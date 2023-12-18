package kind_number

import (
	"context"
	"reflect"
	"testing"

	"github.com/d4l-data4life/mex/mex/shared/solr"

	fieldUtils "github.com/d4l-data4life/mex/mex/shared/fields"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
)

func TestKindContact_ValidateDefinition(t *testing.T) {
	tests := []struct {
		name    string
		request *fieldUtils.FieldDef
		want    fields.BaseFieldDef
		wantErr bool
	}{
		{
			name: "non-number field causes an error",
			request: &fieldUtils.FieldDef{
				Name:     "someName",
				Kind:     "string",
				IndexDef: &fieldUtils.IndexDef{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kind := &KindNumber{}
			got, err := kind.ValidateDefinition(context.TODO(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateDefinition() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GenerateSolrFields(t *testing.T) {
	genericTypeName := solr.KnownLanguagesCategoryMap[solr.GenericLangAbbrev]
	tests := []struct {
		name  string
		input *fieldUtils.FieldDef

		wantErr bool
		want    solr.FieldCategoryToSolrFieldDefsMap
	}{
		{
			name: "invalid field name: disallowed characters",
			input: &fieldUtils.FieldDef{
				Name:     "invalid.name",
				Kind:     KindName,
				IndexDef: &fieldUtils.IndexDef{},
			},
			wantErr: true,
		},
		{
			name: "invalid kind",
			input: &fieldUtils.FieldDef{
				Name:     "test",
				Kind:     "xxx",
				IndexDef: &fieldUtils.IndexDef{},
			},
			wantErr: true,
		},

		{
			name: "generates a single number field with the right properties",
			input: &fieldUtils.FieldDef{
				Name:     "test",
				Kind:     KindName,
				IndexDef: &fieldUtils.IndexDef{},
			},
			want: solr.FieldCategoryToSolrFieldDefsMap{
				genericTypeName: solr.FieldDef{
					Name:         "test",
					Type:         solr.DefaultSolrNumberFieldType,
					Stored:       true,
					Indexed:      false,
					MultiValued:  false,
					DocValues:    false,
					Uninvertible: false,
				},
			},
		},
	}

	for _, tt := range tests {
		kind := KindNumber{}
		t.Run(tt.name, func(t *testing.T) {
			fieldDef, _ := kind.ValidateDefinition(context.TODO(), tt.input)

			got, err := kind.GenerateSolrFields(context.TODO(), fieldDef)
			if tt.wantErr == (err == nil) {
				t.Errorf("GenerateSolrFields: wanted error (true/false) = %v, got error %v", tt.wantErr, err)
				return
			}
			if err != nil {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("solrFieldDefs: wanted %v but got %v", tt.want, got)
			}
		})
	}
}
