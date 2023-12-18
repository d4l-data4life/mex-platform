package kind_link // revive:disable

import (
	"context"
	"fmt"

	"github.com/d4l-data4life/mex/mex/shared/solr"

	"google.golang.org/protobuf/types/known/anypb"

	fieldUtils "github.com/d4l-data4life/mex/mex/shared/fields"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/datamodel"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
)

const KindName = "link"

type KindLink struct{}

type linkFieldDef struct {
	fields.BaseFieldDef

	relationType       string
	linkedTargetFields []string
}

type LinkFieldDef interface {
	fields.BaseFieldDef

	RelationType() string
	LinkedTargetFields() []string
}

func (def *linkFieldDef) RelationType() string         { return def.relationType }
func (def *linkFieldDef) LinkedTargetFields() []string { return def.linkedTargetFields }

func (kind *KindLink) ValidateDefinition(_ context.Context, fieldDef *fieldUtils.FieldDef) (fields.BaseFieldDef, error) {
	if fieldDef == nil {
		return nil, fmt.Errorf("cannot generate backing fields from empty field definition")
	}
	err := fieldUtils.ValidateName(fieldDef.Name)
	if err != nil {
		return nil, err
	}

	if fieldDef.Kind != KindName {
		return nil, fmt.Errorf("kind is not %s: %s", KindName, fieldDef.Kind)
	}

	extLink, err := fields.GetFirstLinkExt(fieldDef.IndexDef)
	if err != nil {
		return nil, err
	}

	return &linkFieldDef{
		BaseFieldDef: fields.NewBaseFieldDef(fieldDef.Name, fieldDef.Kind, fieldDef.DisplayId, false, fields.BaseIndexDef{
			MultiValued: fieldDef.IndexDef.MultiValued,
		}),
		relationType:       extLink.RelationType,
		linkedTargetFields: extLink.LinkedTargetFields,
	}, nil
}

func (kind *KindLink) MustValidateDefinition(ctx context.Context, request *fieldUtils.FieldDef) fields.BaseFieldDef {
	fieldDef, err := kind.ValidateDefinition(ctx, request)
	if err != nil {
		panic(err)
	}
	return fieldDef
}

func (kind *KindLink) MarshalToProtobufFormat(_ context.Context, fieldDef fields.BaseFieldDef) (*fieldUtils.FieldDef, error) {
	if lFieldDef, ok := (fieldDef).(LinkFieldDef); ok {
		ext, err := anypb.New(&fieldUtils.IndexDefExtLink{
			RelationType:       lFieldDef.RelationType(),
			LinkedTargetFields: lFieldDef.LinkedTargetFields(),
		})
		if err != nil {
			return nil, err
		}

		return &fieldUtils.FieldDef{
			Name:      fieldDef.Name(),
			Kind:      fieldDef.Kind(),
			DisplayId: fieldDef.DisplayID(),
			IndexDef: &fieldUtils.IndexDef{
				MultiValued: fieldDef.MultiValued(),
				Ext:         []*anypb.Any{ext},
			},
		}, nil
	}

	return nil, fmt.Errorf("field definition object is not a LinkFieldDef")
}

func (kind *KindLink) ValidateFieldValue(_ context.Context, _ fields.BaseFieldDef, _ string) error {
	return nil
}

func (kind *KindLink) GenerateSolrFields(_ context.Context, fieldDef fields.BaseFieldDef) (solr.FieldCategoryToSolrFieldDefsMap, error) {
	if fieldDef == nil {
		return nil, fmt.Errorf("cannot generate backing fields from empty field definition")
	}
	genericTypeName := solr.KnownLanguagesCategoryMap[solr.GenericLangAbbrev]
	solrFields := solr.FieldCategoryToSolrFieldDefsMap{
		genericTypeName: solr.GetStandardPrimaryBackingField(fieldDef.Name(), solr.DefaultSolrStringFieldType, fieldDef.MultiValued()),
	}
	return solrFields, nil
}

func (*KindLink) GenerateXMLFieldTags(_ context.Context, _ fields.BaseFieldDef, itemValue datamodel.CurrentItemValue) ([]string, error) {
	return []string{fmt.Sprintf("<field name=\"%s\">%s</field>\n", itemValue.FieldName, itemValue.FieldValue)}, nil
}

func (kind *KindLink) EnrichFacetBucket(_ context.Context, bucket *solr.FacetBucket, _ fields.BaseFieldDef) (*solr.FacetBucket, error) {
	return bucket, nil
}

func (kind *KindLink) ResetCaches() {}
