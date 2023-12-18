package search

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/d4l-data4life/mex/mex/shared/bi"
	"github.com/d4l-data4life/mex/mex/shared/constants"
	"github.com/d4l-data4life/mex/mex/shared/errstat"
	"github.com/d4l-data4life/mex/mex/shared/known/statuspb"
	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/searchconfig"
	sharedSolr "github.com/d4l-data4life/mex/mex/shared/solr"
	"github.com/d4l-data4life/mex/mex/shared/telemetry"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields/hooks"

	"github.com/d4l-data4life/mex/mex/services/query/endpoints/search/pb"
	"github.com/d4l-data4life/mex/mex/services/query/solr"
)

type Service struct {
	ServiceTag string
	Log        L.Logger

	Redis                 *redis.Client
	Solr                  sharedSolr.ClientAPI
	SolrCollection        string
	TolerantErrorHandling bool

	FieldRepo        fields.FieldRepo
	SearchConfigRepo searchconfig.SearchConfigRepo

	// Field lifecycle hooks
	PostQueryHooks hooks.PostQueryHooks

	TelemetryService *telemetry.Service

	pb.UnimplementedSearchServer
}

// Search handles search queries
func (svc *Service) Search(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error) {
	engineOpts := solr.QueryEngineOptions{
		Log:                   svc.Log,
		FieldRepo:             svc.FieldRepo,
		SearchConfigRepo:      svc.SearchConfigRepo,
		PostQueryHooks:        svc.PostQueryHooks,
		TolerantErrorHandling: svc.TolerantErrorHandling,
	}
	queryOpts := solr.QueryOptions{
		SearchFocusName: request.SearchFocus,
		MaxEditDistance: request.MaxEditDistance,
		UseNgramField:   request.UseNgramField,
	}
	queryEngine, err := solr.QueryEngineFactory(ctx, queryOpts, engineOpts)
	if err != nil {
		svc.Log.Error(ctx, L.Messagef("failed to create query engine: %s", err.Error()))
		return nil, err
	}

	svc.Log.BIEvent(ctx, L.BIActivity("search-request"), L.BIData(bi.SearchRequestInfo{
		QueryLength:          len(request.Query),
		AxisConstraintsCount: len(request.AxisConstraints),
	}))

	// Get ranges for date fields used for year-range faceting
	svc.Log.Info(ctx, L.Message("get ranges of fields used for year-range facets"))
	dateFieldRanges, ignoredRangeQueryErrors, err := svc.getDateRanges(ctx, request, queryEngine)
	if err != nil {
		svc.Log.Error(ctx, L.Messagef("error getting min-max ranges: %s", err.Error()))
		return nil, errstat.MakeMexStatus(errstat.QueryConstructionErrorInternal, fmt.Sprintf("could not get min-max ranges for field with year-range facets: %s", err.Error())).Err()
	}

	// Build and execute main query
	svc.Log.Info(ctx, L.Message("building Solr query"))
	solrQueryBody, queryDiagnostics, err := queryEngine.CreateSolrQuery(ctx, request, dateFieldRanges)
	if err != nil {
		svc.Log.Error(ctx, L.Messagef("error creating main Solr query: %s", err.Error()))
		if queryDiagnostics != nil {
			svc.Log.Error(ctx, L.Messagef("parsing errors: %s", queryDiagnostics.ParsingErrors))
			svc.Log.Error(ctx, L.Messagef("ignored parsing errors: %s", queryDiagnostics.IgnoredErrors))
		}
		return nil, errstat.MakeGRPCStatus(errstat.CodeFrom(err), "could not create Solr query", errstat.Cause(err)).Err()
	}
	svc.Log.Info(ctx, L.Message("executing main Solr query"))
	solrResponse, statusCode, err := svc.Solr.DoJSONQuery(ctx, nil, solrQueryBody)
	if err != nil {
		svc.Log.Error(ctx, L.Messagef("error executing main Solr query: %s", err.Error()))
		return nil, errstat.MakeMexStatus(errstat.SolrQueryFailedInternal, fmt.Sprintf("solr query failed: %s", err.Error())).Err()
	}
	if statusCode != http.StatusOK {
		errMsg := fmt.Sprintf("main solr query failed with status code %d", statusCode)
		extendedErrMsg := svc.getExtendedErrorMsg(errMsg, solrResponse)
		if svc.TolerantErrorHandling {
			// If using relaxed error handling, only report error and still attempt to parse body
			queryDiagnostics.IgnoredErrors = append(queryDiagnostics.IgnoredErrors, solr.MainQueryPartialSolrFailureWarning)
			warnMsg := fmt.Sprintf("Ignoring non-200 Solr response: %s", extendedErrMsg)
			svc.Log.Warn(ctx, L.Message(warnMsg))
		} else {
			// Return error if using strict error handling
			svc.Log.Error(ctx, L.Message(extendedErrMsg))
			return nil, status.Error(codes.Internal, errMsg)
		}
	}
	// Append any ignored errors from the range queries
	queryDiagnostics.IgnoredErrors = append(queryDiagnostics.IgnoredErrors, ignoredRangeQueryErrors...)

	svc.Log.BIEvent(ctx, L.BIActivity("search-response"), L.BIData(bi.SearchResponseInfo{
		ItemsFound: int(solrResponse.Response.NumFound),
	}))

	svc.Log.Info(ctx, L.Message("extracting result from Solr response"))
	return queryEngine.CreateResponse(ctx, solrResponse, request.Facets, queryDiagnostics)
}

// getDateRanges returns the min-max ranges of all datetime fields for which year-range facets was requested
func (svc *Service) getDateRanges(ctx context.Context, request *pb.SearchRequest, queryEngine *solr.QueryEngine) (*sharedSolr.StringFieldRanges, []string, error) {
	noConstraintYearRangeFacets, constrainedYearRangeFacets, err := solr.GetRangeStatRequestFacets(request)
	if err != nil {
		errMsg := fmt.Sprintf("failed to check for year-range facets: %s", err.Error())
		svc.Log.Error(ctx, L.Message(errMsg))
		return nil, nil, errstat.MakeMexStatus(errstat.QueryConstructionErrorInternal, errMsg).Err()
	}
	if len(noConstraintYearRangeFacets) == 0 && len(constrainedYearRangeFacets) == 0 {
		return nil, nil, nil
	}

	/*
		When finding the min-max range of a given field, all axis constraints on it must be ignored since
		the faceting itself will also ignore them (multi-select faceting).
	*/
	var noConstraintDateFieldRanges sharedSolr.StringFieldRanges
	var ignoredErrors []string
	if len(noConstraintYearRangeFacets) > 0 {
		// All facet fields with no axis constraints can be handled with a single query
		svc.Log.Info(ctx, L.Message(fmt.Sprintf("getting year-ranges for %d facets WITHOUT matching axis constraints",
			len(noConstraintYearRangeFacets))))
		var unconstrainedFacets []*sharedSolr.Facet
		// Combine all min-max pairs of facets into a single slice of facets
		for _, facetPair := range noConstraintYearRangeFacets {
			unconstrainedFacets = append(unconstrainedFacets, facetPair...)
		}
		noConstraintDateFieldRanges, ignoredErrors, err = svc.buildAndRunRangeRequest(ctx, request, queryEngine, unconstrainedFacets, "")
		if err != nil {
			svc.Log.Error(ctx, L.Message(fmt.Sprintf("could not build or run range request: %s", err.Error())))
			return nil, nil, err
		}
	}

	// For facet fields with one or more axis constraints, one query per field is needed since constraints differ
	var constrainedDateFieldRanges []sharedSolr.StringFieldRanges
	if len(constrainedYearRangeFacets) > 0 {
		svc.Log.Info(ctx, L.Message(fmt.Sprintf("getting year-ranges for %d facets WITH matching axis constraints",
			len(constrainedYearRangeFacets))))
		for curField, curFacetPair := range constrainedYearRangeFacets {
			// Get min-max range for field from the current facet, ignoring any axis constraints on it
			curDateFieldRanges, addedIgnoredErrors, curErr := svc.buildAndRunRangeRequest(ctx, request, queryEngine, curFacetPair, curField)
			if curErr != nil {
				svc.Log.Error(ctx, L.Message(fmt.Sprintf("could not build or run range request: %s", curErr.Error())))
				return nil, nil, status.Error(codes.Internal, fmt.Sprintf("building or running range request failed: %s", curErr.Error()))
			}
			constrainedDateFieldRanges = append(constrainedDateFieldRanges, curDateFieldRanges)
			ignoredErrors = append(ignoredErrors, addedIgnoredErrors...)
		}
	}

	// Merge all min-max ranges
	svc.Log.Info(ctx, L.Message("combining all ranges"))
	combinedRange, err := solr.CombineRanges(append(constrainedDateFieldRanges, noConstraintDateFieldRanges))
	if err != nil {
		errMsg := fmt.Sprintf("error combining min-max ranges: %s", err.Error())
		svc.Log.Error(ctx, L.Message(errMsg))
		return nil, nil, errstat.MakeMexStatus(errstat.QueryConstructionErrorInternal, errMsg).Err()
	}

	return &combinedRange, ignoredErrors, nil
}

// buildAndRunRangeRequest retrieves the min-max range for the facet axes in a set of year-range facets
func (svc *Service) buildAndRunRangeRequest(ctx context.Context, request *pb.SearchRequest, qe *solr.QueryEngine, yrFacets []*sharedSolr.Facet, ignoreAxis string) (sharedSolr.StringFieldRanges, []string, error) { //nolint:lll
	var yrAxis []string
	countMap := make(map[string]bool)
	for _, fa := range yrFacets {
		if _, ok := countMap[fa.GetAxis()]; !ok {
			yrAxis = append(yrAxis, fa.GetAxis())
			countMap[fa.GetAxis()] = true
		}
	}
	svc.Log.Info(ctx, L.Message(fmt.Sprintf("creating Solr query to get min-max range of axes (%s) for faceting, "+
		"ignoring constraints on axis '%s' (if found)", strings.Join(yrAxis, ", "), ignoreAxis)))
	yearRangeQueryBody, queryDiagnostics, err := qe.CreateYearRangeQuery(ctx, request, yrFacets, ignoreAxis)
	if err != nil {
		svc.Log.Error(ctx, L.Message(fmt.Sprintf("could not create Solr query for date ranges: %s", err.Error())))
		return nil, nil, fmt.Errorf("could not create Solr query for date ranges: %s", err.Error())
	}
	ignoredErrors := queryDiagnostics.IgnoredErrors

	svc.Log.Info(ctx, L.Message("executing Solr query to get min-max range of ordinal axis for faceting"))
	yearRangeSolrResponse, statusCode, err := svc.Solr.DoJSONQuery(ctx, nil, yearRangeQueryBody)
	if err != nil {
		svc.Log.Error(ctx, L.Message(fmt.Sprintf("solr query for date-ranges failed: %s", err.Error())))
		return nil, nil, status.Error(codes.Internal, fmt.Sprintf("solr query for date-ranges failed: %s", err.Error()))
	}
	if statusCode != http.StatusOK {
		errMsg := fmt.Sprintf("solr query for data range failed with status code %d", statusCode)
		extendedErrMsg := svc.getExtendedErrorMsg(errMsg, yearRangeSolrResponse)
		if svc.TolerantErrorHandling {
			// If using relaxed error handling, only report error and still attempt to parse body
			ignoredErrors = append(ignoredErrors, solr.RangeQueryPartialSolrFailureWarning)
			warnMsg := fmt.Sprintf("Ignoring non-200 Solr response: %s", extendedErrMsg)
			svc.Log.Warn(ctx, L.Message(warnMsg))
		} else {
			// Return error if using strict error handling
			svc.Log.Error(ctx, L.Message(extendedErrMsg))
			return nil, nil, status.Error(codes.Internal, errMsg)
		}
	}

	svc.Log.Info(ctx, L.Message("extracting min-max ranges from Solr response"))
	dateRanges, err := solr.GetDateFieldRangesFromResponse(yearRangeSolrResponse, yrFacets)
	if err != nil {
		svc.Log.Error(ctx, L.Message(fmt.Sprintf("failed to get date-ranges from Solr response: %s", err.Error())))
		return nil, nil, errstat.MakeMexStatus(errstat.SolrResponseProcessingInternal, fmt.Sprintf("failed to get date-ranges from Solr response: %s", err.Error())).Err()
	}

	return dateRanges, ignoredErrors, nil
}

// getExtendedErrorMsg appends any error msg returned in the body of the Solr response to another error msg
func (svc *Service) getExtendedErrorMsg(baseMsg string, response *sharedSolr.QueryResponse) string {
	solrMsg := "[none]"
	if msg, ok := response.Error[sharedSolr.ErrMesssageKey].(string); ok {
		solrMsg = msg
	}
	return fmt.Sprintf("%s - Solr msg: %s", baseMsg, solrMsg)
}

// This method makes the Service an rdb.TopicSubscriber
func (svc *Service) Message(ctx context.Context, topic string, configHash string) {
	if !strings.HasSuffix(topic, constants.ConfigUpdateChannelNameSuffix) {
		return
	}

	_ = svc.FieldRepo.Purge(context.Background())
	_ = svc.SearchConfigRepo.Purge(context.Background())

	svc.TelemetryService.SetStatus(statuspb.Color_GREEN, configHash)
}
