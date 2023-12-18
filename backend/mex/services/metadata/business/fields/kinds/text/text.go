package kind_text // revive:disable

import (
	"context"
	"fmt"

	fieldUtils "github.com/d4l-data4life/mex/mex/shared/fields"
	"github.com/d4l-data4life/mex/mex/shared/solr"
	"github.com/d4l-data4life/mex/mex/shared/utils"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/datamodel"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
)

const KindName = "text"

type KindText struct{}

func (kind *KindText) ValidateDefinition(_ context.Context, request *fieldUtils.FieldDef) (fields.BaseFieldDef, error) {
	err := fieldUtils.ValidateName(request.Name)
	if err != nil {
		return nil, err
	}

	if request.Kind != KindName {
		return nil, fmt.Errorf("kind is not %s: %s", KindName, request.Kind)
	}

	return fields.NewBaseFieldDef(request.Name, request.Kind, request.DisplayId, false, fields.BaseIndexDef{
		MultiValued: request.IndexDef.MultiValued,
	}), nil
}

func (kind *KindText) MustValidateDefinition(ctx context.Context, request *fieldUtils.FieldDef) fields.BaseFieldDef {
	fieldDef, err := kind.ValidateDefinition(ctx, request)
	if err != nil {
		panic(err)
	}
	return fieldDef
}

func (kind *KindText) MarshalToProtobufFormat(_ context.Context, fieldDef fields.BaseFieldDef) (*fieldUtils.FieldDef, error) {
	return &fieldUtils.FieldDef{
		Name:      fieldDef.Name(),
		Kind:      fieldDef.Kind(),
		DisplayId: fieldDef.DisplayID(),
		IndexDef: &fieldUtils.IndexDef{
			MultiValued: fieldDef.MultiValued(),
		},
	}, nil
}

func (kind *KindText) ValidateFieldValue(_ context.Context, _ fields.BaseFieldDef, _ string) error {
	return nil
}

func (kind *KindText) GenerateSolrFields(_ context.Context, fieldDef fields.BaseFieldDef) (solr.FieldCategoryToSolrFieldDefsMap, error) {
	if fieldDef == nil {
		return nil, fmt.Errorf("cannot generate backing fields from empty field definition")
	}
	backingFieldsMap := make(solr.FieldCategoryToSolrFieldDefsMap)
	mexToSolrMap := solr.GetTextCoreBackingFieldInfo(fieldDef.Name(), false)
	var isMultivalued bool
	for category, coreFieldInfo := range mexToSolrMap {
		if category == solr.PrefixContentBaseFieldCategory {
			isMultivalued = true
		} else {
			isMultivalued = fieldDef.MultiValued()
		}
		backingFieldsMap[category] = solr.GetStandardPrimaryBackingField(coreFieldInfo.SolrName, coreFieldInfo.SolrType, isMultivalued)
	}

	return backingFieldsMap, nil
}

func (*KindText) GenerateXMLFieldTags(_ context.Context, _ fields.BaseFieldDef, itemValue datamodel.CurrentItemValue) ([]string, error) {
	// Default is generic language field
	langTargetField, err := solr.GetLangSpecificFieldName(itemValue.FieldName, solr.GenericLangAbbrev)
	if err != nil {
		return nil, err
	}
	// If language is set & recognized, copy to language-specific field
	if itemValue.Language.Valid {
		langCode := itemValue.Language.String
		if _, ok := solr.KnownLanguagesFieldTypeMap[langCode]; ok {
			fn, err := solr.GetLangSpecificFieldName(itemValue.FieldName, langCode)
			if err != nil {
				return nil, err
			}
			langTargetField = fn
		}
	}
	cleanedValue := utils.SanitizeXML(itemValue.FieldValue)
	normalizedValue := utils.NormalizeString(itemValue.FieldValue)
	cleanedNormalizedValue := utils.SanitizeXML(normalizedValue)
	returnTags := []string{fmt.Sprintf("<field name=\"%s\">%s</field>", langTargetField, cleanedValue)}
	// Also add data to prefix, normalized, and unanalyzed fields
	returnTags = append(returnTags, fmt.Sprintf("<field name=\"%s\">%s</field>", solr.GetPrefixBackingFieldName(itemValue.FieldName), cleanedValue))
	returnTags = append(returnTags, fmt.Sprintf("<field name=\"%s\">%s</field>", solr.GetRawBackingFieldName(itemValue.FieldName), cleanedValue))
	returnTags = append(returnTags, fmt.Sprintf("<field name=\"%s\">%s</field>", solr.GetNormalizedBackingFieldName(itemValue.FieldName), cleanedNormalizedValue))
	return returnTags, nil
}

func (kind *KindText) EnrichFacetBucket(_ context.Context, bucket *solr.FacetBucket, _ fields.BaseFieldDef) (*solr.FacetBucket, error) {
	return bucket, nil
}

func (kind *KindText) ResetCaches() {}
