package fields

import (
	"context"

	"github.com/d4l-data4life/mex/mex/shared/fields"
	"github.com/d4l-data4life/mex/mex/shared/solr"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/datamodel"
)

type FieldRepo interface {
	// GetFieldDefNames returns the names of all defined fields, including linked fields
	GetFieldDefNames(ctx context.Context) ([]string, error)
	GetFieldDefByName(ctx context.Context, fieldName string) (BaseFieldDef, error)
	GetFieldDefsByKind(ctx context.Context, fieldKind string) ([]BaseFieldDef, error)
	ListFieldDefs(ctx context.Context) ([]BaseFieldDef, error)
	Purge(ctx context.Context) error
}

type BaseFieldDef interface {
	Name() string
	Kind() string
	DisplayID() string

	MultiValued() bool
}

// +--------------+  MarshalToProtobufFormat   +-------------------+
// | (interface)  | -------------------------> |  (Protobuf)       |
// | BaseFieldDef |                            | *metaV0.FieldDef  |
// |              | <------------------------- |                   |
// +--------------+   ValidateDefinition       +-------------------+

type LifecycleFieldDefinitionHook interface {
	ValidateDefinition(ctx context.Context, request *fields.FieldDef) (BaseFieldDef, error)
	MustValidateDefinition(ctx context.Context, request *fields.FieldDef) BaseFieldDef
	MarshalToProtobufFormat(ctx context.Context, fieldDef BaseFieldDef) (*fields.FieldDef, error)
}

type LifecycleItemCreationHook interface {
	ValidateFieldValue(ctx context.Context, fieldDef BaseFieldDef, fieldValue string) error
}

type LifecycleSolrFieldCreationHook interface {
	// Get map from language code to generated Solr backing fields
	GenerateSolrFields(ctx context.Context, fieldDef BaseFieldDef) (solr.FieldCategoryToSolrFieldDefsMap, error)
}

type LifecycleSolrDataLoadHook interface {
	ResetCaches()
	GenerateXMLFieldTags(ctx context.Context, fieldDef BaseFieldDef, itemValue datamodel.CurrentItemValue) ([]string, error)
}

type LifecyclePostQueryHook interface {
	EnrichFacetBucket(ctx context.Context, bucket *solr.FacetBucket, fieldDef BaseFieldDef) (*solr.FacetBucket, error)
}
