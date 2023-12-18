package index

import (
	"fmt"

	"github.com/d4l-data4life/mex/mex/shared/searchconfig"
	"github.com/d4l-data4life/mex/mex/shared/solr"

	kindCoding "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/coding"
	kindHierarchy "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/hierarchy"
	kindLink "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/link"
	kindNumber "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/number"
	kindString "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/string"
	kindText "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/text"
	kindTimestamp "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/timestamp"
)

// KindWiringMap contains the backing field wiring map for a single MEx kind,
// organized by the category of the underlying backing field
type KindWiringMap map[string][]string

// WiringMap contains the backing field wiring maps organized by MEx field kind (keys)
type WiringMap map[string]KindWiringMap

// ordinalAxisWiringMap contains the wiring logic for ordinal axes
var ordinalAxisWiringMap = WiringMap{
	kindCoding.KindName: {
		// For MEx coding fields, the generic field is used for sorting (faceting/filtering is not supported)
		solr.GenericLangBaseFieldCategory: []string{solr.SortFunctionCategory},
	},
	kindLink.KindName: {
		// For MEx link fields, the generic field is used for faceting/filtering (sorting is not supported)
		solr.GenericLangBaseFieldCategory: []string{solr.FacetAndFilterFunctionCategory},
	},
	kindNumber.KindName: {
		// For MEx number fields, the generic field is used for faceting and sorting
		solr.GenericLangBaseFieldCategory: []string{solr.FacetAndFilterFunctionCategory, solr.SortFunctionCategory},
	},
	kindString.KindName: {
		// For MEx string fields, the generic field is used for faceting and the normalized one for sorting
		solr.GenericLangBaseFieldCategory: []string{solr.FacetAndFilterFunctionCategory},
		solr.NormalizedBaseFieldCategory:  []string{solr.SortFunctionCategory},
	},
	kindText.KindName: {
		// For MEx text fields, all language-specific field are used for faceting and the normalized field is used for sorting
		solr.GenericLangBaseFieldCategory: []string{solr.FacetAndFilterFunctionCategory},
		solr.GermanLangBaseFieldCategory:  []string{solr.FacetAndFilterFunctionCategory},
		solr.EnglishLangBaseFieldCategory: []string{solr.FacetAndFilterFunctionCategory},
		solr.NormalizedBaseFieldCategory:  []string{solr.SortFunctionCategory},
	},
	kindTimestamp.KindName: {
		// For MEx timestamp fields, the generic field is used for faceting and sorting
		solr.GenericLangBaseFieldCategory: []string{solr.FacetAndFilterFunctionCategory, solr.SortFunctionCategory},
	},
}

// hierarchyAxisWiringMap contains the wiring logic for ordinal axes
var hierarchyAxisWiringMap = WiringMap{
	kindHierarchy.KindName: {
		// For MEx hierarchy fields, the single code is added to the single value & the sort axis, and the child codes to the normal facet field
		solr.GenericLangBaseFieldCategory: []string{solr.SingleValueFacetFunctionCategory, solr.SortFunctionCategory},
		solr.ParentCodesBaseFieldCategory: []string{solr.FacetAndFilterFunctionCategory},
	},
}

// searchFocusWiringMap contains the wiring logic for search foci
// Note that the solr.RawSearchFunctionCategory target category must be present to allow pure phrase searches.
var searchFocusWiringMap = WiringMap{
	kindCoding.KindName: {
		// For MEx coding fields, DE and EN labels are used for search (language-specific + prefix & raw)
		solr.GermanLangBaseFieldCategory:  []string{solr.GermanLangSearchFunctionCategory, solr.PrefixSearchFunctionCategory, solr.RawSearchFunctionCategory},
		solr.EnglishLangBaseFieldCategory: []string{solr.EnglishLangSearchFunctionCategory, solr.PrefixSearchFunctionCategory, solr.RawSearchFunctionCategory},
	},
	kindHierarchy.KindName: {
		// For MEx hierarchy fields, DE and EN labels are used for search (language-specific + prefix & raw)
		solr.GermanLangBaseFieldCategory:  []string{solr.GermanLangSearchFunctionCategory, solr.PrefixSearchFunctionCategory, solr.RawSearchFunctionCategory},
		solr.EnglishLangBaseFieldCategory: []string{solr.EnglishLangSearchFunctionCategory, solr.PrefixSearchFunctionCategory, solr.RawSearchFunctionCategory},
	},
	kindNumber.KindName: {
		// For MEx number fields, only the generic field is used for search (base only)
		solr.GenericLangBaseFieldCategory: []string{solr.GenericLangSearchFunctionCategory},
	},
	kindString.KindName: {
		// For MEx string fields, only the generic field is used for search
		solr.GenericLangBaseFieldCategory: []string{solr.GenericLangSearchFunctionCategory, solr.PrefixSearchFunctionCategory, solr.RawSearchFunctionCategory},
	},
	kindText.KindName: {
		// For MEx text fields, the backing field match those of the search focus to ensure functioning highlighting.
		// Accordingly, content is simply copied into the matching field.
		solr.GenericLangBaseFieldCategory:   []string{solr.GenericLangSearchFunctionCategory},
		solr.GermanLangBaseFieldCategory:    []string{solr.GermanLangSearchFunctionCategory},
		solr.EnglishLangBaseFieldCategory:   []string{solr.EnglishLangSearchFunctionCategory},
		solr.PrefixContentBaseFieldCategory: []string{solr.PrefixSearchFunctionCategory},
		solr.RawContentBaseFieldCategory:    []string{solr.RawSearchFunctionCategory},
	},
	kindTimestamp.KindName: {
		// For MEx timestamp fields, only the *raw* (string) timestamp field is used for search (base and unanalyzed only)
		solr.RawContentBaseFieldCategory: []string{solr.GenericLangSearchFunctionCategory, solr.RawSearchFunctionCategory},
	},
}

// GetOrdinalAxisSolrCopyFields returns the copy fields needed to fill the backing fields for the ordinal axis
func GetOrdinalAxisSolrCopyFields(mexMap solr.MexFieldBackingInfoMap, targetFieldMap map[string]string, scElem *searchconfig.SearchConfigObject,
) ([]solr.CopyFieldDef, error) {
	return getSearchConfigSolrCopyFields(mexMap, targetFieldMap, scElem, ordinalAxisWiringMap)
}

// GetHierarchyAxisSolrCopyFields returns the copy fields needed to fill the backing fields for the ordinal axis
func GetHierarchyAxisSolrCopyFields(mexMap solr.MexFieldBackingInfoMap, targetFieldMap map[string]string, scElem *searchconfig.SearchConfigObject,
) ([]solr.CopyFieldDef, error) {
	return getSearchConfigSolrCopyFields(mexMap, targetFieldMap, scElem, hierarchyAxisWiringMap)
}

// GetSearchFocusSolrCopyFields returns the copy fields needed to fill the backing fields for the search focus
func GetSearchFocusSolrCopyFields(mexMap solr.MexFieldBackingInfoMap, targetFieldMap map[string]string, scElem *searchconfig.SearchConfigObject,
) ([]solr.CopyFieldDef, error) {
	return getSearchConfigSolrCopyFields(mexMap, targetFieldMap, scElem, searchFocusWiringMap)
}

// getSearchConfigSolrCopyFields returns the copy fields needed to fill the backing fields for ta given search configs
func getSearchConfigSolrCopyFields(mexMap solr.MexFieldBackingInfoMap, targetFieldMap map[string]string, scElem *searchconfig.SearchConfigObject, wiringMap WiringMap,
) ([]solr.CopyFieldDef, error) {
	var copyFields []solr.CopyFieldDef
	// For each MEx field in search config...
	for _, baseFieldName := range scElem.Fields {
		mexFieldBackingInfo, ok := mexMap[baseFieldName]
		if !ok {
			return nil, fmt.Errorf("could not find Solr backing fields for MEx field '%s' listed in search config field '%s'", baseFieldName, scElem.Name)
		}
		// ... and for each backing field...
		for _, backingFieldInfo := range mexFieldBackingInfo.BackingFields {
			// ... create copy fields from the MEx backing field into the relevant backing fields
			newCopyFields, err := getCopyFields(mexFieldBackingInfo.MexType, backingFieldInfo.Name, backingFieldInfo.Category, targetFieldMap, wiringMap)
			copyFields = append(copyFields, newCopyFields...)
			if err != nil {
				return nil, err
			}
		}
	}

	return copyFields, nil
}

// getCopyFields contains the logic for wiring up a given source field with all relevant functional backing fields (using Solr copy directives)
func getCopyFields(mexType string, source string, sourceCategory string, targetFieldByCategory map[string]string, wiringMap WiringMap) ([]solr.CopyFieldDef, error) {
	var copyFields []solr.CopyFieldDef
	if kindMap, typeOk := wiringMap[mexType]; typeOk {
		if functionTypeKeys, categoryOk := kindMap[sourceCategory]; categoryOk {
			for _, functionType := range functionTypeKeys {
				target, funcTypeOk := targetFieldByCategory[functionType]
				if !funcTypeOk {
					return nil, fmt.Errorf("axis function type '%s' is not valid", functionType)
				}
				copyFields = append(copyFields, getCopyField(source, []string{target}))
			}
		}
	} else {
		return nil, fmt.Errorf("MEx field type '%s' not allowed in search config element", mexType)
	}

	return copyFields, nil
}

func getCopyField(source string, targets []string) solr.CopyFieldDef {
	return solr.CopyFieldDef{
		Source:      source,
		Destination: targets,
	}
}
