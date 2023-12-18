package solr

import "time"

const (
	DefaultSolrConfigSet = "mex_rki"

	// Constants used when merging ingested items
	MergeIDPostfix   = "merged"
	MergeIDSeparator = "#"

	// Pre-configured fields in Solr (assumed to always be present)
	DefaultUniqueKey    = "id"
	ItemEntityNameField = "entityName"
	ItemCreatedAtField  = "createdAt"
	ItemBusinessIDField = "businessId"
	// This is defined by Solr itself, but is not used in MEx
	DefaultVersionKey = "_version_"

	// ErrMesssageKey is the property name for the message in the error object in a search response from Solr
	ErrMesssageKey = "msg"

	// Not exported - only used elsewhere in this file
	defaultSearchFieldName = "_text_"

	// Standard Solr field types
	DefaultSolrStringFieldType       = "string"
	DefaultSolrSortableTextFieldType = "string"
	DefaultSolrTimestampFieldType    = "pdate"
	DefaultSolrTextFieldType         = "text_mex_general"
	DefaultSolrNumberFieldType       = "pfloat"
	DefaultDeSolrTextFieldType       = "text_mex_de"
	DefaultEnSolrTextFieldType       = "text_mex_en"
	DefaultPrefixSolrTextFieldType   = "text_mex_prefix"
	DefaultRawSolrTextFieldType      = "text_mex_minimal" // This should NOT be "string" since that will prevent matching only part of the text

	// Solr post- and prefixes
	FocusPostfix                 = "search_focus"
	PrefixFocusPostfix           = "prefix"
	NormFocusPostfix             = "normalized"
	RawValuePostfix              = "unanalyzed"
	TransitiveHullPostfix        = "trhull"
	TransitiveHullDisplayPostfix = "trhull_display"
	AxisFacetPostfix             = "ordinal_facet_axis"
	AxisSortPostfix              = "ordinal_sort_axis"
	SingleValuePostfix           = "single_value"
	RawValTimestampPostfix       = "raw_value"
	LongSeparator                = "___"
	LinkedFieldSeparator         = "__"
	ExactPostfix                 = "exact"

	// Allowed MEx facet types - these are the types exposed to clients
	MexExactFacetType      = "exact"
	MexYearRangeFacetType  = "yearRange"
	MexStringStatFacetType = "stringStat"

	// Allowed MEx search configuration types
	MexSearchFocusType   = "searchFocus"
	MexOrdinalAxisType   = "ordinalAxis"
	MexHierarchyAxisType = "hierarchyAxis"

	MexDefaultSearchFocusName = "default"

	// MEx Boolean operators
	AndSeparator = " && "
	OrSeparator  = " || "

	// MEx stat facet functions
	MinOperator              = "min"
	MaxOperator              = "max"
	MexExactAxisConstraint   = "exact"
	MexStringRangeConstraint = "stringRange"

	// MEx query settings
	MaxEditDistance = 2
	// EditLowerCutoff and EditUpperCutoff control how the edit distance depends on the word length
	EditLowerCutoff       = 4
	EditUpperCutoff       = 10
	MaxDocLimit           = 1000
	MaxFacetLimit         = 1000
	FacetPrefix           = "facet"
	TagPostfix            = "tag"
	HighlightAlgorithm    = "unified"
	HighlightSnippetNo    = 10
	HighlightFragmentSize = 100
	HighlightStartChar    = "\ue000"
	HighlightStopChar     = "\ue001"
	// We give the boost factors as strings to avoid precision issues
	UnanalyzedBoostFactor = "5.0" // > 1 since exact matches should be rewarded over stemmed ones
	PrefixBoostFactor     = "0.5" // < 1 since fuzzzines will typically produce multiple matches here

	// Language abbreviations
	GenericLangAbbrev     = ""
	GenericLanguageSuffix = "generic"
	GermanLangAbbrev      = "de"
	EnglishLangAbbrev     = "en"

	SimpleAggregation          = "simple"
	SourcePartitionAggregation = "source_partition"

	KeepAllDuplicates        = "keepall"
	RemoveAllDuplicates      = "removeall"
	DefaultDuplicateStrategy = KeepAllDuplicates

	// Size of indexing batches
	IndexBatchSize        = 256
	DefaultSolrCommitTime = 1 * time.Second
	DefaultSolrBatchSize  = 25

	// Standard extension elements in field definition
	HierarchyExtID = "mex.v0.IndexDefExtHierarchy"
	LinkExtID      = "mex.v0.IndexDefExtLink"
)

var AllowedSortOrders = []string{"asc", "desc"}

// Set of fields that are always created in Solr.
var CoreMexFieldNames = []string{
	DefaultUniqueKey,
	ItemEntityNameField,
	ItemCreatedAtField,
	ItemBusinessIDField,
}

// These are Solr fields that will no be touched when cleaning the schema
// This is needed since they are mentioned in static Solr config files, so that
// removing them can lead to start-up issues.
var ProtectedSolrFields = []string{
	DefaultUniqueKey,
	DefaultVersionKey,
	defaultSearchFieldName,
}

// KnownLanguagesFieldTypeMap maps the known (non-generic) language codes to the corresponding Solr field types
var KnownLanguagesFieldTypeMap = map[string]string{
	GenericLangAbbrev: DefaultSolrTextFieldType,
	GermanLangAbbrev:  DefaultDeSolrTextFieldType,
	EnglishLangAbbrev: DefaultEnSolrTextFieldType,
}

var knownLanguagesFieldFunctionCategoryMap = map[string]string{
	GenericLangAbbrev: GenericLangSearchFunctionCategory, // Generic
	GermanLangAbbrev:  GermanLangSearchFunctionCategory,
	EnglishLangAbbrev: EnglishLangSearchFunctionCategory,
}

// KnownLanguagesCategoryMap maps the known language codes to the corresponding field categories
var KnownLanguagesCategoryMap = map[string]string{
	GenericLangAbbrev: GenericLangBaseFieldCategory, // Generic
	GermanLangAbbrev:  GermanLangBaseFieldCategory,
	EnglishLangAbbrev: EnglishLangBaseFieldCategory,
}

// This is the "enum" for categories of the (primary) Solr backing fields for MEx fields
const (
	GenericLangBaseFieldCategory   = "GENERIC_BASE_FIELD_CATEGORY"
	GermanLangBaseFieldCategory    = "GERMAN_LANGBASE_FIELD_CATEGORY"
	EnglishLangBaseFieldCategory   = "ENGLISH_LANG_BASE_FIELD_CATEGORY"
	ParentCodesBaseFieldCategory   = "PARENT_CODES_BASE_FIELD_CATEGORY"
	NormalizedBaseFieldCategory    = "NORMALIZED_BASE_FIELD_CATEGORY"
	RawContentBaseFieldCategory    = "RAW_CONTENT_BASE_FIELD_CATEGORY"
	PrefixContentBaseFieldCategory = "PREFIX_CONTENT_BASE_FIELD_CATEGORY"
)

// This is the "enum" for the functional categories of Solr field backing axes or foci
const (
	GenericLangSearchFunctionCategory = "GENERIC_LANG_SEARCH_FUNCTION_CATEGORY"
	GermanLangSearchFunctionCategory  = "GERMAN_LANG_SEARCH_FUNCTION_CATEGORY"
	EnglishLangSearchFunctionCategory = "ENGLISH_LANG_SEARCH_FUNCTION_CATEGORY"
	PrefixSearchFunctionCategory      = "PREFIX_SEARCH_FUNCTION_CATEGORY"
	RawSearchFunctionCategory         = "RAW_SEARCH_FUNCTION_CATEGORY"
	FacetAndFilterFunctionCategory    = "FACET_AND_FILTER_FUNCTION_CATEGORY"
	SingleValueFacetFunctionCategory  = "SINGLE_VALUE_FACET_FUNCTION_CATEGORY"
	SortFunctionCategory              = "SORT_FUNCTION_CATEGORY"
)
