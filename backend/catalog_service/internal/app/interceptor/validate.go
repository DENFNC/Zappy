package interceptor

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type validator interface {
	Validate() error
}

func ValidateArgsInterceptor(ctx context.Context, log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		const op = "interceptor.ValidateArgsInterceptor"

		log = log.With("op", op)

		if val, ok := req.(validator); ok {
			if err := val.Validate(); err != nil {
				log.Debug(
					"Invalid argument",
					slog.String("method", info.FullMethod),
				)
				return nil, status.Error(
					codes.InvalidArgument,
					"invalid argument",
				)
			}
		}

		return handler(ctx, req)
	}
}
