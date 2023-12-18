package jobs

import (
	"context"

	"github.com/d4l-data4life/mex/mex/shared/known/jobspb"
)

type Jobber interface {
	CreateJob(ctx context.Context, request *jobspb.CreateJobRequest) (*jobspb.CreateJobResponse, error)
	GetJob(ctx context.Context, request *jobspb.GetJobRequest) (*jobspb.GetJobResponse, error)
	JobExists(ctx context.Context, jobID string) bool

	GetLogs(ctx context.Context, request *jobspb.GetJobLogsRequest) (*jobspb.GetJobLogsResponse, error)
	AddLogs(ctx context.Context, request *jobspb.AddJobLogsRequest) (*jobspb.AddJobLogsResponse, error)

	GetItems(ctx context.Context, request *jobspb.GetJobItemsRequest) (*jobspb.GetJobItemsResponse, error)
	AddItems(ctx context.Context, request *jobspb.AddJobItemsRequest) (*jobspb.AddJobItemsResponse, error)

	SetStatus(ctx context.Context, request *jobspb.SetJobStatusRequest) (*jobspb.SetJobStatusResponse, error)
	SetStatusRunning(ctx context.Context, jobID string) error
	SetStatusDone(ctx context.Context, jobID string) error

	SetError(ctx context.Context, request *jobspb.SetJobErrorRequest) (*jobspb.SetJobErrorResponse, error)

	AcquireLock(ctx context.Context, resourceName string) (string, error)
	ReleaseLock(ctx context.Context, resourceName string, handle string) error
}
