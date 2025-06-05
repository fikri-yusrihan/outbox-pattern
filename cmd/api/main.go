package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fikri-yusrihan/outbox-project/internal/consumer"
	"github.com/fikri-yusrihan/outbox-project/internal/db"
	"github.com/fikri-yusrihan/outbox-project/internal/publisher"
)

func main() {
	dbConn, err := db.Connect()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// Context to control shutdown of goroutines
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the Kafka publisher in a separate goroutine
	go publisher.StartKafkaPublisher(ctx, dbConn, "A")
	go publisher.StartKafkaPublisher(ctx, dbConn, "B")

	// Start the Kafka consumer in a separate goroutine
	go consumer.StartKafkaListener("kafka:9092", "user.created")
	go consumer.StartKafkaListener("kafka:9092", "user.created")

	// HTTP server setup
	server := &http.Server{
		Addr: ":8080",
		Handler: nil, // Use default handler
	}

	// Graceful shutdown listener
	idleConnsClosed := make(chan struct{})
	go func() {
		// Listen for termination signals
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig

		log.Println("Shutdown signal received...")

		// Cancel context → stop publisher
		cancel()

		// Graceful HTTP server shutdown
		ctxTimeout, shutdownCancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(ctxTimeout); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}

		close(idleConnsClosed)
	}()

	// Register HTTP handlers
	http.HandleFunc("/users", CreateUserHandler(dbConn))

	if err := http.ListenAndServe(":8080", nil); err != nil && err != http.ErrServerClosed {
		log.Fatal("Failed to start server:", err)
	}

	<- idleConnsClosed // Wait for shutdown to complete
	log.Println("Server gracefully stopped")
}
