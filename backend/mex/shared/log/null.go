package log

import (
	"context"
	"net/http"
	"time"
)

type NullLogger struct{}

func (*NullLogger) Error(_ context.Context, _ ...Opt) {}
func (*NullLogger) Warn(_ context.Context, _ ...Opt)  {}
func (*NullLogger) Info(_ context.Context, _ ...Opt)  {}
func (*NullLogger) Trace(_ context.Context, _ ...Opt) {}

func (*NullLogger) BIEvent(ctx context.Context, opts ...Opt) {}

func (*NullLogger) LogHTTPRequest(r *http.Request) {}
func (*NullLogger) LogHTTPResponse(r *http.Request, status int, bodySize int, startTime time.Time, logRequestURL bool) {
}
