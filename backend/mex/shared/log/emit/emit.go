package emit

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"

	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/try"
)

type RedisEmitter struct {
	sync.Mutex
	*redis.Client
}

func (emitter *RedisEmitter) Emit(e L.Entry) {
	emitter.Lock()
	defer emitter.Unlock()

	if e.JobID() == "" {
		return
	}

	key := fmt.Sprintf("jobs:%s:logs", e.JobID())
	_ = emitter.RPush(context.Background(), key, string(e.Bytes()))
}

type WriterEmitter struct {
	sync.Mutex
	io.Writer
}

func (emitter *WriterEmitter) Emit(e L.Entry) {
	emitter.Lock()
	defer emitter.Unlock()

	_, _ = emitter.Write(e.Bytes())
}

type eventsFilterEmitter struct {
	Origin string
	Path   string
	Secret string

	q   *try.Requeue
	log L.LevelLogger
}

const (
	defaultAttemptCount = 3
	defaultPause        = 100 * time.Millisecond
)

func NewEventsFilterEmitter(ctx context.Context, log L.LevelLogger, origin string, path string, secret string) (L.Emitter, <-chan struct{}) {
	if log == nil {
		log = &L.NullLogger{}
	}
	q, done := try.NewRequeue(ctx, log, defaultAttemptCount, defaultPause)
	emitter := eventsFilterEmitter{
		q:      q,
		log:    log,
		Origin: origin,
		Path:   path,
		Secret: secret,
	}
	return &emitter, done
}

func (emitter *eventsFilterEmitter) Emit(e L.Entry) {
	if e.Type() != L.LogEntryTypeBI {
		return
	}

	buf := e.Bytes()
	if buf != nil {
		emitter.q.Push(func() error {
			req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", emitter.Origin, emitter.Path), bytes.NewReader(buf))
			if err != nil {
				emitter.log.Warn(context.Background(), L.Messagef("error creating POST request to BI events filter service: %s", err.Error()))
				return err
			}
			req.Header.Set("content-type", "application/json")
			req.Header.Set("authorization", emitter.Secret)

			resp, err := http.DefaultClient.Do(req) // FIXME: use custom client!
			if err != nil {
				emitter.log.Warn(context.Background(), L.Messagef("error sending POST request to BI events filter service: %s", err.Error()))
				return err
			}
			defer resp.Body.Close()
			return nil
		})
	}
}
