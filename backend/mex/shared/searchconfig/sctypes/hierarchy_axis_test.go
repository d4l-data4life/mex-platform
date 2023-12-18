package sctypes

import (
	"sort"
	"testing"

	sharedSearchConfig "github.com/d4l-data4life/mex/mex/shared/searchconfig"
	"github.com/d4l-data4life/mex/mex/shared/solr"
)

func TestHierarchyAxisType_GetSolrBackingFieldDefs(t *testing.T) {
	deHullLabelName, _ := solr.GetLangSpecificFieldName(solr.GetTransitiveHullFieldName("unitCode"), solr.GermanLangAbbrev)
	enHullLabelName, _ := solr.GetLangSpecificFieldName(solr.GetTransitiveHullFieldName("unitCode"), solr.EnglishLangAbbrev)

	tests := []struct {
		name            string
		searchFocusElem *sharedSearchConfig.SearchConfigObject
		mexFields       solr.MexFieldBackingInfoMap
		fieldKindMap    map[string]string
	}{
		{
			name: "The name of the actually generated backing fields are returned",
			searchFocusElem: &sharedSearchConfig.SearchConfigObject{
				Type: solr.MexHierarchyAxisType,
				Name: "testAxis",
				Fields: []string{
					"unitCode",
				},
			},
			mexFields: solr.MexFieldBackingInfoMap{
				"unitCode": {
					MexType: "hierarchy",
					BackingFields: []solr.MexBackingFieldWiringInfo{
						{
							Name:     "unitCode",
							Category: solr.GenericLangBaseFieldCategory,
						},
						{
							Name:     solr.GetTransitiveHullFieldName("unitCode"),
							Category: solr.ParentCodesBaseFieldCategory,
						},
						{
							Name:     deHullLabelName,
							Category: solr.GermanLangBaseFieldCategory,
						},
						{
							Name:     enHullLabelName,
							Category: solr.EnglishLangBaseFieldCategory,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asType := &HierarchyAxisType{}
			gotFields, _, _ := asType.GetSolrBackingFields(tt.searchFocusElem, tt.mexFields)
			var generatedNames []string
			for _, f := range gotFields {
				generatedNames = append(generatedNames, f.Name)
			}
			sort.Slice(generatedNames, func(i, j int) bool {
				return generatedNames[i] < generatedNames[j]
			})
		})
	}
}

func TestHierarchyAxisType_GetSolrBackingFields(t *testing.T) {
	deHullLabelName, _ := solr.GetLangSpecificFieldName(solr.GetTransitiveHullFieldName("unitCode"), solr.GermanLangAbbrev)
	enHullLabelName, _ := solr.GetLangSpecificFieldName(solr.GetTransitiveHullFieldName("unitCode"), solr.EnglishLangAbbrev)
	otherDeHullLabelName, _ := solr.GetLangSpecificFieldName(solr.GetTransitiveHullFieldName("otherUnitCode"), solr.GermanLangAbbrev)
	otherEnHullLabelName, _ := solr.GetLangSpecificFieldName(solr.GetTransitiveHullFieldName("otherUnitCode"), solr.EnglishLangAbbrev)

	tests := []struct {
		name            string
		searchFocusElem *sharedSearchConfig.SearchConfigObject
		mexFields       solr.MexFieldBackingInfoMap
		wantFields      []solr.FieldDef
		wantCopyFields  []solr.CopyFieldDef
		wantErr         bool
	}{
		{
			name: "An error is returned if axis contains a field not listed in the field map",
			searchFocusElem: &sharedSearchConfig.SearchConfigObject{
				Type:   solr.MexHierarchyAxisType,
				Name:   "testAxis",
				Fields: []string{"unknown"},
			},
			mexFields: solr.MexFieldBackingInfoMap{
				"unitCode": {
					MexType: "hierarchy",
					BackingFields: []solr.MexBackingFieldWiringInfo{
						{
							Name:     "unitCode",
							Category: solr.GenericLangBaseFieldCategory,
						},
						{
							Name:     solr.GetTransitiveHullFieldName("unitCode"),
							Category: solr.ParentCodesBaseFieldCategory,
						},
						{
							Name:     deHullLabelName,
							Category: solr.GermanLangBaseFieldCategory,
						},
						{
							Name:     enHullLabelName,
							Category: solr.EnglishLangBaseFieldCategory,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Multiple fields in axis is allowed",
			searchFocusElem: &sharedSearchConfig.SearchConfigObject{
				Type: solr.MexHierarchyAxisType,
				Name: "testAxis",
				Fields: []string{
					"unitCode",
					"otherUnitCode",
				},
			},
			mexFields: solr.MexFieldBackingInfoMap{
				"unitCode": {
					MexType: "hierarchy",
					BackingFields: []solr.MexBackingFieldWiringInfo{
						{
							Name:     "unitCode",
							Category: solr.GenericLangBaseFieldCategory,
						},
						{
							Name:     solr.GetTransitiveHullFieldName("unitCode"),
							Category: solr.ParentCodesBaseFieldCategory,
						},
						{
							Name:     deHullLabelName,
							Category: solr.GermanLangBaseFieldCategory,
						},
						{
							Name:     enHullLabelName,
							Category: solr.EnglishLangBaseFieldCategory,
						},
					},
				},
				"otherUnitCode": {
					MexType: "hierarchy",
					BackingFields: []solr.MexBackingFieldWiringInfo{
						{
							Name:     "otherUnitCode",
							Category: solr.GenericLangBaseFieldCategory,
						},
						{
							Name:     solr.GetTransitiveHullFieldName("otherUnitCode"),
							Category: solr.ParentCodesBaseFieldCategory,
						},
						{
							Name:     otherDeHullLabelName,
							Category: solr.GermanLangBaseFieldCategory,
						},
						{
							Name:     otherEnHullLabelName,
							Category: solr.EnglishLangBaseFieldCategory,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "An error is returned if axis contains a non-hierarchy field",
			searchFocusElem: &sharedSearchConfig.SearchConfigObject{
				Type:   solr.MexHierarchyAxisType,
				Name:   "testAxis",
				Fields: []string{"category"},
			},
			mexFields: solr.MexFieldBackingInfoMap{
				"category": {
					MexType: "string",
					BackingFields: []solr.MexBackingFieldWiringInfo{
						{
							Name:     "category",
							Category: solr.GenericLangBaseFieldCategory,
						},
						{
							Name:     solr.GetNormalizedBackingFieldName("category"),
							Category: solr.NormalizedBaseFieldCategory,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Three backing fields are generated (single and parent-codes as string field, plus sorting field as sortable text)",
			searchFocusElem: &sharedSearchConfig.SearchConfigObject{
				Type: solr.MexHierarchyAxisType,
				Name: "testAxis",
				Fields: []string{
					"unitCode",
				},
			},
			mexFields: solr.MexFieldBackingInfoMap{
				"unitCode": {
					MexType: "hierarchy",
					BackingFields: []solr.MexBackingFieldWiringInfo{
						{
							Name:     "unitCode",
							Category: solr.GenericLangBaseFieldCategory,
						},
						{
							Name:     solr.GetTransitiveHullFieldName("unitCode"),
							Category: solr.ParentCodesBaseFieldCategory,
						},
						{
							Name:     deHullLabelName,
							Category: solr.GermanLangBaseFieldCategory,
						},
						{
							Name:     enHullLabelName,
							Category: solr.EnglishLangBaseFieldCategory,
						},
					},
				},
			},
			wantFields: []solr.FieldDef{
				{
					Name:         solr.GetOrdinalAxisFacetAndFilterFieldName("testAxis"),
					Type:         solr.DefaultSolrStringFieldType,
					Stored:       false,
					Indexed:      true,
					MultiValued:  true,
					DocValues:    true,
					Uninvertible: false,
				},
				{
					Name:         solr.GetOrdinalAxisSortFieldName("testAxis"),
					Type:         solr.DefaultSolrSortableTextFieldType,
					Stored:       false,
					Indexed:      true,
					MultiValued:  true,
					DocValues:    true,
					Uninvertible: false,
				},
				{
					Name:         solr.GetSingleNodeAxisFieldName(solr.GetOrdinalAxisFacetAndFilterFieldName("testAxis")),
					Type:         solr.DefaultSolrStringFieldType,
					Stored:       false,
					Indexed:      true,
					MultiValued:  true,
					DocValues:    false,
					Uninvertible: false,
				},
			},
		},
		{
			name: "Copy fields are generated: single value to single value and sort, transitive hull to faceting",
			searchFocusElem: &sharedSearchConfig.SearchConfigObject{
				Type: solr.MexHierarchyAxisType,
				Name: "testAxis",
				Fields: []string{
					"unitCode",
				},
			},
			mexFields: solr.MexFieldBackingInfoMap{
				"unitCode": {
					MexType: "hierarchy",
					BackingFields: []solr.MexBackingFieldWiringInfo{
						{
							Name:     "unitCode",
							Category: solr.GenericLangBaseFieldCategory,
						},
						{
							Name:     solr.GetTransitiveHullFieldName("unitCode"),
							Category: solr.ParentCodesBaseFieldCategory,
						},
						{
							Name:     deHullLabelName,
							Category: solr.GermanLangBaseFieldCategory,
						},
						{
							Name:     enHullLabelName,
							Category: solr.EnglishLangBaseFieldCategory,
						},
					},
				},
			},
			wantCopyFields: []solr.CopyFieldDef{
				{
					Source:      solr.GetTransitiveHullFieldName("unitCode"),
					Destination: []string{solr.GetOrdinalAxisFacetAndFilterFieldName("testAxis")},
				},
				{
					Source:      "unitCode",
					Destination: []string{solr.GetOrdinalAxisSortFieldName("testAxis")},
				},
				{
					Source:      "unitCode",
					Destination: []string{solr.GetSingleNodeAxisFieldName(solr.GetOrdinalAxisFacetAndFilterFieldName("testAxis"))},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asType := &HierarchyAxisType{}
			gotFields, gotCopyFields, err := asType.GetSolrBackingFields(tt.searchFocusElem, tt.mexFields)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSolrBackingFields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantFields != nil {
				checkSolrFields(t, gotFields, tt.wantFields)
			}
			if tt.wantCopyFields != nil {
				checkSolrCopyFields(t, gotCopyFields, tt.wantCopyFields)
			}
		})
	}
}

func TestHierarchyAxisType_GetSolrSearchFieldNames(t *testing.T) {
	t.Run("GetSolrSearchFieldNames() returns error", func(t *testing.T) {
		scType := &HierarchyAxisType{}
		_, err := scType.GetSolrSearchFieldNames("testAxis", false, false)
		if err == nil {
			t.Errorf("GetSolrBackingFields() should return an error but did not")
		}
	})
}
