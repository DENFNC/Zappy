package auth

import (
	"github.com/DENFNC/Zappy/proto/gen/v1/v1"
	"google.golang.org/grpc"
)

type serverAPI struct {
	v1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	v1.RegisterAuthServer(gRPC, &serverAPI{})
}
