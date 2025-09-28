package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gegxkss/wbL0/internal/models"
	"github.com/gegxkss/wbL0/kafka"
	"github.com/google/uuid"
)

func main() {
	kafkaAddresses := []string{"localhost:9091", "localhost:9092", "localhost:9093"}

	producer, err := kafka.NewProducer(kafkaAddresses)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	log.Println("Producer started. Sending test orders to Kafka...")
	counter := 1

	for {
		order := generateTestOrder(counter)
		message, _ := json.Marshal(order)

		err = producer.Produce(string(message), "order", order.OrderUID)
		if err != nil {
			log.Printf("Error producing message: %v", err)
		} else {
			log.Printf("Order %d sent: %s", counter, order.OrderUID)
		}

		counter++
		time.Sleep(5 * time.Second)
	}
}

func generateTestOrder(id int) models.Order {
	orderUID := uuid.New().String()
	products := []string{"Mascaras", "Lipstick", "Eyeshadow", "Foundation", "Concealer"}
	brands := []string{"Vivienne Sabo", "Darling", "RAD", "Maybelline", "Influence", "L'Oreal", "NYX", "MAC"}

	productIndex := rand.Intn(len(products))
	brandIndex := rand.Intn(len(brands))

	return models.Order{
		OrderUID:          orderUID,
		TrackNumber:       fmt.Sprintf("WBILMTESTTRACK%d", id),
		Entry:             "WBIL",
		Locale:            "en",
		InternalSignature: "",
		CustomerId:        "test",
		DeliveryService:   "meest",
		ShardKey:          "9",
		SmId:              99,
		DateCreated:       time.Now(),
		OofShard:          "1",
		Delivery: models.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: models.Payment{
			Transaction:  orderUID,
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817 + id*10,
			PaymentDt:    int(time.Now().Unix()),
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317 + id*5,
		},
		Items: []models.Items{
			{
				ChrtId:      9934930 + id,
				Tracknumber: fmt.Sprintf("WBILMTESTTRACK%d", id),
				Price:       453 + id*2,
				Rid:         fmt.Sprintf("ab4219087a764ae0btest%d", id),
				Name:        products[productIndex],
				Sale:        30,
				Size:        "0",
				TotalPrice:  317 + id*5,
				NmId:        2389212 + id,
				Brand:       brands[brandIndex],
				Status:      202,
			},
		},
	}
}
