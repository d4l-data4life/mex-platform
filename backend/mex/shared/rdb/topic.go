package rdb

import (
	"context"
	"sync"

	"github.com/go-redis/redis/v8"

	L "github.com/d4l-data4life/mex/mex/shared/log"
)

type TopicSubscriber interface {
	Message(ctx context.Context, topic string, msg string)
}

type Topic struct {
	mu sync.Mutex

	subscribers []TopicSubscriber
	topic       string
	pubsub      *redis.PubSub
}

func (cl *Topic) Subscribe(l TopicSubscriber) {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	cl.subscribers = append(cl.subscribers, l)
}

func (cl *Topic) Unsubscribe(ctx context.Context) {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	_ = cl.pubsub.Unsubscribe(ctx)
	_ = cl.pubsub.Close()
}

func (cl *Topic) informAll(ctx context.Context, hash string) {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	for _, l := range cl.subscribers {
		l.Message(ctx, cl.topic, hash)
	}
}

func NewTopic(ctx context.Context, log L.Logger, redisClient *redis.Client, topic string) *Topic {
	cl := Topic{
		topic:  topic,
		pubsub: redisClient.Subscribe(ctx, topic),
	}

	redisChannel := cl.pubsub.Channel()

	go func() {
		for {
			select {
			case msg := <-redisChannel:
				if msg == nil {
					log.Warn(ctx, L.Message("Redis message is nil"))
					continue
				}

				log.Info(ctx, L.Messagef("message on %q: %q", topic, msg.Payload))
				cl.informAll(ctx, msg.Payload)

			case <-ctx.Done():
				log.Info(ctx, L.Messagef("imminent shutdown; unsubscribing from topic '%s'", topic), L.Phase("shutdown"))
				cl.Unsubscribe(ctx)
				return
			}
		}
	}()

	return &cl
}

func (cl *Topic) Topic() string {
	return cl.topic
}
