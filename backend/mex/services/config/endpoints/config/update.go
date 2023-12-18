package config

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/go-git/go-billy/v5/memfs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/d4l-data4life/mex/mex/shared/constants"
	E "github.com/d4l-data4life/mex/mex/shared/errstat"
	"github.com/d4l-data4life/mex/mex/shared/known/jobspb"
	"github.com/d4l-data4life/mex/mex/shared/known/statuspb"
	L "github.com/d4l-data4life/mex/mex/shared/log"

	pbConfig "github.com/d4l-data4life/mex/mex/services/config/endpoints/config/pb"
)

type SubscriberFunc func(ctx context.Context, topic string, message string)

func (f SubscriberFunc) Message(ctx context.Context, topic string, message string) {
	f(ctx, topic, message)
}

const roundDuration = 2 * time.Second

func (svc *Service) UpdateConfig(ctx context.Context, request *pbConfig.UpdateConfigRequest) (*pbConfig.UpdateConfigResponse, error) {
	if svc.RepoName == "" {
		return nil, E.MakeGRPCStatus(codes.Internal, "no repo name configured; cannot clone; test mode only").Err()
	}

	// Check request type, that is, return error if not one of the two supported ones.
	switch ty := request.UpdateType.(type) {
	case *pbConfig.UpdateConfigRequest_RefName:
	case *pbConfig.UpdateConfigRequest_CannedConfig:
		// happy case: drop out below switch
	default:
		return nil, E.MakeGRPCStatus(codes.InvalidArgument, fmt.Sprintf("unknown type: %t", ty)).Err()
	}

	lock, err := svc.Jobber.AcquireLock(ctx, ConfigResourceName)
	if err != nil {
		return nil, E.MakeGRPCStatus(codes.AlreadyExists, "failed to acquire config lock; other job might be running", request).Err()
	}

	job, err := svc.Jobber.CreateJob(ctx, &jobspb.CreateJobRequest{Title: "Repopulate Solr index"})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failure creating job:  %s", err.Error()))
	}

	svc.TelemetryService.SetStatus(statuspb.Color_AMBER, EmptyConfigHash)
	_ = svc.TelemetryService.PublishStatus(ctx)

	ctxJob := constants.NewContextWithValues(ctx, job.JobId)

	go func(ctx context.Context) {
		logJobError := func(message string) {
			svc.Log.Error(ctx, L.Message(message))
			_, err := svc.Jobber.SetError(ctx, &jobspb.SetJobErrorRequest{Error: message, JobId: job.JobId})
			if err != nil {
				svc.Log.Warn(ctx, L.Messagef("could not set job error: %s", err.Error()))
			}
		}

		svc.Jobber.SetStatusRunning(ctx, job.JobId) //nolint:errcheck
		hash := ""

		switch ty := request.UpdateType.(type) {
		case *pbConfig.UpdateConfigRequest_RefName:
			hash, err = svc.updateConfigFromRefName(ctx, ty.RefName)
		case *pbConfig.UpdateConfigRequest_CannedConfig:
			hash, err = svc.updateConfigFromCannedConfig(ctx, ty.CannedConfig)
		}
		if err != nil {
			logJobError(fmt.Sprintf("config update failed: %s", err.Error()))
			svc.Jobber.SetStatusDone(ctx, job.JobId)              //nolint:errcheck
			svc.Jobber.ReleaseLock(ctx, ConfigResourceName, lock) //nolint:errcheck

			svc.TelemetryService.SetStatus(statuspb.Color_RED, EmptyConfigHash)
			return
		}

		svc.announceConfigChange(ctx, hash)

		// Wait some time for all other services to internalize the new config.
		rounds := int(svc.UpdateTimeout / roundDuration)
		svc.Log.Info(ctx, L.Messagef("waiting: for all services (max. %d rounds)", rounds))
		success := false
	L:
		for k := 0; k < rounds; k++ {
			time.Sleep(roundDuration)
			switch svc.checkServices(ctx, hash, defaultMaxAge) {
			case statuspb.Color_RED:
				success = false
				break L
			case statuspb.Color_GREEN:
				success = true
				break L
			// would be the default behavior, but let's make it explicit
			case statuspb.Color_AMBER:
				continue
			}
		}
		if success {
			svc.Log.Info(ctx, L.Message("waited:  for all services (success)"))
		} else {
			svc.Log.Warn(ctx, L.Message("waited:  for all services (timeout)"))

			//nolint:errcheck
			svc.Jobber.SetError(ctx, &jobspb.SetJobErrorRequest{
				JobId: job.JobId,
				Error: fmt.Sprintf("no successful service states for config hash %q after waiting %v", hash, svc.UpdateTimeout),
			})
		}

		svc.Jobber.SetStatusDone(ctx, job.JobId)              //nolint:errcheck
		svc.Jobber.ReleaseLock(ctx, ConfigResourceName, lock) //nolint:errcheck

		if success {
			svc.TelemetryService.SetStatus(statuspb.Color_GREEN, hash)
		} else {
			svc.TelemetryService.SetStatus(statuspb.Color_RED, hash)
		}

		// Job lock released via defer.
	}(ctxJob)

	return &pbConfig.UpdateConfigResponse{JobId: job.JobId}, nil
}

// This function executes the following logic:
//
// - Get all service replicas' statuses which are younger than maxAge.
// - Iterate over all service replicas except your own and return nil only if all config hashes equal the given hash and the status is GREEN.
//   - Return RED: if one status is RED (or the replica statuses could not be retrieved in the first place)
//   - Return AMBER: if one status is AMBER
//   - Return GREEN: only if all statuses are GREEN
func (svc *Service) checkServices(ctx context.Context, hash string, maxAge time.Duration) statuspb.Color {
	m, err := svc.getAllServiceStatuses(ctx, maxAge)
	if err != nil {
		return statuspb.Color_RED
	}

	for serviceTag, bucket := range m {
		// Ignore config service itself as it will become ready only after all others are ready.
		if serviceTag != svc.ServiceTag {
			for _, b := range bucket {
				if b.ConfigHash != hash {
					return statuspb.Color_AMBER
				}
				if b.Color != statuspb.Color_GREEN {
					return b.Color
				}
			}
		}
	}

	return statuspb.Color_GREEN
}

// Announce config change to other services via Redis topic channel
func (svc *Service) announceConfigChange(ctx context.Context, headHash string) {
	svc.Log.Info(ctx, L.Messagef("announcing config change, topic: %q / hash: %q", svc.BroadcastTopicName, headHash))
	if err := svc.Redis.Publish(ctx, svc.BroadcastTopicName, headHash).Err(); err != nil {
		svc.Log.Warn(ctx, L.Messagef("could not publish to topic %q: %s", svc.BroadcastTopicName, err.Error()))
	}
}

func (svc *Service) updateConfigFromRefName(ctx context.Context, refName string) (string, error) {
	if refName == "" {
		return "", E.MakeGRPCStatus(codes.InvalidArgument, "ref name must be specified").Err()
	}

	svc.Log.Info(ctx, L.Messagef("UpdateConfig: %s", refName))

	svc.mu.Lock()
	defer svc.mu.Unlock()

	if svc.currentRepo == nil {
		svc.Log.Info(ctx, L.Messagef("nothing cloned yet; cloning %s", svc.RepoName))
		err := svc.cloneRepo(svc.RepoName)
		if err != nil {
			return "", err
		}
	}

	svc.Log.Info(ctx, L.Messagef("root: %s", svc.fs.Root()))
	svc.Log.Info(ctx, L.Messagef("checking out: %s", refName))

	err := svc.checkout("origin", refName)
	if err != nil {
		return "", err
	}

	svc.Log.Info(ctx, L.Messagef("switched: config ref: %s", refName))

	head, err := svc.currentRepo.Head()
	if err != nil {
		return "", err
	}

	return head.Hash().String(), nil
}

func (svc *Service) updateConfigFromCannedConfig(ctx context.Context, canned *pbConfig.CannedConfig) (string, error) {
	svc.currentRepo = nil
	svc.fs = memfs.New()

	tr := tar.NewReader(bytes.NewBuffer(canned.TarData))
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // end of archive
		}
		if err != nil {
			return "", err
		}
		svc.Log.Info(ctx, L.Messagef("- %s (%v)", hdr.Name, hdr.FileInfo().IsDir()))

		if hdr.FileInfo().IsDir() {
			continue
		}

		//nolint:gomnd
		f, err := svc.fs.OpenFile(hdr.Name, os.O_CREATE|os.O_WRONLY, 0o0666)
		if err != nil {
			return "", err
		}
		//nolint:gosec
		_, err = io.Copy(f, tr)
		if err != nil {
			return "", err
		}
		_ = f.Close()
	}

	return canned.ConfigHash, nil
}
