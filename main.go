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
	"github.com/gegxkss/wbL0/internal/models"
	"github.com/gegxkss/wbL0/kafka"
	"github.com/gegxkss/wbL0/migrations"
	"gorm.io/gorm"
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
	cache := cache.NewCache()
	restoreCacheFromDB(config.DB, cache)

	consumer, _ := kafka.NewConsumer(kafkaAddresses, topic, groupID, config.DB, cache)
	go consumer.Start()
	defer consumer.Stop()

	handlers.SetupRoutes(cache, config.DB)

	go func() {
		log.Println("Starting HTTP server on :8081")
		http.ListenAndServe(":8081", nil)
	}()

	log.Println("Application started successfully!")
	waitForShutdown()
}

func restoreCacheFromDB(db *gorm.DB, cache *cache.Cache) {
	log.Println("Restoring cache from database...")

	var orders []models.Order
	db.Preload("Delivery").Preload("Payment").Preload("Items").Limit(1000).Find(&orders)

	for _, order := range orders {
		cache.Set(order.OrderUID, &order)
	}

	log.Printf("Cache restored with %d orders", len(orders))
}

func waitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down...")
}
