package frepo

import (
	"context"
	"fmt"
	"sync"

	L "github.com/d4l-data4life/mex/mex/shared/log"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
)

type fieldRepoCached struct {
	mu sync.RWMutex

	log      L.Logger
	delegate fields.FieldRepo

	cachedFieldDefs []fields.BaseFieldDef
}

type NewCachedFieldRepoParams struct {
	Log      L.Logger
	Delegate fields.FieldRepo
}

func NewCachedFieldRepo(_ context.Context, params NewCachedFieldRepoParams) fields.FieldRepo {
	return &fieldRepoCached{
		log:      params.Log,
		delegate: params.Delegate,
	}
}

// GetFieldDefNames returns the names of all defined fields, including linked fields
func (repo *fieldRepoCached) GetFieldDefNames(ctx context.Context) ([]string, error) {
	err := repo.ensureFieldDefsLoadedFromDelegate(ctx, false)
	if err != nil {
		return nil, err
	}

	repo.mu.RLock()
	defer repo.mu.RUnlock()

	fieldNames := make([]string, len(repo.cachedFieldDefs))
	for i, fieldDef := range repo.cachedFieldDefs {
		fieldNames[i] = fieldDef.Name()
	}

	return fieldNames, nil
}

func (repo *fieldRepoCached) GetFieldDefByName(ctx context.Context, fieldName string) (fields.BaseFieldDef, error) {
	err := repo.ensureFieldDefsLoadedFromDelegate(ctx, false)
	if err != nil {
		return nil, err
	}

	repo.mu.RLock()
	defer repo.mu.RUnlock()

	for _, fieldDef := range repo.cachedFieldDefs {
		if fieldDef.Name() == fieldName {
			return fieldDef, nil
		}
	}

	return nil, fmt.Errorf("no field definition found with name: %s", fieldName)
}

func (repo *fieldRepoCached) GetFieldDefsByKind(ctx context.Context, fieldKind string) ([]fields.BaseFieldDef, error) {
	err := repo.ensureFieldDefsLoadedFromDelegate(ctx, false)
	if err != nil {
		return nil, err
	}

	repo.mu.RLock()
	defer repo.mu.RUnlock()

	var fieldDefs []fields.BaseFieldDef
	for _, fieldDef := range repo.cachedFieldDefs {
		if fieldDef.Kind() == fieldKind {
			fieldDefs = append(fieldDefs, fieldDef)
		}
	}

	return fieldDefs, nil
}

func (repo *fieldRepoCached) ListFieldDefs(ctx context.Context) ([]fields.BaseFieldDef, error) {
	err := repo.ensureFieldDefsLoadedFromDelegate(ctx, false)
	if err != nil {
		return nil, err
	}

	repo.mu.RLock()
	defer repo.mu.RUnlock()

	return repo.cachedFieldDefs, nil
}

// ensureFieldDefsLoadedFromDelegate loads field configurations from another fields repo.
func (repo *fieldRepoCached) ensureFieldDefsLoadedFromDelegate(ctx context.Context, purge bool) error {
	var err error

	if purge || repo.cachedFieldDefs == nil {
		repo.mu.Lock()
		defer repo.mu.Unlock()

		repo.cachedFieldDefs, err = repo.delegate.ListFieldDefs(ctx)
		if err != nil {
			return err
		}
		repo.log.Info(ctx, L.Messagef("refreshed field definition cache, now contains %d fields", len(repo.cachedFieldDefs)))
	}

	if repo.cachedFieldDefs == nil {
		panic("list still nil")
	}

	return nil
}

func (repo *fieldRepoCached) Purge(ctx context.Context) error {
	return repo.ensureFieldDefsLoadedFromDelegate(ctx, true)
}
