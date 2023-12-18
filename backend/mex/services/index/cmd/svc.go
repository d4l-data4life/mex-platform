package main

import (
	"context"
	"fmt"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/d4l-data4life/mex/mex/shared/blobs/pglo"
	"github.com/d4l-data4life/mex/mex/shared/cfg"
	"github.com/d4l-data4life/mex/mex/shared/codings"
	"github.com/d4l-data4life/mex/mex/shared/codings/csrepo"
	"github.com/d4l-data4life/mex/mex/shared/codings/mesh"
	"github.com/d4l-data4life/mex/mex/shared/entities"
	"github.com/d4l-data4life/mex/mex/shared/entities/erepo"
	"github.com/d4l-data4life/mex/mex/shared/interceptors"
	sharedJobs "github.com/d4l-data4life/mex/mex/shared/jobs"
	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/log/emit"
	"github.com/d4l-data4life/mex/mex/shared/searchconfig"
	"github.com/d4l-data4life/mex/mex/shared/searchconfig/screpo"
	"github.com/d4l-data4life/mex/mex/shared/solr"
	"github.com/d4l-data4life/mex/mex/shared/svcutils"
	"github.com/d4l-data4life/mex/mex/shared/utils"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields/frepo"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields/hooks"
	"github.com/d4l-data4life/mex/mex/services/metadata/endpoints/items"
	"github.com/d4l-data4life/mex/mex/services/metadata/endpoints/jobs"
	"github.com/d4l-data4life/mex/mex/services/metadata/migrations/solr_configset"

	"github.com/d4l-data4life/mex/mex/services/index/endpoints/index"
	"github.com/d4l-data4life/mex/mex/services/index/endpoints/index/pb"
	"github.com/d4l-data4life/mex/mex/services/index/indexer"
)

var (
	build     = "develop" // set during build
	buildDate = "now"     // set during build
)

const (
	serviceName = "MEx Index"
	serviceTag  = "index"
)

func main() {
	ctx := context.Background()

	log, err := L.New(serviceName, build, &emit.WriterEmitter{Writer: os.Stdout})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	config := cfg.MexConfig{}
	dump, err := cfg.InitConfig(log, &cfg.OSEnvs{}, "MEX", serviceTag, &config)
	if err != nil {
		log.Error(ctx, L.Messagef("could not determine config: %s", err.Error()))
		os.Exit(2)
	}
	config.Version = &cfg.Version{Build: build, Desc: serviceName, BuildDate: buildDate}
	log.Info(ctx, L.Messagef("effective config:\n%s", dump))

	setupFunc := func(ctx context.Context, opts svcutils.SetupOpts) error {
		return setup(ctx, opts, config.Strictness.StrictJsonParsing.Config)
	}

	err = svcutils.Run(ctx, svcutils.RunOpts{
		ServiceTag: serviceTag,
		Log:        log,
		Setup:      setupFunc,
		Config:     &config,

		Support: svcutils.Support{
			Postgres:       true,
			Solr:           true,
			TokenValidator: true,
		},

		AdditionalServeMuxOpts: []runtime.ServeMuxOption{
			runtime.WithMarshalerOption("*", interceptors.NewMarshaler(config.Strictness.StrictJsonParsing.Index)),
		},
	})
	if err != nil {
		log.Error(ctx, L.Message(err.Error()))
		os.Exit(1)
	}
}

func setup(ctx context.Context, opts svcutils.SetupOpts, strictConfigParsing bool) error {
	if cfg.StringIsEmpty(opts.Config.Services.Config.Origin) {
		return fmt.Errorf("no origin for config service configured")
	}

	err := solr_configset.Bootstrap(ctx, opts.Log, opts.Solr, opts.Config.Solr.ConfigsetName)
	if err != nil {
		return err
	}

	collections, err := opts.Solr.GetCollections(ctx)
	if err != nil {
		return err
	}
	opts.Log.Info(ctx, L.Messagef("existing Solr collections: %v", collections))
	if !utils.Contains(collections, opts.Config.Solr.Collection) {
		configSet := solr.DefaultSolrConfigSet
		opts.Log.Info(ctx, L.Messagef("creating Solr collection '%s' from configset '%s'", opts.Config.Solr.Collection, configSet))
		err = opts.Solr.CreateCollection(ctx, opts.Config.Solr.Collection, configSet, opts.Config.Solr.ReplicationFactor)
		if err != nil {
			return err
		}
		opts.Log.Info(ctx, L.Messagef("Solr collection created: %s", opts.Config.Solr.Collection))
	}

	jobService := jobs.Service{
		Jobber: sharedJobs.RedisJobber{
			Redis:      opts.Redis,
			Expiration: opts.Config.Jobs.Expiration.AsDuration(),
		},
	}

	// Field hooks
	fieldDefinitionHooks, err := hooks.NewFieldDefinitionHooks(hooks.FieldDefinitionHooksConfig{
		DB: opts.DBPool,
	})
	if err != nil {
		return fmt.Errorf("failed to init field hooks: %w", err)
	}

	blobStore := pglo.PostgresLargeObjectStore{
		DB:              opts.DBPool,
		MasterTableName: opts.Config.Services.Blobs.MasterTableName,
		Log:             opts.Log,
	}
	csrepo.InstallLoader(&codings.BlobStoreCodingsetSourceConfig{}, mesh.NewBlobStoreLoader(&blobStore))

	codingsetRepo, err := csrepo.NewCodingsetsRepo(ctx, csrepo.NewCodingsetsRepoParams{
		Log:                 opts.Log,
		Topic:               opts.TopicConfigChange,
		OriginCMS:           opts.Config.Services.Config.Origin,
		StrictConfigParsing: strictConfigParsing,
	})
	if err != nil {
		return err
	}
	opts.Log.Info(ctx, L.Messagef("coding stores: %v", codingsetRepo.GetNames()))

	solrFieldCreationHooks, err := hooks.NewSolrFieldCreationHooks(hooks.SolrFieldCreationHooksConfig{})
	if err != nil {
		return fmt.Errorf("failed to init field hooks: %w", err)
	}

	solrDataLoadHooks, err := hooks.NewSolrDataLoadHooks(hooks.SolrDataLoadHooksConfig{
		DB:            opts.DBPool,
		CodingsetRepo: codingsetRepo,
	})
	if err != nil {
		return fmt.Errorf("failed to init field hooks: %w", err)
	}

	var fieldRepo fields.FieldRepo
	switch opts.Config.FieldDefs.RepoType {
	case cfg.RepoType_DIRECT:
		fieldRepo = frepo.NewDirectCMSFieldDefsRepo(opts.Config.Services.Config.Origin, fieldDefinitionHooks, strictConfigParsing)
	case cfg.RepoType_CACHED:
		fieldRepo = frepo.NewCachedFieldRepo(ctx, frepo.NewCachedFieldRepoParams{
			Log:      opts.Log,
			Delegate: frepo.NewDirectCMSFieldDefsRepo(opts.Config.Services.Config.Origin, fieldDefinitionHooks, strictConfigParsing),
		})
	default:
		return fmt.Errorf("unknown fields repo type: %s", opts.Config.FieldDefs.RepoType)
	}

	var entityRepo entities.EntityRepo
	switch opts.Config.EntityTypes.RepoType {
	case cfg.RepoType_DIRECT:
		entityRepo = erepo.NewDirectCMSEntityTypesRepo(opts.Config.Services.Config.Origin, strictConfigParsing)
	case cfg.RepoType_CACHED:
		entityRepo = erepo.NewCachedEntityTypesRepo(ctx, erepo.NewCachedEntityTypesRepoParams{
			Log:      opts.Log,
			Delegate: erepo.NewDirectCMSEntityTypesRepo(opts.Config.Services.Config.Origin, strictConfigParsing),
		})
	}

	var searchConfigRepo searchconfig.SearchConfigRepo
	switch opts.Config.SearchConfig.RepoType {
	case cfg.RepoType_DIRECT:
		searchConfigRepo = screpo.NewDirectCMSSearchConfigRepo(opts.Config.Services.Config.Origin, strictConfigParsing)
	case cfg.RepoType_CACHED:
		searchConfigRepo = screpo.NewCachedSearchConfigRepo(ctx,
			screpo.NewCachedSearchConfigRepoParams{
				Log:      opts.Log,
				Delegate: screpo.NewDirectCMSSearchConfigRepo(opts.Config.Services.Config.Origin, strictConfigParsing),
			})
	default:
		return fmt.Errorf("unknown search config repo type: %s", opts.Config.SearchConfig.RepoType)
	}

	indexService := index.Service{
		ServiceTag: serviceTag,
		Log:        opts.Log,

		DB:             opts.DBPool,
		Redis:          opts.Redis,
		Solr:           opts.Solr,
		SolrCollection: opts.Config.Solr.Collection,

		JobService: &jobService,

		FieldRepo:        fieldRepo,
		EntityRepo:       entityRepo,
		SearchConfigRepo: searchConfigRepo,

		SolrFieldCreationHooks: solrFieldCreationHooks,
		SolrDataLoadHooks:      solrDataLoadHooks,
		CodingsetRepo:          codingsetRepo,

		TelemetryService: opts.TelemetryService,

		CollectionName:    opts.Config.Solr.Collection,
		ReplicationFactor: opts.Config.Solr.ReplicationFactor,
	}
	opts.TopicConfigChange.Subscribe(&indexService)

	autoIndexer := indexer.NewAutoIndexer(ctx, indexer.AutoIndexerConfig{
		Log:           opts.Log,
		Redis:         opts.Redis,
		IndexService:  &indexService,
		SetExpiration: opts.Config.AutoIndexer.SetExpiration.AsDuration(),

		TechnicalIDsTopicName: opts.Config.Redis.PubSubPrefix + "/" + items.MetadataItemUpdateByItemIDChannelName,
		BusinessIDsTopicName:  opts.Config.Redis.PubSubPrefix + "/" + items.MetadataItemUpdateByBusinessIDChannelName,
		BusinessIDSetName:     opts.Config.Redis.PubSubPrefix + "/" + items.MetadataItemUpdateByBusinessIDSetName,
	})
	autoIndexer.StartPeriodicIndexer()

	pb.RegisterIndexServer(opts.GRPCServer, &indexService)

	err = pb.RegisterIndexHandlerFromEndpoint(ctx, opts.HTTPMux, opts.Config.Web.GrpcHost, opts.GRPCOpts)
	if err != nil {
		return err
	}

	return nil
}
