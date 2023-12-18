package kind_string

import (
	"context"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"

	fieldUtils "github.com/d4l-data4life/mex/mex/shared/fields"
	"github.com/d4l-data4life/mex/mex/shared/solr"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/datamodel"
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
			name: "non-string field causes an error",
			request: &fieldUtils.FieldDef{
				Name:     "someName",
				Kind:     "number",
				IndexDef: &fieldUtils.IndexDef{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kind := &KindString{}
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

func Test_SolrFieldDefsGeneration(t *testing.T) {
	testFieldName := "test"
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
				Name:     testFieldName,
				Kind:     "xxx",
				IndexDef: &fieldUtils.IndexDef{},
			},
			wantErr: true,
		},

		{
			name: "generates string backing fields for the standard and the normalized value (with the correct properties)",
			input: &fieldUtils.FieldDef{
				Name:     testFieldName,
				Kind:     KindName,
				IndexDef: &fieldUtils.IndexDef{},
			},
			want: solr.FieldCategoryToSolrFieldDefsMap{
				solr.GenericLangBaseFieldCategory: solr.FieldDef{
					Name:         testFieldName,
					Type:         solr.DefaultSolrStringFieldType,
					Stored:       true,
					Indexed:      false,
					MultiValued:  false,
					DocValues:    false,
					Uninvertible: false,
					Required:     false,
				},
				solr.NormalizedBaseFieldCategory: solr.FieldDef{
					Name:         solr.GetNormalizedBackingFieldName(testFieldName),
					Type:         solr.DefaultSolrStringFieldType,
					Stored:       true,
					Indexed:      false,
					MultiValued:  false,
					DocValues:    false,
					Uninvertible: false,
					Required:     false,
				},
			},
		},
		{
			name: "a multivalued MEx field leads to multivalued backing fields",
			input: &fieldUtils.FieldDef{
				Name: testFieldName,
				Kind: KindName,
				IndexDef: &fieldUtils.IndexDef{
					MultiValued: true,
				},
			},
			want: solr.FieldCategoryToSolrFieldDefsMap{
				solr.GenericLangBaseFieldCategory: solr.FieldDef{
					Name:         testFieldName,
					Type:         solr.DefaultSolrStringFieldType,
					Stored:       true,
					Indexed:      false,
					MultiValued:  true,
					DocValues:    false,
					Uninvertible: false,
					Required:     false,
				},
				solr.NormalizedBaseFieldCategory: solr.FieldDef{
					Name:         solr.GetNormalizedBackingFieldName(testFieldName),
					Type:         solr.DefaultSolrStringFieldType,
					Stored:       true,
					Indexed:      false,
					MultiValued:  true,
					DocValues:    false,
					Uninvertible: false,
					Required:     false,
				},
			},
		},
	}

	for _, tt := range tests {
		kind := KindString{}
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

func TestKindString_GenerateXMLFieldTags(t *testing.T) {
	testFieldName := "test"

	tests := []struct {
		name      string
		fieldDef  fields.BaseFieldDef
		itemValue datamodel.CurrentItemValue
		want      []string
		wantErr   bool
	}{
		{
			name:     "If no language is specified, places data in standard field & normalized content in normalized field",
			fieldDef: nil,
			itemValue: datamodel.CurrentItemValue{
				ID:         "abc123",
				ItemID:     "def456",
				FieldName:  testFieldName,
				FieldValue: "hellö",
			},
			want: []string{
				"<field name=\"" + testFieldName + "\">hellö</field>",
				"<field name=\"" + solr.GetNormalizedBackingFieldName(testFieldName) + "\">hello</field>",
			},
		},
		{
			name:     "If unknown language is specified, places data in standard field & normalized content in normalized field",
			fieldDef: nil,
			itemValue: datamodel.CurrentItemValue{
				ID:         "abc123",
				ItemID:     "def456",
				FieldName:  testFieldName,
				FieldValue: "hellö",
				Language: pgtype.Text{
					String: "dk",
					Valid:  true,
				},
			},
			want: []string{
				"<field name=\"" + testFieldName + "\">hellö</field>",
				"<field name=\"" + solr.GetNormalizedBackingFieldName(testFieldName) + "\">hello</field>",
			},
		},
		{
			name:     "If known language is specified, places data in standard field & normalized content in normalized field",
			fieldDef: nil,
			itemValue: datamodel.CurrentItemValue{
				ID:         "abc123",
				ItemID:     "def456",
				FieldName:  testFieldName,
				FieldValue: "hellö",
				Language: pgtype.Text{
					String: solr.GermanLangAbbrev,
					Valid:  true,
				},
			},
			want: []string{
				"<field name=\"" + testFieldName + "\">hellö</field>",
				"<field name=\"" + solr.GetNormalizedBackingFieldName(testFieldName) + "\">hello</field>",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ki := &KindString{}
			got, err := ki.GenerateXMLFieldTags(context.Background(), tt.fieldDef, tt.itemValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateXMLFieldTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateXMLFieldTags() got = %v, want %v", got, tt.want)
			}
		})
	}
}
