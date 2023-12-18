package erepo

import (
	"context"
	"fmt"

	"github.com/d4l-data4life/mex/mex/shared/entities"
)

type mockedEntitiesRepo struct {
	entityObjects []*entities.EntityType
}

func NewMockedEntityTypesRepo(ents []*entities.EntityType) entities.EntityRepo {
	return &mockedEntitiesRepo{
		entityObjects: ents,
	}
}

func (repo *mockedEntitiesRepo) GetEntityTypeNames(_ context.Context, focalOnly bool) ([]string, error) {
	retTypeNames := make([]string, len(repo.entityObjects))
	for i, entityType := range repo.entityObjects {
		if !focalOnly || entityType.Config.IsFocal {
			retTypeNames[i] = entityType.Name
		}
	}
	return retTypeNames, nil
}

func (repo *mockedEntitiesRepo) GetEntityType(_ context.Context, entityTypeName string) (*entities.EntityType, error) {
	for _, entityType := range repo.entityObjects {
		if entityType.Name == entityTypeName {
			return entityType, nil
		}
	}
	return nil, fmt.Errorf("no entity type found with name: %s", entityTypeName)
}

func (repo *mockedEntitiesRepo) ListEntityTypes(_ context.Context, focalOnly bool) ([]*entities.EntityType, error) {
	var retEntityTypes []*entities.EntityType
	for _, entityType := range repo.entityObjects {
		if !focalOnly || entityType.Config.IsFocal {
			retEntityTypes = append(retEntityTypes, entityType)
		}
	}
	return retEntityTypes, nil
}

func (*mockedEntitiesRepo) Purge(ctx context.Context) error { return nil }
