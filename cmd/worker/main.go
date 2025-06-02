package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pdrhp/ms-voto-processor-go/internal/config"
)

func main() {
	cfg := config.Load()

	log.Printf("Starting vote processor worker in %s environment", cfg.App.Environment)
	log.Printf("Database: %s:%s/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.Database)
	log.Printf("Kafka: %v, Topic: %s", cfg.Kafka.Brokers, cfg.Kafka.Topic)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Worker started")
	<-quit
	log.Println("Shutting down worker...")
}
