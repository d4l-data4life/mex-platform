package web

import (
	"context"
	"net/http"
	"strconv"

	"github.com/d4l-data4life/mex/mex/shared/constants"
	"github.com/d4l-data4life/mex/mex/shared/uuid"
)

type HandlerWrapper func(h http.Handler) http.Handler

func NewRequestIDMiddleware() HandlerWrapper {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), constants.ContextKeyRequestID, uuid.MustNewV4())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func NewTraceIDMiddleware(headerName string) HandlerWrapper {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), constants.ContextKeyTraceID, r.Header.Get(headerName))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func NewTracingDecisionMiddleware(globalTraceEnabled bool, traceSecret string) HandlerWrapper {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var ctx context.Context

			// We could stick a proper bool into the context (rather than a string).
			// However, we will feed this value to the gRPC service via its metadata key-value pairs which only accept string data.
			// Since we access this below context value at different locations in the code, we would need to pay attention to which context and
			// thus which type it is. Having it a string everywhere makes things easier (but a little less nice).
			if globalTraceEnabled {
				ctx = context.WithValue(r.Context(), constants.ContextKeyTraceThis, "true")
			} else {
				ctx = context.WithValue(r.Context(), constants.ContextKeyTraceThis, strconv.FormatBool(r.Header.Get(constants.HTTPHeaderTraceSecret) == traceSecret))
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func NewResponseHeaderMiddleware(next http.Handler, header string, value string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(header, value)
		next.ServeHTTP(w, r)
	})
}
