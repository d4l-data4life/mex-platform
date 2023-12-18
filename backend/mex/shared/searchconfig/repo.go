package searchconfig

import (
	"context"
)

//revive:disable-next-line:exported
type SearchConfigRepo interface {
	ListSearchConfigs(ctx context.Context) (*SearchConfigList, error)
	ListSearchConfigsOfType(ctx context.Context, objType string) (*SearchConfigList, error)
	GetSearchConfigObject(ctx context.Context, id string) (*SearchConfigObject, error)
	GetFieldsForSearchFocus(ctx context.Context, requestedSearchFocus string) ([]string, error)
	GetFieldsForAxis(ctx context.Context, requestedOrdinalAxis string) ([]string, error)
	Purge(ctx context.Context) error
}
