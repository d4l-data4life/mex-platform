package index

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/d4l-data4life/mex/mex/shared/constants"
	"github.com/d4l-data4life/mex/mex/shared/errstat"
	"github.com/d4l-data4life/mex/mex/shared/hints"
	"github.com/d4l-data4life/mex/mex/shared/index"
	"github.com/d4l-data4life/mex/mex/shared/known/jobspb"
	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/searchconfig"
	"github.com/d4l-data4life/mex/mex/shared/searchconfig/sctypes"
	"github.com/d4l-data4life/mex/mex/shared/solr"
	"github.com/d4l-data4life/mex/mex/shared/utils"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields/hooks"
	kind_hier "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/hierarchy"

	"github.com/d4l-data4life/mex/mex/services/index/endpoints/index/pb"
)

// CreateIndex updates the Solr schema and the query engine configuration,
// based on the metadata configuration in the DB.
// The data indexed in Solr is not touched.
func (svc *Service) CreateIndex(ctx context.Context, request *pb.CreateIndexRequest) (*pb.CreateIndexResponse, error) {
	lock, err := svc.JobService.AcquireLock(ctx, SvcResourceName)
	if err != nil {
		return nil, errstat.MakeGRPCStatus(codes.AlreadyExists, "create: failed to acquire index lock; other job might be running", request).Err()
	}

	job, err := svc.JobService.CreateJob(ctx, &jobspb.CreateJobRequest{
		Title: "Create Solr index",
	})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failure creating job:  %s", err.Error()))
	}

	// Create new independent context for the job copying the relevant values.
	// (The request's ctx will go out of scope before the job is done, so we cannot use it directly or as a parent context.)
	ctxJob := constants.NewContextWithValues(ctx, job.JobId)

	go func(ctx context.Context) {
		svc.Log.Info(ctx, L.Messagef("schema regeneration: job started (%s)", job.JobId), L.Phase("job"))

		svc.JobService.SetStatusRunning(ctx, job.JobId)              //nolint:errcheck
		defer svc.JobService.SetStatusDone(ctx, job.JobId)           //nolint:errcheck
		defer svc.JobService.ReleaseLock(ctx, SvcResourceName, lock) //nolint:errcheck

		err = svc.doSchemaRebuild(ctx, svc.TelemetryService)
		if err != nil {
			svc.Log.Error(ctx, L.Message(err.Error()))
			_, err := svc.JobService.SetError(ctx, &jobspb.SetJobErrorRequest{Error: err.Error(), JobId: job.JobId})
			if err != nil {
				svc.Log.Warn(ctx, L.Messagef("could not set job error: %s", err.Error()))
			}
		}

		svc.Log.Info(ctx, L.Messagef("schema regeneration: job done (%s)", job.JobId), L.Phase("job"))
		svc.TelemetryService.Done()
	}(ctxJob)

	hints.HintHTTPStatusCode(ctx, http.StatusCreated)
	return &pb.CreateIndexResponse{JobId: job.JobId}, nil
}

func (svc *Service) doSchemaRebuild(ctx context.Context, progressor utils.Progressor) error {
	// Retrieve metadata configuration
	svc.Log.Info(ctx, L.Message("fetching field configuration"))
	fieldDefs, err := svc.FieldRepo.ListFieldDefs(ctx)
	if err != nil {
		return fmt.Errorf("error reading field configurations: %s", err.Error())
	}
	svc.Log.Info(ctx, L.Message("fetching search configuration elements"))
	searchConfigElements, err := svc.SearchConfigRepo.ListSearchConfigs(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve search configuration: %s", err.Error())
	}

	svc.Log.Info(ctx, L.Message("processing search focus field configuration"))
	solrSchemaUpdates, err := generateSolrSchema(ctx, svc.SolrFieldCreationHooks, fieldDefs, searchConfigElements)
	if err != nil {
		return fmt.Errorf("failed to generate Solr schema: %s", err.Error())
	}

	// Clear all fields except core fields like 'id'
	svc.Log.Info(ctx, L.Message("clearing Solr schema"))
	err = index.ClearSchema(ctx, svc.Solr)
	if err != nil {
		return fmt.Errorf("failure when trying to remove Solr schema fields: %s", err.Error())
	}

	// Apply schema update
	svc.Log.Info(ctx, L.Message("re-generating Solr schema"))
	err = index.UploadSchemaUpdates(ctx, solrSchemaUpdates, svc.Solr, progressor)
	if err != nil {
		return fmt.Errorf("error trying to update Solr schema: %s", err.Error())
	}

	return nil
}

type hierarchyFieldConfig struct {
	codeSystemNameOrNodeEntityType string
	linkFieldName                  string
	displayFieldName               string
}

// generateSolrSchema generates all Solr field definitions for the full Solr schema
func generateSolrSchema(ctx context.Context, solrFieldCreationHooks hooks.SolrFieldCreationHooks,
	fieldDefs []fields.BaseFieldDef, searchConfigElements *searchconfig.SearchConfigList,
) (*solr.SchemaUpdates, error) {
	solrSchemaUpdates := &solr.SchemaUpdates{
		FieldDefs:        []solr.FieldDef{},
		CopyFieldDefs:    []solr.CopyFieldDef{},
		DynamicFieldDefs: []solr.DynamicFieldDef{},
	}

	fieldKindMap := make(map[string]string)
	for _, fd := range fieldDefs {
		fieldKindMap[fd.Name()] = fd.Kind()
	}
	focusFieldsMap := make(map[string][]string)
	axisFieldsMap := make(map[string][]string)
	for _, elem := range searchConfigElements.GetSearchConfigs() {
		switch elem.Type {
		case solr.MexSearchFocusType:
			focusFieldsMap[elem.Name] = elem.Fields
		case solr.MexOrdinalAxisType:
			axisFieldsMap[elem.Name] = elem.Fields
		case solr.MexHierarchyAxisType:
			axisFieldsMap[elem.Name] = elem.Fields
		default:
			// Should never happen
			return nil, fmt.Errorf("invalid search configuration type: %s", elem.Type)
		}
	}

	// Loop over all configured fields and create the needed Solr fields for backing them
	mexFieldMap := make(solr.MexFieldBackingInfoMap)
	hierarchyFieldMaps := make(map[string]hierarchyFieldConfig)
	for _, fieldDef := range fieldDefs {
		if fieldDef.Name() == solr.DefaultUniqueKey {
			// Skip the field 'id' which is hardcoded into the base schema
			continue
		}
		hook := solrFieldCreationHooks.GetHook(fieldDef.Kind())
		if hook == nil {
			return nil, fmt.Errorf("hook not found for kind: %s", fieldDef.Kind())
		}

		// We store the hierarchy-specific configuration for all hierarchy fields to be able to check consistency
		if fieldDef.Kind() == kind_hier.KindName {
			if hFieldDef, ok := (fieldDef).(kind_hier.HierarchyFieldDef); ok {
				hierarchyFieldMaps[hFieldDef.Name()] = hierarchyFieldConfig{
					codeSystemNameOrNodeEntityType: hFieldDef.CodeSystemNameOrEntityType(),
					linkFieldName:                  hFieldDef.LinkFieldName(),
					displayFieldName:               hFieldDef.DisplayFieldName(),
				}
			}
		}

		// Add fields needed to back these fields in Solr
		solrFieldsMap, err := hook.GenerateSolrFields(ctx, fieldDef)
		if err != nil {
			return nil, fmt.Errorf(fmt.Sprintf("solr field generation failed: %s", err.Error()))
		}

		// Currently, every generated field is captured in this map
		var mexBackingFields []solr.MexBackingFieldWiringInfo
		for category, fDef := range solrFieldsMap {
			solrSchemaUpdates.FieldDefs = append(solrSchemaUpdates.FieldDefs, fDef)
			mexBackingFields = append(mexBackingFields, solr.MexBackingFieldWiringInfo{
				Name:     fDef.Name,
				Category: category,
			})
		}
		mexFieldMap[fieldDef.Name()] = solr.MexFieldWiringInfo{
			MexType:       fieldDef.Kind(),
			BackingFields: mexBackingFields,
		}
	}

	// Add backing Solr fields for search config elements
	searchConfigBackingFields := []solr.FieldDef{}
	searchConfigCopyFields := []solr.CopyFieldDef{}
	for _, searchConfigElem := range searchConfigElements.GetSearchConfigs() {
		// For hierarchy axes with multiple fields, ensure that all fields use the same hierarchy configuration
		if searchConfigElem.Type == solr.MexHierarchyAxisType && len(searchConfigElem.Fields) > 1 {
			firstFieldName := searchConfigElem.Fields[0]
			firstFieldConfig, ok := hierarchyFieldMaps[firstFieldName]
			if !ok {
				return nil, fmt.Errorf("the hierarchy field '%s' used in the axis '%s' could not be found", firstFieldName, searchConfigElem.Name)
			}
			for i := 1; i < len(searchConfigElem.Fields); i++ {
				curFieldConfig, okNew := hierarchyFieldMaps[searchConfigElem.Fields[i]]
				if !okNew {
					return nil, fmt.Errorf("the hierarchy field '%s' used in the axis '%s' could not be found", curFieldConfig, searchConfigElem.Name)
				}
				if curFieldConfig.linkFieldName != firstFieldConfig.linkFieldName || curFieldConfig.codeSystemNameOrNodeEntityType != firstFieldConfig.codeSystemNameOrNodeEntityType ||
					curFieldConfig.displayFieldName != firstFieldConfig.displayFieldName {
					return nil, fmt.Errorf("the hierarchy configurations of the fields in the hierarchy axis '%s' are not identical", searchConfigElem.Name)
				}
			}
		}

		hook, ok := sctypes.SearchConfigHooks[searchConfigElem.Type]
		if !ok {
			return nil, fmt.Errorf("hook not found for search config element '%s' of type '%s",
				searchConfigElem.Name, searchConfigElem.Type)
		}
		newBackingFields, newCopyFields, err := hook.GetSolrBackingFields(searchConfigElem, mexFieldMap)
		if err != nil {
			return nil, fmt.Errorf("could not generate backing fields for search config element '%s' of type '%s': '%s'",
				searchConfigElem.Name, searchConfigElem.Type, err.Error())
		}
		searchConfigBackingFields = append(searchConfigBackingFields, newBackingFields...)
		searchConfigCopyFields = append(searchConfigCopyFields, newCopyFields...)
	}

	solrSchemaUpdates.FieldDefs = append(solrSchemaUpdates.FieldDefs, searchConfigBackingFields...)
	solrSchemaUpdates.CopyFieldDefs = append(solrSchemaUpdates.CopyFieldDefs, searchConfigCopyFields...)

	return solrSchemaUpdates, nil
}
