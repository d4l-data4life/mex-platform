package index

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/d4l-data4life/mex/mex/shared/errstat"
	"github.com/d4l-data4life/mex/mex/shared/hints"

	"github.com/d4l-data4life/mex/mex/services/index/endpoints/index/pb"
)

// DeleteIndex deletes the Solr index (data stored in Solr) without changing the Solr schema
func (svc *Service) DeleteIndex(ctx context.Context, _ *pb.DeleteIndexRequest) (*pb.DeleteIndexResponse, error) {
	lock, err := svc.JobService.AcquireLock(ctx, SvcResourceName)
	if err != nil {
		return nil, errstat.MakeGRPCStatus(codes.AlreadyExists, "delete: failed to acquire index lock; other job might be running").Err()
	}

	defer svc.JobService.ReleaseLock(ctx, SvcResourceName, lock) //nolint:errcheck

	err = svc.Solr.DropIndex(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failure when trying to delete existing Solr data: %s", err.Error()))
	}

	hints.HintHTTPStatusCode(ctx, http.StatusNoContent)
	return &pb.DeleteIndexResponse{}, nil
}
