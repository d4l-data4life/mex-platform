package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/d4l-data4life/mex/mex/shared/cfg"
	"github.com/d4l-data4life/mex/mex/shared/constants"
	sharedJobs "github.com/d4l-data4life/mex/mex/shared/jobs"
	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/log/emit"
	"github.com/d4l-data4life/mex/mex/shared/svcutils"

	"github.com/d4l-data4life/mex/mex/services/config/endpoints/config"
	pbConfig "github.com/d4l-data4life/mex/mex/services/config/endpoints/config/pb"
)

var (
	build     = "develop" // set during build
	buildDate = "now"     // set during build
)

const (
	serviceName = "MEx Config"
	serviceTag  = "config"
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

	err = svcutils.Run(ctx, svcutils.RunOpts{
		ServiceTag: serviceTag,
		Log:        log,
		Setup:      setup,
		Config:     &mexConfig,

		Support: svcutils.Support{},

		AdditionalTokenValidationExcludePatterns: []string{"/d4l.mex.config.Config/GetFile"},
		AdditionalServeMuxOpts: []runtime.ServeMuxOption{
			runtime.WithMarshalerOption("*", config.NewFileContentMarshaler(mexConfig.Strictness.StrictJsonParsing.Config)),
		},
	})
	if err != nil {
		log.Error(ctx, L.Message(err.Error()))
		os.Exit(1)
	}
}

func setup(ctx context.Context, opts svcutils.SetupOpts) error {
	if len(opts.Config.Services.Config.ApiKeys) == 0 {
		return fmt.Errorf("no API keys configured; please set at least one API key")
	}

	for k, key := range opts.Config.Services.Config.ApiKeys {
		if strings.TrimSpace(key) != key {
			return fmt.Errorf("API key #%d has leading or trailing whitespace; please remove", k+1)
		}
	}

	opts.Log.Info(ctx, L.Messagef("number of API keys: %d", len(opts.Config.Services.Config.ApiKeys)))

	broadcastTopicName := fmt.Sprintf("%s/%s", opts.Config.Redis.PubSubPrefix, constants.ConfigUpdateChannelNameSuffix)
	opts.Log.Info(ctx, L.Messagef("topic: %s", broadcastTopicName))

	configService := config.Service{
		ServiceTag: serviceTag,
		Log:        opts.Log,

		Redis:              opts.Redis,
		BroadcastTopicName: broadcastTopicName,

		RepoName:          opts.Config.Services.Config.Github.RepoName,
		DefaultBranchName: opts.Config.Services.Config.Github.DefaultBranchName,
		EnvPath:           opts.Config.Services.Config.EnvPath,
		UpdateTimeout:     opts.Config.Services.Config.UpdateTimeout.AsDuration(),

		TelemetryService: opts.TelemetryService,
		Jobber: sharedJobs.RedisJobber{
			Redis:      opts.Redis,
			Expiration: opts.Config.Jobs.Expiration.AsDuration(),
		},
	}

	if !cfg.StringIsEmpty(opts.Config.Services.Config.Github.RepoName) {
		err := configService.InitDeployKey(ctx, opts.Config.Services.Config.Github.DeployKeyPem)
		if err != nil {
			return err
		}
		_, err = configService.UpdateConfig(ctx, &pbConfig.UpdateConfigRequest{
			UpdateType: &pbConfig.UpdateConfigRequest_RefName{
				RefName: opts.Config.Services.Config.Github.DefaultBranchName,
			},
		},
		)
		if err != nil {
			return err
		}
	} else {
		opts.Log.Warn(ctx, L.Message("repo name is empty; no Github actions possible (test mode only)"))
	}

	pbConfig.RegisterConfigServer(opts.GRPCServer, &configService)

	err := pbConfig.RegisterConfigHandlerFromEndpoint(ctx, opts.HTTPMux, opts.Config.Web.GrpcHost, opts.GRPCOpts)
	if err != nil {
		return err
	}

	return nil
}
