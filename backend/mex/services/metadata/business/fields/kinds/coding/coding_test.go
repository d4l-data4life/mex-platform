package kind_coding

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	fieldUtils "github.com/d4l-data4life/mex/mex/shared/fields"
	"github.com/d4l-data4life/mex/mex/shared/solr"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
)

func toAnySlice(msg ...proto.Message) []*anypb.Any {
	anys := make([]*anypb.Any, len(msg))

	for i, m := range msg {
		anyOne, err := anypb.New(m)
		if err != nil {
			panic(err)
		}

		anys[i] = anyOne
	}

	return anys
}

func TestKindCoding_ValidateDefinition(t *testing.T) {
	tests := []struct {
		name    string
		request *fieldUtils.FieldDef
		want    fields.BaseFieldDef
		wantErr bool
	}{
		{
			name: "non-coding field causes an error",
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
			kind := &KindCoding{}
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

func TestKindCoding_GenerateSolrFields(t *testing.T) {

	testFieldName := "test"
	baseDisplayFieldNameDe := solr.GetDisplayFieldName(testFieldName, solr.GermanLangAbbrev)
	baseDisplayFieldNameEn := solr.GetDisplayFieldName(testFieldName, solr.EnglishLangAbbrev)

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
			name: "invalid kind causes error",
			input: &fieldUtils.FieldDef{
				Name:     "test",
				Kind:     "xxx",
				IndexDef: &fieldUtils.IndexDef{},
			},
			wantErr: true,
		},
		{
			name: "generates the correct fields (code, DE label, EN label) with the right properties (including label fields always multivalued)",
			input: &fieldUtils.FieldDef{
				Name: "test",
				Kind: KindName,
				IndexDef: &fieldUtils.IndexDef{
					MultiValued: false,
					Ext:         toAnySlice(&fieldUtils.IndexDefExtCoding{CodingsetNames: []string{"mesh"}}),
				},
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
				},
				solr.GermanLangBaseFieldCategory: solr.FieldDef{
					Name:         baseDisplayFieldNameDe,
					Type:         solr.DefaultDeSolrTextFieldType,
					Stored:       true,
					Indexed:      false,
					MultiValued:  true,
					DocValues:    false,
					Uninvertible: false,
				},
				solr.EnglishLangBaseFieldCategory: solr.FieldDef{
					Name:         baseDisplayFieldNameEn,
					Type:         solr.DefaultEnSolrTextFieldType,
					Stored:       true,
					Indexed:      false,
					MultiValued:  true,
					DocValues:    false,
					Uninvertible: false,
				},
			},
		},
		{
			name: "if the MEx field is multi-valued, so are all backing fields",
			input: &fieldUtils.FieldDef{
				Name: "test",
				Kind: KindName,
				IndexDef: &fieldUtils.IndexDef{
					MultiValued: true,
					Ext:         toAnySlice(&fieldUtils.IndexDefExtCoding{CodingsetNames: []string{"mesh"}}),
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
				},
				solr.GermanLangBaseFieldCategory: solr.FieldDef{
					Name:         baseDisplayFieldNameDe,
					Type:         solr.DefaultDeSolrTextFieldType,
					Stored:       true,
					Indexed:      false,
					MultiValued:  true,
					DocValues:    false,
					Uninvertible: false,
				},
				solr.EnglishLangBaseFieldCategory: solr.FieldDef{
					Name:         baseDisplayFieldNameEn,
					Type:         solr.DefaultEnSolrTextFieldType,
					Stored:       true,
					Indexed:      false,
					MultiValued:  true,
					DocValues:    false,
					Uninvertible: false,
				},
			},
		},
	}

	for _, tt := range tests {
		kind0 := KindCoding{}
		kind1 := KindCoding{}
		t.Run(tt.name, func(t *testing.T) {
			fieldDef, err := kind0.ValidateDefinition(context.TODO(), tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateDefinition: error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			x, err := kind1.GenerateSolrFields(context.TODO(), fieldDef)
			if tt.wantErr {
				if err == nil {
					t.Errorf("GenerateSolrFields: error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if !reflect.DeepEqual(x, tt.want) {
				t.Errorf("solrFieldDefs: wanted %v but got %v", tt.want, x)
			}
		})
	}
}

func TestDNumberPattern(t *testing.T) {
	tests := []struct {
		input   string
		heading string
	}{
		{input: "", heading: ""},
		{input: "xxx", heading: ""},
		{input: "D1234", heading: "D1234"},
		{input: "DD1234", heading: "D1234"},
		{input: "foo/D1234", heading: "D1234"},
		{input: " D0001234", heading: "D0001234"},
		{input: "http://foo.bar.de/mesh/D001234666", heading: "D001234666"},
		{input: "http://foo.bar.de/mesh/D001234666x", heading: ""},
	}

	for _, tt := range tests {
		require.Equal(t, tt.heading, extractHeadingCode(tt.input))
	}
}
