package hints

import (
	"context"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func HintHTTPStatusCode(ctx context.Context, code int) {
	_ = grpc.SetHeader(ctx, metadata.Pairs("Mex-Hinted-Status-Code", strconv.Itoa(code)))
}

func HintHTTPLocation(ctx context.Context, location string) {
	_ = grpc.SetHeader(ctx, metadata.Pairs("Mex-Hinted-Location", location))
}

func HintHTTPCacheControl(ctx context.Context, value string) {
	_ = grpc.SetHeader(ctx, metadata.Pairs("Mex-Hinted-Cache-Control", value))
}

func HintHTTPPragma(ctx context.Context, value string) {
	_ = grpc.SetHeader(ctx, metadata.Pairs("Mex-Hinted-Pragma", value))
}
