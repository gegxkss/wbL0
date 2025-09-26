package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gegxkss/wbL0/kafka"
	"github.com/google/uuid"
)

type OrderData struct {
	OrderUID          string   `json:"order_uid"`
	TrackNumber       string   `json:"track_number"`
	Entry             string   `json:"entry"`
	Delivery          Delivery `json:"delivery"`
	Payment           Payment  `json:"payment"`
	Items             []Item   `json:"items"`
	Locale            string   `json:"locale"`
	InternalSignature string   `json:"internal_signature"`
	CustomerId        string   `json:"customer_id"`
	DeliveryService   string   `json:"delivery_service"`
	ShardKey          string   `json:"shardkey"`
	SmId              int      `json:"sm_id"`
	DateCreated       string   `json:"date_created"`
	OofShard          string   `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	ChrtId      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmId        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func main() {
	kafkaAddresses := []string{"localhost:9091", "localhost:9092", "localhost:9093"}

	producer, err := kafka.NewProducer(kafkaAddresses)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	log.Println("🚀 Producer started. Sending test orders to Kafka...")
	log.Println("Press Ctrl+C to stop")

	// Бесконечный цикл отправки заказов
	counter := 1
	for {
		order := generateTestOrder(counter)

		message, err := json.Marshal(order)
		if err != nil {
			log.Printf("❌ Error marshaling order: %v", err)
			continue
		}

		err = producer.Produce(string(message), "order", order.OrderUID)
		if err != nil {
			log.Printf("❌ Error producing message: %v", err)
		} else {
			log.Printf("✅ Order %d sent: %s", counter, order.OrderUID)
		}

		counter++
		time.Sleep(5 * time.Second) // Отправка каждые 5 секунд
	}
}

func generateTestOrder(id int) OrderData {
	now := time.Now().Format(time.RFC3339)

	// Генерация случайных данных для разнообразия
	names := []string{"Test Testov", "Иван Иванов", "John Doe", "Maria Garcia"}
	cities := []string{"Kiryat Mozkin", "Moscow", "New York", "London"}
	products := []string{"Mascaras", "Laptop", "Smartphone", "Book", "Shoes"}

	return OrderData{
		OrderUID:          uuid.New().String(),
		TrackNumber:       fmt.Sprintf("WBILMTESTTRACK%d", id),
		Entry:             "WBIL",
		Locale:            "en",
		InternalSignature: "",
		CustomerId:        "test",
		DeliveryService:   "meest",
		ShardKey:          "9",
		SmId:              99,
		DateCreated:       now,
		OofShard:          "1",
		Delivery: Delivery{
			Name:    names[rand.Intn(len(names))],
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    cities[rand.Intn(len(cities))],
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: Payment{
			Transaction:  uuid.New().String(),
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       rand.Intn(5000) + 1000,
			PaymentDt:    int(time.Now().Unix()),
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   rand.Intn(500) + 100,
			CustomFee:    0,
		},
		Items: []Item{
			{
				ChrtId:      9934930 + id,
				TrackNumber: fmt.Sprintf("WBILMTESTTRACK%d", id),
				Price:       rand.Intn(1000) + 100,
				Rid:         uuid.New().String(),
				Name:        products[rand.Intn(len(products))],
				Sale:        rand.Intn(50),
				Size:        "0",
				TotalPrice:  rand.Intn(500) + 100,
				NmId:        2389212 + id,
				Brand:       "Brand " + string(rune(65+rand.Intn(26))), // A-Z
				Status:      202,
			},
		},
	}
}
