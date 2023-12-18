package kind_timestamp // revive:disable

import (
	"context"
	"fmt"
	"time"

	fieldUtils "github.com/d4l-data4life/mex/mex/shared/fields"
	"github.com/d4l-data4life/mex/mex/shared/solr"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/datamodel"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
)

const KindName = "timestamp"

type KindTimestamp struct{}

func (kind *KindTimestamp) ValidateDefinition(_ context.Context, request *fieldUtils.FieldDef) (fields.BaseFieldDef, error) {
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

func (kind *KindTimestamp) MustValidateDefinition(ctx context.Context, request *fieldUtils.FieldDef) fields.BaseFieldDef {
	fieldDef, err := kind.ValidateDefinition(ctx, request)
	if err != nil {
		panic(err)
	}
	return fieldDef
}

func (kind *KindTimestamp) MarshalToProtobufFormat(_ context.Context, fieldDef fields.BaseFieldDef) (*fieldUtils.FieldDef, error) {
	return &fieldUtils.FieldDef{
		Name:      fieldDef.Name(),
		Kind:      fieldDef.Kind(),
		DisplayId: fieldDef.DisplayID(),
		IndexDef: &fieldUtils.IndexDef{
			MultiValued: fieldDef.MultiValued(),
		},
	}, nil
}

func (kind *KindTimestamp) ValidateFieldValue(_ context.Context, _ fields.BaseFieldDef, _ string) error {
	return nil
}

func (kind *KindTimestamp) GenerateSolrFields(_ context.Context, fieldDef fields.BaseFieldDef) (solr.FieldCategoryToSolrFieldDefsMap, error) {
	if fieldDef == nil {
		return nil, fmt.Errorf("cannot generate backing fields from empty field definition")
	}
	solrFields := solr.FieldCategoryToSolrFieldDefsMap{
		solr.GenericLangBaseFieldCategory: solr.GetStandardPrimaryBackingField(fieldDef.Name(), solr.DefaultSolrTimestampFieldType, fieldDef.MultiValued()),
		solr.RawContentBaseFieldCategory:  solr.GetStandardPrimaryBackingField(solr.GetRawValTimestampName(fieldDef.Name()), solr.DefaultSolrStringFieldType, fieldDef.MultiValued()),
	}
	return solrFields, nil
}

func (*KindTimestamp) GenerateXMLFieldTags(_ context.Context, _ fields.BaseFieldDef, itemValue datamodel.CurrentItemValue) ([]string, error) {
	// Parse timestamp string, allowing for different levels of precision
	t, err := time.Parse("2006-01-02T15:04:05Z", itemValue.FieldValue)
	if err != nil {
		t, err = time.Parse("2006-01-02", itemValue.FieldValue)
		if err != nil {
			t, err = time.Parse("2006-01", itemValue.FieldValue)
			if err != nil {
				t, err = time.Parse("2006", itemValue.FieldValue)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return []string{
		fmt.Sprintf("<field name=\"%s\">%s</field>", itemValue.FieldName, t.Format("2006-01-02T15:04:05Z")),
		fmt.Sprintf("<field name=\"%s\">%s</field>", solr.GetRawValTimestampName(itemValue.FieldName),
			itemValue.FieldValue),
	}, nil
}

func (kind *KindTimestamp) EnrichFacetBucket(_ context.Context, bucket *solr.FacetBucket, _ fields.BaseFieldDef) (*solr.FacetBucket, error) {
	return bucket, nil
}

func (kind *KindTimestamp) ResetCaches() {}
