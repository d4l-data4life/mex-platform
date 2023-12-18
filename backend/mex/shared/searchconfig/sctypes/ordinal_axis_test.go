package sctypes

import (
	"sort"
	"testing"

	sharedSearchConfig "github.com/d4l-data4life/mex/mex/shared/searchconfig"
	"github.com/d4l-data4life/mex/mex/shared/solr"
)

func getExpectedTestFieldsForAxis(axisName string, fieldType string) []solr.FieldDef {
	fn := solr.GetOrdinalAxisFacetAndFilterFieldName(axisName)
	sn := solr.GetOrdinalAxisSortFieldName(axisName)
	expectedFields := []solr.FieldDef{
		{
			Name:         fn,
			Type:         fieldType,
			Stored:       false,
			Indexed:      true,
			MultiValued:  true,
			DocValues:    true,
			Uninvertible: false,
		},
		{
			Name:         sn,
			Type:         fieldType,
			Stored:       false,
			Indexed:      true,
			MultiValued:  true,
			DocValues:    true,
			Uninvertible: false,
		},
	}
	return expectedFields
}

func TestOrdinalAxisType_GetSolrBackingFieldDefs(t *testing.T) {
	tests := []struct {
		name            string
		searchFocusElem *sharedSearchConfig.SearchConfigObject
		mexFields       solr.MexFieldBackingInfoMap
	}{
		{
			name: "The name of the actually generated backing field is returned",
			searchFocusElem: &sharedSearchConfig.SearchConfigObject{
				Type:   solr.MexOrdinalAxisType,
				Name:   "testAxis",
				Fields: []string{"category", "created"},
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
				"created": {
					MexType: "timestamp",
					BackingFields: []solr.MexBackingFieldWiringInfo{
						{
							Name:     "created",
							Category: solr.GenericLangBaseFieldCategory,
						},
						{
							Name:     solr.GetRawValTimestampName("category"),
							Category: solr.RawContentBaseFieldCategory,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asType := &OrdinalAxisType{}
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

func TestOrdinalAxisType_GetSolrBackingFields(t *testing.T) {
	_, _ = solr.GetLangSpecificFieldName("seeThis", solr.GenericLangAbbrev)
	descriptionNameGeneric, _ := solr.GetLangSpecificFieldName("description", solr.GenericLangAbbrev)
	descriptionNameDe, _ := solr.GetLangSpecificFieldName("description", solr.GermanLangAbbrev)
	descriptionNameEn, _ := solr.GetLangSpecificFieldName("description", solr.EnglishLangAbbrev)
	descriptionNameNormalized := solr.GetNormalizedBackingFieldName("description")
	deHullLabelName, _ := solr.GetLangSpecificFieldName(solr.GetTransitiveHullFieldName("unitCode"), solr.GermanLangAbbrev)
	enHullLabelName, _ := solr.GetLangSpecificFieldName(solr.GetTransitiveHullFieldName("unitCode"), solr.EnglishLangAbbrev)

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
				Type:   solr.MexOrdinalAxisType,
				Name:   "testAxis",
				Fields: []string{"unknown"},
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
			name: "An error is returned if an ordinal axis contains a field of kind hierarchy",
			searchFocusElem: &sharedSearchConfig.SearchConfigObject{
				Type:   solr.MexOrdinalAxisType,
				Name:   "hierarchyAxis",
				Fields: []string{"unitCode"},
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
			name: "If the axis contains link fields only, the base field is copied to the facet/filter field (but not to sort field)",
			searchFocusElem: &sharedSearchConfig.SearchConfigObject{
				Type:   solr.MexOrdinalAxisType,
				Name:   "testAxis",
				Fields: []string{"seeThis"},
			},
			mexFields: solr.MexFieldBackingInfoMap{
				"seeThis": {
					MexType: "link",
					BackingFields: []solr.MexBackingFieldWiringInfo{
						{
							Name:     "seeThis",
							Category: solr.GenericLangBaseFieldCategory,
						},
					},
				},
			},
			wantFields: getExpectedTestFieldsForAxis("testAxis", solr.DefaultSolrSortableTextFieldType),
			wantCopyFields: []solr.CopyFieldDef{
				{
					Source:      "seeThis",
					Destination: []string{solr.GetOrdinalAxisFacetAndFilterFieldName("testAxis")},
				},
			},
		},
		{
			name: "If the axis contains string fields only, sortable text facet and sort backing fields are generated and the raw and normalized content, respectively, is copied into them",
			searchFocusElem: &sharedSearchConfig.SearchConfigObject{
				Type:   solr.MexOrdinalAxisType,
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
			wantFields: getExpectedTestFieldsForAxis("testAxis", solr.DefaultSolrSortableTextFieldType),
			wantCopyFields: []solr.CopyFieldDef{
				{
					Source:      "category",
					Destination: []string{solr.GetOrdinalAxisFacetAndFilterFieldName("testAxis")},
				},
				{
					Source:      solr.GetNormalizedBackingFieldName("category"),
					Destination: []string{solr.GetOrdinalAxisSortFieldName("testAxis")},
				},
			},
		},
		{
			name: "If the axis contains text fields only, sortable text facet and sort backing fields are generated and the generic-language, DE, and EN content is copied into the former and the normalized content into the latter",
			searchFocusElem: &sharedSearchConfig.SearchConfigObject{
				Type:   solr.MexOrdinalAxisType,
				Name:   "testAxis",
				Fields: []string{"description"},
			},
			mexFields: solr.MexFieldBackingInfoMap{
				"description": {
					MexType: "text",
					BackingFields: []solr.MexBackingFieldWiringInfo{
						{
							Name:     descriptionNameGeneric,
							Category: solr.GenericLangBaseFieldCategory,
						},
						{
							Name:     descriptionNameDe,
							Category: solr.GermanLangBaseFieldCategory,
						},
						{
							Name:     descriptionNameEn,
							Category: solr.EnglishLangBaseFieldCategory,
						},
						{
							Name:     descriptionNameNormalized,
							Category: solr.NormalizedBaseFieldCategory,
						},
					},
				},
			},
			wantFields: getExpectedTestFieldsForAxis("testAxis", solr.DefaultSolrSortableTextFieldType),
			wantCopyFields: []solr.CopyFieldDef{
				{
					Source:      descriptionNameGeneric,
					Destination: []string{solr.GetOrdinalAxisFacetAndFilterFieldName("testAxis")},
				},
				{
					Source:      descriptionNameDe,
					Destination: []string{solr.GetOrdinalAxisFacetAndFilterFieldName("testAxis")},
				},
				{
					Source:      descriptionNameEn,
					Destination: []string{solr.GetOrdinalAxisFacetAndFilterFieldName("testAxis")},
				},
				{
					Source:      descriptionNameNormalized,
					Destination: []string{solr.GetOrdinalAxisSortFieldName("testAxis")},
				},
			},
		},
		{
			name: "If the axis contains link fields only, number facet and sort backing fields are generated and the raw content is copied into them",
			searchFocusElem: &sharedSearchConfig.SearchConfigObject{
				Type:   solr.MexOrdinalAxisType,
				Name:   "testAxis",
				Fields: []string{"count"},
			},
			mexFields: solr.MexFieldBackingInfoMap{
				"count": {
					MexType: "number",
					BackingFields: []solr.MexBackingFieldWiringInfo{
						{
							Name:     "count",
							Category: solr.GenericLangBaseFieldCategory,
						},
					},
				},
			},
			wantFields: getExpectedTestFieldsForAxis("testAxis", solr.DefaultSolrNumberFieldType),
			wantCopyFields: []solr.CopyFieldDef{
				{
					Source:      "count",
					Destination: []string{solr.GetOrdinalAxisFacetAndFilterFieldName("testAxis")},
				},
				{
					Source:      "count",
					Destination: []string{solr.GetOrdinalAxisSortFieldName("testAxis")},
				},
			},
		},
		{
			name: "If the axis contains coding fields only, a string sort backing fields is generated and the raw content is copied into it",
			searchFocusElem: &sharedSearchConfig.SearchConfigObject{
				Type:   solr.MexOrdinalAxisType,
				Name:   "codingAxis",
				Fields: []string{"meshCode"},
			},
			mexFields: solr.MexFieldBackingInfoMap{
				"meshCode": {
					MexType: "coding",
					BackingFields: []solr.MexBackingFieldWiringInfo{
						{
							Name:     "meshCode",
							Category: solr.GenericLangBaseFieldCategory,
						},
						{
							Name:     solr.GetDisplayFieldName("meshCode", solr.GermanLangAbbrev),
							Category: solr.GermanLangBaseFieldCategory,
						},
						{
							Name:     solr.GetDisplayFieldName("meshCode", solr.EnglishLangAbbrev),
							Category: solr.EnglishLangBaseFieldCategory,
						},
					},
				},
			},
			wantFields: getExpectedTestFieldsForAxis("codingAxis", solr.DefaultSolrSortableTextFieldType),
			wantCopyFields: []solr.CopyFieldDef{
				{
					Source:      "meshCode",
					Destination: []string{solr.GetOrdinalAxisSortFieldName("codingAxis")},
				},
			},
		},
		{
			name: "If the axis contains timestamp fields only, timestamp facet and sort backing fields are generated and the raw content is copied into them",
			searchFocusElem: &sharedSearchConfig.SearchConfigObject{
				Type:   solr.MexOrdinalAxisType,
				Name:   "testAxis",
				Fields: []string{"created"},
			},
			mexFields: solr.MexFieldBackingInfoMap{
				"created": {
					MexType: "timestamp",
					BackingFields: []solr.MexBackingFieldWiringInfo{
						{
							Name:     "created",
							Category: solr.GenericLangBaseFieldCategory,
						},
						{
							Name:     solr.GetRawValTimestampName("created"),
							Category: solr.RawContentBaseFieldCategory,
						},
					},
				},
			},
			wantFields: getExpectedTestFieldsForAxis("testAxis", solr.DefaultSolrTimestampFieldType),
			wantCopyFields: []solr.CopyFieldDef{
				{
					Source:      "created",
					Destination: []string{solr.GetOrdinalAxisFacetAndFilterFieldName("testAxis")},
				},
				{
					Source:      "created",
					Destination: []string{solr.GetOrdinalAxisSortFieldName("testAxis")},
				},
			},
		},
		{
			name: "If the axis contains fields of different MEx types, a sortable text facet and sort backing fields are generated and the relevant content is copied into them",
			searchFocusElem: &sharedSearchConfig.SearchConfigObject{
				Type:   solr.MexOrdinalAxisType,
				Name:   "testAxis",
				Fields: []string{"category", "created"},
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
				"created": {
					MexType: "timestamp",
					BackingFields: []solr.MexBackingFieldWiringInfo{
						{
							Name:     "created",
							Category: solr.GenericLangBaseFieldCategory,
						},
						{
							Name:     solr.GetRawValTimestampName("created"),
							Category: solr.RawContentBaseFieldCategory,
						},
					},
				},
			},
			wantFields: getExpectedTestFieldsForAxis("testAxis", solr.DefaultSolrSortableTextFieldType),
			wantCopyFields: []solr.CopyFieldDef{
				{
					Source:      "category",
					Destination: []string{solr.GetOrdinalAxisFacetAndFilterFieldName("testAxis")},
				},
				{
					Source:      solr.GetNormalizedBackingFieldName("category"),
					Destination: []string{solr.GetOrdinalAxisSortFieldName("testAxis")},
				},
				{
					Source:      "created",
					Destination: []string{solr.GetOrdinalAxisFacetAndFilterFieldName("testAxis")},
				},
				{
					Source:      "created",
					Destination: []string{solr.GetOrdinalAxisSortFieldName("testAxis")},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asType := &OrdinalAxisType{}
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

func TestOrdinalAxisType_GetSolrSearchFieldNames(t *testing.T) {
	t.Run("GetSolrSearchFieldNames() panics", func(t *testing.T) {
		scType := &OrdinalAxisType{}
		_, err := scType.GetSolrSearchFieldNames("testAxis", false, false)
		if err == nil {
			t.Errorf("GetSolrBackingFields() should return an error but did not")
		}
	})
}
