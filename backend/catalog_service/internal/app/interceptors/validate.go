package interceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type validator interface {
	Validate() error
}

func ValidateArgsInterceptor(ctx context.Context) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		if val, ok := req.(validator); ok {
			if err := val.Validate(); err != nil {
				return nil, status.Error(
					codes.InvalidArgument,
					"invalid argument",
				)
			}
		}

		return handler(ctx, req)
	}
}
