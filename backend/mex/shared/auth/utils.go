package auth

import (
	"context"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/d4l-data4life/mex/mex/shared/constants"
)

func GetMexUser(ctx context.Context) (*Claims, error) {
	v := ctx.Value(constants.ContextKeyUserClaims)
	if v == nil {
		return nil, status.New(codes.Unauthenticated, "no user claims in request").Err()
	}

	claims, ok := v.(*Claims)
	if !ok {
		return nil, status.New(codes.Unauthenticated, "invalid user claims in request").Err()
	}

	return claims, nil
}

func GetUserID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	rdbRaw := ctx.Value(constants.ContextKeyUserClaims)
	if rdbRaw == nil {
		return ""
	}

	rdb, ok := rdbRaw.(*Claims)
	if !ok {
		return ""
	}

	return rdb.UserId
}

func GetTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	rawTraceID := ctx.Value(constants.ContextKeyTraceID)
	if rawTraceID == nil {
		return ""
	}

	traceID, ok := rawTraceID.(string)
	if !ok {
		return ""
	}

	return traceID
}

func GetJobID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	rawJobID := ctx.Value(constants.ContextKeyJobID)
	if rawJobID == nil {
		return ""
	}

	jobID, ok := rawJobID.(string)
	if !ok {
		return ""
	}

	return jobID
}

func AuthnHeaderReader(ctx context.Context, r *http.Request) metadata.MD {
	md := metadata.Pairs("authorization", r.Header.Get("authorization"))
	return md
}
