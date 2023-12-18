package sctypes

import (
	"github.com/d4l-data4life/mex/mex/services/query/parser"
	"github.com/d4l-data4life/mex/mex/shared/searchconfig"
	"github.com/d4l-data4life/mex/mex/shared/solr"
)

type ConfigHook interface {
	// GetSolrBackingFields returns the definitions of the Solr fields and copy fields needed to back a search config element.
	GetSolrBackingFields(configElem *searchconfig.SearchConfigObject, mexFieldMap solr.MexFieldBackingInfoMap) ([]solr.FieldDef, []solr.CopyFieldDef, error)
	// GetSolrSearchFieldNames returns the names of the Solr field to be searched for this element
	// (Boolean flag decides if the fuzzy search field should also be returned)
	GetSolrSearchFieldNames(searchConfigName string, useFuzzyField bool, isPhraseOnlyQuery bool) ([]string, error)
	// GetMatchingOpsConfig returns the configuration for the various field that should be searched
	GetMatchingOpsConfig(searchConfigName string, maxEditDistance uint32, usePrefixField bool) (parser.MatchingOpsConfig, error)
}

// SearchConfigHooks maps the config element type to the matching hook object
var SearchConfigHooks = map[string]ConfigHook{
	solr.MexOrdinalAxisType:   &OrdinalAxisType{},
	solr.MexSearchFocusType:   &SearchFocusType{},
	solr.MexHierarchyAxisType: &HierarchyAxisType{},
}
