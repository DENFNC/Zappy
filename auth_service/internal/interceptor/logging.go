package interceptor

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
)

func NewLoggingInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		timeStart := time.Now()

		resp, err := handler(ctx, req)

		if err != nil {
			log.Error(
				"Error processing request",
				slog.String("method", info.FullMethod),
				slog.Any("error", err.Error()),
				slog.Duration("time_ms", time.Duration(time.Since(timeStart).Microseconds())),
			)
		} else {
			log.Info(
				"Request processed",
				slog.String("method", info.FullMethod),
				slog.Duration("time_ms", time.Duration(time.Since(timeStart).Microseconds())),
			)
		}
		return resp, err
	}
}
