package hooks

import (
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"

	kind_coding "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/coding"
	kind_hierarchy "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/hierarchy"
	kind_link "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/link"
	kind_number "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/number"
	kind_string "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/string"
	kind_text "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/text"
	kind_timestamp "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/timestamp"
)

type SolrFieldCreationHooks map[string]fields.LifecycleSolrFieldCreationHook

type SolrFieldCreationHooksConfig struct{}

func NewSolrFieldCreationHooks(_ SolrFieldCreationHooksConfig) (SolrFieldCreationHooks, error) {
	hooks := make(SolrFieldCreationHooks)

	hooks[kind_number.KindName] = &kind_number.KindNumber{}
	hooks[kind_string.KindName] = &kind_string.KindString{}
	hooks[kind_text.KindName] = &kind_text.KindText{}
	hooks[kind_timestamp.KindName] = &kind_timestamp.KindTimestamp{}
	hooks[kind_hierarchy.KindName] = &kind_hierarchy.KindHierarchy{}
	hooks[kind_link.KindName] = &kind_link.KindLink{}

	hooks[kind_coding.KindName] = &kind_coding.KindCoding{}

	return hooks, nil
}

func (hooks SolrFieldCreationHooks) GetHook(fieldKind string) fields.LifecycleSolrFieldCreationHook {
	return hooks[fieldKind]
}
