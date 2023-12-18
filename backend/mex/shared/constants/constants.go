package constants

import "context"

type keyType string

const (
	ContextKeyTraceSecret keyType = "mex-trace-secret"
	ContextKeyTraceThis   keyType = "mex-trace-this"
	ContextKeyTraceID     keyType = "mex-trace-id"
	ContextKeyUserID      keyType = "mex-user-id"
	ContextKeyUserClaims  keyType = "mex-user-claims"
	ContextKeyJobID       keyType = "mex-job-id"
	ContextKeyTenantID    keyType = "mex-tenant-id"
	ContextKeyRequestID   keyType = "mex-request-id"
)

const (
	HTTPHeaderJobID       = "x-mex-job-id"
	HTTPHeaderTraceID     = "trace-id"
	HTTPHeaderTraceSecret = "x-mex-trace"
	HTTPHeaderRealIP      = "x-real-ip"
)

func GetContextValueDefault[T any](ctx context.Context, key interface{}, defaultValue T) T {
	if ctx == nil {
		return defaultValue
	}

	v := ctx.Value(key)
	if v == nil {
		return defaultValue
	}

	if s, ok := v.(T); ok {
		return s
	}
	return defaultValue
}

func NewContextWithValues(ctxSource context.Context, jobID string) context.Context {
	ctxNew := context.Background()

	ctxNew = context.WithValue(ctxNew, ContextKeyJobID, jobID)
	ctxNew = context.WithValue(ctxNew, ContextKeyRequestID, ctxSource.Value(ContextKeyRequestID))
	ctxNew = context.WithValue(ctxNew, ContextKeyTraceID, ctxSource.Value(ContextKeyTraceID))
	ctxNew = context.WithValue(ctxNew, ContextKeyTenantID, ctxSource.Value(ContextKeyTenantID))
	ctxNew = context.WithValue(ctxNew, ContextKeyUserID, ctxSource.Value(ContextKeyUserID))
	ctxNew = context.WithValue(ctxNew, ContextKeyUserClaims, ctxSource.Value(ContextKeyUserClaims))
	ctxNew = context.WithValue(ctxNew, ContextKeyTraceSecret, ctxSource.Value(ContextKeyTraceSecret))

	return ctxNew
}

const (
	ConfigUpdateChannelNameSuffix = "mex-config-update"
	MergedItemToFragmentRelation  = "dataOriginatesFrom"
)
