package items

import (
	"context"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/d4l-data4life/mex/mex/shared/auth"
	"github.com/d4l-data4life/mex/mex/shared/cfg"
	"github.com/d4l-data4life/mex/mex/shared/constants"
	"github.com/d4l-data4life/mex/mex/shared/entities"
	"github.com/d4l-data4life/mex/mex/shared/index"
	"github.com/d4l-data4life/mex/mex/shared/jobs"
	"github.com/d4l-data4life/mex/mex/shared/known/jobspb"
	"github.com/d4l-data4life/mex/mex/shared/known/statuspb"
	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/telemetry"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields/hooks"
	itemspb "github.com/d4l-data4life/mex/mex/services/metadata/endpoints/items/pb"
)

const (
	MetadataItemUpdateByItemIDChannelName     = "mex-metadata-item-update-by-item-id"
	MetadataItemUpdateByBusinessIDChannelName = "mex-metadata-item-update-by-business-id"
	MetadataItemUpdateByBusinessIDSetName     = "mex-metadata-item-update-by-business-id-set"
)

type Service struct {
	ServiceTag string
	Log        L.Logger

	DB    *pgxpool.Pool
	Redis *redis.Client

	Jobber jobs.Jobber

	FieldRepo  fields.FieldRepo
	EntityRepo entities.EntityRepo

	// Field lifecycle hooks
	ItemCreationHooks      hooks.ItemCreationHooks
	SolrFieldCreationHooks hooks.SolrFieldCreationHooks
	SolrDataLoadHooks      hooks.SolrDataLoadHooks

	Announcer        index.Announcer
	TelemetryService *telemetry.Service

	DuplicateDetectionAlgorithm cfg.DuplicateDetectionAlgorithm

	itemspb.UnimplementedItemsServer
}

const SvcResourceName = "items"

func AddJobItemIDs(ctx context.Context, j jobs.Jobber, itemIDs []string) {
	jobID := auth.GetJobID(ctx)
	if jobID != "" {
		_, _ = j.AddItems(ctx, &jobspb.AddJobItemsRequest{
			JobId:   jobID,
			ItemIds: itemIDs,
		})
	}
}

// This method makes the Service an rdb.TopicSubscriber
func (svc *Service) Message(ctx context.Context, topic string, configHash string) {
	if !strings.HasSuffix(topic, constants.ConfigUpdateChannelNameSuffix) {
		return
	}

	_ = svc.EntityRepo.Purge(context.Background())
	_ = svc.FieldRepo.Purge(context.Background())

	svc.TelemetryService.SetStatus(statuspb.Color_GREEN, configHash)
}
