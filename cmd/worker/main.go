package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pdrhp/ms-voto-processor-go/internal/config"
	"github.com/pdrhp/ms-voto-processor-go/internal/container"
)

func main() {
	cfg := config.Load()

	app := container.NewContainer(cfg)
	defer app.Close()

	if err := app.Build(); err != nil {
		log.Fatalf("Failed to build application: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := app.Start(ctx); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Worker started")
	<-quit

	cancel()
	app.Stop()

	log.Println("Worker stopped")
}
