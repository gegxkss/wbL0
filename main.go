package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

var kafkaAddresses = []string{"localhost:9091", "localhost:9092"}

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

	go simulateOrders(cache, config.DB) //////////////

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
func simulateOrders(cache *cache.Cache, db *gorm.DB) {
	time.Sleep(5 * time.Second)

	log.Println("Creating demo orders for testing...")
	testOrder1 := models.Order{
		OrderUID:          "b563feb7b2b84b6test",
		TrackNumber:       "WBILMTESTTRACK",
		Entry:             "WBIL",
		Locale:            "en",
		InternalSignature: "",
		CustomerId:        "test",
		DeliveryService:   "meest",
		ShardKey:          "9",
		SmId:              99,
		DateCreated:       time.Now(),
		OofShard:          "1",
	}

	if err := cache.Set(testOrder1.OrderUID, &testOrder1); err != nil {
		log.Printf("Failed to add demo order to cache: %v", err)
	} else {
		log.Printf("Demo order 1 created: %s", testOrder1.OrderUID)
	}

	testOrder2 := models.Order{
		OrderUID:          "test-order-456",
		TrackNumber:       "TESTTRACK456",
		Entry:             "TEST",
		Locale:            "ru",
		InternalSignature: "",
		CustomerId:        "demo",
		DeliveryService:   "pochta",
		ShardKey:          "5",
		SmId:              50,
		DateCreated:       time.Now(),
		OofShard:          "2",
	}

	if err := cache.Set(testOrder2.OrderUID, &testOrder2); err != nil {
		log.Printf("Failed to add demo order to cache: %v", err)
	} else {
		log.Printf("Demo order 2 created: %s", testOrder2.OrderUID)
	}

	log.Println("Demo orders ready:")
	log.Println("   - b563feb7b2b84b6test")
	log.Println("   - test-order-456")
}

func waitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down...")
}
