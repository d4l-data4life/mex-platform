package interceptors

import (
	"context"
	"net/http"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func NewHeaderInterceptor(key string, key2 any) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			value := md.Get(key)
			if len(value) > 0 {
				return handler(context.WithValue(ctx, key2, value[0]), req)
			}
		}

		return handler(ctx, req)
	}
}

func HeaderReader(key string, header string) func(context.Context, *http.Request) metadata.MD {
	return func(ctx context.Context, r *http.Request) metadata.MD {
		return metadata.Pairs(key, r.Header.Get(header))
	}
}

func ContextConnect(key string, key2 any) func(context.Context, *http.Request) metadata.MD {
	return func(ctx context.Context, r *http.Request) metadata.MD {
		// could also used ctx, but r.Context() is more explicit
		value := r.Context().Value(key2)
		if value == nil {
			return metadata.Pairs()
		}

		switch v := value.(type) {
		case string:
			return metadata.Pairs(key, v)
		case bool:
			return metadata.Pairs(key, strconv.FormatBool(v))
		default:
			panic("unsupported context value type")
		}
	}
}
