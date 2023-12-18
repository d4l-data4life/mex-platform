package svcutils

import (
	"context"
	"crypto/tls"
	"expvar"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"

	"github.com/d4l-data4life/mex/mex/shared/auth"
	"github.com/d4l-data4life/mex/mex/shared/auth/authn"
	"github.com/d4l-data4life/mex/mex/shared/certs"
	"github.com/d4l-data4life/mex/mex/shared/cfg"
	"github.com/d4l-data4life/mex/mex/shared/constants"
	"github.com/d4l-data4life/mex/mex/shared/db"
	"github.com/d4l-data4life/mex/mex/shared/errstat/errhandler"
	"github.com/d4l-data4life/mex/mex/shared/interceptors"
	"github.com/d4l-data4life/mex/mex/shared/keys"
	"github.com/d4l-data4life/mex/mex/shared/known/securitypb"
	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/log/emit"
	"github.com/d4l-data4life/mex/mex/shared/log/grpc_adapt"
	"github.com/d4l-data4life/mex/mex/shared/rdb"
	"github.com/d4l-data4life/mex/mex/shared/solr"
	"github.com/d4l-data4life/mex/mex/shared/telemetry"
	"github.com/d4l-data4life/mex/mex/shared/telemetry/pingers"
	"github.com/d4l-data4life/mex/mex/shared/try"
	"github.com/d4l-data4life/mex/mex/shared/web"
)

// This is the common service run function of all services.
// The business logic is implemented as gRPC endpoints.
// The definitions are contained in respective proto files.
//
// After a successful start, the process listens on two ports:
//
//   - [A] The gRPC server is running on port 9000 by default.
//
//   - [B] We also generate an HTTP REST gateway from the proto files.
//     It maps annotated gRPC endpoints to REST endpoints.
//     It is just a relay, all logic (including authentication)
//     is done in the gRPC server/endpoints.
//
//           O H2               O HTTP 1.1
//           | (:9000)          | (:3000)
//           | [A]              | [B]
//           |                  |
//     +-----|------------------|----------------------------+
//     |     |                  |                            |
//     |     |        +---------V---------+                  |
//     |     |        | HTTP REST gateway |                  |
//     |     |        | (generated)       |                  |
//     |     |        +---------+---------+                  |
//     |     |                  |                            |
//     |     |                  |                            |
//     |  +--V------------------V---------+                  |
//     |  |                               |     .---------.  |
//     |  | gRPC service endpoints        <--->( Databases ) |
//     |  |                               |     `---------'  |
//     |  +-------------------------------+                  |
//     |                                                     |
//     +-----------------------------------------------------+

const PromMexNamespace = "mex"

type SetupOpts struct {
	Config            *cfg.MexConfig
	Log               L.Logger
	Redis             *redis.Client
	DBPool            *pgxpool.Pool
	Solr              solr.ClientAPI
	GRPCServer        *grpc.Server
	GRPCOpts          []grpc.DialOption
	TokenValidator    keys.TokenValidator
	HTTPMux           *runtime.ServeMux
	TelemetryService  *telemetry.Service
	TopicConfigChange *rdb.Topic
	PromRegistry      *prometheus.Registry
}

type SetupFunc = func(ctx context.Context, opts SetupOpts) error

type Support struct {
	Postgres       bool
	Solr           bool
	TokenValidator bool
}

type RunOpts struct {
	ServiceTag string

	Log   *L.MexLogger
	Setup SetupFunc

	Config *cfg.MexConfig

	Support Support

	AdditionalServeMuxOpts                   []runtime.ServeMuxOption
	AdditionalHTTPHandlerWrappers            []web.HandlerWrapper
	AdditionalTokenValidationExcludePatterns []string
}

type InitParams struct {
	EnvPrefix   string
	EnvFileName string
	PrintHelp   bool
	Log         L.Logger
	Config      interface{}
}

func Run(ctx context.Context, opts RunOpts) error {
	var err error
	ctx, cancel := context.WithCancel(ctx)
	// We do not defer cancel() here as it would be idiomatic otherwise.
	// If we did, then the cancel call would be the last deferred call of this function.
	// However, we use the cancellation signal to e.g. unsubscribe cleanly from Redis pub/sub topics.
	// The Redis connection is also closed using a defer call below.
	// Since it comes below, it would be called before the cancel call, thus closing
	// the Redis connection before de-subscriptions can happen.

	// ------------------------------------------------------------------------------------------------------
	opts.Log.SetRedact(opts.Config.Logging.RedactPersonalFields)
	opts.Log.SetQueryRedactions(opts.Config.Logging.RedactQueryParams)

	if cfg.StringIsEmpty(opts.Config.Services.BiEventsFilter.Origin) || cfg.StringIsEmpty(opts.Config.Services.BiEventsFilter.Secret) {
		cancel()
		return fmt.Errorf("misconfigured BI events filter")
	}

	efEmitter, efDone := emit.NewEventsFilterEmitter(ctx, opts.Log,
		opts.Config.Services.BiEventsFilter.Origin,
		opts.Config.Services.BiEventsFilter.Path,
		opts.Config.Services.BiEventsFilter.Secret,
	)
	opts.Log.AddEmitter(efEmitter)

	// ------------------------------------------------------------------------------------------------------
	// Service start
	opts.Log.Info(ctx, L.Messagef("starting service: [%s / %s / %s]",
		opts.Config.Version.Desc, opts.Config.Version.Build, opts.Config.Version.BuildDate), L.PhaseStartup)

	expvar.NewString("build").Set(opts.Config.Version.Build)

	// ------------------------------------------------------------------------------------------------------
	// Load (optional) additional CA certificates
	rootCAs, err := certs.LoadCerts(ctx, certs.LoadParams{
		AdditionalCACertsFiles: opts.Config.Web.CaCerts.AdditionalCaCertsFiles,
		AdditionalCACertsPEM:   opts.Config.Web.CaCerts.AdditionalCaCertsPem,
		ConnectionAttempts:     opts.Config.Web.CaCerts.AccessAttempts,
		ConnectionPause:        opts.Config.Web.CaCerts.GetAccessPause().AsDuration(),
		Log:                    opts.Log,
	})
	if err != nil {
		cancel()
		return err
	}

	// ------------------------------------------------------------------------------------------------------
	// Set up base metrics
	promRegistry := prometheus.NewRegistry()

	m0 := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: PromMexNamespace,
		Name:      "build_info",
		Help:      "MEx build information exposed as labels of a dummy gauge",
		ConstLabels: map[string]string{
			"buildVersion": opts.Config.Version.Build,
			"buildDate":    opts.Config.Version.BuildDate,
			"service":      opts.Config.Version.Desc,
		},
	})
	m0.Set(1)

	m1 := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: PromMexNamespace,
		Name:      "http_status_codes",
		Help:      "counter of HTTP status codes in case of errors",
	}, []string{"service", "code"})

	statusCounter := func(status int) {
		m1.WithLabelValues(opts.ServiceTag, fmt.Sprintf("%03d", status)).Inc()
	}

	promRegistry.MustRegister(m0)
	promRegistry.MustRegister(m1)

	// ------------------------------------------------------------------------------------------------------
	// Redis connection
	var redisClient *redis.Client
	var topicConfigChange *rdb.Topic
	opts.Log.Info(ctx, L.Message("--[ Redis ]------------------------------------------------------"), L.PhaseStartup)
	opts.Log.Info(ctx, L.Messagef("init: Redis support, host: %s:%d", opts.Config.Redis.Hostname, opts.Config.Redis.Port), L.PhaseStartup)

	redisSupport := &rdb.InitRedisSupport{
		RootCAs:  rootCAs,
		Hostname: opts.Config.Redis.Hostname,
		Port:     opts.Config.Redis.Port,
		Password: opts.Config.Redis.Password,
		DB:       opts.Config.Redis.Db,
		UseTLS:   opts.Config.Redis.UseTls,

		RetryStrategy: try.NewMaxAttemptsConstantPauseStrategy(opts.Config.Redis.ConnectionAttempts, opts.Config.Redis.ConnectionPause.AsDuration()),
	}

	redisClient, err = redisSupport.Init(ctx, opts.Log)
	if err != nil {
		cancel()
		return fmt.Errorf("connecting to Redis: %w", err)
	}

	if cfg.StringIsEmpty(opts.Config.Redis.PubSubPrefix) {
		cancel()
		return fmt.Errorf("config Redis pub/sub prefix is empty")
	}
	topicConfigChange = rdb.NewTopic(ctx, opts.Log, redisClient, fmt.Sprintf("%s/%s", opts.Config.Redis.PubSubPrefix, constants.ConfigUpdateChannelNameSuffix))

	defer func() {
		opts.Log.Info(ctx, L.Message("stopping: Redis support"), L.PhaseShutdown)

		// Give the Redis pub/sub listeners time to unsubscribe before closing the Redis connection.
		time.Sleep(opts.Config.Redis.ShutdownGracePeriod.AsDuration())
		redisClient.Close()
		opts.Log.Info(ctx, L.Message("stopped: Redis support"), L.PhaseShutdown)
	}()

	opts.Log.AddEmitter(&emit.RedisEmitter{Client: redisClient})

	// Telemetry service
	telemetryService := telemetry.New(opts.Log, opts.ServiceTag, redisClient, opts.Config.Telemetry.StatusUpdateInterval.AsDuration(), ctx.Done())
	telemetryService.AddPinger(pingers.NewRedisPinger(redisClient, opts.Config.Telemetry.PingerUpdateInterval.AsDuration()))

	// ------------------------------------------------------------------------------------------------------
	// Solr Cloud setup
	var solrClient solr.ClientAPI
	if opts.Support.Solr {
		opts.Log.Info(ctx, L.Message("--[ Solr ]------------------------------------------------------"), L.PhaseStartup)
		opts.Log.Info(ctx, L.Messagef("init: Solr support, hostname: %s", opts.Config.Solr.Origin), L.PhaseStartup)

		solrClient = solr.NewClient(opts.Config.Solr.Origin, opts.Config.Solr.Collection,
			solr.WithLogger(opts.Log),
			solr.WithCertificates(rootCAs),
			solr.WithBatchSize(opts.Config.Solr.IndexBatchSize),
			solr.WithCommitWithin(opts.Config.Solr.CommitWithin.AsDuration()),
			// No-op when not configured
			solr.WithBasicAuth(opts.Config.Solr.BasicauthUser, opts.Config.Solr.BasicauthPassword),
		)
		_, err = try.Try(ctx, try.Task[struct{}]{
			Desc:          "Solr connection",
			Phase:         "startup",
			PauseStrategy: try.NewMaxAttemptsConstantPauseStrategy(opts.Config.Solr.ConnectionAttempts, opts.Config.Solr.ConnectionPause.AsDuration()),
			Log:           opts.Log,
			Func: func() (struct{}, error) {
				return struct{}{}, solrClient.Ping(ctx)
			},
		})
		if err != nil {
			cancel()
			return fmt.Errorf("connection to Solr failed: %w", err)
		}

		telemetryService.AddPinger(pingers.NewSolrPinger(solrClient, opts.Config.Telemetry.PingerUpdateInterval.AsDuration()))
	}

	// ------------------------------------------------------------------------------------------------------
	// Relational database connection
	var pgClient *pgxpool.Pool
	if opts.Support.Postgres {
		opts.Log.Info(ctx, L.Message("--[ Postgres ]------------------------------------------------------"), L.PhaseStartup)
		opts.Log.Info(ctx, L.Messagef("init: database support, hostname: %s", opts.Config.Db.Hostname), L.PhaseStartup)

		dbSupport := &db.InitDatabaseSupport{
			RootCAs:       rootCAs,
			User:          opts.Config.Db.User,
			Password:      opts.Config.Db.Password,
			Hostname:      opts.Config.Db.Hostname,
			Port:          opts.Config.Db.Port,
			Name:          opts.Config.Db.Name,
			SearchPath:    opts.Config.Db.SearchPath,
			SSLMode:       opts.Config.Db.SslMode,
			RetryStrategy: try.NewMaxAttemptsConstantPauseStrategy(opts.Config.Db.ConnectionAttempts, opts.Config.Db.ConnectionPause.AsDuration()),
		}

		pgClient, err = dbSupport.Init(ctx, opts.Log)
		if err != nil {
			cancel()
			return fmt.Errorf("database connection failed: %w", err)
		}

		defer func() {
			opts.Log.Info(ctx, L.Message("stopping: database support"), L.PhaseShutdown)
			pgClient.Close()
			opts.Log.Info(ctx, L.Message("stopped: database support"), L.PhaseShutdown)
		}()

		postgresPinger := pingers.NewPostgresPinger(pgClient, opts.Config.Telemetry.PingerUpdateInterval.AsDuration())
		postgresPinger.RegisterMetrics(PromMexNamespace, promRegistry)
		telemetryService.AddPinger(postgresPinger)
	}

	// ------------------------------------------------------------------------------------------------------
	// Authnz
	apiKeyAuthn, err := authn.NewAPIKeyAuthenticator(opts.Config.TenantId, opts.Config.Auth.ApiKeysRoles)
	if err != nil {
		cancel()
		return err
	}

	opts.Log.Info(ctx, L.Messagef("configured API keys: %d", apiKeyAuthn.Count()))
	if apiKeyAuthn.Count() == 0 {
		opts.Log.Info(ctx, L.Message("no API keys configured; service will not accept requests with API keys"))
	}

	authnRegistry := auth.RequestAuthenticatorRegistry{
		securitypb.AuthenticationType_NONE:    authn.NewNoneAuthenticator(opts.Config.TenantId),
		securitypb.AuthenticationType_API_KEY: apiKeyAuthn,
	}

	// ------------------------------------------------------------------------------------------------------
	// Key stores
	var tokenValidator keys.TokenValidator
	if opts.Support.TokenValidator {
		opts.Log.Info(ctx, L.Message("--[ Key store ]------------------------------------------------------"), L.PhaseStartup)
		opts.Log.Info(ctx, L.Message("init: key store"), L.PhaseStartup)

		//nolint:gosec // gosec complains about this issue: https://github.com/go-redis/redis/issues/1553
		keystoreTLSConfig := &tls.Config{
			InsecureSkipVerify: false,
			RootCAs:            rootCAs,
		}

		tokenValidator, err = try.Try(context.Background(), try.Task[keys.TokenValidator]{
			Desc:          "JWKS query",
			Phase:         "startup",
			Log:           opts.Log,
			PauseStrategy: try.NewMaxAttemptsConstantPauseStrategy(opts.Config.Jwks.ConnectionAttempts, opts.Config.Jwks.ConnectionPause.AsDuration()),
			Func: func() (keys.TokenValidator, error) {
				return keys.NewRemoteKeyStore(context.Background(), keys.RemoteKeyStoreOptions{
					Log:                         opts.Log,
					URI:                         opts.Config.Jwks.RemoteKeysUri,
					TLSConfig:                   keystoreTLSConfig,
					InternalAuthServiceHostname: opts.Config.Oauth.InternalAuthServiceHostname,
				})
			},
		})
		if err != nil {
			cancel()
			return err
		}

		authJWT, err := authn.NewJWTAuthenticator(
			tokenValidator,
			opts.Config.Oauth.ClientId,
			opts.Config.Oauth.ConsumerGroupId,
			opts.Config.Oauth.ProducerGroupId,
		)
		if err != nil {
			cancel()
			return err
		}
		authnRegistry[securitypb.AuthenticationType_BEARER_TOKEN] = authJWT

		telemetryService.AddPinger(pingers.NewKeystorePinger(tokenValidator, opts.Config.Telemetry.PingerUpdateInterval.AsDuration()))
	}

	opts.Log.Info(ctx, L.Message("--[ Services ]------------------------------------------------------"), L.PhaseStartup)

	// ------------------------------------------------------------------------------------------------------
	// Shutdown channels

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// ------------------------------------------------------------------------------------------------------
	// Start gRPC service
	grpclog.SetLoggerV2(grpc_adapt.GRPCLogger{Logger: opts.Log, Level: grpc_adapt.ParseLevel(opts.Config.Logging.LogLevelGrpc)})

	ints := []grpc.UnaryServerInterceptor{
		interceptors.NewHeaderInterceptor(string(constants.ContextKeyTraceThis), constants.ContextKeyTraceThis),
		interceptors.NewHeaderInterceptor(string(constants.ContextKeyTraceID), constants.ContextKeyTraceID),
		interceptors.NewHeaderInterceptor(string(constants.ContextKeyJobID), constants.ContextKeyJobID),
		interceptors.NewHeaderInterceptor(string(constants.ContextKeyRequestID), constants.ContextKeyRequestID),
		L.NewLogInterceptor(opts.Log, []string{
			"/d4l.mex.telemetry.Telemetry/ReadinessProbe",
			"/d4l.mex.telemetry.Telemetry/LivenessProbe",
		}),
		auth.NewInterceptor(authnRegistry, auth.NewPrivMgr()),
	}

	grpcServer := grpc.NewServer(grpc.MaxRecvMsgSize(int(opts.Config.Web.MaxBodyBytes)), grpc.ChainUnaryInterceptor(ints...))
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", opts.Config.Web.GrpcHost)
	if err != nil {
		cancel()
		return fmt.Errorf("failed to listen: %w", err)
	}

	// ------------------------------------------------------------------------------------------------------
	// REST Gateway
	opts.Log.Info(ctx, L.Message("initializing REST gateway support"), L.PhaseStartup)

	// Create the mux
	muxOptions := []runtime.ServeMuxOption{
		runtime.WithErrorHandler(errhandler.CustomErrorHandler(opts.Log, statusCounter)),
		runtime.WithMetadata(auth.AuthnHeaderReader),
		runtime.WithMetadata(interceptors.HeaderReader(string(constants.ContextKeyTraceSecret), constants.HTTPHeaderTraceSecret)),
		runtime.WithMetadata(interceptors.HeaderReader(string(constants.ContextKeyTraceID), constants.HTTPHeaderTraceID)),
		runtime.WithMetadata(interceptors.HeaderReader(string(constants.ContextKeyJobID), constants.HTTPHeaderJobID)),
		runtime.WithMetadata(interceptors.ContextConnect(string(constants.ContextKeyRequestID), constants.ContextKeyRequestID)),
		runtime.WithMetadata(interceptors.ContextConnect(string(constants.ContextKeyTraceThis), constants.ContextKeyTraceThis)),
		runtime.WithForwardResponseOption(interceptors.RemoveResponseHeaders([]string{"Grpc-Metadata-Content-Type"})),
		runtime.WithForwardResponseOption(interceptors.RewriteHintedHeaders()),
		runtime.WithForwardResponseOption(interceptors.SetResponseHeader("X-Content-Type-Options", "nosniff")),
		runtime.WithForwardResponseOption(interceptors.RewriteStatusCode()),
	}
	muxOptions = append(muxOptions, opts.AdditionalServeMuxOpts...)
	mux := runtime.NewServeMux(muxOptions...)

	grpcOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	telemetry.RegisterTelemetryServer(grpcServer, telemetryService)

	err = telemetry.RegisterTelemetryHandlerFromEndpoint(ctx, mux, opts.Config.Web.GrpcHost, grpcOpts)
	if err != nil {
		cancel()
		return err
	}

	// Service-specific setup code
	err = opts.Setup(ctx, SetupOpts{
		Log:               opts.Log,
		Config:            opts.Config,
		Redis:             redisClient,
		DBPool:            pgClient,
		Solr:              solrClient,
		GRPCServer:        grpcServer,
		GRPCOpts:          grpcOpts,
		TokenValidator:    tokenValidator,
		HTTPMux:           mux,
		TelemetryService:  telemetryService,
		TopicConfigChange: topicConfigChange,
		PromRegistry:      promRegistry,
	})
	if err != nil {
		cancel()
		return err
	}

	// Hook up the metrics endpoint
	promHandler := promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{Registry: promRegistry})
	err = mux.HandlePath(http.MethodGet, opts.Config.Web.MetricsPath, func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		promHandler.ServeHTTP(w, r)
	})
	if err != nil {
		cancel()
		return err
	}

	go func() {
		opts.Log.Info(ctx, L.Messagef("gRPC service started: %s", opts.Config.Web.GrpcHost), L.PhaseStartup)
		serverErrors <- grpcServer.Serve(lis)
	}()

	go func() {
		<-ctx.Done()
		opts.Log.Info(ctx, L.Message("stopping: gRPC server"), L.PhaseShutdown)
		grpcServer.GracefulStop()
		opts.Log.Info(ctx, L.Message("stopped: gRPC server"), L.PhaseShutdown)
	}()

	h := L.NewRequestResponseLoggingMiddleware(opts.Log, []string{
		// Endpoints excluded from logging
		"/probes/liveness",
		"/probes/readiness",
		"/metrics",
	}, true)(mux)

	for _, handler := range opts.AdditionalHTTPHandlerWrappers {
		h = handler(h)
	}

	h = web.NewResponseHeaderMiddleware(h, "X-Mex-Build", opts.Config.Version.Build)

	h = web.NewRequestIDMiddleware()(h)
	h = web.NewTraceIDMiddleware(constants.HTTPHeaderTraceID)(h)
	h = web.NewTracingDecisionMiddleware(opts.Config.Logging.TraceEnabled, opts.Config.Logging.TraceSecret)(h)

	if opts.Config.Web.IpFilter.Enabled {
		opts.Log.Info(context.Background(), L.Message("IP filter enabled"))
		h = web.NewIPFilter(opts.Log, opts.Config.Web.IpFilter.AllowedIps)(h)
	}

	if opts.Config.Web.RateLimiting.Enabled {
		rate := limiter.Rate{
			Period: opts.Config.Web.RateLimiting.Period.AsDuration(),
			Limit:  opts.Config.Web.RateLimiting.Limit,
		}
		opts.Log.Info(context.Background(), L.Messagef("rate limiter enabled: period: %v, limit: %d", rate.Period, rate.Limit), L.PhaseStartup)

		store := memory.NewStore()
		lim := limiter.New(store, rate, limiter.WithClientIPHeader(opts.Config.Web.RateLimiting.ClientIpHeader))
		h = stdlib.NewMiddleware(lim).Handler(h)
	}

	if opts.Config.Web.MaxBodyBytes > 0 {
		h = http.MaxBytesHandler(h, opts.Config.Web.MaxBodyBytes)
	}

	prodServer := &http.Server{
		Addr:              opts.Config.Web.ApiHost,
		ReadTimeout:       opts.Config.Web.ReadTimeout.AsDuration(),
		ReadHeaderTimeout: opts.Config.Web.ReadTimeout.AsDuration(),
		WriteTimeout:      opts.Config.Web.WriteTimeout.AsDuration(),
		IdleTimeout:       opts.Config.Web.IdleTimeout.AsDuration(),
		MaxHeaderBytes:    int(opts.Config.Web.MaxHeaderBytes),
		Handler:           h,
	}

	// Start the service and listen for API requests.
	go func() {
		opts.Log.Info(ctx, L.Messagef("started: HTTP gateway (%s, PID %d)", opts.Config.Web.ApiHost, os.Getpid()), L.PhaseStartup)
		serverErrors <- prodServer.ListenAndServe()
	}()

	go func() {
		<-ctx.Done()
		opts.Log.Info(ctx, L.Message("stopping: HTTP server"), L.PhaseShutdown)
		err = prodServer.Shutdown(ctx)
		if err != nil {
			opts.Log.Warn(ctx, L.Messagef("prod server shutdown error: %s", err.Error()))
		}
		opts.Log.Info(ctx, L.Message("stopped: HTTP server"), L.PhaseShutdown)
	}()

	// ------------------------------------------------------------------------------------------------------
	// Shutdown logic

	var errReturn error

	// Block and wait for shutdown.
	select {
	case err := <-serverErrors:
		errReturn = fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		opts.Log.Info(ctx, L.Messagef("shutdown initiated, signal: %v", sig), L.PhaseShutdown)
	}

	// Cancelling the context ctx signals resources (who actively support that) to shutdown/disconnect.
	cancel()

	// The above context cancellation will trigger the flushing of the event filter queue.
	// In order to guarantee issuance of all messages, we wait for its completion.
	opts.Log.Info(ctx, L.Message("flushing: event filter queue"), L.PhaseShutdown)
	<-efDone
	opts.Log.Info(ctx, L.Message("flushed: event filter queue"), L.PhaseShutdown)

	return errReturn
}
