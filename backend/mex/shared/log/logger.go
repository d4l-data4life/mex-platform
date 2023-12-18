package log

import (
	"context"

	"google.golang.org/grpc"
)

func NewLogInterceptor(log Logger, excludePatterns []string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// Bypass logs for excluded paths
		for _, pattern := range excludePatterns {
			if pattern == info.FullMethod {
				return handler(ctx, req)
			}
		}

		log.Info(ctx, Messagef(">>> %s", info.FullMethod))

		resp, err := handler(ctx, req)
		if err != nil {
			log.Warn(ctx, Message(err.Error()))
		}

		log.Info(ctx, Messagef("<<< %s", info.FullMethod))

		return resp, err
	}
}
