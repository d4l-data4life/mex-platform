package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/d4l-data4life/mex/mex/shared/uuid"

	"github.com/d4l-data4life/mex/mex/shared/known/jobspb"
)

const (
	PropTitle     = "title"
	PropStatus    = "status"
	PropError     = "error"
	PropCreatedAt = "createdAt"

	StatusCreated = "CREATED"
	StatusRunning = "RUNNING"
	StatusDone    = "DONE"
)

var scriptReleaseLock = redis.NewScript(`
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end
`)

func hashNameJobs(jobID string) string {
	return fmt.Sprintf("jobs:%s:data", jobID)
}

func hashNameJobLogs(jobID string) string {
	return fmt.Sprintf("jobs:%s:logs", jobID)
}

func hashNameJobItems(jobID string) string {
	return fmt.Sprintf("jobs:%s:items", jobID)
}

func hashNameLocks(s string) string {
	return fmt.Sprintf("locks:%s", s)
}

type RedisJobber struct {
	Redis      *redis.Client
	Expiration time.Duration
}

func (j RedisJobber) CreateJob(ctx context.Context, request *jobspb.CreateJobRequest) (*jobspb.CreateJobResponse, error) {
	jobID := uuid.MustNewV4()
	hashName := hashNameJobs(jobID)
	hashNameLogs := hashNameJobLogs(jobID)
	hashNameItems := hashNameJobItems(jobID)

	setCmd := j.Redis.HMSet(ctx, hashName, map[string]interface{}{
		PropTitle:     request.Title,
		PropStatus:    StatusCreated,
		PropError:     "",
		PropCreatedAt: time.Now().Format(time.RFC3339),
	})
	if setCmd.Err() != nil {
		return nil, setCmd.Err()
	}

	pushCmd := j.Redis.RPush(ctx, hashNameLogs, "{}")
	if pushCmd.Err() != nil {
		return nil, pushCmd.Err()
	}

	pushCmd = j.Redis.RPush(ctx, hashNameItems, "{}")
	if pushCmd.Err() != nil {
		return nil, pushCmd.Err()
	}

	j.Redis.Expire(ctx, hashName, j.Expiration)
	j.Redis.Expire(ctx, hashNameLogs, j.Expiration)
	j.Redis.Expire(ctx, hashNameItems, j.Expiration)

	return &jobspb.CreateJobResponse{JobId: jobID}, nil
}

func (j RedisJobber) GetLogs(ctx context.Context, request *jobspb.GetJobLogsRequest) (*jobspb.GetJobLogsResponse, error) {
	if !j.JobExists(ctx, request.JobId) {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("job not found: %s", request.JobId))
	}

	hashNameLogs := hashNameJobLogs(request.JobId)

	cmdSlice := j.Redis.LRange(ctx, hashNameLogs, 0, -1)
	if cmdSlice.Err() != nil {
		return nil, cmdSlice.Err()
	}

	return &jobspb.GetJobLogsResponse{
		JobId: request.JobId,
		Logs:  cmdSlice.Val(),
	}, nil
}

func (j RedisJobber) AddLogs(ctx context.Context, request *jobspb.AddJobLogsRequest) (*jobspb.AddJobLogsResponse, error) {
	if !j.JobExists(ctx, request.JobId) {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("job not found: %s", request.JobId))
	}

	hashNameLogs := hashNameJobLogs(request.JobId)
	var count int32

	for _, log := range request.Logs {
		cmd := j.Redis.RPush(ctx, hashNameLogs, log)
		if cmd.Err() != nil {
			return nil, cmd.Err()
		}
		count = int32(cmd.Val())
	}

	return &jobspb.AddJobLogsResponse{
		JobId:    request.JobId,
		LogCount: count,
	}, nil
}

func (j RedisJobber) GetItems(ctx context.Context, request *jobspb.GetJobItemsRequest) (*jobspb.GetJobItemsResponse, error) {
	if !j.JobExists(ctx, request.JobId) {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("job not found: %s", request.JobId))
	}

	hashNameItems := hashNameJobItems(request.JobId)

	cmdSlice := j.Redis.LRange(ctx, hashNameItems, 0, -1)
	if cmdSlice.Err() != nil {
		return nil, cmdSlice.Err()
	}

	return &jobspb.GetJobItemsResponse{
		JobId:   request.JobId,
		ItemIds: cmdSlice.Val(),
	}, nil
}

func (j RedisJobber) AddItems(ctx context.Context, request *jobspb.AddJobItemsRequest) (*jobspb.AddJobItemsResponse, error) {
	if !j.JobExists(ctx, request.JobId) {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("job not found: %s", request.JobId))
	}

	hashNameItems := hashNameJobItems(request.JobId)
	var count int32

	for _, itemID := range request.ItemIds {
		cmd := j.Redis.RPush(ctx, hashNameItems, itemID)
		if cmd.Err() != nil {
			return nil, cmd.Err()
		}
		count = int32(cmd.Val())
	}

	return &jobspb.AddJobItemsResponse{
		JobId:     request.JobId,
		ItemCount: count,
	}, nil
}

func (j RedisJobber) SetStatus(ctx context.Context, request *jobspb.SetJobStatusRequest) (*jobspb.SetJobStatusResponse, error) {
	if !j.JobExists(ctx, request.JobId) {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("job not found: %s", request.JobId))
	}

	hashName := hashNameJobs(request.JobId)

	cmd := j.Redis.HSet(ctx, hashName, PropStatus, request.Status)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	return &jobspb.SetJobStatusResponse{
		JobId:  request.JobId,
		Status: request.Status,
	}, nil
}

func (j RedisJobber) SetError(ctx context.Context, request *jobspb.SetJobErrorRequest) (*jobspb.SetJobErrorResponse, error) {
	if !j.JobExists(ctx, request.JobId) {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("job not found: %s", request.JobId))
	}

	hashName := hashNameJobs(request.JobId)

	cmd := j.Redis.HSet(ctx, hashName, PropError, request.Error)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	return &jobspb.SetJobErrorResponse{
		JobId: request.JobId,
	}, nil
}

func (j RedisJobber) GetJob(ctx context.Context, request *jobspb.GetJobRequest) (*jobspb.GetJobResponse, error) {
	if !j.JobExists(ctx, request.JobId) {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("job not found: %s", request.JobId))
	}

	hashName := hashNameJobs(request.JobId)
	hashNameLogs := hashNameJobLogs(request.JobId)

	cmd := j.Redis.HGetAll(ctx, hashName)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	var response jobspb.GetJobResponse
	response.JobId = request.JobId

	response.Title = cmd.Val()[PropTitle]
	response.CreatedAt = cmd.Val()[PropCreatedAt]
	response.Status = cmd.Val()[PropStatus]
	response.Error = cmd.Val()[PropError]

	cmdLength := j.Redis.LLen(ctx, hashNameLogs)
	if cmdLength.Err() != nil {
		return nil, cmdLength.Err()
	}

	response.LogCount = int32(cmdLength.Val())

	return &response, nil
}

func (j RedisJobber) AcquireLock(ctx context.Context, resourceName string) (string, error) {
	key := hashNameLocks(resourceName)
	handle := uuid.MustNewV4()

	cmd := j.Redis.SetNX(ctx, key, handle, j.Expiration)
	if cmd.Err() != nil {
		return "", fmt.Errorf("error acquiring lock for resource %s (%w)", resourceName, cmd.Err())
	}

	if cmd.Val() {
		return handle, nil
	}

	return "", fmt.Errorf("could not acquire lock for resource %s", resourceName)
}

func (j RedisJobber) ReleaseLock(ctx context.Context, resourceName string, handle string) error {
	key := hashNameLocks(resourceName)

	_, err := scriptReleaseLock.Run(ctx, j.Redis, []string{key}, handle).Result()
	if err != nil {
		return fmt.Errorf("error releasing lock for resource %s (%w)", resourceName, err)
	}

	return nil
}

func (j RedisJobber) JobExists(ctx context.Context, jobID string) bool {
	hashName := hashNameJobs(jobID)

	cmd := j.Redis.Exists(ctx, hashName)
	if cmd.Err() != nil {
		return false
	}

	return cmd.Val() == 1
}

func (j RedisJobber) SetStatusRunning(ctx context.Context, jobID string) error {
	if !j.JobExists(ctx, jobID) {
		return fmt.Errorf("job not found: %s", jobID)
	}

	_, err := j.SetStatus(ctx, &jobspb.SetJobStatusRequest{JobId: jobID, Status: StatusRunning})
	return err
}

func (j RedisJobber) SetStatusDone(ctx context.Context, jobID string) error {
	if !j.JobExists(ctx, jobID) {
		return fmt.Errorf("job not found: %s", jobID)
	}

	_, err := j.SetStatus(ctx, &jobspb.SetJobStatusRequest{JobId: jobID, Status: StatusDone})
	return err
}
