package csrepo

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/d4l-data4life/mex/mex/shared/codings"
	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/rdb"
	"github.com/d4l-data4life/mex/mex/shared/utils"
	"github.com/d4l-data4life/mex/mex/shared/utils/async"
)

const CodingsetsourcesUpdateChannelNameSuffix = "mex-codingsetsources-update"

type CodingsetRepo interface {
	GetNames() []string
	GetCodingset(name string) (codings.Codingset, error)
	Purge(ctx context.Context) error
}

type codingsetRepo struct {
	mu sync.RWMutex

	log L.Logger

	originCMS           string
	codingsetsMap       map[string]async.Promise[codings.Codingset]
	strictConfigParsing bool
}

const apiPath = "/api/v0/config/files"

type NewCodingsetsRepoParams struct {
	Log                 L.Logger
	Topic               *rdb.Topic
	OriginCMS           string
	StrictConfigParsing bool
}

func NewCodingsetsRepo(ctx context.Context, params NewCodingsetsRepoParams) (CodingsetRepo, error) {
	repo := codingsetRepo{
		log:                 params.Log,
		originCMS:           params.OriginCMS,
		codingsetsMap:       make(map[string]async.Promise[codings.Codingset]),
		strictConfigParsing: params.StrictConfigParsing,
	}

	if err := repo.update(ctx); err != nil {
		return nil, err
	}

	return &repo, nil
}

func (repo *codingsetRepo) update(ctx context.Context) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.log.Info(ctx, L.Message("update codingsets"), L.Phase("codingsets"))
	repo.codingsetsMap = make(map[string]async.Promise[codings.Codingset])

	client := &http.Client{}
	resp, err := client.Get(fmt.Sprintf("%s%s/codingset_sources", repo.originCMS, apiPath))
	if err != nil {
		repo.log.Warn(ctx, L.Messagef("could not query codingsets: %s)", err.Error()))
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		repo.log.Warn(ctx, L.Messagef("could not fetch codingset sources from CMS - got response status code %d", resp.StatusCode))
		return nil
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var codingsetSources codings.CodingsetSources
	discardUnknown := !repo.strictConfigParsing
	err = protojson.UnmarshalOptions{DiscardUnknown: discardUnknown}.Unmarshal(data, &codingsetSources)
	if err != nil {
		return err
	}

	for _, codingsetSource := range codingsetSources.CodingsetSources {
		repo.log.Info(ctx, L.Messagef("- %s (%s)", codingsetSource.Name, codingsetSource.Config.TypeUrl))

		loaderFunc := GetLoader(codingsetSource.Config)
		if loaderFunc == nil {
			repo.log.Warn(ctx, L.Messagef("unknown/unsupported source: %s", codingsetSource.Config.TypeUrl))
			continue
		}

		repo.codingsetsMap[codingsetSource.Name] = loaderFunc(codingsetSource.Config)
	}

	return nil
}

func (repo *codingsetRepo) GetCodingset(name string) (codings.Codingset, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	if codingsetPromise, ok := repo.codingsetsMap[name]; ok {
		return codingsetPromise.Await()
	}
	return nil, fmt.Errorf("codingset not found: %s", name)
}

func (repo *codingsetRepo) GetNames() []string {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	return utils.KeysOfMap(repo.codingsetsMap)
}

func (repo *codingsetRepo) Purge(ctx context.Context) error {
	return repo.update(ctx)
}
