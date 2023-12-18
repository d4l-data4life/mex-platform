package pingers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/d4l-data4life/mex/mex/shared/utils"
)

type RedisPinger struct {
	lastError error
	quit      chan<- struct{}
}

func (p *RedisPinger) LastError() error {
	return p.lastError
}

func (p *RedisPinger) Stop() {
	close(p.quit)
}

func NewRedisPinger(client *redis.Client, updateInterval time.Duration) *RedisPinger {
	if client == nil {
		panic("Redis client is nil")
	}

	p := RedisPinger{
		lastError: fmt.Errorf("Redis pinger did not run yet"),
	}

	p.quit = utils.ExponentialTicker(time.Second, updateInterval, func() {
		statusVal := client.Ping(context.Background())
		p.lastError = statusVal.Err()
	})

	return &p
}
