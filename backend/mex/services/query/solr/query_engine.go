package solr

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/d4l-data4life/mex/mex/shared/errstat"
	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/searchconfig"
	sharedSearchConfig "github.com/d4l-data4life/mex/mex/shared/searchconfig"
	"github.com/d4l-data4life/mex/mex/shared/searchconfig/sctypes"
	"github.com/d4l-data4life/mex/mex/shared/solr"
	"github.com/d4l-data4life/mex/mex/shared/utils"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields/hooks"
	kindHierarchy "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/hierarchy"
	kindTimestamp "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/timestamp"

	"github.com/d4l-data4life/mex/mex/services/query/endpoints/search/pb"
	"github.com/d4l-data4life/mex/mex/services/query/parser"
)

// Ignored error strings - should be considered enum values
const (
	SkippedFacetWarning                 = "FACET_SKIPPED"
	MainQueryPartialSolrFailureWarning  = "MAIN_QUERY_PARTIAL_SOLR_FAILURE"
	RangeQueryPartialSolrFailureWarning = "RANGE_QUERY_PARTIAL_SOLR_FAILURE"
)

/*
QueryEngine the central entity used for mediating between MEx clients and Solr for searches.
*/
type QueryEngine struct {
	log                   L.Logger
	TolerantErrorHandling bool

	Converter parser.QueryConverter // Converts MEx query string to SOlr query strings

	fieldRepo         fields.FieldRepo              // Access to field configuration info
	searchConfigRepo  searchconfig.SearchConfigRepo // Access to ordinal axes and search foci
	returnFieldMapper *fieldMapper                  // Maps MEx field requested by clients (return fields & highlighting) to the underlying Solr fields
	postQueryHooks    hooks.PostQueryHooks          // Type-specific tasks to be run before returning result to client
}

type QueryOptions struct {
	SearchFocusName string
	MaxEditDistance uint32
	UseNgramField   bool
}

type QueryEngineOptions struct {
	Log L.Logger

	TolerantErrorHandling bool

	FieldRepo        fields.FieldRepo              // Access to field configuration info
	SearchConfigRepo searchconfig.SearchConfigRepo // Access to ordinal axes and search foci
	PostQueryHooks   hooks.PostQueryHooks          // Type-specific tasks to be run before returning result to client
}

func QueryEngineFactory(ctx context.Context, queryOpts QueryOptions, engineOpts QueryEngineOptions) (*QueryEngine, error) {
	searchFocusHooks := sctypes.SearchConfigHooks[solr.MexSearchFocusType]
	effectiveSearchFocusName := GetEffectiveSearchFocus(queryOpts.SearchFocusName)
	_, err := engineOpts.SearchConfigRepo.GetSearchConfigObject(ctx, effectiveSearchFocusName)
	if err != nil {
		errText := fmt.Sprintf("invalid search focus requested ('%s', resolved to '%s')", queryOpts.SearchFocusName, effectiveSearchFocusName)
		return nil, errstat.MakeMexStatus(errstat.InvalidClientQuery, errText).Err()
	}
	matchingOpConfig, err := searchFocusHooks.GetMatchingOpsConfig(effectiveSearchFocusName, queryOpts.MaxEditDistance, queryOpts.UseNgramField)
	if err != nil {
		return nil, errstat.MakeMexStatus(errstat.InvalidClientQuery, fmt.Sprintf("invalid request: %s", err.Error())).Err()
	}
	listener, err := parser.NewListener(matchingOpConfig)
	if err != nil {
		return nil, errstat.MakeMexStatus(errstat.QueryEngineCreationFailedInternal, fmt.Sprintf("could not create query parser: %s", err.Error())).Err()
	}
	queryConverter := parser.NewMexParser(listener)
	queryEngine, err := newQueryEngine(queryConverter, engineOpts)
	if err != nil {
		return nil, errstat.MakeMexStatus(errstat.QueryEngineCreationFailedInternal, fmt.Sprintf("could not create Solr query engine: %s", err.Error())).Err()
	}
	return queryEngine, err
}

func newQueryEngine(converter parser.QueryConverter, opts QueryEngineOptions) (*QueryEngine, error) {
	mapper, err := newFieldMapper(context.Background(), opts.FieldRepo)
	if err != nil {
		return nil, err
	}
	return &QueryEngine{
		log:                   opts.Log,
		TolerantErrorHandling: opts.TolerantErrorHandling,
		Converter:             converter,

		fieldRepo:         opts.FieldRepo,
		searchConfigRepo:  opts.SearchConfigRepo,
		returnFieldMapper: mapper,
		postQueryHooks:    opts.PostQueryHooks,
	}, nil
}

// CreateSolrQuery turns a MEx search request into the corresponding Solr request
func (qe *QueryEngine) CreateSolrQuery(ctx context.Context, searchRequest *pb.SearchRequest, dateFieldRanges *solr.StringFieldRanges) (*solr.QueryBody, *solr.Diagnostics, error) {
	mexFieldToMexKindMap := make(map[string]string)
	fieldsDefs, fieldErr := qe.fieldRepo.ListFieldDefs(ctx)
	if fieldErr != nil {
		return nil, &solr.Diagnostics{}, errstat.MakeMexStatus(errstat.InvalidConfigurationClient,
			"could not retrieve field configuration").Err()
	}
	for _, fd := range fieldsDefs {
		mexFieldToMexKindMap[fd.Name()] = fd.Kind()
	}
	queryBody := getBaseQueryBody()
	setPaging(queryBody, searchRequest)
	diagnostics, isPhraseOnlyQuery, queryErr := qe.setQuery(queryBody, searchRequest, qe.Converter)
	if queryErr != nil {
		return nil, diagnostics, queryErr
	}
	solrFieldTags, constraintErr := qe.setConstraints(ctx, queryBody, searchRequest)
	if constraintErr != nil {
		return nil, nil, constraintErr
	}
	sortErr := qe.setSorting(ctx, queryBody, searchRequest)
	if sortErr != nil {
		return nil, nil, sortErr
	}
	fieldsErr := qe.setFields(queryBody, searchRequest)
	if fieldsErr != nil {
		return nil, nil, fieldsErr
	}
	ignoredErrors, facetErr := qe.setFacets(ctx, queryBody, searchRequest, dateFieldRanges, solrFieldTags, mexFieldToMexKindMap)
	if facetErr != nil {
		return nil, nil, facetErr
	}
	diagnostics.IgnoredErrors = append(diagnostics.IgnoredErrors, ignoredErrors...)
	highlightErr := qe.setHighlighting(ctx, queryBody, searchRequest, isPhraseOnlyQuery)
	if highlightErr != nil {
		return nil, nil, highlightErr
	}

	return queryBody, diagnostics, nil
}

// getBaseQueryBody returns a base Solr query body with the fixed parameters set
func getBaseQueryBody() *solr.QueryBody {
	queryBody := solr.QueryBody{}
	// Add parameters that would normally go in query URL
	queryBody.Params = solr.ParamObj{
		OmitHeader: true,
		QOp:        "AND",
		DefType:    "edismax",
	}
	return &queryBody
}

// setPaging adds paging information to the passed Solr search request
func setPaging(queryBody *solr.QueryBody, searchRequest *pb.SearchRequest) {
	// Cap limit at max allowed value
	limit := math.Min(float64(searchRequest.GetLimit()), solr.MaxDocLimit)
	queryBody.Limit = uint32(limit)
	// Here, the Go default (0) happens to be correct
	queryBody.Offset = searchRequest.GetOffset()
}

// setPaging adds sorting information to the passed Solr search request
func (qe *QueryEngine) setSorting(ctx context.Context, queryBody *solr.QueryBody, searchRequest *pb.SearchRequest) error {
	if searchRequest.GetSorting() != nil {
		requestedAxis := searchRequest.GetSorting().GetAxis()
		if requestedAxis == "" {
			return errstat.MakeMexStatus(errstat.InvalidClientQuery, "no ordinal axis given for sorting").Err()
		}
		sortOrder := "asc"
		requestedOrder := searchRequest.GetSorting().GetOrder()
		if requestedOrder != "" {
			if !utils.Contains(solr.AllowedSortOrders, requestedOrder) {
				return errstat.MakeMexStatus(
					errstat.InvalidClientQuery,
					"the sort order must be 'asc' or 'desc'",
					errstat.DevMessagef("invalid the sort order: '%s' (must be 'asc' or 'desc')", requestedOrder),
				).Err()
			}
			sortOrder = requestedOrder
		}

		// Check if requested ordinal axis is valid
		_, err := qe.searchConfigRepo.GetFieldsForAxis(ctx, requestedAxis)
		if err != nil {
			return errstat.MakeMexStatus(
				errstat.InvalidClientQuery,
				"requested ordinal axis is not configured").Err()
		}
		// Get the Solr sort backing field for the chosen axis
		solrBackingFieldName := solr.GetOrdinalAxisSortFieldName(requestedAxis)
		queryBody.Sort = fmt.Sprintf("%s %s", solrBackingFieldName, sortOrder)
	}
	return nil
}

// setFields adds information about the fields to return to the passed Solr search request
func (qe *QueryEngine) setFields(queryBody *solr.QueryBody, searchRequest *pb.SearchRequest) error {
	// Set the fields requested from Solr
	var fieldsRequested []string
	if len(searchRequest.GetFields()) > 0 {
		for _, fn := range searchRequest.GetFields() {
			// We return data from all underlying Solr backing field configured for mapping
			res, err := qe.returnFieldMapper.getReturnBackingFieldNamesFromMexName(fn)
			if err != nil {
				// Should really not happen since we should have fields for all supported languages
				return fmt.Errorf("could not set requested fields: %s", err.Error())
			}
			fieldsRequested = append(fieldsRequested, res...)
		}
	}
	// Always include technical ID
	if !utils.Contains(fieldsRequested, solr.DefaultUniqueKey) {
		fieldsRequested = append(fieldsRequested, solr.DefaultUniqueKey)
	}
	// Always include entity SolrName
	if !utils.Contains(fieldsRequested, solr.ItemEntityNameField) {
		fieldsRequested = append(fieldsRequested, solr.ItemEntityNameField)
	}

	queryBody.Fields = fieldsRequested
	return nil
}

// setFacets adds faceting information to the passed Solr search request
func (qe *QueryEngine) setFacets(ctx context.Context, queryBody *solr.QueryBody, req *pb.SearchRequest,
	ranges *solr.StringFieldRanges, solrFieldTags map[string][]string, mexFieldToMexKindMap map[string]string,
) ([]string, error) {
	// Add facet information requested from Solr
	var ignoredErrors []string
	if len(req.GetFacets()) > 0 {
		queryBody.Facet = map[string]solr.SolrFacet{}
		facetNamesSeen := map[string]bool{}
		var valErr error
		for _, f := range req.GetFacets() {
			mexAxisFields, axisErr := qe.searchConfigRepo.GetFieldsForAxis(ctx, f.GetAxis())
			if axisErr != nil {
				errorStr := "ordinal axis used for faceting is not configured"
				if qe.TolerantErrorHandling {
					// Skip facet
					warnStr := fmt.Sprintf("Skpping facet: %s", errorStr)
					qe.log.Warn(ctx, L.Message(warnStr))
					ignoredErrors = append(ignoredErrors, SkippedFacetWarning)
					continue
				}
				return nil, errstat.MakeMexStatus(errstat.InvalidClientQuery, errorStr).Err()
			}
			facetNamesSeen, valErr = validateFacet(f, mexAxisFields, mexFieldToMexKindMap, facetNamesSeen)
			if valErr != nil {
				errorStr := fmt.Sprintf("invalid facet found: %s", valErr.Error())
				if qe.TolerantErrorHandling {
					// Skip facet
					warnStr := fmt.Sprintf("Skpping facet: %s", errorStr)
					qe.log.Warn(ctx, L.Message(warnStr))
					ignoredErrors = append(ignoredErrors, SkippedFacetWarning)
					continue
				}
				return nil, errstat.MakeMexStatus(errstat.InvalidConfigurationClient, errorStr).Err()
			}
			facetName := createFacetName(f.GetAxis(), f.GetStatName())
			curFacet, facetErr := qe.getFacetForField(ctx, f, ranges, solrFieldTags)
			if facetErr != nil {
				errorStr := fmt.Sprintf("could not created requested facet on ordinal axis: %s", facetErr.Error())
				if qe.TolerantErrorHandling {
					// Skip facet
					warnStr := fmt.Sprintf("Skpping facet: %s", errorStr)
					qe.log.Warn(ctx, L.Message(warnStr))
					ignoredErrors = append(ignoredErrors, SkippedFacetWarning)
					continue
				}
				return nil, errstat.MakeMexStatus(errstat.InvalidClientQuery, errorStr).Err()
			}
			if curFacet != nil {
				queryBody.Facet[facetName] = *curFacet
			}
		}
	}
	return ignoredErrors, nil
}

/*
getFacetForField generates the Solr facet settings for the different facet types. The returned data type (SolrFacet)
is a sort of "union-type" in that it contains fields for all Solr facet types. Only the required ones are filled
and custom JSON serialization for the SolrFacetSet-type ensures that (depending on facet type) only the correct
fields are serialized to JSON.
*/
func (qe *QueryEngine) getFacetForField(ctx context.Context, f *solr.Facet, ranges *solr.StringFieldRanges, tags map[string][]string) (*solr.SolrFacet, error) {
	solrFieldName := solr.GetOrdinalAxisFacetAndFilterFieldName(f.GetAxis())
	var curFacet *solr.SolrFacet
	switch f.GetType() {
	case solr.MexExactFacetType:
		curFacet = &solr.SolrFacet{
			DetailedType: solr.SolrTermsFacetType,
			NumBuckets:   true,
			Field:        solrFieldName,
		}
		if f.GetLimit() > 0 {
			curFacet.Limit = uint32(math.Min(float64(f.GetLimit()), solr.MaxFacetLimit))
		}
		if f.GetOffset() > 0 {
			curFacet.Offset = f.GetOffset()
		}
		if tags, ok := tags[f.GetAxis()]; ok {
			curFacet.ExcludeTags = tags
		}
	case solr.MexYearRangeFacetType:
		dateRange, dateOk := (*ranges)[f.GetAxis()]
		if !dateOk {
			// If there is no date range for a year-range axis,
			// we take it to mean that there are no hits and drop the facet
			qe.log.Info(ctx, L.Messagef("no date range available for ordinal axis '%s' - dropping corresponding facet",
				f.GetAxis()))
			return nil, nil
		}
		// To do a range facet, we first need to get the min and max years
		startDate, endDate, dateErr := getFacetingYearRange(dateRange)
		if dateErr != nil {
			return nil, fmt.Errorf("could not determine date range for faceting from the dates")
		}

		curFacet = &solr.SolrFacet{
			DetailedType: solr.SolrStringRangeFacetType,
			Field:        solrFieldName,
			StartString:  startDate,
			EndString:    endDate,
			GapString:    "+1YEARS",
		}
		if tags, ok := tags[f.GetAxis()]; ok {
			curFacet.ExcludeTags = tags
		}
	case solr.MexStringStatFacetType:
		curFacet = &solr.SolrFacet{
			DetailedType: solr.SolrStringStatFacetType,
			Field:        solrFieldName,
			StatOp:       f.GetStatOp(),
		}
	default:
		return nil, fmt.Errorf("received invalid MEx facet type")
	}
	return curFacet, nil
}

// validateFacet checks a requested facet for validity,
// adding it to the passed facetNamesSeen array if it is a valid string facet
func validateFacet(facet *solr.Facet, mexOrdinalAxisFields []string, mexFieldToMexKindMap map[string]string, facetNamesSeen map[string]bool) (map[string]bool, error) {
	switch facet.GetType() {
	case solr.MexExactFacetType:
		// No special checks
	case solr.MexYearRangeFacetType:
		axisType, err := sctypes.GetOrdinalAxisFieldType(mexOrdinalAxisFields, mexFieldToMexKindMap)
		if err != nil {
			return nil, fmt.Errorf("could not get ordinal axis configuration for the request axis")
		}
		if axisType != solr.DefaultSolrTimestampFieldType {
			return facetNamesSeen, fmt.Errorf("facets of type '%s' are only possible for ordinal axis of underlying kind '%s'", solr.MexYearRangeFacetType, kindTimestamp.KindName)
		}
	case solr.MexStringStatFacetType:
		axisType, err := sctypes.GetOrdinalAxisFieldType(mexOrdinalAxisFields, mexFieldToMexKindMap)
		if err != nil {
			return facetNamesSeen, fmt.Errorf("ordinal axis configuration for the request axis is invalid")
		}
		if axisType != solr.DefaultSolrTimestampFieldType && axisType != solr.DefaultSolrStringFieldType {
			return facetNamesSeen, fmt.Errorf("facets of type '%s' are only possible for ordinal axis with underlying types '%s'",
				solr.MexStringStatFacetType,
				strings.Join([]string{solr.DefaultSolrTimestampFieldType, solr.DefaultSolrStringFieldType}, " & "))
		}
		if facet.GetStatOp() != solr.MinOperator && facet.GetStatOp() != solr.MaxOperator {
			return facetNamesSeen, fmt.Errorf("a stat operator for facetting was not recognized")
		}
		if facet.GetStatName() == "" {
			return facetNamesSeen, fmt.Errorf("no facet name given for facet on the ordinal axis")
		}
		if _, ok := facetNamesSeen[facet.GetStatName()]; ok {
			return facetNamesSeen, fmt.Errorf("a stat facet name was used multiple time")
		}
		facetNamesSeen[facet.GetStatName()] = true
	default:
		return facetNamesSeen, fmt.Errorf("requested an undefined facet type")
	}
	return facetNamesSeen, nil
}

// setHighlighting adds highlighting information to the passed Solr search request
func (qe *QueryEngine) setHighlighting(ctx context.Context, queryBody *solr.QueryBody, searchRequest *pb.SearchRequest, isPhraseOnlyQuery bool) error {
	var err error
	var mexHighlightFields []string
	// Get fields to highlight
	switch {
	case searchRequest.GetAutoHighlight():
		requestedSearchFocus := GetEffectiveSearchFocus(searchRequest.SearchFocus)
		mexHighlightFields, err = qe.searchConfigRepo.GetFieldsForSearchFocus(ctx, requestedSearchFocus)
		if err != nil {
			return err
		}
	case len(searchRequest.GetHighlightFields()) > 0:
		mexHighlightFields = searchRequest.GetHighlightFields()
	}
	if len(mexHighlightFields) == 0 {
		return nil
	}

	// Map MEx highlight fields to underlying Solr fields
	var mappedNames []string
	var highlightFields []string
	for _, fn := range mexHighlightFields {
		mappedNames, err = qe.returnFieldMapper.getHighlightBackingFieldNamesFromMexName(fn, isPhraseOnlyQuery)
		if err != nil {
			// Should really not happen since we should have fields for all supported languages
			return fmt.Errorf("could not set requested highlight fields: %s", err.Error())
		}
		highlightFields = append(highlightFields, mappedNames...)
	}
	// If any highlighting is needed, set params to switch it on
	if len(highlightFields) > 0 {
		setHighlightingParams(queryBody, highlightFields)
	}

	return nil
}

// GetEffectiveSearchFocus returns the requested search focus if set, otherwise the default
func GetEffectiveSearchFocus(requestedSearchFocus string) string {
	effectiveSearchFocus := requestedSearchFocus
	if effectiveSearchFocus == "" {
		effectiveSearchFocus = solr.MexDefaultSearchFocusName
	}
	return effectiveSearchFocus
}

func setHighlightingParams(queryBody *solr.QueryBody, highlightFields []string) {
	queryBody.Params.Hl = true
	queryBody.Params.HlFl = strings.Join(highlightFields, ",")
	// Default highlighting parameters are set
	queryBody.Params.HlMethod = solr.HighlightAlgorithm
	queryBody.Params.HlSnippets = solr.HighlightSnippetNo
	queryBody.Params.HlFragsize = solr.HighlightFragmentSize
	// The used tags around highlighted terms are set.
	// Here unicode code points of the 'Private Use Area' are used to allow unambiguous identification in front end.
	queryBody.Params.HlTagPre = solr.HighlightStartChar
	queryBody.Params.HlTagPost = solr.HighlightStopChar
}

// setQuery generates the full Solr query string
func (qe *QueryEngine) setQuery(queryBody *solr.QueryBody, searchRequest *pb.SearchRequest, converter parser.QueryConverter) (*solr.Diagnostics, bool, error) {
	var diagnostics *solr.Diagnostics
	queryParseResult, err := converter.ConvertToSolrQuery(searchRequest.GetQuery())
	if err != nil {
		switch err.GetType() {
		case parser.ParserErrorType:
			// Parsing error - capture diagnostic information and return error
			diagnostics = &solr.Diagnostics{
				ParsingSucceeded: false,
				ParsingErrors:    err.GetMessages(),
			}
			return diagnostics, false, errstat.MakeMexStatus(errstat.InvalidClientQuery, err.Error()).Err()
		case parser.QueryConstructionErrorType:
			// Our own logic failed - return error
			return nil, false, errstat.MakeMexStatus(errstat.QueryConstructionErrorInternal, err.Error()).Err()
		case parser.InvalidArgumentErrorType:
			return nil, false, errstat.MakeMexStatus(errstat.InvalidClientQuery, err.Error()).Err()
		default:
			// Unknown error - return error
			return nil, false, errstat.MakeMexStatus(errstat.OtherError, err.Error()).Err()
		}
	} else {
		diagnostics = &solr.Diagnostics{
			ParsingSucceeded: true,
			CleanedQuery:     queryParseResult.CleanedQuery,
			QueryWasCleaned:  queryParseResult.QueryWasCleaned,
		}
	}

	queryBody.Query = queryParseResult.SolrQuery

	return diagnostics, queryParseResult.PhrasesOnlyQuery, nil
}

func (qe *QueryEngine) setConstraints(ctx context.Context, queryBody *solr.QueryBody, searchRequest *pb.SearchRequest) (map[string][]string, error) {
	if len(searchRequest.GetAxisConstraints()) == 0 {
		return nil, nil
	}

	solrFieldTags := make(map[string][]string)
	// Each entry in axisConstraintClauses represents a constraint on a single ordinal axis
	var axisConstraintClauses []string
	for _, constraint := range searchRequest.GetAxisConstraints() {
		axisName := constraint.GetAxis()
		axisConfig, axisErr := qe.searchConfigRepo.GetSearchConfigObject(ctx, axisName)
		if axisErr != nil {
			return nil, axisErr
		}
		if axisConfig.Type != solr.MexOrdinalAxisType && axisConfig.Type != solr.MexHierarchyAxisType {
			return nil, fmt.Errorf("invalid axis name in constraint")
		}
		valConstraintClause, tag, fieldConstErr := getConstraintClauseForAxis(constraint, axisConfig)
		if fieldConstErr != nil {
			return nil, fieldConstErr
		}
		axisConstraintClauses = append(axisConstraintClauses, valConstraintClause)
		solrFieldTags[axisName] = []string{tag}
	}
	queryBody.Filter = axisConstraintClauses
	return solrFieldTags, nil
}

// getConstraintClauseForAxis turns the constraint for a single axis into a single Solr search clause
func getConstraintClauseForAxis(constraint *solr.AxisConstraint, axisConfig *sharedSearchConfig.SearchConfigObject) (string, string, error) {
	var valConstraintsForAxis []string
	axisBackingFieldName := solr.GetOrdinalAxisFacetAndFilterFieldName(axisConfig.Name)
	targetFieldName := parser.SanitizeTerm(axisBackingFieldName)
	switch constraint.GetType() {
	case solr.MexExactAxisConstraint:
		if axisConfig.Type == solr.MexOrdinalAxisType && len(constraint.GetValues()) == 0 {
			return "", "", errstat.MakeMexStatus(errstat.InvalidClientQuery,
				"no constraint values given for exact axis constraint on ordinal axis").Err()
		}
		if axisConfig.Type == solr.MexHierarchyAxisType && len(constraint.GetValues()) == 0 && len(constraint.GetSingleNodeValues()) == 0 {
			return "", "", errstat.MakeMexStatus(errstat.InvalidClientQuery,
				"no constraint values (sub-tree or single-node) given for exact axis constraint on hierarchy axis").Err()
		}
		for _, val := range constraint.GetValues() {
			escapedVal := parser.SanitizeTerm(val)
			valConstraintsForAxis = append(valConstraintsForAxis, fmt.Sprintf(`%s:"%s"`, targetFieldName,
				escapedVal))
		}
		// Handle single-node constraints for hierarchy axis
		if axisConfig.Type == solr.MexHierarchyAxisType {
			singleNodeTargetFieldName := solr.GetSingleNodeAxisFieldName(targetFieldName)
			for _, val := range constraint.GetSingleNodeValues() {
				escapedVal := parser.SanitizeTerm(val)
				valConstraintsForAxis = append(valConstraintsForAxis, fmt.Sprintf(`%s:"%s"`, singleNodeTargetFieldName,
					escapedVal))
			}
		}
	case solr.MexStringRangeConstraint:
		if len(constraint.GetStringRanges()) == 0 {
			return "", "", errstat.MakeMexStatus(errstat.InvalidClientQuery,
				"no constraint values given for string range axis constraint on ordinal axis").Err()
		}
		for _, strRange := range constraint.GetStringRanges() {
			strRangeConstraint, strRangeErr := getStringRangeConstraint(strRange, targetFieldName)
			if strRangeErr != nil {
				return "", "", errstat.MakeGRPCStatus(errstat.CodeFrom(strRangeErr), "error creating string range constraint", errstat.Cause(strRangeErr)).Err()
			}
			valConstraintsForAxis = append(valConstraintsForAxis, strRangeConstraint)
		}
	default:
		return "", "", errstat.MakeMexStatus(errstat.InvalidConfigurationClient,
			"an axis constraint on an ordinal axis has an unknown type").Err()
	}
	// Different values for a single field are either ORed (default) or ANDed together
	combOp := solr.OrSeparator
	if len(constraint.GetCombineOperator()) > 0 {
		separatorForOp := map[string]string{
			"and": solr.AndSeparator,
			"or":  solr.OrSeparator,
		}
		op, ok := separatorForOp[strings.ToLower(constraint.GetCombineOperator())]
		if !ok {
			return "", "", errstat.MakeMexStatus(errstat.InvalidClientQuery,
				fmt.Sprintf("unsupported Boolean operator used for combining value-constraints on"+
					" fields")).Err()
		}
		combOp = op
	}
	valConstraintClause := strings.Join(valConstraintsForAxis, combOp)
	// Add tag to allow excluding this constraint for faceting
	valConstraintClause, tag := solr.TagExpr(valConstraintClause, axisConfig.Name)
	return valConstraintClause, tag, nil
}

// getStringRangeConstraint
func getStringRangeConstraint(strRange *solr.StringRange, escapedField string) (string, error) {
	if strRange.Min == "" && strRange.Max == "" {
		return "", errstat.MakeMexStatus(errstat.InvalidClientQuery, "string range must have a min or max value (or both)").Err()
	}
	lowerLimit := "*"
	upperLimit := "*"
	if strRange.Min != "" {
		lowerLimit = fmt.Sprintf(`"%s"`, strRange.Min)
	}
	if strRange.Max != "" {
		upperLimit = fmt.Sprintf(`"%s"`, strRange.Max)
	}
	return fmt.Sprintf(`%s:[%s TO %s]`, escapedField, lowerLimit, upperLimit), nil
}

// getFacetingYearRange calculates the whole-year faceting range based on the actual max and min dates in the data
func getFacetingYearRange(dateRange *solr.StringRange) (string, string, error) {
	// NOTE: We do not check that the passed datetime is actually valid
	yearMatcher := regexp.MustCompile(`^(\d{4})-`)
	startYear := yearMatcher.FindString(dateRange.Min)
	endYear := yearMatcher.FindString(dateRange.Max)
	if startYear == "" || endYear == "" {
		return "", "", fmt.Errorf("unable to identify start and end years from date strings")
	}
	minDate := fmt.Sprintf("%s01-01T00:00:00.000Z", startYear)
	maxDate := fmt.Sprintf("%s12-31T23:59:59.999Z", endYear)

	return minDate, maxDate, nil
}

/*
GetRangeStatRequestFacets extracts the year-range facets from a request and creates corresponding facets
for getting the min-max range, dividing the result into two groups depending on whether there is an
axis constraint on the relevant field or not. The return values are maps from field names to 2-element
arrays containing the min and max facets for that field.
*/
func GetRangeStatRequestFacets(searchRequest *pb.SearchRequest) (map[string][]*solr.Facet, map[string][]*solr.Facet, error) {
	noConstraintFacets := make(map[string][]*solr.Facet)
	constrainedFacets := make(map[string][]*solr.Facet)
	if len(searchRequest.GetFacets()) > 0 {
		axisConstraintFields := map[string]bool{}
		for _, c := range searchRequest.GetAxisConstraints() {
			if _, ok := axisConstraintFields[c.GetAxis()]; !ok {
				axisConstraintFields[c.GetAxis()] = true
			}
		}
		for _, facet := range searchRequest.GetFacets() {
			curAxis := facet.GetAxis()
			if curAxis == "" {
				return nil, nil, fmt.Errorf("found facet with no axis specified")
			}
			if facet.GetType() == solr.MexYearRangeFacetType {
				curMinMaxFacets := make([]*solr.Facet, 2)
				for i, op := range [2]string{solr.MinOperator, solr.MaxOperator} {
					facetName, nameErr := getStatNameForAxisAndOp(curAxis, op)
					if nameErr != nil {
						return nil, nil, fmt.Errorf("error generating facet name: %s", nameErr.Error())
					}
					curMinMaxFacets[i] = &solr.Facet{
						Type:     solr.MexStringStatFacetType,
						Axis:     curAxis,
						StatName: facetName,
						StatOp:   op,
					}
				}
				if _, ok := axisConstraintFields[curAxis]; ok {
					constrainedFacets[curAxis] = curMinMaxFacets
				} else {
					noConstraintFacets[curAxis] = curMinMaxFacets
				}
			}
		}
	}
	return noConstraintFacets, constrainedFacets, nil
}

/*
CreateYearRangeQuery returns a Solr query with based on the passed request in which (1) the same query and facet
constraints are applied EXCEPT axis constraints on ignoreField, and (2) just the min & max facets for the listed
axes are returned.
*/
func (qe *QueryEngine) CreateYearRangeQuery(ctx context.Context, searchRequest *pb.SearchRequest, rangeFacets []*solr.Facet, ignoreAxis string,
) (*solr.QueryBody, *solr.Diagnostics, error) {
	if len(rangeFacets) == 0 {
		return nil, nil, fmt.Errorf("cannot construct year-range query from request with no year-range facets")
	}
	// Build a new request with the same constraints, but the new facets
	adaptedAxisConstraints := cloneAxisConstraints(searchRequest.GetAxisConstraints(), ignoreAxis)
	rangeRequest := &pb.SearchRequest{
		Query:           searchRequest.GetQuery(),
		Limit:           0,
		Facets:          rangeFacets,
		AxisConstraints: adaptedAxisConstraints,
	}
	// Get the Solr request for this new query
	queryBody, queryDiagnostics, queryErr := qe.CreateSolrQuery(ctx, rangeRequest, nil)
	return queryBody, queryDiagnostics, queryErr
}

/*
cloneAxisConstraints deep-clones an array of AxisConstraint pointers, dropping constraints on the
specified field. Passing an empty string for dropField means that no fields should be ignored.
*/
func cloneAxisConstraints(constraints []*solr.AxisConstraint, dropAxis string) []*solr.AxisConstraint {
	if len(constraints) == 0 {
		return nil
	}
	var clonedConstraints []*solr.AxisConstraint
	for _, c := range constraints {
		// Drop constraint on the axis that should be ignored ("" implies copying all constraints)
		if dropAxis != "" && c.GetAxis() == dropAxis {
			continue
		}
		var clonedStringRanges []*solr.StringRange
		if len(c.GetStringRanges()) > 0 {
			clonedStringRanges = []*solr.StringRange{}
			for _, r := range c.GetStringRanges() {
				clonedStringRanges = append(clonedStringRanges, &solr.StringRange{
					Min: r.Min,
					Max: r.Max,
				})
			}
		}
		newConstraint := &solr.AxisConstraint{
			Type:             c.GetType(),
			Axis:             c.GetAxis(),
			Values:           c.GetValues(),
			SingleNodeValues: c.GetSingleNodeValues(),
			StringRanges:     clonedStringRanges,
			CombineOperator:  c.GetCombineOperator(),
		}
		clonedConstraints = append(clonedConstraints, newConstraint)
	}
	return clonedConstraints
}

/*
CreateResponse converts a Solr response into a MEx core service response.

Although the Solr response in principle contains all info needed to extract the facets,
we pass the fields we faceted by to be able to pick them out. In the JSON returned by Solr,
the facet information is contained in an object that has a "count" property (giving the full
count) and then one property for each facet, with the property SolrName being the SolrName given
to the corresponding facet in the query (in MEx always set to the SolrName of the facet return).
What is more, if one names a facet "count", the Solr facet response JSON will contain
*two* properties called "count", one being the overall count and the other being the
information for the facet with that SolrName.
*/
func (qe *QueryEngine) CreateResponse(ctx context.Context, solrResponse *solr.QueryResponse, facets []*solr.Facet, diagnostics *solr.Diagnostics,
) (*pb.SearchResponse, error) {
	if solrResponse == nil {
		return nil, errstat.MakeMexStatus(errstat.SolrResponseProcessingInternal, "Solr response object is nil").Err()
	}

	response := pb.NewSearchResponse()
	response.Diagnostics = diagnostics
	response.NumFound = solrResponse.Response.NumFound
	response.NumFoundExact = solrResponse.Response.NumFoundExact
	response.Start = solrResponse.Response.Start
	response.MaxScore = solrResponse.Response.MaxScore

	// Format search matches
	for _, doc := range solrResponse.Response.Docs {
		docItem, docErr := qe.makeDocItem(doc)
		if docErr != nil {
			return nil, docErr
		}
		if docItem != nil {
			response.AddDocItem(docItem)
		}
	}

	// Format facets
	collectedFacets, facetErr := qe.getResponseFacets(ctx, facets, solrResponse)
	if facetErr != nil {
		return nil, errstat.MakeMexStatus(errstat.SolrResponseProcessingInternal,
			fmt.Sprintf("could not parse Solr facet response: %s", facetErr.Error())).Err()
	}
	if len(collectedFacets) > 0 {
		response.Facets = collectedFacets
	}

	// Format highlight information
	for docID, highlightInfo := range solrResponse.Highlighting {
		highlight, highlightErr := qe.makeHighlight(docID, highlightInfo)
		if highlightErr != nil {
			return nil, errstat.MakeMexStatus(errstat.SolrResponseProcessingInternal,
				fmt.Sprintf("could not parse Solr highlight response: %s", highlightErr.Error())).Err()
		}
		if highlight != nil {
			response.Highlights = append(response.Highlights, highlight)
		}
	}
	if len(response.Highlights) > 0 {
		// Sort by item ID
		sort.Slice(response.Highlights, func(i, j int) bool {
			return response.Highlights[i].ItemId < response.Highlights[j].ItemId
		})
	}

	return response, nil
}

// getResponseFacets maps facets returned by Solr to the MEx facet return format
func (qe *QueryEngine) getResponseFacets(ctx context.Context, facets []*solr.Facet, solrResponse *solr.QueryResponse) ([]*solr.FacetResult, error) {
	// Format facet information
	var collectedFacets []*solr.FacetResult
	for _, curFacet := range facets {
		qe.log.Info(ctx, L.Messagef("facet field: %s", curFacet.Axis), L.Phase("query"))

		facetName := createFacetName(curFacet.Axis, curFacet.GetStatName())
		// codeSystemName := getCodeSystemIfAny(fieldConfigs, curFacet.Field)

		facet := &solr.FacetResult{
			Type: curFacet.GetType(),
			Axis: curFacet.GetAxis(),
		}
		if facetResponse, ok := solrResponse.Facets[facetName]; ok {
			if curFacet.GetType() == solr.MexStringStatFacetType {
				// For stat facets, the returned facet result is simply a string, not an object
				strVal, valErr := getStringResultFromFacet(map[string]interface{}{
					facetName: facetResponse,
				}, curFacet)
				if valErr != nil {
					return nil, valErr
				}
				facet.StatName = facetName
				facet.StringStatResult = strVal
			} else {
				// For non-stat facets, the result is an object
				var facetErr error
				facet, facetErr = qe.makeFacet(ctx, curFacet, facetResponse)
				if facetErr != nil {
					return nil, status.Error(codes.Internal, "could not parse Solr facet response: "+facetErr.Error())
				}
			}
		} else {
			// No Solr response for this facet --> return empty facet
			if curFacet.GetType() == solr.MexStringStatFacetType {
				facet.StatName = facetName
				facet.StringStatResult = ""
			} else {
				facet.BucketNo = 0
				facet.Buckets = []*solr.FacetBucket{}
			}
		}
		if facet != nil {
			collectedFacets = append(collectedFacets, facet)
		}
	}
	return collectedFacets, nil
}

// makeDocItem converts a generic JSON object representing a Solr document into an array of key-value pairs
func (qe *QueryEngine) makeDocItem(solrDoc solr.GenericObject) (*solr.DocItem, error) {
	// Safe access to ID
	id, ok := solrDoc[solr.DefaultUniqueKey].(string)
	// Documents without ID are dropped
	if !ok {
		return nil, nil
	}

	docItem := solr.NewDocItem(id)
	// Safe access to entity type
	entity, ok := solrDoc[solr.ItemEntityNameField].(string)
	if ok {
		docItem.EntityType = entity
	}

	// Reduce the object to key-value pairs with all values strings
	for k, v := range solrDoc {
		mappedFieldInfo := qe.returnFieldMapper.getMexNameFromBackingFieldName(k)
		// ID and entity type are handled separately
		if mappedFieldInfo.Name == solr.DefaultUniqueKey || mappedFieldInfo.Name == solr.ItemEntityNameField {
			continue
		}
		switch s := v.(type) {
		case []interface{}:
			// Arrays elements are split into separate key-value pairs
			for _, arrVal := range s {
				// Everything is turned into a string, even if it is an array or an object
				docItem.AddValue(mappedFieldInfo, fmt.Sprintf("%v", arrVal))
			}
		default:
			// All non-array elements are turned into a string, even if it is an object
			docItem.AddValue(mappedFieldInfo, fmt.Sprintf("%v", s))
		}
	}

	sort.Sort(solr.ByDocValue(docItem.Values))

	return docItem, nil
}

// makeFacet translates the facet information returned by Solr into a MEx facet response
func (qe *QueryEngine) makeFacet(ctx context.Context, facet *solr.Facet, solrFacetResult interface{}) (*solr.FacetResult, error) {
	axisName := facet.GetAxis()
	returnFacet := solr.NewFacetResult(facet.GetType(), axisName)

	// TODO: The manual JSON parsing below could perhaps be simplified using serialization to objects
	// partly) replaced by de-serialization into a new facet result type
	typedFacet, facetOk := solrFacetResult.(map[string]interface{})
	if !facetOk {
		return nil, fmt.Errorf("could not parse Solr facet")
	}
	bucketNo, bucketNoErr := getBucketNo(typedFacet, facet.GetType())
	if bucketNoErr != nil {
		return nil, bucketNoErr
	}
	returnFacet.BucketNo = bucketNo
	buckets, bucketOk := typedFacet["buckets"].([]interface{})
	if !bucketOk {
		return nil, fmt.Errorf("could not parse Solr facet bucket for facet axis")
	}

	ordinalAxisConfigs, err := qe.searchConfigRepo.ListSearchConfigsOfType(ctx, solr.MexOrdinalAxisType)
	if err != nil {
		return nil, fmt.Errorf("failed to get ordinal axis configs")
	}
	var currentAxisConfig *sharedSearchConfig.SearchConfigObject
	for _, oaC := range ordinalAxisConfigs.GetSearchConfigs() {
		if oaC.Name == axisName {
			currentAxisConfig = oaC
			break
		}
	}
	// If no ordinal axis with this name was found, try the hierarchy axes
	if currentAxisConfig == nil {
		hierarchyAxisConfigs, err := qe.searchConfigRepo.ListSearchConfigsOfType(ctx, solr.MexHierarchyAxisType)
		if err != nil {
			return nil, fmt.Errorf("failed to get hierarchy axis configs")
		}
		for _, haC := range hierarchyAxisConfigs.GetSearchConfigs() {
			if haC.Name == axisName {
				currentAxisConfig = haC
				break
			}
		}
	}
	if currentAxisConfig == nil {
		return nil, fmt.Errorf("could not find config for ordinal axis")
	}

	var fieldDefRep fields.BaseFieldDef
	var fieldHook fields.LifecyclePostQueryHook
	// If we have a hierarchical axis, we need to get the hierarchy field hook for enrichment
	if currentAxisConfig.Type == solr.MexHierarchyAxisType {
		if len(currentAxisConfig.Fields) == 0 {
			return nil, fmt.Errorf("hierarchy axis '%s' contains no fields", currentAxisConfig.Name)
		}
		// Since all fields in a hierarchical axis must have identical configurations w.r.t. linked hierarchy.
		// we simply take the configuration from the first field
		fieldDefRep, err = qe.fieldRepo.GetFieldDefByName(ctx, currentAxisConfig.Fields[0])
		if err != nil {
			return nil, fmt.Errorf("no configuration available for field used for in a hierarchical axis")
		}
		if fieldDefRep.Kind() != kindHierarchy.KindName {
			return nil, fmt.Errorf("hierarchy axis used for faceting contains a non-hierarchy field")
		}
		fieldHook = qe.postQueryHooks.GetHook(kindHierarchy.KindName)
		if fieldHook == nil {
			return nil, fmt.Errorf("no post query hook for field kind %s ", kindHierarchy.KindName)
		}
	}

	for _, bucket := range buckets {
		bucketMap, mapOk := bucket.(map[string]interface{})
		if !mapOk {
			return nil, fmt.Errorf("could not parse Solr facet bucket content")
		}
		bucketCount, countOk := bucketMap["count"].(float64)
		if !countOk {
			return nil, fmt.Errorf("could not parse Solr facet bucket count for facet axis")
		}
		bucketVal, valOk := bucketMap["val"].(string)
		if !valOk {
			return nil, fmt.Errorf("could not parse Solr facet bucket value for facet axis")
		}

		returnBucket := &solr.FacetBucket{Value: bucketVal, Count: uint32(bucketCount)}
		// Enrich hierarchy facets
		if currentAxisConfig.Type == solr.MexHierarchyAxisType && fieldHook != nil {
			returnBucket, err = fieldHook.EnrichFacetBucket(ctx, returnBucket, fieldDefRep)
			if err != nil {
				return nil, fmt.Errorf("enrichment of facet bucket failed: %s", err.Error())
			}
		}

		returnFacet.Buckets = append(returnFacet.Buckets, returnBucket)
	}

	return returnFacet, nil
}

// GetDateFieldRangesFromResponse extracts the axis min-max ranges from the response to a min-max range query
func GetDateFieldRangesFromResponse(solrResponse *solr.QueryResponse, reqFacets []*solr.Facet) (solr.StringFieldRanges, error) {
	resultRanges := solr.StringFieldRanges{}
	for _, facet := range reqFacets {
		if facet.GetType() == solr.MexStringStatFacetType {
			solrField, op, parseErr := getSolrFieldAndOpFromStatName(facet.GetStatName())
			if parseErr != nil {
				return nil, fmt.Errorf("could not interpret facet as a min or max")
			}
			val, err := getStringResultFromFacet(solrResponse.Facets, facet)
			if err != nil {
				return nil, err
			}
			// An empty result indicates that nothing was returned --> skip
			if val == "" {
				continue
			}
			switch op {
			case solr.MinOperator:
				if curRange, ok := resultRanges[solrField]; ok {
					if curRange.Min != "" {
						return nil, fmt.Errorf("encountered more than one min value for the solrField '%s'",
							solrField)
					}
					curRange.Min = val
					resultRanges[solrField] = curRange
				} else {
					resultRanges[solrField] = &solr.StringRange{
						Min: val,
					}
				}
			case solr.MaxOperator:
				if curRange, ok := resultRanges[solrField]; ok {
					if curRange.Max != "" {
						return nil, fmt.Errorf("encountered more than one max value for the solrField '%s'",
							solrField)
					}
					curRange.Max = val
					resultRanges[solrField] = curRange
				} else {
					resultRanges[solrField] = &solr.StringRange{
						Max: val,
					}
				}
			default:
				return nil, fmt.Errorf("operator for the facet must be 'min' or 'max'")
			}
		}
	}
	for _, strRange := range resultRanges {
		if strRange.Min == "" || strRange.Max == "" {
			return nil, fmt.Errorf("either min or max value is missing for an axis")
		}
	}
	return resultRanges, nil
}

// getStringResultFromFacet extracts a string-valued result (if possible) for the given facet from the returned response
func getStringResultFromFacet(resultFacets map[string]interface{}, facet *solr.Facet) (string, error) {
	rawVal, ok := resultFacets[facet.GetStatName()]
	if !ok {
		return "", nil
	}
	val, isString := rawVal.(string)
	if !isString {
		return "", fmt.Errorf("value returned for facet is not a string")
	}
	return val, nil
}

// getBucketNo extract the total facet bucket no. if present and otherwise sets it to 0
func getBucketNo(typedFacet map[string]interface{}, mexFacetType string) (uint32, error) {
	// Nos. from JSON are encoded as float64
	bucketNo := float64(0)
	// The bucket count is not returned for range facets and so not always present
	if bucketNumRaw, bucketNumOk := typedFacet["numBuckets"]; bucketNumOk {
		var convOK bool
		bucketNo, convOK = bucketNumRaw.(float64)
		if !convOK {
			return 0, fmt.Errorf("could not convert no. of bucket to float64")
		}
	} else if mexFacetType == solr.MexExactFacetType {
		return 0, fmt.Errorf("total bucket no. missing for terms facet")
	}
	return uint32(bucketNo), nil
}

// makeHighlight constructs the highlight part of the response from the Solr response
func (qe *QueryEngine) makeHighlight(docID string, highlightInfo interface{}) (*solr.Highlight, error) {
	highlight := solr.NewHighlight(docID)

	typedHighlightInfo, highlightOk := highlightInfo.(map[string]interface{})
	if !highlightOk {
		return nil, fmt.Errorf("could not parse Solr highlights")
	}
	matchMap := make(map[string]map[string]string)
	for fieldName, matches := range typedHighlightInfo {
		mappedFieldInfo := qe.returnFieldMapper.getMexNameFromBackingFieldName(fieldName)
		if _, ok := matchMap[mappedFieldInfo.Name]; !ok {
			matchMap[mappedFieldInfo.Name] = make(map[string]string)
		}
		typedMatches, matchOk := matches.([]interface{})
		if !matchOk {
			return nil, fmt.Errorf("could not parse Solr field highlights")
		}
		// skip adding matches if no matches were found
		if len(typedMatches) == 0 {
			continue
		}
		for _, snippet := range typedMatches {
			stringSnippet, snippetOk := snippet.(string)
			if !snippetOk {
				return nil, fmt.Errorf("could not parse Solr highlight snippet")
			}
			// New snippet --> store, repeat snippet previously only see for generic language --> overwrite language
			if curLang, ok := matchMap[mappedFieldInfo.Name][stringSnippet]; !ok || curLang == solr.GenericLangAbbrev {
				matchMap[mappedFieldInfo.Name][stringSnippet] = mappedFieldInfo.Language
			}
		}
	}
	for mappedFieldName, snippetsLangMap := range matchMap {
		// Invert map to organize snippets by language
		snippetsByLang := make(map[string][]string)
		for snippet, langCode := range snippetsLangMap {
			if _, ok := snippetsByLang[langCode]; !ok {
				snippetsByLang[langCode] = []string{snippet}
			} else {
				snippetsByLang[langCode] = append(snippetsByLang[langCode], snippet)
			}
		}
		// Format the snippets for return
		for lc, snippets := range snippetsByLang {
			mf := solr.FieldWithLanguage{
				Name:     mappedFieldName,
				Language: lc,
			}
			highlight.AddMatch(mf, snippets)
		}
	}
	// return nil if no highlights were found
	if len(highlight.Matches) == 0 {
		return nil, nil
	}
	return highlight, nil
}

// CombineRanges combines different sets of StringFieldRanges
func CombineRanges(stringRanges []solr.StringFieldRanges) (solr.StringFieldRanges, error) {
	combinedRange := solr.StringFieldRanges{}
	for _, cRange := range stringRanges {
		for key, val := range cRange {
			if _, ok := combinedRange[key]; ok {
				return nil, fmt.Errorf("A min-max range for an axis appeared more than once")
			}
			combinedRange[key] = val
		}
	}
	return combinedRange, nil
}

/*
	createFacetName return a facet SolrName, taking a pre-specified one if available and otherwise making one.

For simplicity, we generate the SolrName of the facet from the axis SolrName in the latter case. However,
we do not simply set it _equal_ to the axis name since an axis called "count" would then cause a
facet JSON object containing _two_ properties called "count" in the response (the other count being
the overall facet count)!
*/
func createFacetName(axisName string, statName string) string {
	if statName != "" {
		return statName
	}
	return fmt.Sprintf("%s_%s", solr.FacetPrefix, axisName)
}

func getStatNameForAxisAndOp(axis string, op string) (string, error) {
	if strings.Contains(op, "_") {
		return "", fmt.Errorf("a used operator is invalid as it contains a low-dash ('_')")
	}
	if axis == "" || op == "" {
		return "", fmt.Errorf("can only create stat facet for non-empty axis names and string")
	}
	return fmt.Sprintf("%s_%s", op, axis), nil
}

func getSolrFieldAndOpFromStatName(name string) (string, string, error) {
	parts := strings.SplitN(name, "_", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("a used facet cannot be split into an axis name and an operator")
	}
	return parts[1], parts[0], nil
}
