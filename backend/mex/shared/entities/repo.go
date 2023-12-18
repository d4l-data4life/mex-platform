package entities

import (
	"context"
)

type EntityRepo interface {
	GetEntityTypeNames(ctx context.Context, focalOnly bool) ([]string, error)
	GetEntityType(ctx context.Context, entityTypeName string) (*EntityType, error)
	ListEntityTypes(ctx context.Context, focalOnly bool) ([]*EntityType, error)
	Purge(ctx context.Context) error
}
