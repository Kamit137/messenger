package main

import (
	"fmt"
	"log"
	"messenger/internal/app"
	"messenger/internal/config"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := *config.MustLoad()
	fmt.Println(cfg)
	logger := log.New(os.Stdout, "[GRPC] ", log.LstdFlags|log.Lshortfile)

	application := app.New(logger, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)
	go application.GRPCSERVER.Run()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	signal := <-stop
	application.GRPCSERVER.Stop(signal)
	log.Print("ALL STOPED")
}
