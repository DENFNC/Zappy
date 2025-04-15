package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type ServiceRegistrar interface {
	Register(server *grpc.Server)
}

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(
	log *slog.Logger,
	port int,
	services ...ServiceRegistrar,
) *App {
	server := grpc.NewServer()

	for _, service := range services {
		service.Register(server)
	}

	return &App{
		log:        log,
		gRPCServer: server,
		port:       port,
	}
}

func (a *App) MustRun() {
	const op = "grpcapp.App.MustRun"

	log := a.log.With("op", op)
	if err := a.start(); err != nil {
		log.Error("Failed to start gRPC server", "error", err)
		panic(fmt.Errorf("failed to start gRPC server: %w", err))
	}
}

func (a *App) Stop() {
	const op = "grpcapp.App.Stop"

	log := a.log.With("op", op)
	log.Info("Stopping gRPC server")
	a.gRPCServer.GracefulStop()
}

func (a *App) start() error {
	const op = "grpcapp.App.start"

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
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
