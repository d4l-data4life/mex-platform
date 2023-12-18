package kind_number // revive:disable

import (
	"context"
	"fmt"
	"strconv"

	fieldUtils "github.com/d4l-data4life/mex/mex/shared/fields"
	"github.com/d4l-data4life/mex/mex/shared/solr"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/datamodel"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
)

const KindName = "number"

type KindNumber struct{}

func (kind *KindNumber) ValidateDefinition(_ context.Context, request *fieldUtils.FieldDef) (fields.BaseFieldDef, error) {
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

func (kind *KindNumber) MustValidateDefinition(ctx context.Context, request *fieldUtils.FieldDef) fields.BaseFieldDef {
	fieldDef, err := kind.ValidateDefinition(ctx, request)
	if err != nil {
		panic(err)
	}
	return fieldDef
}

func (kind *KindNumber) MarshalToProtobufFormat(_ context.Context, fieldDef fields.BaseFieldDef) (*fieldUtils.FieldDef, error) {
	return &fieldUtils.FieldDef{
		Name:      fieldDef.Name(),
		Kind:      fieldDef.Kind(),
		DisplayId: fieldDef.DisplayID(),
		IndexDef: &fieldUtils.IndexDef{
			MultiValued: fieldDef.MultiValued(),
		},
	}, nil
}

func (kind *KindNumber) ValidateFieldValue(_ context.Context, _ fields.BaseFieldDef, fieldValue string) error {
	_, err := strconv.ParseFloat(fieldValue, 64)
	return err
}

func (kind *KindNumber) GenerateSolrFields(_ context.Context, fieldDef fields.BaseFieldDef) (solr.FieldCategoryToSolrFieldDefsMap, error) {
	if fieldDef == nil {
		return nil, fmt.Errorf("cannot generate backing fields from empty field definition")
	}
	genericTypeName := solr.KnownLanguagesCategoryMap[solr.GenericLangAbbrev]
	solrFields := solr.FieldCategoryToSolrFieldDefsMap{
		genericTypeName: solr.GetStandardPrimaryBackingField(fieldDef.Name(), solr.DefaultSolrNumberFieldType, fieldDef.MultiValued()),
	}
	return solrFields, nil
}

func (*KindNumber) GenerateXMLFieldTags(_ context.Context, _ fields.BaseFieldDef, itemValue datamodel.CurrentItemValue) ([]string, error) {
	return []string{fmt.Sprintf("<field name=\"%s\">%s</field>\n", itemValue.FieldName, itemValue.FieldValue)}, nil
}

func (kind *KindNumber) EnrichFacetBucket(_ context.Context, bucket *solr.FacetBucket, _ fields.BaseFieldDef) (*solr.FacetBucket, error) {
	return bucket, nil
}

func (kind *KindNumber) ResetCaches() {}
