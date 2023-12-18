package linked

import (
	"fmt"

	"github.com/d4l-data4life/mex/mex/shared/solr"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
	kindHierarchy "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/hierarchy"
	kindLink "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/link"
)

// GetLinkedFieldDefs returns all linked fields
func GetLinkedFieldDefs(fieldDefs []fields.BaseFieldDef) ([]fields.BaseFieldDef, error) {
	var linkedFieldDefs []fields.BaseFieldDef
	for _, fdBase := range fieldDefs {
		// Skip fields that do not represent links (i.e. are not link or hierarchy fields)
		fd, isLinkType := fdBase.(kindLink.LinkFieldDef)
		if !isLinkType {
			fd, isLinkType = fdBase.(kindHierarchy.HierarchyFieldDef)
			if !isLinkType {
				continue
			}
		}

		for _, targetFieldName := range fd.LinkedTargetFields() {
			targetFieldDef, err := getFieldDefByName(fieldDefs, targetFieldName)
			if err != nil {
				// Ignore linked field if target cannot be found
				continue
			}
			/*
				Linked field is configured identically to the target field, except
				1. their name (always different from target)
				2. their display ID (always different from target)
				3. whether they are multivalued (depends on both link and target fields)
			*/
			linkedFieldName := solr.GetLinkedFieldName(fd.Name(), targetFieldName)
			linkedFieldDisplayID := solr.GetLinkedDisplayID(linkedFieldName)
			linkedFieldMultiValued := fd.MultiValued() || targetFieldDef.MultiValued()
			linkedFieldIndexDef := fields.BaseIndexDef{
				MultiValued: linkedFieldMultiValued,
			}
			newLinkedFieldDef := fields.NewBaseFieldDef(linkedFieldName, targetFieldDef.Kind(), linkedFieldDisplayID,
				true, linkedFieldIndexDef)
			linkedFieldDefs = append(linkedFieldDefs, newLinkedFieldDef)
		}
	}
	return linkedFieldDefs, nil
}

// getFieldDefByName return the field definition with the given name, otherwise error
func getFieldDefByName(fieldDefs []fields.BaseFieldDef, targetName string) (fields.BaseFieldDef, error) {
	for _, fd := range fieldDefs {
		if fd.Name() == targetName {
			return fd, nil
		}
	}
	return nil, fmt.Errorf("no field with name '%s' found", targetName)
}
