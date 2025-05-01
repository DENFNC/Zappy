package interceptor

import (
	"context"
	"log/slog"
	"runtime/debug"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewPanicInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		startTime := time.Now()
		defer func() {
			if r := recover(); r != nil {
				log.Error(
					"PANIC recovered",
					slog.Any("panic", r),
					slog.String("stack", string(debug.Stack())),
					slog.String("method", info.FullMethod),
					slog.Duration("duration", time.Since(startTime)),
				)
				err = status.Error(codes.Internal, "internal server error")
			}
		}()

		resp, err = handler(ctx, req)

		return
	}
}
