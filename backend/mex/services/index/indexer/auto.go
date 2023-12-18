package indexer

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"

	L "github.com/d4l-data4life/mex/mex/shared/log"

	"github.com/d4l-data4life/mex/mex/services/index/endpoints/index"
	"github.com/d4l-data4life/mex/mex/services/index/endpoints/index/pb"
)

type AutoIndexer interface {
	AnnounceTechnicalItemID(technicalID string)
	AnnounceBusinessItemID(businessID string)

	StartPeriodicIndexer()
	StopPeriodicIndexer()
}

type indexer struct {
	mu sync.Mutex

	log   L.Logger
	redis *redis.Client

	technicalIDsTopicName string
	businessIDsTopicName  string
	businessIDSetName     string

	indexService *index.Service

	setExpiration time.Duration

	ticker *time.Ticker
	quit   chan struct{}
}

type AutoIndexerConfig struct {
	Log          L.Logger
	Redis        *redis.Client
	IndexService *index.Service

	TechnicalIDsTopicName string
	BusinessIDsTopicName  string
	BusinessIDSetName     string

	SetExpiration time.Duration
}

func NewAutoIndexer(ctx context.Context, config AutoIndexerConfig) AutoIndexer {
	idx := indexer{
		log:   config.Log,
		redis: config.Redis,

		technicalIDsTopicName: config.TechnicalIDsTopicName,
		businessIDsTopicName:  config.BusinessIDsTopicName,
		businessIDSetName:     config.BusinessIDSetName,

		indexService:  config.IndexService,
		setExpiration: config.SetExpiration,
	}

	topicByTechnicalID := idx.redis.Subscribe(ctx, config.TechnicalIDsTopicName)
	channelByTechnicalID := topicByTechnicalID.Channel()

	topicByBusinessID := idx.redis.Subscribe(ctx, config.BusinessIDsTopicName)
	channelByBusinessID := topicByBusinessID.Channel()

	go func() {
		for {
			select {
			case msg := <-channelByTechnicalID:
				idx.log.Trace(context.Background(), L.Messagef("AutoIndexer registered new item with technical ID: %s", msg.Payload), L.Phase("id-update"))

			case msg := <-channelByBusinessID:
				idx.log.Info(ctx, L.Messagef("AutoIndexer registered new item with business ID: %s", msg.Payload), L.Phase("bid-update"))

				go func() {
					// The event of a new business item will be received by every service replica.
					// We want to ensure that only one replica will do an indexing, and not each replica.
					// This is achieved by attempting to remove the business ID from the list first.
					// Since Redis is single-threaded, this will only succeed for one replica, which is the one responsible for indexing.
					// The SREM command conveniently returns the actual number of removed items; if that value is 1, we are the indexing replica.
					sremCmd := idx.redis.SRem(ctx, idx.businessIDSetName, msg.Payload)
					switch {
					case sremCmd.Err() != nil:
						idx.log.Warn(ctx, L.Messagef("error deleting item '%s' from set '%s': %s", msg.Payload, idx.businessIDSetName, sremCmd.Err().Error()), L.Phase("bid-update"))

					case sremCmd.Val() == 1:
						_, err := config.IndexService.IndexLatestItem(ctx, &pb.IndexLatestItemRequest{BusinessId: msg.Payload})
						if err != nil {
							idx.log.Warn(ctx, L.Messagef("error indexing item '%s': %s", msg.Payload, err.Error()), L.Phase("bid-update"))

							// We are the indexing replica, but something went wrong.
							// We stick the business ID back into the set so the periodic indexer can take care of it.
							err = idx.redis.SAdd(ctx, idx.businessIDSetName, msg.Payload).Err()
							if err != nil {
								idx.log.Warn(ctx, L.Messagef("error re-adding business ID '%s' to Redis set '%s': %s", msg.Payload, idx.businessIDSetName, sremCmd.Err().Error()), L.Phase("bid-update"))
							}
						}

					default:
						idx.log.Info(ctx, L.Messagef("wanted to index '%s', but it seems another instance took care of it", msg.Payload), L.Phase("bid-update"))
					}
				}()

			case <-ctx.Done():
				idx.log.Info(ctx, L.Message("imminent shutdown; unsubscribing from topics"), L.Phase("shutdown"))
				_ = topicByTechnicalID.Unsubscribe(ctx)
				_ = topicByTechnicalID.Close()

				_ = topicByBusinessID.Unsubscribe(ctx)
				_ = topicByBusinessID.Close()

				return
			}
		}
	}()

	return &idx
}

func (idx *indexer) AnnounceTechnicalItemID(technicalID string) {
	ctx := context.Background()
	if idx.redis.Publish(ctx, idx.technicalIDsTopicName, technicalID).Err() != nil {
		idx.log.Warn(ctx, L.Messagef("could not publish technical ID: %s", technicalID))
	}
}

func (idx *indexer) AnnounceBusinessItemID(businessID string) {
	ctx := context.Background()

	// Add the item ID to the set so it is recorded for later indexing in case the below (*) broadcast is somehow lost.
	if cmdSAdd := idx.redis.SAdd(ctx, idx.businessIDSetName, businessID); cmdSAdd.Err() != nil {
		idx.log.Warn(ctx, L.Messagef("could not add business ID to indexer set: %s (%s)", businessID, cmdSAdd.Err().Error()))
	}

	if idx.redis.Expire(ctx, idx.businessIDSetName, idx.setExpiration).Err() != nil {
		idx.log.Warn(ctx, L.Message("could not reset indexer expiration time"))
	}

	// (*) Broadcast the item ID.
	if cmdPublish := idx.redis.Publish(ctx, idx.businessIDsTopicName, businessID); cmdPublish.Err() != nil {
		idx.log.Warn(ctx, L.Messagef("could not publish business ID: %s: %s", businessID, cmdPublish.Err().Error()))
	}
}

func (idx *indexer) StartPeriodicIndexer() {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	if idx.ticker != nil {
		return
	}

	idx.ticker = time.NewTicker(idx.setExpiration / 2)
	idx.quit = make(chan struct{})

	go func() {
		for {
			select {
			case <-idx.ticker.C:
				idx.indexOverlookedItems()

			case <-idx.quit:
				idx.ticker.Stop()
				idx.log.Info(context.Background(), L.Message("stopping ticker"))
				return
			}
		}
	}()
}

func (idx *indexer) StopPeriodicIndexer() {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	if idx.quit == nil {
		return
	}

	close(idx.quit)

	idx.quit = nil
	idx.ticker = nil
}

func (idx *indexer) indexOverlookedItems() {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	ctx := context.Background()

	scardCmd := idx.redis.SCard(ctx, idx.businessIDSetName)
	if scardCmd.Err() != nil {
		idx.log.Warn(ctx, L.Messagef("error getting set size of set '%s': %s", idx.businessIDSetName, scardCmd.Err().Error()), L.Phase("bid-update"))
		return
	}

	if scardCmd.Val() == 0 {
		return
	}

	idx.log.Info(ctx, L.Messagef("indexing %d overlooked items", scardCmd.Val()), L.Phase("bid-update"))

	smembersCmd := idx.redis.SMembers(ctx, idx.businessIDSetName)
	if smembersCmd.Err() != nil {
		idx.log.Warn(ctx, L.Messagef("error getting set '%s': %s", idx.businessIDSetName, smembersCmd.Err().Error()), L.Phase("bid-update"))
		return
	}

	for _, businessID := range smembersCmd.Val() {
		_, err := idx.indexService.IndexLatestItem(ctx, &pb.IndexLatestItemRequest{BusinessId: businessID})
		if err == nil {
			sremCmd := idx.redis.SRem(ctx, idx.businessIDSetName, businessID)
			if sremCmd.Err() != nil {
				idx.log.Warn(ctx, L.Messagef("error deleting item %s from set '%s': %s", businessID, idx.businessIDSetName, sremCmd.Err().Error()), L.Phase("bid-update"))
			}
		} else {
			idx.log.Warn(ctx, L.Messagef("could not index item %s: %s", businessID, err.Error()), L.Phase("bid-update"))
		}
	}
}
