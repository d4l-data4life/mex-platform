package index

import (
	"context"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/d4l-data4life/mex/mex/shared/codings/csrepo"
	"github.com/d4l-data4life/mex/mex/shared/constants"
	"github.com/d4l-data4life/mex/mex/shared/entities"
	"github.com/d4l-data4life/mex/mex/shared/known/statuspb"
	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/searchconfig"
	"github.com/d4l-data4life/mex/mex/shared/solr"
	"github.com/d4l-data4life/mex/mex/shared/telemetry"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields/hooks"
	"github.com/d4l-data4life/mex/mex/services/metadata/endpoints/jobs"

	"github.com/d4l-data4life/mex/mex/services/index/endpoints/index/pb"
)

type Service struct {
	ServiceTag string
	Log        L.Logger

	DB             *pgxpool.Pool
	Redis          *redis.Client
	Solr           solr.ClientAPI
	SolrCollection string

	JobService *jobs.Service

	FieldRepo        fields.FieldRepo
	EntityRepo       entities.EntityRepo
	SearchConfigRepo searchconfig.SearchConfigRepo
	CodingsetRepo    csrepo.CodingsetRepo

	// Field lifecycle hooks
	SolrFieldCreationHooks hooks.SolrFieldCreationHooks
	SolrDataLoadHooks      hooks.SolrDataLoadHooks

	TelemetryService *telemetry.Service

	CollectionName    string
	ReplicationFactor uint32

	pb.UnimplementedIndexServer
}

const SvcResourceName = "index"

// This method makes the Service an rdb.TopicSubscriber.
// It gets called when the subscribed topic is triggered.
func (svc *Service) Message(ctx context.Context, topic string, configHash string) {
	if !strings.HasSuffix(topic, constants.ConfigUpdateChannelNameSuffix) {
		return
	}

	_ = svc.FieldRepo.Purge(ctx)
	_ = svc.EntityRepo.Purge(ctx)
	_ = svc.SearchConfigRepo.Purge(ctx)
	_ = svc.CodingsetRepo.Purge(ctx)

	// Trigger collection recreation
	svc.Log.Info(ctx, L.Messagef("received config change message for config commit hash %q", configHash))

	if err := svc.runRecreateCollectionJob(configHash); err != nil {
		svc.Log.Warn(ctx, L.Messagef("index rebuild based on hash %q: error: %s", configHash, err.Error()))
		svc.TelemetryService.SetStatus(statuspb.Color_RED, configHash)
	}
}
