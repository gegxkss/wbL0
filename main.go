package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gegxkss/wbL0/internal/cache"
	"github.com/gegxkss/wbL0/internal/config"
	"github.com/gegxkss/wbL0/internal/handlers"
	"github.com/gegxkss/wbL0/kafka"
	"github.com/gegxkss/wbL0/migrations"
)

var migrate = flag.Bool("m", false, "Run database migrations")

const (
	topic   = "order"
	groupID = "orders-group"
)

var kafkaAddresses = []string{"localhost:9091", "localhost:9092", "localhost:9093"}

func main() {
	flag.Parse()

	if *migrate {
		migrations.Migration()
		log.Println("Migrations completed. Exiting.")
		return
	}
	config.ConnectDB()
	log.Println("Database connection successful!")
	log.Println("Add your Kafka logic here later.")

	cache := cache.NewCache()
	log.Printf("Cache initialized (empty)")

	consumer, err := kafka.NewConsumer(kafkaAddresses, topic, groupID, config.DB, cache)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}

	go func() {
		log.Println("Starting Kafka consumer...")
		consumer.Start()
	}()
	defer consumer.Stop()

	handlers.SetupRoutes(cache, config.DB)

	go func() {
		log.Println("Starting HTTP server on :8081")
		if err := http.ListenAndServe(":8081", nil); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	log.Println("Application started successfully!")
	log.Println("Kafka consumer is listening for orders...")
	log.Println("HTTP API is available on http://localhost:8081")

	waitForShutdown()
}

func waitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down...")
}
