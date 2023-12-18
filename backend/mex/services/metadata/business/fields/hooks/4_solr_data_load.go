package hooks

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/d4l-data4life/mex/mex/shared/codings/csrepo"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"

	kind_coding "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/coding"
	kind_hierarchy "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/hierarchy"
	kind_link "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/link"
	kind_number "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/number"
	kind_string "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/string"
	kind_text "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/text"
	kind_timestamp "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/timestamp"
)

type SolrDataLoadHooks map[string]fields.LifecycleSolrDataLoadHook

type SolrDataLoadHooksConfig struct {
	DB            *pgxpool.Pool
	CodingsetRepo csrepo.CodingsetRepo
}

func NewSolrDataLoadHooks(cfg SolrDataLoadHooksConfig) (SolrDataLoadHooks, error) {
	hooks := make(SolrDataLoadHooks)

	hooks[kind_number.KindName] = &kind_number.KindNumber{}
	hooks[kind_string.KindName] = &kind_string.KindString{}
	hooks[kind_text.KindName] = &kind_text.KindText{}
	hooks[kind_timestamp.KindName] = &kind_timestamp.KindTimestamp{}
	hooks[kind_link.KindName] = &kind_link.KindLink{}

	kindHierarchy, err := kind_hierarchy.NewKindHierarchy(cfg.DB)
	if err != nil {
		return nil, err
	}
	hooks[kind_hierarchy.KindName] = kindHierarchy

	hooks[kind_coding.KindName] = &kind_coding.KindCoding{
		CodingsetRepo: cfg.CodingsetRepo,
	}

	return hooks, nil
}

func (hooks SolrDataLoadHooks) GetHook(fieldKind string) fields.LifecycleSolrDataLoadHook {
	return hooks[fieldKind]
}
