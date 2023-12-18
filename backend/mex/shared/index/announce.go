package index

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"

	L "github.com/d4l-data4life/mex/mex/shared/log"
)

type Announcer interface {
	AnnounceTechnicalItemID(technicalID string)
	AnnounceBusinessItemID(businessID string)
}

type RedisAnnouncer struct {
	Log   L.Logger
	Redis *redis.Client

	TechnicalIDsTopicName string
	BusinessIDsTopicName  string
	BusinessIDSetName     string

	SetExpiration time.Duration
}

func (idx *RedisAnnouncer) AnnounceTechnicalItemID(technicalID string) {
	ctx := context.Background()
	if idx.Redis.Publish(ctx, idx.TechnicalIDsTopicName, technicalID).Err() != nil {
		idx.Log.Warn(ctx, L.Messagef("could not publish technical ID: %s", technicalID))
	}
}

func (idx *RedisAnnouncer) AnnounceBusinessItemID(businessID string) {
	ctx := context.Background()

	// Add the item ID to the set so it is recorded for later indexing in case the below (*) broadcast is somehow lost.
	if cmdSAdd := idx.Redis.SAdd(ctx, idx.BusinessIDSetName, businessID); cmdSAdd.Err() != nil {
		idx.Log.Warn(ctx, L.Messagef("could not add business ID to indexer set: %s (%s)", businessID, cmdSAdd.Err().Error()))
	}

	if idx.Redis.Expire(ctx, idx.BusinessIDSetName, idx.SetExpiration).Err() != nil {
		idx.Log.Warn(ctx, L.Message("could not reset indexer expiration time"))
	}

	// (*) Broadcast the item ID.
	if cmdPublish := idx.Redis.Publish(ctx, idx.BusinessIDsTopicName, businessID); cmdPublish.Err() != nil {
		idx.Log.Warn(ctx, L.Messagef("could not publish business ID: %s: %s", businessID, cmdPublish.Err().Error()))
	}
}
