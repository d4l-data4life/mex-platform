package frepo

import (
	"fmt"

	"github.com/d4l-data4life/mex/mex/shared/solr"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
	kindString "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/string"
	kindTimestamp "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/timestamp"
)

func getPredefinedFields() []fields.BaseFieldDef {
	result := []fields.BaseFieldDef{}
	// We loop over the list of predefined field names to ensure the output is aligned with that list.
	for _, coreFieldName := range solr.CoreMexFieldNames {
		switch coreFieldName {
		case solr.DefaultUniqueKey:
			// Unique item ID field
			result = append(result, fields.NewBaseFieldDef(
				solr.DefaultUniqueKey,
				kindString.KindName,
				"FIELD_ID",
				false,
				fields.BaseIndexDef{
					MultiValued: false,
				},
			))
		case solr.ItemEntityNameField:
			// Entity type of item
			result = append(result, fields.NewBaseFieldDef(
				solr.ItemEntityNameField,
				kindString.KindName,
				"FIELD_ENTITY_NAME",
				false,
				fields.BaseIndexDef{
					MultiValued: false,
				},
			),
			)
		case solr.ItemCreatedAtField:
			// Time at which the item was created in MEx
			result = append(result, fields.NewBaseFieldDef(
				solr.ItemCreatedAtField,
				kindTimestamp.KindName,
				"FIELD_CREATED_AT",
				false,
				fields.BaseIndexDef{
					MultiValued: false,
				},
			),
			)
		case solr.ItemBusinessIDField:
			// Business ID of item
			result = append(result, fields.NewBaseFieldDef(
				solr.ItemBusinessIDField,
				kindString.KindName,
				"FIELD_BUSINESS_ID",
				false,
				fields.BaseIndexDef{
					MultiValued: false,
				},
			))
		default:
			panic(fmt.Errorf("default Solr field configs do not align with the list of predefined Solr field names"))
		}
	}
	if len(result) != len(solr.CoreMexFieldNames) {
		panic(fmt.Errorf("default Solr field configs do not align with the list of predefined Solr field names"))
	}
	return result
}
