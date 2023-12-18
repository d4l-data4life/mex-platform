package index

import (
	"context"
	"reflect"
	"sort"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	sharedFields "github.com/d4l-data4life/mex/mex/shared/fields"
	sharedSearchConfig "github.com/d4l-data4life/mex/mex/shared/searchconfig"
	"github.com/d4l-data4life/mex/mex/shared/solr"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields/frepo"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields/hooks"
	kind_hierarchy "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/hierarchy"
	kind_string "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/string"
	kind_text "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/text"
	kind_timestamp "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/timestamp"
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

func checkSolrFields(t *testing.T, got []solr.FieldDef, expected []solr.FieldDef) {
	if len(got) != len(expected) {
		t.Errorf("Expected %d fields but got %d", len(expected), len(got))
		return
	}
	sort.Slice(got, func(i, j int) bool {
		return got[i].Name < got[j].Name
	})
	sort.Slice(expected, func(i, j int) bool {
		return expected[i].Name < expected[j].Name
	})

	for i, gotF := range got {
		if !reflect.DeepEqual(gotF, expected[i]) {
			t.Errorf("Mismatch for field with index %d (after sorting): got %v but wanted %v d", i, gotF,
				expected[i])
		}
	}
}

func checkSolrCopyFields(t *testing.T, got []solr.CopyFieldDef, expected []solr.CopyFieldDef) {
	if len(got) != len(expected) {
		t.Errorf("Expected %d copy fields but got %d", len(expected), len(got))
		return
	}
	sort.Slice(got, func(i, j int) bool {
		if got[i].Source == got[j].Source {
			return got[i].Destination[0] < got[j].Destination[0]
		}
		return got[i].Source < got[j].Source
	})
	sort.Slice(expected, func(i, j int) bool {
		if expected[i].Source == expected[j].Source {
			return expected[i].Destination[0] < expected[j].Destination[0]
		}
		return expected[i].Source < expected[j].Source
	})

	for i, gotCF := range got {
		if !reflect.DeepEqual(gotCF, expected[i]) {
			t.Errorf("Mismatch for copy field with index %d (after sorting): got %v but wanted %v d", i, gotCF,
				expected[i])
		}
	}
}

func Test_generateSolrSchema(t *testing.T) {

	// Field sets used for testing

	noFieldsRepo := frepo.NewMockedFieldRepo([]fields.BaseFieldDef{})
	someFieldsRepo := frepo.NewMockedFieldRepo([]fields.BaseFieldDef{
		(&kind_text.KindText{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{
			Name: "label",
			Kind: "text",
			IndexDef: &sharedFields.IndexDef{
				MultiValued: true,
			}}),
		(&kind_string.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{
			Name: "keyword",
			Kind: "string",
			IndexDef: &sharedFields.IndexDef{
				MultiValued: true,
			}}),
		(&kind_timestamp.KindTimestamp{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{
			Name: "created",
			Kind: "timestamp",
			IndexDef: &sharedFields.IndexDef{
				MultiValued: false,
			}}),
	})

	hierField1, _ := kind_hierarchy.NewKindHierarchy(nil)
	hierField2, _ := kind_hierarchy.NewKindHierarchy(nil)
	brokenHierarchyFieldsRepo := frepo.NewMockedFieldRepo([]fields.BaseFieldDef{
		(&kind_text.KindText{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{
			Name: "label",
			Kind: "text",
			IndexDef: &sharedFields.IndexDef{
				MultiValued: true,
			}}),
		(&kind_string.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{
			Name: "keyword",
			Kind: "string",
			IndexDef: &sharedFields.IndexDef{
				MultiValued: true,
			}}),
		(&kind_timestamp.KindTimestamp{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{
			Name: "created",
			Kind: "timestamp",
			IndexDef: &sharedFields.IndexDef{
				MultiValued: false,
			}}),
		hierField1.MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{
			Name: "hier1",
			Kind: "hierarchy",
			IndexDef: &sharedFields.IndexDef{
				MultiValued: false,
				Ext: toAnySlice(
					&sharedFields.IndexDefExtLink{
						RelationType: "someType1",
					},
					&sharedFields.IndexDefExtHierarchy{
						CodeSystemNameOrNodeEntityType: "A",
						LinkFieldName:                  "linkField",
						DisplayFieldName:               "keyword",
					}),
			}}),
		hierField2.MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{
			Name: "hier2",
			Kind: "hierarchy",
			IndexDef: &sharedFields.IndexDef{
				MultiValued: false,
				Ext: toAnySlice(
					&sharedFields.IndexDefExtLink{
						RelationType: "someType2",
					},
					&sharedFields.IndexDefExtHierarchy{
						CodeSystemNameOrNodeEntityType: "A",
						LinkFieldName:                  "linkField",
						DisplayFieldName:               "label",
					}),
			}}),
	})

	emptySearchConfig := &sharedSearchConfig.SearchConfigList{}
	focusSearchConfig := &sharedSearchConfig.SearchConfigList{
		SearchConfigs: []*sharedSearchConfig.SearchConfigObject{
			{
				Type: solr.MexSearchFocusType,
				Name: "testFocus",
				Fields: []string{
					"label",
					"keyword",
				},
			},
		},
	}
	axisSearchConfig := &sharedSearchConfig.SearchConfigList{
		SearchConfigs: []*sharedSearchConfig.SearchConfigObject{
			{
				Type: solr.MexOrdinalAxisType,
				Name: "testAxis",
				Fields: []string{
					"label",
					"keyword",
				},
			},
		},
	}
	focusAndAxisSearchConfig := &sharedSearchConfig.SearchConfigList{
		SearchConfigs: []*sharedSearchConfig.SearchConfigObject{
			{
				Type: solr.MexSearchFocusType,
				Name: "testFocus",
				Fields: []string{
					"label",
					"keyword",
				},
			},
			{
				Type: solr.MexOrdinalAxisType,
				Name: "testAxis",
				Fields: []string{
					"label",
					"keyword",
				},
			},
		},
	}

	focusAndAxisSearchConfigWithHierarchyAxis := &sharedSearchConfig.SearchConfigList{
		SearchConfigs: []*sharedSearchConfig.SearchConfigObject{
			{
				Type: solr.MexSearchFocusType,
				Name: "default",
				Fields: []string{
					"label",
					"keyword",
				},
			},
			{
				Type: solr.MexHierarchyAxisType,
				Name: "hierarchyAxis",
				Fields: []string{
					"hier1",
					"hier2",
				},
			},
		},
	}

	tests := []struct {
		name                 string
		fieldRepo            fields.FieldRepo
		searchConfigElements *sharedSearchConfig.SearchConfigList
		want                 *solr.SchemaUpdates
		wantErr              bool
	}{
		{
			name:                 "If no fields have been manually configured, only backing fields for the predefined fields 'entityName' and 'createdAt' are created",
			fieldRepo:            noFieldsRepo,
			searchConfigElements: emptySearchConfig,
			want: &solr.SchemaUpdates{
				FieldDefs: []solr.FieldDef{
					{
						Name:        "entityName",
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						MultiValued: false,
						Stored:      true,
						DocValues:   false,
					},
					{
						Name:        "entityName___normalized",
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						MultiValued: false,
						Stored:      true,
						DocValues:   false,
					},
					{
						Name:        "businessId",
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						Stored:      true,
						MultiValued: false,
						DocValues:   false,
					},
					{
						Name:        "businessId___normalized",
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						Stored:      true,
						MultiValued: false,
						DocValues:   false,
					},
					{
						Name:        "createdAt_raw_value",
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						Stored:      true,
						MultiValued: false,
						DocValues:   false,
					},
					{
						Name:        "createdAt",
						Type:        solr.DefaultSolrTimestampFieldType,
						Indexed:     false,
						Stored:      true,
						MultiValued: false,
						DocValues:   false,
					},
				},
				CopyFieldDefs:    []solr.CopyFieldDef{},
				DynamicFieldDefs: []solr.DynamicFieldDef{},
			},
		},
		{
			name:                 "If no fields are defined but a search focus, an error is returned",
			fieldRepo:            noFieldsRepo,
			searchConfigElements: focusSearchConfig,
			wantErr:              true,
		},
		{
			name:                 "Correctly generates schema when only fields are defined (no axes or foci)",
			fieldRepo:            someFieldsRepo,
			searchConfigElements: emptySearchConfig,
			want: &solr.SchemaUpdates{
				FieldDefs: []solr.FieldDef{
					{
						Name:        "entityName",
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						MultiValued: false,
						Stored:      true,
						DocValues:   false,
					},
					{
						Name:        solr.GetNormalizedBackingFieldName("entityName"),
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						MultiValued: false,
						Stored:      true,
						DocValues:   false,
					},
					{
						Name:        "createdAt",
						Type:        solr.DefaultSolrTimestampFieldType,
						Indexed:     false,
						Stored:      true,
						MultiValued: false,
						DocValues:   false,
					},
					{
						Name:        "createdAt_raw_value",
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						Stored:      true,
						MultiValued: false,
						DocValues:   false,
					},
					{
						Name:        "businessId",
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						Stored:      true,
						MultiValued: false,
						DocValues:   false,
					},
					{
						Name:        solr.GetNormalizedBackingFieldName("businessId"),
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						Stored:      true,
						MultiValued: false,
						DocValues:   false,
					},
					{
						Name:        "label___generic",
						Type:        solr.DefaultSolrTextFieldType,
						MultiValued: true,
						Stored:      true,
					},
					{
						Name:        "label___de",
						Type:        solr.DefaultDeSolrTextFieldType,
						MultiValued: true,
						Stored:      true,
					},
					{
						Name:        "label___en",
						Type:        solr.DefaultEnSolrTextFieldType,
						MultiValued: true,
						Stored:      true,
					},
					{
						Name:        "label___" + solr.PrefixFocusPostfix,
						Type:        solr.DefaultPrefixSolrTextFieldType,
						MultiValued: true,
						Stored:      true,
					},
					{
						Name:        "label___" + solr.RawValuePostfix,
						Type:        solr.DefaultRawSolrTextFieldType,
						MultiValued: true,
						Stored:      true,
					},
					{
						Name:        "label___" + solr.NormFocusPostfix,
						Type:        solr.DefaultSolrStringFieldType,
						MultiValued: true,
						Stored:      true,
					},
					{
						Name:        "keyword",
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						MultiValued: true,
						Stored:      true,
						DocValues:   false,
					},
					{
						Name:        solr.GetNormalizedBackingFieldName("keyword"),
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						MultiValued: true,
						Stored:      true,
						DocValues:   false,
					},
					{
						Name:        "created",
						Type:        solr.DefaultSolrTimestampFieldType,
						Indexed:     false,
						Stored:      true,
						MultiValued: false,
						DocValues:   false,
					},
					{
						Name:        "created_raw_value",
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						Stored:      true,
						MultiValued: false,
						DocValues:   false,
					},
				},
				CopyFieldDefs:    []solr.CopyFieldDef{},
				DynamicFieldDefs: []solr.DynamicFieldDef{},
			},
		},
		{
			name:                 "Correctly generates schema when fields and a search focus are defined",
			fieldRepo:            someFieldsRepo,
			searchConfigElements: focusSearchConfig,
			want: &solr.SchemaUpdates{
				FieldDefs: []solr.FieldDef{
					{
						Name:        "entityName",
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						MultiValued: false,
						Stored:      true,
						DocValues:   false,
					},
					{
						Name:        solr.GetNormalizedBackingFieldName("entityName"),
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						MultiValued: false,
						Stored:      true,
						DocValues:   false,
					},
					{
						Name:        "createdAt",
						Type:        solr.DefaultSolrTimestampFieldType,
						Indexed:     false,
						Stored:      true,
						MultiValued: false,
						DocValues:   false,
					},
					{
						Name:        "createdAt_raw_value",
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						Stored:      true,
						MultiValued: false,
						DocValues:   false,
					},
					{
						Name:        "businessId",
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						Stored:      true,
						MultiValued: false,
						DocValues:   false,
					},
					{
						Name:        solr.GetNormalizedBackingFieldName("businessId"),
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						Stored:      true,
						MultiValued: false,
						DocValues:   false,
					},
					{
						Name:        "label___generic",
						Type:        solr.DefaultSolrTextFieldType,
						MultiValued: true,
						Stored:      true,
					},
					{
						Name:        "label___de",
						Type:        solr.DefaultDeSolrTextFieldType,
						MultiValued: true,
						Stored:      true,
					},
					{
						Name:        "label___en",
						Type:        solr.DefaultEnSolrTextFieldType,
						MultiValued: true,
						Stored:      true,
					},
					{
						Name:        "label___" + solr.PrefixFocusPostfix,
						Type:        solr.DefaultPrefixSolrTextFieldType,
						MultiValued: true,
						Stored:      true,
					},
					{
						Name:        "label___" + solr.RawValuePostfix,
						Type:        solr.DefaultRawSolrTextFieldType,
						MultiValued: true,
						Stored:      true,
					},
					{
						Name:        "label___" + solr.NormFocusPostfix,
						Type:        solr.DefaultSolrStringFieldType,
						MultiValued: true,
						Stored:      true,
					},
					{
						Name:        "keyword",
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						MultiValued: true,
						Stored:      true,
						DocValues:   false,
					},
					{
						Name:        solr.GetNormalizedBackingFieldName("keyword"),
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						MultiValued: true,
						Stored:      true,
						DocValues:   false,
					},
					{
						Name:        "created",
						Type:        solr.DefaultSolrTimestampFieldType,
						Indexed:     false,
						Stored:      true,
						MultiValued: false,
						DocValues:   false,
					},
					{
						Name:        "created_raw_value",
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						Stored:      true,
						MultiValued: false,
						DocValues:   false,
					},
					{
						Name:        "testFocus_search_focus___generic",
						Type:        solr.DefaultSolrTextFieldType,
						MultiValued: true,
						Indexed:     true,
					},
					{
						Name:        "testFocus_search_focus___de",
						Type:        solr.DefaultDeSolrTextFieldType,
						MultiValued: true,
						Indexed:     true,
					},
					{
						Name:        "testFocus_search_focus___en",
						Type:        solr.DefaultEnSolrTextFieldType,
						MultiValued: true,
						Indexed:     true,
					},
					{
						Name:        "testFocus_search_focus___prefix",
						Type:        solr.DefaultPrefixSolrTextFieldType,
						MultiValued: true,
						Indexed:     true,
					},
					{
						Name:        "testFocus_search_focus___unanalyzed",
						Type:        solr.DefaultRawSolrTextFieldType,
						MultiValued: true,
						Indexed:     true,
					},
				},
				CopyFieldDefs: []solr.CopyFieldDef{
					{
						Source:      "label___generic",
						Destination: []string{"testFocus_search_focus___generic"},
					},
					{
						Source:      "label___de",
						Destination: []string{"testFocus_search_focus___de"},
					},
					{
						Source:      "label___en",
						Destination: []string{"testFocus_search_focus___en"},
					},
					{
						Source:      "label___prefix",
						Destination: []string{"testFocus_search_focus___prefix"},
					},
					{
						Source:      "label___" + solr.RawValuePostfix,
						Destination: []string{"testFocus_search_focus___" + solr.RawValuePostfix},
					},
					{
						Source:      "keyword",
						Destination: []string{"testFocus_search_focus___generic"},
					},
					{
						Source:      "keyword",
						Destination: []string{"testFocus_search_focus___prefix"},
					},
					{
						Source:      "keyword",
						Destination: []string{"testFocus_search_focus___unanalyzed"},
					},
				},
				DynamicFieldDefs: []solr.DynamicFieldDef{},
			},
		},
		{
			name:                 "Correctly generates schema when fields and an ordinal axis are defined",
			fieldRepo:            someFieldsRepo,
			searchConfigElements: axisSearchConfig,
			want: &solr.SchemaUpdates{
				FieldDefs: []solr.FieldDef{
					{
						Name:        "entityName",
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						MultiValued: false,
						Stored:      true,
						DocValues:   false,
					},
					{
						Name:        solr.GetNormalizedBackingFieldName("entityName"),
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						MultiValued: false,
						Stored:      true,
						DocValues:   false,
					},
					{
						Name:        "createdAt",
						Type:        solr.DefaultSolrTimestampFieldType,
						Indexed:     false,
						Stored:      true,
						MultiValued: false,
						DocValues:   false,
					},
					{
						Name:        "createdAt_raw_value",
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						Stored:      true,
						MultiValued: false,
						DocValues:   false,
					},
					{
						Name:        "businessId",
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						Stored:      true,
						MultiValued: false,
						DocValues:   false,
					},
					{
						Name:        solr.GetNormalizedBackingFieldName("businessId"),
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						Stored:      true,
						MultiValued: false,
						DocValues:   false,
					},
					{
						Name:        "label___generic",
						Type:        solr.DefaultSolrTextFieldType,
						MultiValued: true,
						Stored:      true,
					},
					{
						Name:        "label___de",
						Type:        solr.DefaultDeSolrTextFieldType,
						MultiValued: true,
						Stored:      true,
					},
					{
						Name:        "label___en",
						Type:        solr.DefaultEnSolrTextFieldType,
						MultiValued: true,
						Stored:      true,
					},
					{
						Name:        "label___" + solr.PrefixFocusPostfix,
						Type:        solr.DefaultPrefixSolrTextFieldType,
						MultiValued: true,
						Stored:      true,
					},
					{
						Name:        "label___" + solr.RawValuePostfix,
						Type:        solr.DefaultRawSolrTextFieldType,
						MultiValued: true,
						Stored:      true,
					},
					{
						Name:        "label___" + solr.NormFocusPostfix,
						Type:        solr.DefaultSolrStringFieldType,
						MultiValued: true,
						Stored:      true,
					},
					{
						Name:        "keyword",
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						MultiValued: true,
						Stored:      true,
						DocValues:   false,
					},
					{
						Name:        solr.GetNormalizedBackingFieldName("keyword"),
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						MultiValued: true,
						Stored:      true,
						DocValues:   false,
					},
					{
						Name:        "created",
						Type:        solr.DefaultSolrTimestampFieldType,
						Indexed:     false,
						Stored:      true,
						MultiValued: false,
						DocValues:   false,
					},
					{
						Name:        "created_raw_value",
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						Stored:      true,
						MultiValued: false,
						DocValues:   false,
					},
					{
						Name:        "testAxis_ordinal_facet_axis",
						Type:        solr.DefaultSolrSortableTextFieldType,
						MultiValued: true,
						Indexed:     true,
						DocValues:   true,
					},
					{
						Name:        "testAxis_ordinal_sort_axis",
						Type:        solr.DefaultSolrSortableTextFieldType,
						MultiValued: true,
						Indexed:     true,
						DocValues:   true,
					},
				},
				CopyFieldDefs: []solr.CopyFieldDef{
					{
						Source:      "keyword",
						Destination: []string{"testAxis_ordinal_facet_axis"},
					},
					{
						Source:      "keyword___normalized",
						Destination: []string{"testAxis_ordinal_sort_axis"},
					},
					{
						Source:      "label___generic",
						Destination: []string{"testAxis_ordinal_facet_axis"},
					},
					{
						Source:      "label___de",
						Destination: []string{"testAxis_ordinal_facet_axis"},
					},
					{
						Source:      "label___en",
						Destination: []string{"testAxis_ordinal_facet_axis"},
					},
					{
						Source:      "label___normalized",
						Destination: []string{"testAxis_ordinal_sort_axis"},
					},
				},
				DynamicFieldDefs: []solr.DynamicFieldDef{},
			},
		},
		{
			name:                 "Correctly generates schema using fields, search foci, and ordinal axes",
			fieldRepo:            someFieldsRepo,
			searchConfigElements: focusAndAxisSearchConfig,
			want: &solr.SchemaUpdates{
				FieldDefs: []solr.FieldDef{
					{
						Name:        "entityName",
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						MultiValued: false,
						Stored:      true,
						DocValues:   false,
					},
					{
						Name:        solr.GetNormalizedBackingFieldName("entityName"),
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						MultiValued: false,
						Stored:      true,
						DocValues:   false,
					},
					{
						Name:        "createdAt",
						Type:        solr.DefaultSolrTimestampFieldType,
						Indexed:     false,
						Stored:      true,
						MultiValued: false,
						DocValues:   false,
					},
					{
						Name:        "createdAt_raw_value",
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						Stored:      true,
						MultiValued: false,
						DocValues:   false,
					},
					{
						Name:        "businessId",
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						Stored:      true,
						MultiValued: false,
						DocValues:   false,
					},
					{
						Name:        solr.GetNormalizedBackingFieldName("businessId"),
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						Stored:      true,
						MultiValued: false,
						DocValues:   false,
					},
					{
						Name:        "label___generic",
						Type:        solr.DefaultSolrTextFieldType,
						MultiValued: true,
						Stored:      true,
					},
					{
						Name:        "label___de",
						Type:        solr.DefaultDeSolrTextFieldType,
						MultiValued: true,
						Stored:      true,
					},
					{
						Name:        "label___en",
						Type:        solr.DefaultEnSolrTextFieldType,
						MultiValued: true,
						Stored:      true,
					},
					{
						Name:        "label___" + solr.PrefixFocusPostfix,
						Type:        solr.DefaultPrefixSolrTextFieldType,
						MultiValued: true,
						Stored:      true,
					},
					{
						Name:        "label___" + solr.RawValuePostfix,
						Type:        solr.DefaultRawSolrTextFieldType,
						MultiValued: true,
						Stored:      true,
					},
					{
						Name:        "label___" + solr.NormFocusPostfix,
						Type:        solr.DefaultSolrStringFieldType,
						MultiValued: true,
						Stored:      true,
					},
					{
						Name:        "keyword",
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						MultiValued: true,
						Stored:      true,
						DocValues:   false,
					},
					{
						Name:        solr.GetNormalizedBackingFieldName("keyword"),
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						MultiValued: true,
						Stored:      true,
						DocValues:   false,
					},
					{
						Name:        "created",
						Type:        solr.DefaultSolrTimestampFieldType,
						Indexed:     false,
						Stored:      true,
						MultiValued: false,
						DocValues:   false,
					},
					{
						Name:        "created_raw_value",
						Type:        solr.DefaultSolrStringFieldType,
						Indexed:     false,
						Stored:      true,
						MultiValued: false,
						DocValues:   false,
					},
					{
						Name:        "testAxis_ordinal_facet_axis",
						Type:        solr.DefaultSolrSortableTextFieldType,
						MultiValued: true,
						Indexed:     true,
						DocValues:   true,
					},
					{
						Name:        "testAxis_ordinal_sort_axis",
						Type:        solr.DefaultSolrSortableTextFieldType,
						MultiValued: true,
						Indexed:     true,
						DocValues:   true,
					},
					{
						Name:        "testFocus_search_focus___generic",
						Type:        solr.DefaultSolrTextFieldType,
						MultiValued: true,
						Indexed:     true,
					},
					{
						Name:        "testFocus_search_focus___de",
						Type:        solr.DefaultDeSolrTextFieldType,
						MultiValued: true,
						Indexed:     true,
					},
					{
						Name:        "testFocus_search_focus___en",
						Type:        solr.DefaultEnSolrTextFieldType,
						MultiValued: true,
						Indexed:     true,
					},
					{
						Name:        "testFocus_search_focus___prefix",
						Type:        solr.DefaultPrefixSolrTextFieldType,
						MultiValued: true,
						Indexed:     true,
					},
					{
						Name:        "testFocus_search_focus___unanalyzed",
						Type:        solr.DefaultRawSolrTextFieldType,
						MultiValued: true,
						Indexed:     true,
					},
				},
				CopyFieldDefs: []solr.CopyFieldDef{
					{
						Source:      "keyword",
						Destination: []string{"testAxis_ordinal_facet_axis"},
					},
					{
						Source:      "keyword___normalized",
						Destination: []string{"testAxis_ordinal_sort_axis"},
					},
					{
						Source:      "label___generic",
						Destination: []string{"testAxis_ordinal_facet_axis"},
					},
					{
						Source:      "label___de",
						Destination: []string{"testAxis_ordinal_facet_axis"},
					},
					{
						Source:      "label___en",
						Destination: []string{"testAxis_ordinal_facet_axis"},
					},
					{
						Source:      "label___normalized",
						Destination: []string{"testAxis_ordinal_sort_axis"},
					},
					{
						Source:      "label___generic",
						Destination: []string{"testFocus_search_focus___generic"},
					},
					{
						Source:      "label___de",
						Destination: []string{"testFocus_search_focus___de"},
					},
					{
						Source:      "label___en",
						Destination: []string{"testFocus_search_focus___en"},
					},
					{
						Source:      "label___prefix",
						Destination: []string{"testFocus_search_focus___prefix"},
					},
					{
						Source:      "label___" + solr.RawValuePostfix,
						Destination: []string{"testFocus_search_focus___" + solr.RawValuePostfix},
					},
					{
						Source:      "keyword",
						Destination: []string{"testFocus_search_focus___generic"},
					},
					{
						Source:      "keyword",
						Destination: []string{"testFocus_search_focus___prefix"},
					},
					{
						Source:      "keyword",
						Destination: []string{"testFocus_search_focus___unanalyzed"},
					},
				},
				DynamicFieldDefs: []solr.DynamicFieldDef{},
			},
		},
		{
			name:                 "Returns an error if a hierarchy axis uses fields with different configurations",
			fieldRepo:            brokenHierarchyFieldsRepo,
			searchConfigElements: focusAndAxisSearchConfigWithHierarchyAxis,
			wantErr:              true,
		},
	}
	solrFieldCreationHooks, _ := hooks.NewSolrFieldCreationHooks(hooks.SolrFieldCreationHooksConfig{})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fieldDefs, _ := tt.fieldRepo.ListFieldDefs(context.TODO())
			got, err := generateSolrSchema(context.TODO(), solrFieldCreationHooks, fieldDefs,
				tt.searchConfigElements)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateSolrSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err == nil {
				checkSolrFields(t, got.FieldDefs, tt.want.FieldDefs)
				checkSolrCopyFields(t, got.CopyFieldDefs, tt.want.CopyFieldDefs)
			}
		})
	}
}
