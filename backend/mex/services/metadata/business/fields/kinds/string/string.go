package kind_string // revive:disable

import (
	"context"
	"fmt"

	fieldUtils "github.com/d4l-data4life/mex/mex/shared/fields"
	"github.com/d4l-data4life/mex/mex/shared/solr"
	"github.com/d4l-data4life/mex/mex/shared/utils"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/datamodel"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
)

const KindName = "string"

type KindString struct{}

func (kind *KindString) ValidateDefinition(_ context.Context, request *fieldUtils.FieldDef) (fields.BaseFieldDef, error) {
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

func (kind *KindString) MustValidateDefinition(ctx context.Context, request *fieldUtils.FieldDef) fields.BaseFieldDef {
	fieldDef, err := kind.ValidateDefinition(ctx, request)
	if err != nil {
		panic(err)
	}
	return fieldDef
}

func (kind *KindString) MarshalToProtobufFormat(_ context.Context, fieldDef fields.BaseFieldDef) (*fieldUtils.FieldDef, error) {
	return &fieldUtils.FieldDef{
		Name:      fieldDef.Name(),
		Kind:      fieldDef.Kind(),
		DisplayId: fieldDef.DisplayID(),
		IndexDef: &fieldUtils.IndexDef{
			MultiValued: fieldDef.MultiValued(),
		},
	}, nil
}

func (kind *KindString) ValidateFieldValue(context.Context, fields.BaseFieldDef, string) error {
	return nil
}

func (kind *KindString) GenerateSolrFields(_ context.Context, fieldDef fields.BaseFieldDef) (solr.FieldCategoryToSolrFieldDefsMap, error) {
	if fieldDef == nil {
		return nil, fmt.Errorf("cannot generate backing fields from empty field definition")
	}
	solrFields := solr.FieldCategoryToSolrFieldDefsMap{
		solr.GenericLangBaseFieldCategory: solr.GetStandardPrimaryBackingField(fieldDef.Name(), solr.DefaultSolrStringFieldType, fieldDef.MultiValued()),
		solr.NormalizedBaseFieldCategory: solr.GetStandardPrimaryBackingField(solr.GetNormalizedBackingFieldName(fieldDef.Name()),
			solr.DefaultSolrSortableTextFieldType, fieldDef.MultiValued()),
	}
	return solrFields, nil
}

func (*KindString) GenerateXMLFieldTags(_ context.Context, _ fields.BaseFieldDef, itemValue datamodel.CurrentItemValue) ([]string, error) {
	cleanedValue := utils.SanitizeXML(itemValue.FieldValue)
	normalizedValue := utils.NormalizeString(itemValue.FieldValue)
	cleanedNormalizedValue := utils.SanitizeXML(normalizedValue)

	returnTags := []string{fmt.Sprintf("<field name=\"%s\">%s</field>", itemValue.FieldName, cleanedValue)}
	returnTags = append(returnTags, fmt.Sprintf("<field name=\"%s\">%s</field>", solr.GetNormalizedBackingFieldName(itemValue.FieldName), cleanedNormalizedValue))

	return returnTags, nil
}

func (kind *KindString) EnrichFacetBucket(_ context.Context, bucket *solr.FacetBucket, _ fields.BaseFieldDef) (*solr.FacetBucket, error) {
	return bucket, nil
}

func (kind *KindString) ResetCaches() {}
