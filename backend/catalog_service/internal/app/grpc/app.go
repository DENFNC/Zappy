package grpcapp

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ServiceRegistrar interface {
	GRPCRegister(server *grpc.Server)
	HTTPRegister(
		ctx context.Context,
		mux *runtime.ServeMux,
	)
}

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	httpServer *runtime.ServeMux
	grpcPort   int
	httpPort   int
}

func New(
	ctx context.Context,
	log *slog.Logger,
	grpcPort int,
	httpPort int,
	services ...ServiceRegistrar,
) *App {
	mux := runtime.NewServeMux()

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	for _, service := range services {
		service.GRPCRegister(grpcServer)
		service.HTTPRegister(ctx, mux)
	}

	return &App{
		log:        log,
		gRPCServer: grpcServer,
		httpServer: mux,
		grpcPort:   grpcPort,
		httpPort:   httpPort,
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

func (a *App) MustRunHttp() {
	const op = "grpcapp.App.MustRunHttp"

	log := a.log.With("op", op)
	if err := a.httpStart(); err != nil {
		log.Error(
			"Failed to start HTTP server",
			slog.String("error", err.Error()),
		)
		panic(err)
	}
}

func (a *App) Stop() {
	const op = "grpcapp.App.Stop"

	log := a.log.With("op", op)
	log.Info("Stopping gRPC server")

	a.gRPCServer.GracefulStop()
}

func (a *App) httpStart() error {
	const op = "grpcapp.App.httpStart"

	log := a.log.With("op", op)
	log.Info(
		"Starting HTTP server",
		"addr", fmt.Sprintf(":%d", a.httpPort),
	)

	if err := http.ListenAndServe(fmt.Sprintf("localhost:%d", a.httpPort), a.httpServer); err != nil {
		return err
	}

	return nil
}

func (a *App) gRPCstart() error {
	const op = "grpcapp.App.gRPCstart"

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.grpcPort))
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
