package screpo

import (
	"context"
	"fmt"
	"sync"

	"github.com/d4l-data4life/mex/mex/shared/errstat"
	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/searchconfig"
	"github.com/d4l-data4life/mex/mex/shared/solr"
)

// searchConfigRepoCached is a search configuration repository caching the data of a delegate repository.
// If the underlying repository accesses a database, using this wrapper can speed up data fetching significantly.
type searchConfigRepoCached struct {
	mu sync.RWMutex

	log      L.Logger
	delegate searchconfig.SearchConfigRepo

	cachedSearchConfigs *searchconfig.SearchConfigList
}

type NewCachedSearchConfigRepoParams struct {
	Log      L.Logger
	Delegate searchconfig.SearchConfigRepo
}

func NewCachedSearchConfigRepo(_ context.Context, params NewCachedSearchConfigRepoParams) searchconfig.SearchConfigRepo {
	return &searchConfigRepoCached{
		log:      params.Log,
		delegate: params.Delegate,
	}
}

// ListSearchConfigs lists all config objects
func (repo *searchConfigRepoCached) ListSearchConfigs(ctx context.Context) (*searchconfig.SearchConfigList, error) {
	err := repo.ensureSearchConfigsLoadedFromDelegate(ctx, false)
	if err != nil {
		return nil, err
	}

	repo.mu.RLock()
	defer repo.mu.RUnlock()

	return repo.cachedSearchConfigs, nil
}

// ListSearchConfigsOfType lists all config objects of a certain type
func (repo *searchConfigRepoCached) ListSearchConfigsOfType(ctx context.Context, objType string) (*searchconfig.SearchConfigList, error) {
	err := repo.ensureSearchConfigsLoadedFromDelegate(ctx, false)
	if err != nil {
		return nil, err
	}

	repo.mu.RLock()
	defer repo.mu.RUnlock()

	configs := searchconfig.SearchConfigList{
		SearchConfigs: []*searchconfig.SearchConfigObject{},
	}
	for _, sc := range repo.cachedSearchConfigs.GetSearchConfigs() {
		if sc.Type == objType {
			configs.SearchConfigs = append(configs.SearchConfigs, sc)
		}
	}
	return &configs, nil
}

// GetSearchConfigObject returns one search config object by id
func (repo *searchConfigRepoCached) GetSearchConfigObject(ctx context.Context, name string) (*searchconfig.SearchConfigObject, error) {
	err := repo.ensureSearchConfigsLoadedFromDelegate(ctx, false)
	if err != nil {
		return nil, err
	}

	repo.mu.RLock()
	defer repo.mu.RUnlock()

	for _, sc := range repo.cachedSearchConfigs.GetSearchConfigs() {
		if sc.Name == name {
			return sc, nil
		}
	}
	return nil, errstat.MakeMexStatus(errstat.ElementNotFound, "search config not found").Err()
}

// GetFieldsForSearchFocus return the requested search focus (error if not found)
func (repo *searchConfigRepoCached) GetFieldsForSearchFocus(ctx context.Context, requestedSearchFocus string) ([]string, error) {
	err := repo.ensureSearchConfigsLoadedFromDelegate(ctx, false)
	if err != nil {
		return nil, err
	}

	repo.mu.RLock()
	defer repo.mu.RUnlock()

	for _, sc := range repo.cachedSearchConfigs.GetSearchConfigs() {
		if sc.Type == solr.MexSearchFocusType && sc.Name == requestedSearchFocus {
			return sc.Fields, nil
		}
	}
	return nil, fmt.Errorf("could not find requested search focus")
}

// GetFieldsForAxis returns fields for an axis (ordinal or hierarchy)
func (repo *searchConfigRepoCached) GetFieldsForAxis(ctx context.Context, requestedAxis string) ([]string, error) {
	err := repo.ensureSearchConfigsLoadedFromDelegate(ctx, false)
	if err != nil {
		return nil, err
	}

	repo.mu.RLock()
	defer repo.mu.RUnlock()

	for _, sc := range repo.cachedSearchConfigs.GetSearchConfigs() {
		if (sc.Type == solr.MexOrdinalAxisType || sc.Type == solr.MexHierarchyAxisType) && sc.Name == requestedAxis {
			return sc.Fields, nil
		}
	}
	return nil, fmt.Errorf("could not find requested ordinal axis")
}

// ensureSearchConfigsLoadedFromDelegate refills the field configuration cache (if needed) from another fields repo
func (repo *searchConfigRepoCached) ensureSearchConfigsLoadedFromDelegate(ctx context.Context, purge bool) error {
	var err error

	if purge || repo.cachedSearchConfigs == nil {
		repo.mu.Lock()
		defer repo.mu.Unlock()

		repo.cachedSearchConfigs, err = repo.delegate.ListSearchConfigs(ctx)
		if err != nil {
			return err
		}
		if repo.cachedSearchConfigs != nil {
			repo.log.Info(ctx, L.Messagef("refreshed search config cache, now contains %d search configs", len(repo.cachedSearchConfigs.SearchConfigs)))
		}
	}

	if repo.cachedSearchConfigs == nil {
		panic("list still nil")
	}

	return nil
}

func (repo *searchConfigRepoCached) Purge(ctx context.Context) error {
	return repo.ensureSearchConfigsLoadedFromDelegate(ctx, true)
}
