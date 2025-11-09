package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	handler "github.com/AchilleasB/southpark-messaging/go-api-service/api"
	"github.com/AchilleasB/southpark-messaging/go-api-service/internal/adapters/rabbitmq"
	"github.com/AchilleasB/southpark-messaging/go-api-service/internal/app"
	"github.com/AchilleasB/southpark-messaging/go-api-service/internal/server"
)

func main() {
	amqpURL := os.Getenv("RABBITMQ_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@localhost:5672/"
	}

	// Create concrete RabbitMQ publisher (adapter)
	pub, err := rabbitmq.NewPublisher(amqpURL, "southpark_messages")
	if err != nil {
		log.Fatalf("failed to create publisher: %v", err)
	}
	// ensure adapter resources are closed on exit
	defer pub.Close()

	// Create application core, injecting the publisher port implementation
	application := app.NewApp(pub)

	// Create HTTP handler layer (adapts HTTP -> application)
	h := handler.NewHandler(application)

	// Create server (accepts handler)
	srv := server.New(":8081", h)

	// Start server in background
	go func() {
		if err := srv.Start(); err != nil {
			log.Printf("server stopped: %v", err)
		}
	}()

	// Wait for interrupt signal and gracefully shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
}
