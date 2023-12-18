package sctypes

/*
This file contains the logic associated with a hierarchical axis.
*/

import (
	"fmt"

	"github.com/d4l-data4life/mex/mex/services/query/parser"
	"github.com/d4l-data4life/mex/mex/shared/index"
	sharedSearchConfig "github.com/d4l-data4life/mex/mex/shared/searchconfig"
	"github.com/d4l-data4life/mex/mex/shared/solr"

	kindHierarchy "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/hierarchy"
)

type HierarchyAxisType struct{}

// GetSolrSearchFieldNames returns the names of the Solr fields backing the search element with a given name
func (haType *HierarchyAxisType) GetSolrSearchFieldNames(string, bool, bool) ([]string, error) {
	return nil, fmt.Errorf("should not ask for fields to search for a hierarchy axis")
}

func (haType *HierarchyAxisType) GetMatchingOpsConfig(string, uint32, bool) (parser.MatchingOpsConfig, error) {
	panic("should not ask for matching configuration for a hierarchy axis")
}

// GetSolrBackingFields returns the Solr fields and copy fields needed for a given hierarchical axis
func (haType *HierarchyAxisType) GetSolrBackingFields(hierarchyAxisElem *sharedSearchConfig.SearchConfigObject, mexFieldMap solr.MexFieldBackingInfoMap,
) ([]solr.FieldDef, []solr.CopyFieldDef, error) {
	// Check that this is actually a hierarchy axis
	if hierarchyAxisElem.Type != solr.MexHierarchyAxisType {
		return nil, nil, fmt.Errorf("hierarchy axis hook was passed a search config element with unknown type")
	}

	for _, fn := range hierarchyAxisElem.Fields {
		bInfo, ok := mexFieldMap[fn]
		if !ok {
			return nil, nil, fmt.Errorf("the field '%s' used in a hierarchy axis could not be found", fn)
		}
		if bInfo.MexType != kindHierarchy.KindName {
			return nil, nil, fmt.Errorf("hierarchy axis contains a non-hierarchy field (field '%s' of type %s)", fn, bInfo.MexType)
		}
	}

	// Create a backing string field for faceting based on code + parent codes
	hierarchyFacetAxisName := solr.GetOrdinalAxisFacetAndFilterFieldName(hierarchyAxisElem.Name)
	// Create a backing string field for constraints based on just the actually assigned code
	singleNodeHierarchyAxisName := solr.GetSingleNodeAxisFieldName(hierarchyFacetAxisName)
	// Create a backing string field for sorting
	sortHierarchyAxisName := solr.GetOrdinalAxisSortFieldName(hierarchyAxisElem.Name)
	axisBackingFields := []solr.FieldDef{
		solr.GetStandardSecondaryBackingField(hierarchyFacetAxisName, solr.DefaultSolrStringFieldType, true),
		solr.GetStandardSecondaryBackingField(singleNodeHierarchyAxisName, solr.DefaultSolrStringFieldType, false),
		solr.GetStandardSecondaryBackingField(sortHierarchyAxisName, solr.DefaultSolrSortableTextFieldType, true),
	}

	// Add copy fields connecting backing fields with the matching hierarchical axis backing fields
	targetFieldByFunctionCategory := map[string]string{
		solr.FacetAndFilterFunctionCategory:   hierarchyFacetAxisName,
		solr.SortFunctionCategory:             sortHierarchyAxisName,
		solr.SingleValueFacetFunctionCategory: singleNodeHierarchyAxisName,
	}
	var haCopyFields []solr.CopyFieldDef
	haCopyFields, err := index.GetHierarchyAxisSolrCopyFields(mexFieldMap, targetFieldByFunctionCategory, hierarchyAxisElem)
	if err != nil {
		return nil, nil, fmt.Errorf(fmt.Sprintf("solr copy field generation for hierarchy axis failed: %s",
			err.Error()))
	}

	return axisBackingFields, haCopyFields, nil
}
