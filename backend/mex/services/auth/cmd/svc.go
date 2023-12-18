package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/d4l-data4life/mex/mex/shared/cfg"
	"github.com/d4l-data4life/mex/mex/shared/interceptors"
	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/log/emit"
	"github.com/d4l-data4life/mex/mex/shared/svcutils"
	"github.com/d4l-data4life/mex/mex/shared/try"

	"github.com/d4l-data4life/mex/mex/services/auth/endpoints/auth"
	pbAuth "github.com/d4l-data4life/mex/mex/services/auth/endpoints/auth/pb"
)

var (
	build     = "develop" // set during build
	buildDate = "now"     // set during build
)

const (
	serviceName = "MEx Auth"
	serviceTag  = "auth"
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

	err = svcutils.Run(ctx, svcutils.RunOpts{
		ServiceTag: serviceTag,
		Log:        log,
		Setup:      setup,
		Config:     &config,

		Support: svcutils.Support{},

		AdditionalServeMuxOpts: []runtime.ServeMuxOption{
			runtime.WithMarshalerOption("*", interceptors.NewMarshaler(config.Strictness.StrictJsonParsing.Auth)),
		},
	})
	if err != nil {
		log.Error(ctx, L.Message(err.Error()))
		os.Exit(1)
	}
}

func setup(ctx context.Context, opts svcutils.SetupOpts) error {
	if !opts.Config.Oauth.Server.Enabled {
		return nil
	}

	if !cfg.StringIsEmpty(opts.Config.Oauth.Server.SigningPrivateKeyFile) && !cfg.BytesAreEmpty(opts.Config.Oauth.Server.SigningPrivateKeyPem) {
		return fmt.Errorf("cannot specify both signing key file and PEM")
	}

	if cfg.StringIsEmpty(opts.Config.Oauth.Server.SigningPrivateKeyFile) && cfg.BytesAreEmpty(opts.Config.Oauth.Server.SigningPrivateKeyPem) {
		return fmt.Errorf("must specify either signing key file or PEM")
	}

	var signingPEMReader io.Reader
	if !cfg.StringIsEmpty(opts.Config.Oauth.Server.SigningPrivateKeyFile) {
		trimmedFilename := cfg.TrimEmpty(opts.Config.Oauth.Server.SigningPrivateKeyFile)

		r, err := try.Try(ctx, try.Task[io.Reader]{
			Desc:  "load PEM file: " + trimmedFilename,
			Phase: "startup",

			PauseStrategy: try.NewMaxAttemptsConstantPauseStrategy(opts.Config.Oauth.Server.KeyFileAccessAttempts, opts.Config.Oauth.Server.KeyFileAccessPause.AsDuration()),

			Log: opts.Log,
			Func: func() (io.Reader, error) {
				return os.Open(trimmedFilename)
			},
		})
		if err != nil {
			return fmt.Errorf("could not load PEM file %s: %w", trimmedFilename, err)
		}
		signingPEMReader = r
	}

	if !cfg.BytesAreEmpty(opts.Config.Oauth.Server.SigningPrivateKeyPem) {
		opts.Log.Info(ctx, L.Messagef("PEM content in env var, length: %d", len(opts.Config.Oauth.Server.SigningPrivateKeyPem)))
		if opts.Config.Oauth.Server.SigningPrivateKeyPem[0] != ' ' {
			signingPEMReader = strings.NewReader(string(opts.Config.Oauth.Server.SigningPrivateKeyPem))
		} else {
			opts.Log.Info(ctx, L.Message("ignoring empty PEM parameter"))
		}
	}

	authService, err := auth.NewService(auth.ServiceOptions{
		ServiceTag: serviceTag,
		Log:        opts.Log,

		Redis: opts.Redis,

		ClientID:        opts.Config.Oauth.ClientId,
		ProducerGroupID: opts.Config.Oauth.ProducerGroupId,
		ConsumerGroupID: opts.Config.Oauth.ConsumerGroupId,

		SigningKeyPEMReader: signingPEMReader,
		SignatureAlg:        opts.Config.Oauth.Server.SignatureAlg,
		ClientSecrets:       opts.Config.Oauth.Server.ClientSecrets,
		RedirectURIs:        opts.Config.Oauth.Server.RedirectUris,
		GrantFlows:          opts.Config.Oauth.Server.GrantFlows,

		AuthCodeValidity:     opts.Config.Oauth.Server.AuthCodeValidity.AsDuration(),
		AccessTokenValidity:  opts.Config.Oauth.Server.AccessTokenValidity.AsDuration(),
		RefreshTokenValidity: opts.Config.Oauth.Server.RefreshTokenValidity.AsDuration(),

		TelemetryService: opts.TelemetryService,
	})
	if err != nil {
		return err
	}

	opts.TopicConfigChange.Subscribe(authService)

	opts.Log.Info(ctx, L.Message("OAuth endpoints enabled"), L.PhaseStartup)
	pbAuth.RegisterAuthServer(opts.GRPCServer, authService)

	err = pbAuth.RegisterAuthHandlerFromEndpoint(ctx, opts.HTTPMux, opts.Config.Web.GrpcHost, opts.GRPCOpts)
	if err != nil {
		return err
	}

	return nil
}
