package index

import (
	"context"
	"testing"

	"github.com/d4l-data4life/mex/mex/shared/solr"
	"github.com/d4l-data4life/mex/mex/shared/testutils"
)

var (
	df1 = solr.DynamicFieldDef{Name: "test*", Type: solr.DefaultDeSolrTextFieldType, MultiValued: true, Indexed: true}
	df2 = solr.DynamicFieldDef{Name: "date*", Type: solr.DefaultSolrTextFieldType, MultiValued: true, Indexed: true}

	mockReturnVals = solr.ReturnVals{
		Fields: []solr.FieldDef{
			{solr.DefaultUniqueKey, solr.DefaultSolrStringFieldType, false, true, true, true, true, false},
			{"default_search_focus___generic", solr.DefaultSolrTextFieldType, false, false, false, false, false,
				false},
			{"test", solr.DefaultSolrTextFieldType, true, false, true, false, false, false},
			{"date", solr.DefaultSolrTimestampFieldType, true, false, true, false, true, false},
		},
		CopyFields: []solr.CopyFieldResponse{
			{"test", "testCopy"},
			{"date", "dateCopy"},
		},
		DynamicFields: []solr.DynamicFieldDef{df1, df2},
	}
)

func TestClearSchema(t *testing.T) {
	tests := []struct {
		name    string
		client  solr.MockClient
		checks  *[]testutils.ClientCheck
		wantErr bool
	}{
		{
			name:    "Fetches current schema fields from Solr",
			client:  solr.NewMockClient(false, "id", mockReturnVals),
			checks:  &[]testutils.ClientCheck{testutils.CheckFieldsFetchedViaClient(true)},
			wantErr: false,
		},
		{
			name:    "Fetches current schema copy fields from Solr",
			client:  solr.NewMockClient(false, "id", mockReturnVals),
			checks:  &[]testutils.ClientCheck{testutils.CheckCopyFieldsFetchedViaClient(true)},
			wantErr: false,
		},
		{
			name:    "Fetches current schema dynamic fields from Solr",
			client:  solr.NewMockClient(false, "id", mockReturnVals),
			checks:  &[]testutils.ClientCheck{testutils.CheckDynamicFieldsFetchedViaClient(true)},
			wantErr: false,
		},
		{
			name:    "Removes all current schema fields from Solr EXCEPT 'id', '_version_', and '_text_'",
			client:  solr.NewMockClient(false, "id", mockReturnVals),
			checks:  &[]testutils.ClientCheck{testutils.CheckFieldsRemovedViaClient([]string{"test", "date", "default_search_focus___generic"})},
			wantErr: false,
		},
		{
			name:    "Removes all current schema copy fields from Solr",
			client:  solr.NewMockClient(false, "id", mockReturnVals),
			checks:  &[]testutils.ClientCheck{testutils.CheckCopyFieldsRemovedViaClient([]string{"test-->testCopy", "date-->dateCopy"})},
			wantErr: false,
		},
		{
			name:    "removes all current schema dynamic fields from Solr",
			client:  solr.NewMockClient(false, "id", mockReturnVals),
			checks:  &[]testutils.ClientCheck{testutils.CheckDynamicFieldsRemovedViaClient([]string{"test*", "date*"})},
			wantErr: false,
		},
		{
			name:    "returns error if client calls fail",
			client:  solr.NewMockClient(true, "id", mockReturnVals),
			checks:  nil,
			wantErr: true,
		},
		{
			name:    "copy fields are removed before normal fields",
			client:  solr.NewMockClient(false, "id", mockReturnVals),
			checks:  &[]testutils.ClientCheck{testutils.CheckClientCallOrder("RemoveSchemaCopyFields", "RemoveSchemaFields")},
			wantErr: false,
		},
		{
			name:    "copy fields are removed before dynamic fields",
			client:  solr.NewMockClient(false, "id", mockReturnVals),
			checks:  &[]testutils.ClientCheck{testutils.CheckDynamicFieldsRemovedViaClient([]string{"test*", "date*"})},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ClearSchema(context.TODO(), &tt.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClearSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.checks != nil {
				for _, check := range *tt.checks {
					check(tt.client, t)
				}
			}
		})
	}
}

func TestUploadSchemaUpdates(t *testing.T) {
	df1 := solr.DynamicFieldDef{Name: "test*", Type: solr.DefaultSolrTextFieldType, MultiValued: true, Indexed: true}
	df2 := solr.DynamicFieldDef{Name: "date*", Type: "pdate", MultiValued: true, Indexed: true}
	tests := []struct {
		name          string
		client        solr.MockClient
		schemaUpdates *solr.SchemaUpdates
		checks        *[]testutils.ClientCheck
		wantErr       bool
	}{
		{
			name:          "Fetches current unique ID",
			client:        solr.NewMockClient(false, "id", mockReturnVals),
			schemaUpdates: &solr.SchemaUpdates{},
			checks:        &[]testutils.ClientCheck{testutils.CheckUniqueIDFetchedViaClient(true)},
			wantErr:       false,
		},
		{
			name:          "Returns error if unique ID does not have required value",
			client:        solr.NewMockClient(false, "fake_id", mockReturnVals),
			schemaUpdates: &solr.SchemaUpdates{},
			checks:        nil,
			wantErr:       true,
		},
		{
			name:   "Adds all fields",
			client: solr.NewMockClient(false, solr.DefaultUniqueKey, mockReturnVals),
			schemaUpdates: &solr.SchemaUpdates{
				FieldDefs: []solr.FieldDef{
					{"test", solr.DefaultSolrTextFieldType, true, false, true, false, false, false},
					{"date", solr.DefaultSolrTimestampFieldType, true, false, true, false, true, false},
				},
				CopyFieldDefs:    []solr.CopyFieldDef{},
				DynamicFieldDefs: []solr.DynamicFieldDef{},
			},
			checks: &[]testutils.ClientCheck{
				testutils.CheckFieldsAddedViaClient([]string{"test", "date"}),
			},
			wantErr: false,
		},
		{
			name:   "Adds all copy fields",
			client: solr.NewMockClient(false, solr.DefaultUniqueKey, mockReturnVals),
			schemaUpdates: &solr.SchemaUpdates{
				FieldDefs: []solr.FieldDef{},
				CopyFieldDefs: []solr.CopyFieldDef{
					{"test", []string{"testCopy"}, 0},
					{"date", []string{"dateCopy"}, 0},
				},
				DynamicFieldDefs: []solr.DynamicFieldDef{},
			},
			checks: &[]testutils.ClientCheck{
				testutils.CheckCopyFieldsAddedViaClient([]string{"test-->testCopy", "date-->dateCopy"}),
			},
			wantErr: false,
		},
		{
			name:   "Adds all dynamic fields",
			client: solr.NewMockClient(false, solr.DefaultUniqueKey, mockReturnVals),
			schemaUpdates: &solr.SchemaUpdates{
				FieldDefs:        []solr.FieldDef{},
				CopyFieldDefs:    []solr.CopyFieldDef{},
				DynamicFieldDefs: []solr.DynamicFieldDef{df1, df2},
			},
			checks: &[]testutils.ClientCheck{
				testutils.CheckDynamicFieldsAddedViaClient([]string{"test*", "date*"}),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UploadSchemaUpdates(context.Background(), tt.schemaUpdates, &tt.client, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadSchemaUpdates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.checks != nil {
				for _, check := range *tt.checks {
					check(tt.client, t)
				}
			}
		})
	}
}
