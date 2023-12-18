package jobs

import (
	"context"
	"net/http"

	E "github.com/d4l-data4life/mex/mex/shared/errstat"
	"github.com/d4l-data4life/mex/mex/shared/hints"
	"github.com/d4l-data4life/mex/mex/shared/jobs"
	"github.com/d4l-data4life/mex/mex/shared/known/jobspb"

	pbJobs "github.com/d4l-data4life/mex/mex/services/metadata/endpoints/jobs/pb"
)

type Service struct {
	jobs.Jobber

	pbJobs.UnimplementedJobsServer
}

func (svc Service) CreateJob(ctx context.Context, request *jobspb.CreateJobRequest) (*jobspb.CreateJobResponse, error) {
	response, err := svc.Jobber.CreateJob(ctx, request)
	if err != nil {
		return nil, E.MakeGRPCStatus(E.CodeFrom(err), "job creation failed", E.DevMessage(err.Error())).Err()
	}

	hints.HintHTTPStatusCode(ctx, http.StatusCreated)
	return response, nil
}

func (svc Service) GetLogs(ctx context.Context, request *jobspb.GetJobLogsRequest) (*jobspb.GetJobLogsResponse, error) {
	response, err := svc.Jobber.GetLogs(ctx, request)
	if err != nil {
		return nil, E.MakeGRPCStatus(E.CodeFrom(err), "failed to get job logs", E.DevMessage(err.Error())).Err()
	}

	return response, nil
}

func (svc Service) AddLogs(ctx context.Context, request *jobspb.AddJobLogsRequest) (*jobspb.AddJobLogsResponse, error) {
	response, err := svc.Jobber.AddLogs(ctx, request)
	if err != nil {
		return nil, E.MakeGRPCStatus(E.CodeFrom(err), "failed to add job logs", E.DevMessage(err.Error())).Err()
	}

	return response, nil
}

func (svc Service) GetItems(ctx context.Context, request *jobspb.GetJobItemsRequest) (*jobspb.GetJobItemsResponse, error) {
	response, err := svc.Jobber.GetItems(ctx, request)
	if err != nil {
		return nil, E.MakeGRPCStatus(E.CodeFrom(err), "failed to get job items", E.DevMessage(err.Error())).Err()
	}

	return response, nil
}

func (svc Service) AddItems(ctx context.Context, request *jobspb.AddJobItemsRequest) (*jobspb.AddJobItemsResponse, error) {
	response, err := svc.Jobber.AddItems(ctx, request)
	if err != nil {
		return nil, E.MakeGRPCStatus(E.CodeFrom(err), "failed to add job items", E.DevMessage(err.Error())).Err()
	}

	return response, nil
}

func (svc Service) SetStatus(ctx context.Context, request *jobspb.SetJobStatusRequest) (*jobspb.SetJobStatusResponse, error) {
	response, err := svc.Jobber.SetStatus(ctx, request)
	if err != nil {
		return nil, E.MakeGRPCStatus(E.CodeFrom(err), "failed to set job status", E.DevMessage(err.Error())).Err()
	}

	return response, nil
}

func (svc Service) SetError(ctx context.Context, request *jobspb.SetJobErrorRequest) (*jobspb.SetJobErrorResponse, error) {
	response, err := svc.Jobber.SetError(ctx, request)
	if err != nil {
		return nil, E.MakeGRPCStatus(E.CodeFrom(err), "failed to set job error", E.DevMessage(err.Error())).Err()
	}

	return response, nil
}

func (svc Service) GetJob(ctx context.Context, request *jobspb.GetJobRequest) (*jobspb.GetJobResponse, error) {
	response, err := svc.Jobber.GetJob(ctx, request)
	if err != nil {
		return nil, E.MakeGRPCStatus(E.CodeFrom(err), "failed to get job", E.DevMessage(err.Error())).Err()
	}

	return response, nil
}
