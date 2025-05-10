// TODO: Прописать всю логику по человечески
// TODO: так-же когда буду переходить на http дополнить вспомогательными функциями

package grpcapp

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ServiceRegistrar interface {
	GRPCRegister(server *grpc.Server)
	// HTTPRegister(
	// 	ctx context.Context,
	// 	mux *runtime.ServeMux,
	// )
}

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	// httpServer *runtime.ServeMux
	gRPCPort int
	// httpPort   int
}

func New(
	ctx context.Context,
	log *slog.Logger,
	grpcPort int,
	// httpPort int,
	services ...ServiceRegistrar,
) *App {

	grpcServer := grpc.NewServer(
	//TODO: Прописать нужные интерцепторы
	)

	reflection.Register(grpcServer)

	return &App{
		log:        log,
		gRPCServer: grpcServer,
		gRPCPort:   grpcPort,
	}
}

func (a *App) MustRunGrpc() {
	const op = "grpcapp.App.MustRunGrpc"

	log := a.log.With("op", op)
	if err := a.gRPCstart(); err != nil {
		log.Error(
			"Failed to start gRPC server",
			"error", err,
		)
		panic(fmt.Errorf("failed to start gRPC server: %w", err))
	}
}

func (a *App) gRPCstart() error {
	const op = "grpcapp.App.gRPCstart"

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.gRPCPort))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	log := a.log.With("op", op)
	log.Info(
		"Starting gRPC server",
		"addr", lis.Addr().String(),
	)

	if err := a.gRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to gRPC serve: %w", err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.App.Stop"

	log := a.log.With("op", op)
	log.Info("Stopping gRPC server")

	a.gRPCServer.GracefulStop()
}
