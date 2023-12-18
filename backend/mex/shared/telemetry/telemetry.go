package telemetry

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/d4l-data4life/mex/mex/shared/known/statuspb"
	L "github.com/d4l-data4life/mex/mex/shared/log"
)

type Pinger interface {
	Stop()
	LastError() error
}

type Service struct {
	mu sync.Mutex

	log        L.Logger
	serviceTag string
	pingers    []Pinger
	redis      *redis.Client
	dayOfMonth int
	status     *statuspb.Status

	UnimplementedTelemetryServer
}

func New(logger L.Logger, serviceTag string, redisClient *redis.Client, statusUpdateInterval time.Duration, quit <-chan struct{}) *Service {
	svc := Service{
		log:        logger,
		serviceTag: serviceTag,
		pingers:    []Pinger{},
		redis:      redisClient,
		status: &statuspb.Status{
			ServiceTag: serviceTag,
			Replica:    hostnameOrLocalhost(),
			Color:      statuspb.Color_GREEN,
			ConfigHash: "∅",
		},
	}

	svc.start(statusUpdateInterval, quit)

	return &svc
}

func (svc *Service) LivenessProbe(_ context.Context, _ *LivenessRequest) (*LivenessResponse, error) {
	return &LivenessResponse{}, nil
}

// Readiness is defined by the availability of the pingers.
//
// These will typically include:
//   - Postgres
//   - Solr
//   - Redis
//   - Key store
func (svc *Service) ReadinessProbe(ctx context.Context, _ *ReadinessRequest) (*statuspb.Status, error) {
	for _, p := range svc.pingers {
		err := p.LastError()
		if err != nil {
			svc.log.Warn(ctx, L.Messagef("pinger: last error: %s", err.Error()))
			return nil, status.Error(codes.Unavailable, "service not ready")
		}
	}

	// Report service as ready
	return svc.status, nil
}

func (svc *Service) AddPinger(p Pinger) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	svc.pingers = append(svc.pingers, p)
}

func (svc *Service) start(interval time.Duration, quit <-chan struct{}) {
	statusTicker := time.NewTicker(interval)

	go func() {
		ctx := context.Background()
		for {
			select {
			case <-statusTicker.C:
				if svc.redis == nil {
					// Don't do anything if Redis not specified.
					continue
				}

				err := svc.PublishStatus(ctx)
				if err != nil {
					svc.log.Error(ctx, L.Messagef("error publishing status to Redis: %s", err.Error()))
				}

			case <-quit:
				svc.log.Info(context.Background(), L.Message("stopping: service status reporter"))
				statusTicker.Stop()
				for _, p := range svc.pingers {
					p.Stop()
				}
				svc.log.Info(context.Background(), L.Message("stopped: service status reporter"))
				return
			}
		}
	}()
}

func hostnameOrLocalhost() string {
	hostname := os.Getenv("HOSTNAME")
	if hostname != "" {
		return hostname
	}
	return "localhost"
}

func (svc *Service) SetStatus(color statuspb.Color, configHash string) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	svc.log.Trace(context.Background(), L.Messagef("telemetry message: new status: color :%v, config hash: %q", color, configHash))

	if svc.status == nil {
		panic("status is nil")
	}

	svc.status.Color = color
	svc.status.ConfigHash = configHash
}

func (svc *Service) GetStatus() (statuspb.Color, string) {
	return svc.status.Color, svc.status.ConfigHash
}

// Implement the utils.Progressor interface
func (svc *Service) Progress(step string, details string) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	svc.status.Progress = &statuspb.Progress{
		Step:    step,
		Details: details,
	}
}

func (svc *Service) Done() {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	svc.status.Progress = nil
}

const StatusHashName = "status"

func (svc *Service) PublishStatus(ctx context.Context) error {
	svc.status.LastReported = timestamppb.Now()

	// Delete the status hash once a day so it does not accumulate stale fields of past replicas.
	// (It will take statusUpdateInterval time to re-populate it, but that is okay.)
	newDayOfMonth := svc.status.LastReported.AsTime().Day()
	if svc.dayOfMonth != newDayOfMonth {
		svc.log.Info(ctx, L.Messagef("deleting: %q Redis hash", StatusHashName))
		err := svc.redis.Del(ctx, StatusHashName).Err()
		if err != nil {
			svc.log.Warn(ctx, L.Messagef("error: could not delete Redis hash %q: %s", StatusHashName, err.Error()))
		} else {
			svc.log.Info(ctx, L.Messagef("deleted:  %q Redis hash", StatusHashName))
		}
		svc.dayOfMonth = newDayOfMonth
	}

	buf, err := protojson.Marshal(svc.status)
	if err != nil {
		return err
	}

	err = svc.redis.HSet(ctx, StatusHashName, fmt.Sprintf("%s:%s", svc.status.ServiceTag, svc.status.Replica), string(buf)).Err()
	if err != nil {
		svc.log.Error(ctx, L.Messagef("error writing status to Redis: %s", err.Error()))
	}

	return svc.redis.Expire(ctx, StatusHashName, time.Minute).Err()
}
