package sctypes

import (
	"fmt"

	kindHierarchy "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/hierarchy"
	"github.com/d4l-data4life/mex/mex/shared/solr"

	kind_number "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/number"
	kind_timestamp "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/timestamp"
)

// GetOrdinalAxisFieldType returns the appropriate Solr field type for a given ordinal axis
func GetOrdinalAxisFieldType(axisFields []string, mexFieldMap map[string]string) (string, error) {
	if len(axisFields) == 0 {
		return "", fmt.Errorf("no ordinal axis fields given")
	}
	// Check if all fields have the same MEx type anf, if so, find that type
	firstKind := ""
	for i, field := range axisFields {
		newKind, ok := mexFieldMap[field]
		if !ok {
			return "", fmt.Errorf("could not find MEx kind of the field '%s' (field probably not configured)", field)
		}
		if newKind == kindHierarchy.KindName {
			return "", fmt.Errorf("ordinal axis cannot contain fields of the kinds hierarchy and coding")
		}
		if i == 0 {
			firstKind = newKind
			continue
		}
		// Different field types --> sortable text field
		if newKind != firstKind {
			return solr.DefaultSolrSortableTextFieldType, nil
		}
	}

	switch firstKind {
	// NOTE THAT ANY SOLR FIELD TYPE RETURNED MUST SUPPORT SORTING
	case kind_timestamp.KindName:
		return solr.DefaultSolrTimestampFieldType, nil
	case kind_number.KindName:
		return solr.DefaultSolrNumberFieldType, nil
	default:
		return solr.DefaultSolrSortableTextFieldType, nil
	}
}
