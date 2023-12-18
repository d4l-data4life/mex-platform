package sctypes

import (
	"reflect"
	"sort"
	"testing"

	sharedSearchConfig "github.com/d4l-data4life/mex/mex/shared/searchconfig"
	"github.com/d4l-data4life/mex/mex/shared/solr"
	"github.com/d4l-data4life/mex/mex/shared/testutils"
)

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

func getExpectedTestFieldsForFocus(focusName string) []solr.FieldDef {
	var expectedFields []solr.FieldDef
	focusFieldName := solr.GetSearchFocusFieldName(focusName)
	for lc, fieldType := range solr.KnownLanguagesFieldTypeMap {
		fn, _ := solr.GetLangSpecificFieldName(focusFieldName, lc)
		expectedFields = append(expectedFields, solr.FieldDef{
			Name:         fn,
			Type:         fieldType,
			Stored:       false,
			Indexed:      true,
			MultiValued:  true,
			DocValues:    false,
			Uninvertible: false,
		},
		)
	}
	// Add prefix search field
	expectedFields = append(expectedFields, solr.FieldDef{
		Name:         solr.GetPrefixBackingFieldName(focusFieldName),
		Type:         solr.DefaultPrefixSolrTextFieldType,
		Stored:       false,
		Indexed:      true,
		MultiValued:  true,
		DocValues:    false,
		Uninvertible: false,
	},
	)
	// Add unanalyzed search field
	expectedFields = append(expectedFields, solr.FieldDef{
		Name:         solr.GetRawBackingFieldName(focusFieldName),
		Type:         solr.DefaultRawSolrTextFieldType,
		Stored:       false,
		Indexed:      true,
		MultiValued:  true,
		DocValues:    false,
		Uninvertible: false,
	},
	)
	return expectedFields
}

func TestSearchFocusType_GetSolrBackingFieldNames(t *testing.T) {
	fieldName := "description"
	descriptionNameDe, _ := solr.GetLangSpecificFieldName(fieldName, solr.GermanLangAbbrev)
	descriptionNameEn, _ := solr.GetLangSpecificFieldName(fieldName, solr.EnglishLangAbbrev)

	focusName := "testFocus"
	searchFocusFieldNameBase := solr.GetSearchFocusFieldName(focusName)
	searchFocusFieldNameGeneric, _ := solr.GetLangSpecificFieldName(searchFocusFieldNameBase, solr.GenericLangAbbrev)
	searchFocusFieldNameDe, _ := solr.GetLangSpecificFieldName(searchFocusFieldNameBase, solr.GermanLangAbbrev)
	searchFocusFieldNameEn, _ := solr.GetLangSpecificFieldName(searchFocusFieldNameBase, solr.EnglishLangAbbrev)
	searchFocusFieldNamePrefix := solr.GetPrefixBackingFieldName(searchFocusFieldNameBase) + "^" + solr.PrefixBoostFactor
	searchFocusFieldNameUnanalyzed := solr.GetRawBackingFieldName(searchFocusFieldNameBase) + "^" + solr.UnanalyzedBoostFactor
	descriptionNameNormalized := solr.GetNormalizedBackingFieldName("description")

	tests := []struct {
		name            string
		searchFocusElem *sharedSearchConfig.SearchConfigObject
		mexFields       solr.MexFieldBackingInfoMap
		searchFocusName string
		includePrefix   bool
		isPhraseOnly    bool
		wantErr         bool
		expectedNames   []string
	}{
		{
			name: "Empty search focus name causes error",
			searchFocusElem: &sharedSearchConfig.SearchConfigObject{
				Type:   solr.MexSearchFocusType,
				Name:   focusName,
				Fields: []string{"description"},
			},
			mexFields: solr.MexFieldBackingInfoMap{
				"description": {
					MexType: "text",
					BackingFields: []solr.MexBackingFieldWiringInfo{
						{
							Name:     "description",
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
			searchFocusName: "",
			includePrefix:   true,
			isPhraseOnly:    false,
			wantErr:         true,
		},
		{
			name: "Query IS NOT phrase-only, prefix-search IS requested: Names of all backing fields returned (with boosts)",
			searchFocusElem: &sharedSearchConfig.SearchConfigObject{
				Type:   solr.MexSearchFocusType,
				Name:   focusName,
				Fields: []string{"description"},
			},
			mexFields: solr.MexFieldBackingInfoMap{
				"description": {
					MexType: "text",
					BackingFields: []solr.MexBackingFieldWiringInfo{
						{
							Name:     "description",
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
			searchFocusName: focusName,
			includePrefix:   true,
			isPhraseOnly:    false,
			expectedNames:   []string{searchFocusFieldNameGeneric, searchFocusFieldNameDe, searchFocusFieldNameEn, searchFocusFieldNamePrefix, searchFocusFieldNameUnanalyzed},
		},
		{
			name: "Query IS NOT phrase-only, prefix search IS NOT requests: Names all fields are returned (with boosts) EXCEPT prefix",
			searchFocusElem: &sharedSearchConfig.SearchConfigObject{
				Type:   solr.MexSearchFocusType,
				Name:   focusName,
				Fields: []string{"description"},
			},
			mexFields: solr.MexFieldBackingInfoMap{
				"description": {
					MexType: "text",
					BackingFields: []solr.MexBackingFieldWiringInfo{
						{
							Name:     "description",
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
			searchFocusName: focusName,
			includePrefix:   false,
			isPhraseOnly:    false,
			expectedNames:   []string{searchFocusFieldNameGeneric, searchFocusFieldNameDe, searchFocusFieldNameEn, searchFocusFieldNameUnanalyzed},
		},
		{
			name: "Queries IS phrase-only, prefix search IS requested: Only the unanalyzed field is returned (with boost)",
			searchFocusElem: &sharedSearchConfig.SearchConfigObject{
				Type:   solr.MexSearchFocusType,
				Name:   focusName,
				Fields: []string{"description"},
			},
			mexFields: solr.MexFieldBackingInfoMap{
				"description": {
					MexType: "text",
					BackingFields: []solr.MexBackingFieldWiringInfo{
						{
							Name:     "description",
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
			searchFocusName: focusName,
			includePrefix:   true,
			isPhraseOnly:    true,
			expectedNames:   []string{searchFocusFieldNameUnanalyzed},
		},
		{
			name: "Queries IS phrase-only, prefix search IS NOT requested: Only the unanalyzed field is returned (with boosts)",
			searchFocusElem: &sharedSearchConfig.SearchConfigObject{
				Type:   solr.MexSearchFocusType,
				Name:   focusName,
				Fields: []string{"description"},
			},
			mexFields: solr.MexFieldBackingInfoMap{
				"description": {
					MexType: "text",
					BackingFields: []solr.MexBackingFieldWiringInfo{
						{
							Name:     "description",
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
			searchFocusName: focusName,
			includePrefix:   false,
			isPhraseOnly:    true,
			expectedNames:   []string{searchFocusFieldNameUnanalyzed},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scType := &SearchFocusType{}
			gotNames, gotErr := scType.GetSolrSearchFieldNames(tt.searchFocusName, tt.includePrefix, tt.isPhraseOnly)
			if tt.wantErr != (gotErr != nil) {
				t.Errorf("Wanted error %v, but the returned error was %v", tt.wantErr, gotErr)
				return
			}
			if !tt.wantErr {
				sort.Slice(tt.expectedNames, func(i, j int) bool {
					return tt.expectedNames[i] < tt.expectedNames[j]
				})
				sort.Slice(gotNames, func(i, j int) bool {
					return gotNames[i] < gotNames[j]
				})
				if !reflect.DeepEqual(gotNames, tt.expectedNames) {
					t.Errorf("Returned names do not match generated names: wanted %v, got %v", tt.expectedNames, gotNames)
				}
			}
		})
	}
}

func TestSearchFocusType_GetSolrBackingFieldDefs(t *testing.T) {
	focusFieldName := solr.GetSearchFocusFieldName("testFocus")
	testFocusNameGeneric, _ := solr.GetLangSpecificFieldName(focusFieldName, solr.GenericLangAbbrev)
	testFocusNameDe, _ := solr.GetLangSpecificFieldName(focusFieldName, solr.GermanLangAbbrev)
	testFocusNameEn, _ := solr.GetLangSpecificFieldName(focusFieldName, solr.EnglishLangAbbrev)
	testFocusNamePrefix := solr.GetPrefixBackingFieldName(focusFieldName)
	testFocusNameRaw := solr.GetRawBackingFieldName(focusFieldName)
	descriptionNameGeneric, _ := solr.GetLangSpecificFieldName("description", solr.GenericLangAbbrev)
	descriptionNameDe, _ := solr.GetLangSpecificFieldName("description", solr.GermanLangAbbrev)
	descriptionNameEn, _ := solr.GetLangSpecificFieldName("description", solr.EnglishLangAbbrev)
	descriptionNameRaw := solr.GetRawBackingFieldName("description")
	descriptionNamePrefix := solr.GetPrefixBackingFieldName("description")
	descriptionNameNormalized := solr.GetNormalizedBackingFieldName("description")
	deHullLabelName, _ := solr.GetLangSpecificFieldName(solr.GetTransitiveHullDisplayFieldName("unitCode"), solr.GermanLangAbbrev)
	enHullLabelName, _ := solr.GetLangSpecificFieldName(solr.GetTransitiveHullDisplayFieldName("unitCode"), solr.EnglishLangAbbrev)

	tests := []struct {
		name            string
		searchFocusElem *sharedSearchConfig.SearchConfigObject
		mexFields       solr.MexFieldBackingInfoMap
		wantFields      []solr.FieldDef
		wantCopyFields  []solr.CopyFieldDef
		wantErr         bool
	}{
		{
			name: "Backing text fields are generated for all known languages, prefixes, and unanalyzed text",
			searchFocusElem: &sharedSearchConfig.SearchConfigObject{
				Type:   solr.MexSearchFocusType,
				Name:   "testFocus",
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
			wantFields:     getExpectedTestFieldsForFocus("testFocus"),
			wantCopyFields: nil,
		},
		{
			name: "For a string field in a search focus, the standard field value is copied into the generic, prefix, and raw axis backing fields",
			searchFocusElem: &sharedSearchConfig.SearchConfigObject{
				Type:   solr.MexSearchFocusType,
				Name:   "testFocus",
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
			wantCopyFields: []solr.CopyFieldDef{
				{
					Source:      "category",
					Destination: []string{testFocusNameGeneric},
				},
				{
					Source:      "category",
					Destination: []string{testFocusNamePrefix},
				},
				{
					Source:      "category",
					Destination: []string{testFocusNameRaw},
				},
			},
		},
		{
			name: "For a timestamp field in a search focus, the raw field value is copied into the generic language and the unanalyzed backing fields",
			searchFocusElem: &sharedSearchConfig.SearchConfigObject{
				Type:   solr.MexSearchFocusType,
				Name:   "testFocus",
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
			wantCopyFields: []solr.CopyFieldDef{
				{
					Source:      "created_raw_value",
					Destination: []string{testFocusNameGeneric},
				},
				{
					Source:      "created_raw_value",
					Destination: []string{testFocusNameRaw},
				},
			},
		},
		{
			name: "For a hierarchy field in a search focus, the DE & EN labels are copied to into the matching language-specific search fields, as well as into the prefix and the unanalyzed backing fields",
			searchFocusElem: &sharedSearchConfig.SearchConfigObject{
				Type:   solr.MexSearchFocusType,
				Name:   "testFocus",
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
			wantCopyFields: []solr.CopyFieldDef{
				{
					Source:      "unitCode_trhull_display___de",
					Destination: []string{testFocusNameDe},
				},
				{
					Source:      "unitCode_trhull_display___en",
					Destination: []string{testFocusNameEn},
				},
				{
					Source:      "unitCode_trhull_display___de",
					Destination: []string{testFocusNamePrefix},
				},
				{
					Source:      "unitCode_trhull_display___en",
					Destination: []string{testFocusNamePrefix},
				},
				{
					Source:      "unitCode_trhull_display___de",
					Destination: []string{testFocusNameRaw},
				},
				{
					Source:      "unitCode_trhull_display___en",
					Destination: []string{testFocusNameRaw},
				},
			},
		},
		{
			name: "For text field belonging to a search focus, the values for each backing field is copied into the corresponding focus backing field",
			searchFocusElem: &sharedSearchConfig.SearchConfigObject{
				Type:   solr.MexSearchFocusType,
				Name:   "testFocus",
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
							Name:     descriptionNamePrefix,
							Category: solr.PrefixContentBaseFieldCategory,
						},
						{
							Name:     descriptionNameRaw,
							Category: solr.RawContentBaseFieldCategory,
						},
						{
							Name:     descriptionNameNormalized,
							Category: solr.NormalizedBaseFieldCategory,
						},
					},
				},
			},
			wantCopyFields: []solr.CopyFieldDef{
				{
					Source:      descriptionNameGeneric,
					Destination: []string{testFocusNameGeneric},
				},
				{
					Source:      descriptionNameDe,
					Destination: []string{testFocusNameDe},
				},
				{
					Source:      descriptionNameEn,
					Destination: []string{testFocusNameEn},
				},
				{
					Source:      descriptionNamePrefix,
					Destination: []string{testFocusNamePrefix},
				},
				{
					Source:      descriptionNameRaw,
					Destination: []string{testFocusNameRaw},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scType := &SearchFocusType{}
			gotFields, gotCopyFields, err := scType.GetSolrBackingFields(tt.searchFocusElem, tt.mexFields)
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

func TestSearchFocusType_GetMatchingOpsConfig(t *testing.T) {

	focusName := "testFocus"
	searchFocusFieldNameBase := solr.GetSearchFocusFieldName(focusName)
	searchFocusFieldNameGeneric, _ := solr.GetLangSpecificFieldName(searchFocusFieldNameBase, solr.GenericLangAbbrev)
	searchFocusFieldNameDe, _ := solr.GetLangSpecificFieldName(searchFocusFieldNameBase, solr.GermanLangAbbrev)
	searchFocusFieldNameEn, _ := solr.GetLangSpecificFieldName(searchFocusFieldNameBase, solr.EnglishLangAbbrev)
	searchFocusFieldNamePrefix := solr.GetPrefixBackingFieldName(searchFocusFieldNameBase)
	searchFocusFieldNameUnanalyzed := solr.GetRawBackingFieldName(searchFocusFieldNameBase)

	tests := []struct {
		name            string
		searchFocusElem *sharedSearchConfig.SearchConfigObject
		includePrefix   bool
		maxEditDistance uint32
		wantErr         bool
		check           testutils.MatchingOpCheck
	}{
		{
			name: "The correct fields are queried when prefix search is requested",
			searchFocusElem: &sharedSearchConfig.SearchConfigObject{
				Type:   solr.MexSearchFocusType,
				Name:   focusName,
				Fields: []string{"description"},
			},
			includePrefix:   true,
			maxEditDistance: 0,
			check: testutils.CheckSearchedFields(
				[]string{
					searchFocusFieldNameUnanalyzed,
					searchFocusFieldNameGeneric,
					searchFocusFieldNameDe,
					searchFocusFieldNameEn,
					searchFocusFieldNamePrefix,
				},
				[]string{searchFocusFieldNameUnanalyzed},
			),
		},
		{
			name: "The correct fields are queries when prefix search is NOT requested",
			searchFocusElem: &sharedSearchConfig.SearchConfigObject{
				Type:   solr.MexSearchFocusType,
				Name:   focusName,
				Fields: []string{"description"},
			},
			includePrefix:   false,
			maxEditDistance: 0,
			check: testutils.CheckSearchedFields(
				[]string{
					searchFocusFieldNameUnanalyzed,
					searchFocusFieldNameGeneric,
					searchFocusFieldNameDe,
					searchFocusFieldNameEn,
				},
				[]string{searchFocusFieldNameUnanalyzed},
			),
		},
		{
			name: "The correct boost factors are ser when prefix search is requested",
			searchFocusElem: &sharedSearchConfig.SearchConfigObject{
				Type:   solr.MexSearchFocusType,
				Name:   focusName,
				Fields: []string{"description"},
			},
			includePrefix:   true,
			maxEditDistance: 0,
			check: testutils.CheckBoostFactor(
				map[string]string{
					searchFocusFieldNameUnanalyzed: solr.UnanalyzedBoostFactor,
					searchFocusFieldNameGeneric:    "",
					searchFocusFieldNameDe:         "",
					searchFocusFieldNameEn:         "",
					searchFocusFieldNamePrefix:     solr.PrefixBoostFactor,
				},
				map[string]string{
					searchFocusFieldNameUnanalyzed: solr.UnanalyzedBoostFactor,
				},
			),
		},
		{
			name: "The correct boost factors are ser when prefix search is NOT requested",
			searchFocusElem: &sharedSearchConfig.SearchConfigObject{
				Type:   solr.MexSearchFocusType,
				Name:   focusName,
				Fields: []string{"description"},
			},
			includePrefix:   false,
			maxEditDistance: 0,
			check: testutils.CheckBoostFactor(
				map[string]string{
					searchFocusFieldNameUnanalyzed: solr.UnanalyzedBoostFactor,
					searchFocusFieldNameGeneric:    "",
					searchFocusFieldNameDe:         "",
					searchFocusFieldNameEn:         "",
				},
				map[string]string{
					searchFocusFieldNameUnanalyzed: solr.UnanalyzedBoostFactor,
				},
			),
		},
		{
			name: "The correct edit distances are ser when prefix search is requested",
			searchFocusElem: &sharedSearchConfig.SearchConfigObject{
				Type:   solr.MexSearchFocusType,
				Name:   focusName,
				Fields: []string{"description"},
			},
			includePrefix:   true,
			maxEditDistance: 1,
			check: testutils.CheckEditDistance(
				map[string]uint32{
					searchFocusFieldNameUnanalyzed: 0,
					searchFocusFieldNameGeneric:    1,
					searchFocusFieldNameDe:         1,
					searchFocusFieldNameEn:         1,
					searchFocusFieldNamePrefix:     0,
				},
				map[string]uint32{
					searchFocusFieldNameUnanalyzed: 0,
				},
			),
		},
		{
			name: "The correct edit distances are ser when prefix search is NOT requested",
			searchFocusElem: &sharedSearchConfig.SearchConfigObject{
				Type:   solr.MexSearchFocusType,
				Name:   focusName,
				Fields: []string{"description"},
			},
			includePrefix:   false,
			maxEditDistance: 1,
			check: testutils.CheckEditDistance(
				map[string]uint32{
					searchFocusFieldNameUnanalyzed: 0,
					searchFocusFieldNameGeneric:    1,
					searchFocusFieldNameDe:         1,
					searchFocusFieldNameEn:         1,
				},
				map[string]uint32{
					searchFocusFieldNameUnanalyzed: 0,
				},
			),
		},
		{
			name: "An error is returned if the edit distance is larger than the allowed max value",
			searchFocusElem: &sharedSearchConfig.SearchConfigObject{
				Type:   solr.MexSearchFocusType,
				Name:   focusName,
				Fields: []string{"description"},
			},
			includePrefix:   true,
			maxEditDistance: solr.MaxEditDistance + 1,
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scType := &SearchFocusType{}

			gotConfig, gotErr := scType.GetMatchingOpsConfig("testFocus", tt.maxEditDistance, tt.includePrefix)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("Want error: %v - but got error %v", tt.wantErr, gotErr)
			}
			if gotErr == nil {
				tt.check(gotConfig, t)
			}
		})
	}
}
