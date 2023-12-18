package solr

import (
	"fmt"
	"regexp"
	"strings"
)

type BackingFieldDefWithCategory struct {
	Def             FieldDef
	FieldCategoryID string // Functional category of field
}

/*
FieldCategoryToSolrFieldDefsMap maps, for a single MEx field, the code for the backing field category to an object with
the corresponding Solr backing field definition and flags indicating whether it is relevant for search/faceting/
sorting.
*/
type FieldCategoryToSolrFieldDefsMap map[string]FieldDef

// Wiring information for a single backing field
type MexBackingFieldWiringInfo struct {
	Name     string
	Category string
}

// MexFieldWiringInfo stores information about a single MEx field that is needed for wiring it up
type MexFieldWiringInfo struct {
	MexType       string
	BackingFields []MexBackingFieldWiringInfo
}

// MexFieldBackingInfoMap maps the MEx name of a field to information about it and its backing fields
type MexFieldBackingInfoMap map[string]MexFieldWiringInfo

// GetSearchFocusFieldName returns the name of the auxiliary field for a given search focus
func GetSearchFocusFieldName(name string) string {
	return fmt.Sprintf("%s_%s", name, FocusPostfix)
}

// GetPrefixBackingFieldName returns the name of the prefix backing field
func GetPrefixBackingFieldName(name string) string {
	return fmt.Sprintf("%s%s%s", name, LongSeparator, PrefixFocusPostfix)
}

// GetNormalizedBackingFieldName returns the name of the normalized backing field
func GetNormalizedBackingFieldName(name string) string {
	return fmt.Sprintf("%s%s%s", name, LongSeparator, NormFocusPostfix)
}

// GetRawBackingFieldName returns the name of the unanalyzed backing field
func GetRawBackingFieldName(name string) string {
	return fmt.Sprintf("%s%s%s", name, LongSeparator, RawValuePostfix)
}

// GetOrdinalAxisFacetAndFilterFieldName returns the name of the auxiliary field for a given ordinal axis
func GetOrdinalAxisFacetAndFilterFieldName(axisName string) string {
	return fmt.Sprintf("%s_%s", axisName, AxisFacetPostfix)
}

// GetOrdinalAxisSortFieldName returns the name of the auxiliary field for a given ordinal axis
func GetOrdinalAxisSortFieldName(axisName string) string {
	return fmt.Sprintf("%s_%s", axisName, AxisSortPostfix)
}

// GetSingleNodeAxisFieldName returns the name of the single-node field for a hierarchical axis
func GetSingleNodeAxisFieldName(axisName string) string {
	return fmt.Sprintf("%s_%s", axisName, SingleValuePostfix)
}

// GetLangSpecificFieldName returns the name of the auxiliary field for holding text for a given field in s specific
// language. If no language is given, the generic language name is returned. If the language is unknown, an error is
// returned.
func GetLangSpecificFieldName(name string, langCode string) (string, error) {
	if langCode == "" {
		return name + LongSeparator + GenericLanguageSuffix, nil
	}
	if _, ok := KnownLanguagesFieldTypeMap[langCode]; ok {
		return name + LongSeparator + langCode, nil
	}
	return "", fmt.Errorf("invalid language code")
}

// GetExactFieldName returns the name of the exact auxiliary field corresponding to a given primary field
func GetExactFieldName(name string) string {
	return fmt.Sprintf("%s_%s", name, ExactPostfix)
}

// GetRawValTimestampName returns the name of the timestamp backing field for a given primary field (should be a date field)
func GetRawValTimestampName(name string) string {
	return fmt.Sprintf("%s_%s", name, RawValTimestampPostfix)
}

// GetLinkedFieldName returns the standard linked field name given a link field name and the name of a field on the
// target
func GetLinkedFieldName(linkFieldName string, targetFieldName string) string {
	return linkFieldName + LinkedFieldSeparator + targetFieldName
}

// GetLinkedDisplayID returns the standard display name for a linked field name
func GetLinkedDisplayID(linkedFieldName string) string {
	isUpperCase := regexp.MustCompile("[[:upper:]]")
	var stringChar string
	normalizedFieldName := ""
	prevStringChar := ""
	for pos, char := range linkedFieldName {
		stringChar = string(char)
		/*
			Insert an underscore if
			1. the next character is uppercase AND
			2. we are not at the start of the string AND
			3. the last character we added was not an underscore
			The last condition ensures that an underscore followed by a capital letter is kept
			as-is and does not lead to a double underscore
		*/
		if isUpperCase.MatchString(stringChar) && pos != 0 && prevStringChar != "_" {
			normalizedFieldName += "_"
		}
		normalizedFieldName += stringChar
		prevStringChar = stringChar
	}
	return strings.ToUpper("FIELD_" + normalizedFieldName)
}

// CreateMergeID generates the business ID for merged item from the business ID of the undelrying fragments
func CreateMergeID(rawID string) string {
	newBusinessID := rawID + MergeIDSeparator + MergeIDPostfix
	return newBusinessID
}

type FieldWithLanguage struct {
	Name     string
	Language string
}

type CoreBackingFieldInfo struct {
	SolrName string
	Language string
	SolrType string
}

func NewFacetResult(facetType string, axis string) *FacetResult {
	return &FacetResult{
		Type:    facetType,
		Axis:    axis,
		Buckets: make([]*FacetBucket, 0),
	}
}

func (h *Highlight) AddMatch(field FieldWithLanguage, matches []string) *Highlight {
	h.Matches = append(h.Matches, &FieldHighlight{
		FieldName: field.Name,
		Language:  field.Language,
		Snippets:  matches,
	})
	return h
}

func NewHighlight(docID string) *Highlight {
	return &Highlight{
		ItemId:  docID,
		Matches: make([]*FieldHighlight, 0),
	}
}

func NewDocItem(id string) *DocItem {
	return &DocItem{
		ItemId: id,
		Values: make([]*DocValue, 0),
	}
}

func (docItem *DocItem) AddValue(field FieldWithLanguage, fieldValue string) *DocItem {
	docItem.Values = append(docItem.Values, &DocValue{
		FieldName:  field.Name,
		FieldValue: fieldValue,
		Language:   field.Language,
	})
	return docItem
}

type ByDocValue []*DocValue

func (d ByDocValue) Len() int { return len(d) }

func (d ByDocValue) Less(i, j int) bool {
	if d[i].FieldName == d[j].FieldName {
		return d[i].FieldValue < d[j].FieldValue
	}
	return d[i].FieldName < d[j].FieldName
}

func (d ByDocValue) Swap(i, j int) { d[i], d[j] = d[j], d[i] }

// TagExpr returns the expression tagged with a tag generated from the second param using getTagNameForField()
func TagExpr(expr string, name string) (string, string) {
	tag := GenerateTagName(name)
	return fmt.Sprintf("{!tag=%s}%s", tag, expr), tag
}

// GetTagNameForField constructs the SolrName of Solr tag based on a string
func GenerateTagName(str string) string {
	return fmt.Sprintf("%s_%s", str, TagPostfix)
}

func GetTransitiveHullFieldName(fieldName string) string {
	return fmt.Sprintf("%s_%s", fieldName, TransitiveHullPostfix)
}

func GetTransitiveHullDisplayFieldName(fieldName string) string {
	return fmt.Sprintf("%s_%s", fieldName, TransitiveHullDisplayPostfix)
}

func GetDisplayFieldName(fieldName string, language string) string {
	return fmt.Sprintf("%s_display___%s", fieldName, language)
}

// GetTextCoreBackingFieldInfo constructs the Solr backing fields needed for a MEx text field
func GetTextCoreBackingFieldInfo(mexName string, isForSearchConfig bool) map[string]CoreBackingFieldInfo {
	mexToSolr := make(map[string]CoreBackingFieldInfo)

	var categoryNameMap map[string]string
	var prefixCat string
	var rawCat string
	if isForSearchConfig {
		categoryNameMap = knownLanguagesFieldFunctionCategoryMap
		prefixCat = PrefixSearchFunctionCategory
		rawCat = RawSearchFunctionCategory
	} else {
		categoryNameMap = KnownLanguagesCategoryMap
		prefixCat = PrefixContentBaseFieldCategory
		rawCat = RawContentBaseFieldCategory
		// Add a  field for normalized content when backing a field
		mexToSolr[NormalizedBaseFieldCategory] = CoreBackingFieldInfo{
			SolrName: GetNormalizedBackingFieldName(mexName),
			// Note that this field need not support search since it is only used for backing fields
			SolrType: DefaultSolrStringFieldType,
		}
	}

	// Language-specific fields
	var categoryName string
	for lc := range KnownLanguagesFieldTypeMap {
		// Since field names for different fields are distinct and cannot contain double underscores,
		// langSpecFieldName is unique for every (field, Language)-combination
		langSpecFieldName, _ := GetLangSpecificFieldName(mexName, lc)
		categoryName = categoryNameMap[lc]
		mexToSolr[categoryName] = CoreBackingFieldInfo{
			SolrName: langSpecFieldName,
			Language: lc,
			SolrType: KnownLanguagesFieldTypeMap[lc],
		}
	}

	// Prefix field
	mexToSolr[prefixCat] = CoreBackingFieldInfo{
		SolrName: GetPrefixBackingFieldName(mexName),
		SolrType: DefaultPrefixSolrTextFieldType,
	}

	// Unanalyzed field
	mexToSolr[rawCat] = CoreBackingFieldInfo{
		SolrName: GetRawBackingFieldName(mexName),
		SolrType: DefaultRawSolrTextFieldType,
	}

	return mexToSolr
}

func GetStandardPrimaryBackingField(name string, solrType string, multivalued bool) FieldDef {
	return FieldDef{
		Name:         name,
		Type:         solrType,
		Stored:       true,
		Indexed:      false,
		MultiValued:  multivalued,
		DocValues:    false,
		Uninvertible: false,
	}
}

func GetStandardSecondaryBackingField(name string, solrType string, useDocValue bool) FieldDef {
	return FieldDef{
		Name:         name,
		Type:         solrType,
		Stored:       false,
		Indexed:      true,
		MultiValued:  true,
		DocValues:    useDocValue,
		Uninvertible: false,
	}
}
