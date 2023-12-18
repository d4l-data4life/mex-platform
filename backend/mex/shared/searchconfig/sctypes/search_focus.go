package sctypes

/*
This file contains the logic associated with a search focus.
*/

import (
	"fmt"

	"github.com/d4l-data4life/mex/mex/services/query/parser"
	"github.com/d4l-data4life/mex/mex/shared/index"
	sharedSearchConfig "github.com/d4l-data4life/mex/mex/shared/searchconfig"
	"github.com/d4l-data4life/mex/mex/shared/solr"
)

type SearchFocusType struct{}

// GetSolrSearchFieldNames returns the names of the Solr field backing the search element with a given name
func (scType *SearchFocusType) GetSolrSearchFieldNames(searchConfigName string, usePrefixField bool, isPhraseOnlyQuery bool) ([]string, error) {
	if searchConfigName == "" {
		return nil, fmt.Errorf("name of search config cannot be empty")
	}
	focusFieldName := solr.GetSearchFocusFieldName(searchConfigName)
	// Always do a (boosted) search of the unanalyzed field
	unanalyzedField := fmt.Sprintf("%s^%s", solr.GetRawBackingFieldName(focusFieldName), solr.UnanalyzedBoostFactor)
	backingFieldNames := []string{unanalyzedField}
	// For pure phrase queries, do not search any other fields to avoid stemming matches
	if isPhraseOnlyQuery {
		return backingFieldNames, nil
	}
	// Add language-specific backing fields
	for lc := range solr.KnownLanguagesFieldTypeMap {
		fn, _ := solr.GetLangSpecificFieldName(focusFieldName, lc)
		backingFieldNames = append(backingFieldNames, fn)
	}
	// Add prefix search backing field if requested
	if usePrefixField {
		prefixField := fmt.Sprintf("%s^%s", solr.GetPrefixBackingFieldName(focusFieldName), solr.PrefixBoostFactor)
		backingFieldNames = append(backingFieldNames, prefixField)
	}
	return backingFieldNames, nil
}

// GetMatchingOpsConfig returns the configuration for the various field that should be searched for this search focus
func (scType *SearchFocusType) GetMatchingOpsConfig(searchConfigName string, maxEditDistance uint32, usePrefixField bool) (parser.MatchingOpsConfig, error) {
	if maxEditDistance > solr.MaxEditDistance {
		return parser.MatchingOpsConfig{}, fmt.Errorf("term matching operator: max edit distance cannot be set to more than %d", solr.MaxEditDistance)
	}
	focusFieldName := solr.GetSearchFocusFieldName(searchConfigName)
	searchFocusBackingFields := solr.GetTextCoreBackingFieldInfo(focusFieldName, true)
	matchingOpsConfig := parser.MatchingOpsConfig{
		Term:   []parser.MatchingFieldConfig{},
		Phrase: []parser.MatchingFieldConfig{},
	}

	for category, fieldInfo := range searchFocusBackingFields {
		switch category {
		case solr.RawSearchFunctionCategory:
			// The unanalyzed content is searched for both terms and phrases
			unanalyzedFieldOp := parser.MatchingFieldConfig{
				FieldName:       fieldInfo.SolrName,
				BoostFactor:     solr.UnanalyzedBoostFactor,
				MaxEditDistance: 0,
			}
			matchingOpsConfig.Term = append(matchingOpsConfig.Term, unanalyzedFieldOp)
			matchingOpsConfig.Phrase = append(matchingOpsConfig.Phrase, unanalyzedFieldOp)
		case solr.PrefixSearchFunctionCategory:
			// Add term search on prefix backing field if requested (not fuzzy, boosted)
			if usePrefixField {
				matchingOpsConfig.Term = append(matchingOpsConfig.Term, parser.MatchingFieldConfig{
					FieldName:       fieldInfo.SolrName,
					BoostFactor:     solr.PrefixBoostFactor,
					MaxEditDistance: 0,
				})
			}
		default:
			// Other backing fields are used for term search (fuzzy, no boost)
			matchingOpsConfig.Term = append(matchingOpsConfig.Term, parser.MatchingFieldConfig{
				FieldName:       fieldInfo.SolrName,
				BoostFactor:     "",
				MaxEditDistance: maxEditDistance,
			})
		}
	}

	return matchingOpsConfig, nil
}

// GetSolrBackingFields returns the Solr fields and copy fields needed for a given search focus
func (scType *SearchFocusType) GetSolrBackingFields(searchFocusElem *sharedSearchConfig.SearchConfigObject, mexFieldMap solr.MexFieldBackingInfoMap,
) ([]solr.FieldDef, []solr.CopyFieldDef, error) {
	// Check that this is actually a search focus
	if searchFocusElem.Type != solr.MexSearchFocusType {
		return nil, nil, fmt.Errorf("search focus hook was passed a search config element with unknown type")
	}

	// Create backing fields
	focusFieldName := solr.GetSearchFocusFieldName(searchFocusElem.Name)
	searchFocusBackingFieldInfo := make(solr.FieldCategoryToSolrFieldDefsMap)
	mexToSolrMap := solr.GetTextCoreBackingFieldInfo(focusFieldName, true)
	for category, coreFieldInfo := range mexToSolrMap {
		searchFocusBackingFieldInfo[category] = solr.GetStandardSecondaryBackingField(coreFieldInfo.SolrName, coreFieldInfo.SolrType, false)
	}

	var focusBackingFields []solr.FieldDef
	targetFieldByFunctionCategory := make(map[string]string)
	for category, fDef := range searchFocusBackingFieldInfo {
		targetFieldByFunctionCategory[category] = fDef.Name
		focusBackingFields = append(focusBackingFields, fDef)
	}

	// Add copy fields connecting backing fields for the matching core fields with the search foci backing fields
	focusCopyFields, err := index.GetSearchFocusSolrCopyFields(mexFieldMap, targetFieldByFunctionCategory, searchFocusElem)
	if err != nil {
		return nil, nil, fmt.Errorf(fmt.Sprintf("solr copy field generation for search focus failed: %s", err.Error()))
	}

	return focusBackingFields, focusCopyFields, nil
}
