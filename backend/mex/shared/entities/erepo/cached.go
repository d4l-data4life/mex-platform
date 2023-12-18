package erepo

import (
	"context"
	"fmt"
	"sync"

	"github.com/d4l-data4life/mex/mex/shared/entities"
	L "github.com/d4l-data4life/mex/mex/shared/log"
)

type entityTypesRepoCached struct {
	mu sync.RWMutex

	log      L.Logger
	delegate entities.EntityRepo

	cachedFocalOnlyEntityTypes []*entities.EntityType
	cachedEntityTypes          []*entities.EntityType
}

type NewCachedEntityTypesRepoParams struct {
	Log      L.Logger
	Delegate entities.EntityRepo
}

func NewCachedEntityTypesRepo(_ context.Context, params NewCachedEntityTypesRepoParams) entities.EntityRepo {
	return &entityTypesRepoCached{
		log:      params.Log,
		delegate: params.Delegate,
	}
}

func (repo *entityTypesRepoCached) GetEntityTypeNames(ctx context.Context, focalOnly bool) ([]string, error) {
	err := repo.ensureEntityTypesLoadedFromDelegate(ctx, false)
	if err != nil {
		return nil, err
	}

	repo.mu.RLock()
	defer repo.mu.RUnlock()

	entityTypeNames := []string{}

	if focalOnly {
		for _, entityType := range repo.cachedFocalOnlyEntityTypes {
			entityTypeNames = append(entityTypeNames, entityType.Name)
		}
	} else {
		for _, entityType := range repo.cachedEntityTypes {
			entityTypeNames = append(entityTypeNames, entityType.Name)
		}
	}

	return entityTypeNames, nil
}

func (repo *entityTypesRepoCached) GetEntityType(ctx context.Context, entityTypeName string) (*entities.EntityType, error) {
	err := repo.ensureEntityTypesLoadedFromDelegate(ctx, false)
	if err != nil {
		return nil, err
	}

	repo.mu.RLock()
	defer repo.mu.RUnlock()

	for _, entityType := range repo.cachedEntityTypes {
		if entityType.Name == entityTypeName {
			return entityType, nil
		}
	}

	return nil, fmt.Errorf("no entity type found with name: '%s'", entityTypeName)
}

func (repo *entityTypesRepoCached) ListEntityTypes(ctx context.Context, focalOnly bool) ([]*entities.EntityType, error) {
	err := repo.ensureEntityTypesLoadedFromDelegate(ctx, false)
	if err != nil {
		return nil, err
	}

	repo.mu.RLock()
	defer repo.mu.RUnlock()

	if focalOnly {
		return repo.cachedFocalOnlyEntityTypes, nil
	}

	return repo.cachedEntityTypes, nil
}

// ensureEntityTypesLoadedFromDelegate loads field configurations from another fields repo.
func (repo *entityTypesRepoCached) ensureEntityTypesLoadedFromDelegate(ctx context.Context, purge bool) error {
	var err error

	if purge || repo.cachedFocalOnlyEntityTypes == nil {
		repo.mu.Lock()
		defer repo.mu.Unlock()

		repo.cachedFocalOnlyEntityTypes, err = repo.delegate.ListEntityTypes(ctx, true)
		if err != nil {
			return err
		}

		repo.cachedEntityTypes, err = repo.delegate.ListEntityTypes(ctx, false)
		if err != nil {
			repo.cachedFocalOnlyEntityTypes = nil
			return err
		}
		repo.log.Info(ctx, L.Messagef("refreshed entity definition cache, now contains %d entity types", len(repo.cachedEntityTypes)))
	}

	if repo.cachedFocalOnlyEntityTypes == nil {
		repo.cachedFocalOnlyEntityTypes = []*entities.EntityType{}
	}

	if repo.cachedEntityTypes == nil {
		repo.cachedEntityTypes = []*entities.EntityType{}
	}

	return nil
}

func (repo *entityTypesRepoCached) Purge(ctx context.Context) error {
	return repo.ensureEntityTypesLoadedFromDelegate(ctx, true)
}
