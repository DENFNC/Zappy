package interceptor

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func TimingInterceptor(ctx context.Context, logger *slog.Logger) grpc.UnaryServerInterceptor {
	const op = "interceptor.TimingInterceptor"

	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		log := logger.With(slog.String("op", op))

		start := time.Now()
		resp, err = handler(ctx, req)
		duration := time.Since(start).Round(time.Millisecond)
		st, _ := status.FromError(err)

		log.InfoContext(
			ctx,
			"gRPC request",
			slog.String("method", info.FullMethod),
			slog.String("status", st.Code().String()),
			slog.String("duration", duration.String()),
		)

		return resp, err
	}
}
