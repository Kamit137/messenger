package grpcapp

import (
	"fmt"
	"log"
	"net"
	"os"

	authgrpc "messenger/internal/grpc"

	"google.golang.org/grpc"
)

type App struct {
	logger     *log.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(logger *log.Logger, port int) *App {
	gRPCServer := grpc.NewServer()
	authgrpc.RegServer(gRPCServer)

	return &App{
		logger:     logger,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"
	a.logger.Printf("%s: starting on port %d", op, a.port)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		a.logger.Fatalf("%s: failed to listen: %v", op, err)
	}

	a.logger.Printf("%s: server is listening", op)

	if err := a.gRPCServer.Serve(lis); err != nil {
		a.logger.Fatalf("%s: server stopped with error: %v", op, err)
	}

	return nil
}

func (a *App)Stop(signal os.Signal){

	log.Printf("Stop grpc server graceful: %v", signal)
	a.gRPCServer.GracefulStop()
}
