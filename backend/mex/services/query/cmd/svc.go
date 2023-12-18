package main

import (
	"context"
	"fmt"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/d4l-data4life/mex/mex/shared/cfg"
	"github.com/d4l-data4life/mex/mex/shared/interceptors"
	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/log/emit"
	"github.com/d4l-data4life/mex/mex/shared/searchconfig"
	"github.com/d4l-data4life/mex/mex/shared/searchconfig/screpo"
	"github.com/d4l-data4life/mex/mex/shared/svcutils"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields/frepo"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields/hooks"

	"github.com/d4l-data4life/mex/mex/services/query/endpoints/search"
	"github.com/d4l-data4life/mex/mex/services/query/endpoints/search/pb"
)

var (
	build     = "develop" // set during build
	buildDate = "now"     // set during build
)

const (
	serviceName = "MEx Query"
	serviceTag  = "query"
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
			runtime.WithMarshalerOption("*", interceptors.NewMarshaler(config.Strictness.StrictJsonParsing.Query)),
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

	// Field hooks
	fieldDefinitionHooks, err := hooks.NewFieldDefinitionHooks(hooks.FieldDefinitionHooksConfig{
		DB: opts.DBPool,
	})
	if err != nil {
		return fmt.Errorf("failed to init field hooks: %w", err)
	}

	postQueryHooks, err := hooks.NewPostQueryHooks(hooks.PostQueryHooksConfig{
		DB: opts.DBPool,
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
	searchService := search.Service{
		ServiceTag: serviceTag,
		Log:        opts.Log,

		Redis:                 opts.Redis,
		Solr:                  opts.Solr,
		SolrCollection:        opts.Config.Solr.Collection,
		TolerantErrorHandling: opts.Config.Strictness.Search.ToleratePartialFailures,

		FieldRepo:        fieldRepo,
		SearchConfigRepo: searchConfigRepo,

		TelemetryService: opts.TelemetryService,

		PostQueryHooks: postQueryHooks,
	}
	opts.TopicConfigChange.Subscribe(&searchService)

	pb.RegisterSearchServer(opts.GRPCServer, &searchService)

	err = pb.RegisterSearchHandlerFromEndpoint(ctx, opts.HTTPMux, opts.Config.Web.GrpcHost, opts.GRPCOpts)
	if err != nil {
		return err
	}

	return nil
}
