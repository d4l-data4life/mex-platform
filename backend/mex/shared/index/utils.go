package index

import (
	"context"
	"fmt"

	"github.com/d4l-data4life/mex/mex/shared/solr"
	"github.com/d4l-data4life/mex/mex/shared/utils"
)

// ClearSchema removes all existing fields, copy fields, and dynamic fields from schema
func ClearSchema(ctx context.Context, client solr.ClientAPI) error {
	// Remove all copy fields - has to happen first since fields cannot be deleted if referenced in a copy field
	var copyFieldsToRemove []solr.RemoveCopyFieldSubBody
	currentCopyFields, copyFieldErr := client.GetSchemaCopyFields(ctx)
	if copyFieldErr != nil {
		return fmt.Errorf("failed to retrieve schema copy fields: %s", copyFieldErr.Error())
	}
	for _, f := range currentCopyFields {
		copyFieldsToRemove = append(copyFieldsToRemove, solr.RemoveCopyFieldSubBody{
			Source: f.Source,
			Dest:   f.Destination,
		})
	}
	cfErr := client.RemoveSchemaCopyFields(ctx, copyFieldsToRemove)
	if cfErr != nil {
		return fmt.Errorf("failed to remove %d schema copy field(s): %s", len(copyFieldsToRemove), cfErr.Error())
	}

	// Remove all fields
	currentFields, fieldErr := client.GetSchemaFields(ctx)
	if fieldErr != nil {
		return fmt.Errorf("failed to retrieve schema fields: %s", fieldErr.Error())
	}

	var fieldsToRemove []string
	for _, f := range currentFields {
		// Leave 'id', '_version_', and default search fields to avoid Solr startup issues
		if utils.Contains(solr.ProtectedSolrFields, f.Name) {
			continue
		}
		fieldsToRemove = append(fieldsToRemove, f.Name)
	}
	fErr := client.RemoveSchemaFields(ctx, fieldsToRemove)
	if fErr != nil {
		return fmt.Errorf("failed to remove %d schema field(s): %s", len(fieldsToRemove), fErr.Error())
	}

	// Remove all dynamic fields
	currentDynamicFields, dynamicFieldErr := client.GetSchemaDynamicFields(ctx)
	if dynamicFieldErr != nil {
		return fmt.Errorf("failed to retrieve schema dynamic fields: %s", dynamicFieldErr.Error())
	}
	fieldsToRemove = []string{}
	for _, f := range currentDynamicFields {
		fieldsToRemove = append(fieldsToRemove, f.Name)
	}
	dfErr := client.RemoveSchemaDynamicFields(ctx, fieldsToRemove)
	if dfErr != nil {
		return fmt.Errorf("failed to remove %d schema dynamic field(s): %s", len(fieldsToRemove), dfErr.Error())
	}

	return nil
}

// UploadSchemaUpdates pushes a set of schema updates to a Solr instance
func UploadSchemaUpdates(ctx context.Context, schema *solr.SchemaUpdates, client solr.ClientAPI, progressor utils.Progressor) error {
	if progressor == nil {
		progressor = &utils.NopProgressor{}
	}

	uniqueKey, keyErr := client.GetSchemaUniqueKey(ctx)
	if keyErr != nil {
		return fmt.Errorf("failed to retrieve schema unique key: %s", keyErr.Error())
	}
	if uniqueKey != solr.DefaultUniqueKey {
		return fmt.Errorf("unique key of Solr schema is not set to the required standard value '%s'", solr.DefaultUniqueKey)
	}

	fieldErr := client.AddSchemaFields(ctx, schema.FieldDefs)
	if fieldErr != nil {
		return fmt.Errorf("failed to create schema fields: %v", fieldErr.Error())
	}
	progressor.Progress("schema update", "fields: done")

	fieldErr = client.AddSchemaCopyFields(ctx, schema.CopyFieldDefs)
	if fieldErr != nil {
		return fmt.Errorf("failed to create schema copy fields: %v", fieldErr.Error())
	}
	progressor.Progress("schema update", "copy fields: done")

	fieldErr = client.AddSchemaDynamicFields(ctx, schema.DynamicFieldDefs)
	if fieldErr != nil {
		return fmt.Errorf("failed to create dynamic fields: %v", fieldErr.Error())
	}
	progressor.Progress("schema update", "dynamic fields: done")

	return nil
}
