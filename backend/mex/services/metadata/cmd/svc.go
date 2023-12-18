package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/d4l-data4life/mex/mex/shared/cfg"
	"github.com/d4l-data4life/mex/mex/shared/db"
	"github.com/d4l-data4life/mex/mex/shared/entities"
	"github.com/d4l-data4life/mex/mex/shared/entities/erepo"
	"github.com/d4l-data4life/mex/mex/shared/index"
	"github.com/d4l-data4life/mex/mex/shared/interceptors"
	sharedJobs "github.com/d4l-data4life/mex/mex/shared/jobs"
	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/log/emit"
	"github.com/d4l-data4life/mex/mex/shared/mail"
	"github.com/d4l-data4life/mex/mex/shared/mail/flowmailer"
	"github.com/d4l-data4life/mex/mex/shared/svcutils"
	"github.com/d4l-data4life/mex/mex/shared/web"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields/frepo"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields/hooks"

	"github.com/d4l-data4life/mex/mex/services/metadata/endpoints/blobs"
	pbBlobs "github.com/d4l-data4life/mex/mex/services/metadata/endpoints/blobs/pb"
	"github.com/d4l-data4life/mex/mex/services/metadata/endpoints/items"
	pbItems "github.com/d4l-data4life/mex/mex/services/metadata/endpoints/items/pb"
	"github.com/d4l-data4life/mex/mex/services/metadata/endpoints/jobs"
	pbJobs "github.com/d4l-data4life/mex/mex/services/metadata/endpoints/jobs/pb"
	"github.com/d4l-data4life/mex/mex/services/metadata/endpoints/notify"
	pbNotify "github.com/d4l-data4life/mex/mex/services/metadata/endpoints/notify/pb"

	"github.com/d4l-data4life/mex/mex/services/metadata/migrations/migrate_database"
)

var (
	build     = "develop" // set during build
	buildDate = "now"     // set during build
)

const (
	serviceName = "MEx Metadata"
	serviceTag  = "metadata"

	contentTypeSendNotification = "application/json+send-notification"
)

func main() {
	ctx := context.Background()

	log, err := L.New(serviceName, build, &emit.WriterEmitter{Writer: os.Stdout})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mexConfig := cfg.MexConfig{}
	dump, err := cfg.InitConfig(log, &cfg.OSEnvs{}, "MEX", serviceTag, &mexConfig)
	if err != nil {
		log.Error(ctx, L.Messagef("could not determine config: %s", err.Error()))
		os.Exit(2)
	}
	mexConfig.Version = &cfg.Version{Build: build, Desc: serviceName, BuildDate: buildDate}
	log.Info(ctx, L.Messagef("effective config:\n%s", dump))

	setupFunc := func(ctx context.Context, opts svcutils.SetupOpts) error {
		return setup(ctx, opts, mexConfig.Strictness.StrictJsonParsing.Config)
	}

	strictJSONParsing := mexConfig.Strictness.StrictJsonParsing.Metadata
	err = svcutils.Run(ctx, svcutils.RunOpts{
		ServiceTag: serviceTag,
		Log:        log,
		Setup:      setupFunc,
		Config:     &mexConfig,

		Support: svcutils.Support{
			Postgres:       true,
			TokenValidator: true,
		},

		AdditionalServeMuxOpts: []runtime.ServeMuxOption{
			runtime.WithMarshalerOption(contentTypeSendNotification, notify.NewSendNotificationMarshaler()),
			runtime.WithMarshalerOption("*", interceptors.NewMarshaler(strictJSONParsing)),
		},

		// Change the Content-Type for the /notify endpoint so the body is parsed by the custom unmarshaler
		AdditionalHTTPHandlerWrappers: []web.HandlerWrapper{
			func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if strings.HasSuffix(r.URL.Path, "/notify") {
						r.Header.Set("Content-Type", contentTypeSendNotification)
					}
					next.ServeHTTP(w, r)
				})
			},
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

	err := db.Migrate(ctx, opts.Log, opts.DBPool, migrate_database.BindataMigrations{})
	if err != nil {
		return fmt.Errorf("database migration failed: %w", err)
	}

	jobber := sharedJobs.RedisJobber{
		Redis:      opts.Redis,
		Expiration: opts.Config.Jobs.Expiration.AsDuration(),
	}

	jobService := jobs.Service{Jobber: jobber}

	// Field hooks
	fieldDefinitionHooks, err := hooks.NewFieldDefinitionHooks(hooks.FieldDefinitionHooksConfig{
		DB: opts.DBPool,
	})
	if err != nil {
		return fmt.Errorf("failed to init field hooks: %w", err)
	}

	itemCreationHooks, err := hooks.NewItemCreationHooks(hooks.ItemCreationHooksConfig{
		DB: opts.DBPool,
	})
	if err != nil {
		return fmt.Errorf("failed to init field hooks: %w", err)
	}

	solrFieldCreationHooks, err := hooks.NewSolrFieldCreationHooks(hooks.SolrFieldCreationHooksConfig{})
	if err != nil {
		return fmt.Errorf("failed to init field hooks: %w", err)
	}

	solrDataLoadHooks, err := hooks.NewSolrDataLoadHooks(hooks.SolrDataLoadHooksConfig{
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

	announcer := index.RedisAnnouncer{
		Log:   opts.Log,
		Redis: opts.Redis,

		SetExpiration: opts.Config.AutoIndexer.SetExpiration.AsDuration(),

		TechnicalIDsTopicName: opts.Config.Redis.PubSubPrefix + "/" + items.MetadataItemUpdateByItemIDChannelName,
		BusinessIDsTopicName:  opts.Config.Redis.PubSubPrefix + "/" + items.MetadataItemUpdateByBusinessIDChannelName,
		BusinessIDSetName:     opts.Config.Redis.PubSubPrefix + "/" + items.MetadataItemUpdateByBusinessIDSetName,
	}

	metadataService := items.Service{
		ServiceTag: serviceTag,
		Log:        opts.Log,

		DB:     opts.DBPool,
		Redis:  opts.Redis,
		Jobber: jobber,

		FieldRepo:  fieldRepo,
		EntityRepo: entityRepo,

		ItemCreationHooks:      itemCreationHooks,
		SolrFieldCreationHooks: solrFieldCreationHooks,
		SolrDataLoadHooks:      solrDataLoadHooks,

		Announcer:        &announcer,
		TelemetryService: opts.TelemetryService,

		DuplicateDetectionAlgorithm: opts.Config.Indexing.DuplicationDetectionAlgorithm,
	}
	opts.TopicConfigChange.Subscribe(&metadataService)

	blobsService := blobs.Service{
		Log:             opts.Log,
		DB:              opts.DBPool,
		MasterTableName: opts.Config.Services.Blobs.MasterTableName,
	}

	var mailer mail.Mailer
	switch opts.Config.Notify.EmailerType {
	case cfg.EmailerType_MOCKMAILER:
		opts.Log.Warn(ctx, L.Message("using mockmailer"))
		mailer = flowmailer.NewMockMailer(opts.Redis)

	case cfg.EmailerType_FLOWMAILER:
		mailer, err = flowmailer.NewSimpleFlowmailer(flowmailer.Params{
			OriginOAuth:         opts.Config.Notify.Flowmailer.OriginOauth,
			OriginAPI:           opts.Config.Notify.Flowmailer.OriginApi,
			ClientID:            opts.Config.Notify.Flowmailer.ClientId,
			ClientSecret:        opts.Config.Notify.Flowmailer.ClientSecret,
			AccountID:           opts.Config.Notify.Flowmailer.AccountId,
			NoReplyEmailAddress: opts.Config.Notify.Flowmailer.NoreplyEmailAddress,
			Timeout:             opts.Config.Web.ReadTimeout.AsDuration(), // reuse
		})
		if err != nil {
			return fmt.Errorf("failed to init Flowmailer client: %w", err)
		}
	default:
		return fmt.Errorf("unsupported emailer type: %d", opts.Config.Notify.EmailerType)
	}

	notifyService := notify.Service{
		Log:   opts.Log,
		DB:    opts.DBPool,
		Redis: opts.Redis,

		Mailer: mailer,

		ConfigServiceOrigin: opts.Config.Services.Config.Origin,
	}

	pbItems.RegisterItemsServer(opts.GRPCServer, &metadataService)
	pbJobs.RegisterJobsServer(opts.GRPCServer, &jobService)
	pbBlobs.RegisterBlobsServer(opts.GRPCServer, &blobsService)
	pbNotify.RegisterNotifyServer(opts.GRPCServer, &notifyService)

	err = pbItems.RegisterItemsHandlerFromEndpoint(ctx, opts.HTTPMux, opts.Config.Web.GrpcHost, opts.GRPCOpts)
	if err != nil {
		return err
	}

	err = pbJobs.RegisterJobsHandlerFromEndpoint(ctx, opts.HTTPMux, opts.Config.Web.GrpcHost, opts.GRPCOpts)
	if err != nil {
		return err
	}

	err = pbBlobs.RegisterBlobsHandlerFromEndpoint(ctx, opts.HTTPMux, opts.Config.Web.GrpcHost, opts.GRPCOpts)
	if err != nil {
		return err
	}

	err = pbNotify.RegisterNotifyHandlerFromEndpoint(ctx, opts.HTTPMux, opts.Config.Web.GrpcHost, opts.GRPCOpts)
	if err != nil {
		return err
	}

	return nil
}
