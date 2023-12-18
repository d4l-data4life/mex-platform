package sctypes

/*
This file contains the logic associated with an ordinal axis.
*/

import (
	"fmt"

	"github.com/d4l-data4life/mex/mex/services/query/parser"
	"github.com/d4l-data4life/mex/mex/shared/index"
	sharedSearchConfig "github.com/d4l-data4life/mex/mex/shared/searchconfig"
	"github.com/d4l-data4life/mex/mex/shared/solr"
)

type OrdinalAxisType struct{}

// GetSolrSearchFieldNames returns the names of the Solr field backing the search element with a given name
func (scType *OrdinalAxisType) GetSolrSearchFieldNames(string, bool, bool) ([]string, error) {
	return nil, fmt.Errorf("should not ask for fields to search for an ordinal axis")
}

func (scType *OrdinalAxisType) GetMatchingOpsConfig(string, uint32, bool) (parser.MatchingOpsConfig, error) {
	panic("should not ask for matching configuration for an ordinal axis")
}

// GetSolrBackingFields returns the Solr fields and copy fields needed for a given ordinal axis
func (scType *OrdinalAxisType) GetSolrBackingFields(ordinalAxisElem *sharedSearchConfig.SearchConfigObject, mexFieldMap solr.MexFieldBackingInfoMap,
) ([]solr.FieldDef, []solr.CopyFieldDef, error) {
	// Check that this is actually an ordinal axis
	if ordinalAxisElem.Type != solr.MexOrdinalAxisType {
		return nil, nil, fmt.Errorf("ordinal axis hook was passed a search config element with unknown type")
	}

	fieldToTypeMap := make(map[string]string)
	for fieldName, backingInfo := range mexFieldMap {
		fieldToTypeMap[fieldName] = backingInfo.MexType
	}
	ordinalAxisFieldType, err := GetOrdinalAxisFieldType(ordinalAxisElem.Fields, fieldToTypeMap)
	if err != nil {
		return nil, nil, err
	}

	// Create a backing field for faceting
	ordinalFacetAxisName := solr.GetOrdinalAxisFacetAndFilterFieldName(ordinalAxisElem.Name)
	ordinalAxisFacetBackingField := sharedSearchConfig.FunctionBackingFieldInfo{
		Def:                solr.GetStandardSecondaryBackingField(ordinalFacetAxisName, ordinalAxisFieldType, true),
		FunctionCategoryID: solr.FacetAndFilterFunctionCategory,
	}

	// Create a backing field for sorting
	ordinalSortAxisName := solr.GetOrdinalAxisSortFieldName(ordinalAxisElem.Name)
	ordinalAxisSortBackingField := sharedSearchConfig.FunctionBackingFieldInfo{
		Def:                solr.GetStandardSecondaryBackingField(ordinalSortAxisName, ordinalAxisFieldType, true),
		FunctionCategoryID: solr.SortFunctionCategory,
	}

	// Add copy fields connecting backing fields with the matching core fields with the search foci backing fields
	targetFieldByFunctionCategory := map[string]string{
		solr.FacetAndFilterFunctionCategory: ordinalAxisFacetBackingField.Def.Name,
		solr.SortFunctionCategory:           ordinalAxisSortBackingField.Def.Name,
	}
	focusCopyFields, err := index.GetOrdinalAxisSolrCopyFields(mexFieldMap, targetFieldByFunctionCategory, ordinalAxisElem)
	if err != nil {
		return nil, nil, fmt.Errorf(fmt.Sprintf("solr copy field generation for ordinal axis failed: %s",
			err.Error()))
	}

	return []solr.FieldDef{ordinalAxisFacetBackingField.Def, ordinalAxisSortBackingField.Def}, focusCopyFields, nil
}
