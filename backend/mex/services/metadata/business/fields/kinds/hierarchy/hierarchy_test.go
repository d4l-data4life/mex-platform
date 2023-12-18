package kind_hierarchy

import (
	"context"
	"reflect"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	fieldUtils "github.com/d4l-data4life/mex/mex/shared/fields"
	"github.com/d4l-data4life/mex/mex/shared/solr"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/hierarchy/codings"
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

func TestKindHierarchy_ValidateDefinition(t *testing.T) {
	tests := []struct {
		name    string
		request *fieldUtils.FieldDef
		wantErr bool
	}{
		{
			name: "non-hierarchy field causes an error",
			request: &fieldUtils.FieldDef{
				Name:     "someName",
				Kind:     "string",
				IndexDef: &fieldUtils.IndexDef{},
			},
			wantErr: true,
		},
		{
			name: "hierarchy field with neither a link nor a hierarchy config extension causers error",
			request: &fieldUtils.FieldDef{
				Name:     "someName",
				Kind:     "hierarchy",
				IndexDef: &fieldUtils.IndexDef{},
			},
			wantErr: true,
		},
		{
			name: "hierarchy field with a link config extension but no hierarchy config extensions causes an error",
			request: &fieldUtils.FieldDef{
				Name: "someName",
				Kind: "hierarchy",
				IndexDef: &fieldUtils.IndexDef{
					Ext: toAnySlice(
						&fieldUtils.IndexDefExtLink{RelationType: "someType"},
					),
				},
			},
			wantErr: true,
		},
		{
			name: "hierarchy field with a hierarchy config extension but no link config extensions causes an error",
			request: &fieldUtils.FieldDef{
				Name: "someName",
				Kind: "hierarchy",
				IndexDef: &fieldUtils.IndexDef{
					Ext: toAnySlice(
						&fieldUtils.IndexDefExtHierarchy{CodeSystemNameOrNodeEntityType: "mesh"},
					),
				},
			},
			wantErr: true,
		},
		{
			name: "hierarchy field with both a link and a hierarchy config extension works",
			request: &fieldUtils.FieldDef{
				Name: "someName",
				Kind: "hierarchy",
				IndexDef: &fieldUtils.IndexDef{
					Ext: toAnySlice(
						&fieldUtils.IndexDefExtLink{RelationType: "someType"},
						&fieldUtils.IndexDefExtHierarchy{CodeSystemNameOrNodeEntityType: "mesh"},
					),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kind := &kindHierarchy{}
			_, err := kind.ValidateDefinition(context.TODO(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestKindHierarchy_GenerateSolrFields(t *testing.T) {

	testFieldName := "test"
	baseTransitiveHullFieldName := solr.GetTransitiveHullFieldName(testFieldName)
	baseDisplayFieldName := solr.GetTransitiveHullDisplayFieldName(testFieldName)
	baseDisplayFieldNameDe, _ := solr.GetLangSpecificFieldName(baseDisplayFieldName, solr.GermanLangAbbrev)
	baseDisplayFieldNameEn, _ := solr.GetLangSpecificFieldName(baseDisplayFieldName, solr.EnglishLangAbbrev)

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
			name: "generates the correct fields (code, parent code set, DE label, EN label) with the right properties (including parent codes & label fields always multivalued)",
			input: &fieldUtils.FieldDef{
				Name: "test",
				Kind: KindName,
				IndexDef: &fieldUtils.IndexDef{
					MultiValued: false,
					Ext: toAnySlice(
						&fieldUtils.IndexDefExtLink{RelationType: "someType"},
						&fieldUtils.IndexDefExtHierarchy{CodeSystemNameOrNodeEntityType: "mesh"},
					),
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
				solr.ParentCodesBaseFieldCategory: solr.FieldDef{
					Name:         baseTransitiveHullFieldName,
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
		{
			name: "if the MEx field is multi-valued, so are all backing fields",
			input: &fieldUtils.FieldDef{
				Name: "test",
				Kind: KindName,
				IndexDef: &fieldUtils.IndexDef{
					MultiValued: true,
					Ext: toAnySlice(
						&fieldUtils.IndexDefExtLink{RelationType: "someType"},
						&fieldUtils.IndexDefExtHierarchy{CodeSystemNameOrNodeEntityType: "mesh"},
					)},
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
				solr.ParentCodesBaseFieldCategory: solr.FieldDef{
					Name:         baseTransitiveHullFieldName,
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

	curCodings := &codings.MockedCodings{Codings: map[string][]codings.Coding{
		"mesh": {},
	}}

	for _, tt := range tests {
		kind0 := kindHierarchy{codings: curCodings}
		kind1 := KindHierarchy{}
		t.Run(tt.name, func(t *testing.T) {
			fieldDef, _ := kind0.ValidateDefinition(context.TODO(), tt.input)

			got, err := kind1.GenerateSolrFields(context.TODO(), fieldDef)
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

func TestKindHierarchy_MarshalToProtobufFormat(t *testing.T) {

	fieldDef := &fieldUtils.FieldDef{
		Name:      "testField",
		Kind:      KindName,
		DisplayId: "SOMETHING",
		IndexDef: &fieldUtils.IndexDef{
			MultiValued: true,
			Ext: toAnySlice(
				&fieldUtils.IndexDefExtHierarchy{
					CodeSystemNameOrNodeEntityType: "mesh",
					LinkFieldName:                  "parentTerm",
					DisplayFieldName:               "acronym",
				},
				&fieldUtils.IndexDefExtLink{
					RelationType:       "someType",
					LinkedTargetFields: []string{},
				},
			)},
	}

	hierarchyField, _ := (&kindHierarchy{}).ValidateDefinition(context.TODO(), fieldDef)

	tests := []struct {
		name     string
		fieldDef fields.BaseFieldDef
		want     *fieldUtils.FieldDef
		wantErr  bool
	}{
		{
			name:     "Correctly includes all fields",
			fieldDef: hierarchyField,
			want:     fieldDef,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kind := &kindHierarchy{}
			got, err := kind.MarshalToProtobufFormat(context.TODO(), tt.fieldDef)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalToProtobufFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalToProtobufFormat() got = %v, want %v", got, tt.want)
			}
		})
	}
}
