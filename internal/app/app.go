package app

import (
	"log"
	grpcapp "messenger/internal/app/grpc"
	"time"
)

type App struct {
	GRPCSERVER *grpcapp.App
}

func New(log *log.Logger, grpcPort int, sroragePath string, tokenTTL time.Duration) *App {
	grpcApp := grpcapp.New(log, grpcPort)
	return &App{
		GRPCSERVER: grpcApp,
	}
}
