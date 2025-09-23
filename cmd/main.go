package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

const (
	topic = "order"
)

var address = []string{"localhost:9092"}

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 10 * time.Millisecond,
			BatchSize:    100,
			MaxAttempts:  3,
			RequiredAcks: kafka.RequireOne,
		},
	}
}

func (p *Producer) Produce(message, key string) error {
	return p.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: []byte(message),
			Time:  time.Now(),
		},
	)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

func main() {
	p := NewProducer(address)
	defer p.Close()

	for i := 0; i < 3; i++ {
		orderUID := uuid.New().String()
		order := map[string]interface{}{
			"order_uid":    orderUID,
			"track_number": "WBILMTESTTRACK",
			"entry":        "WBIL",
			"delivery": map[string]interface{}{
				"name":    "Test Testov",
				"phone":   "+9720000000",
				"zip":     "2639809",
				"city":    "Kiryat Mozkin",
				"address": "Ploshad Mira 15",
				"region":  "Kraiot",
				"email":   "test@gmail.com",
			},
			"payment": map[string]interface{}{
				"transaction":   orderUID,
				"request_id":    "",
				"currency":      "USD",
				"provider":      "wbpay",
				"amount":        1817,
				"payment_dt":    1637907727,
				"bank":          "alpha",
				"delivery_cost": 1500,
				"goods_total":   317,
				"custom_fee":    0,
			},
			"items": []map[string]interface{}{
				{
					"chrt_id":      9934930,
					"track_number": "WBILMTESTTRACK",
					"price":        453,
					"rid":          uuid.New().String(),
					"name":         "Mascaras",
					"sale":         30,
					"size":         "0",
					"total_price":  317,
					"nm_id":        2389212,
					"brand":        "Vivienne Sabo",
					"status":       202,
				},
			},
			"locale":             "en",
			"internal_signature": "",
			"customer_id":        "test",
			"delivery_service":   "meest",
			"shardkey":           "9",
			"sm_id":              99,
			"date_created":       "2021-11-26T06:22:19Z",
			"oof_shard":          "1",
		}

		jsonData, err := json.Marshal(order)
		if err != nil {
			log.Printf("Failed to marshal order to JSON: %v", err)
			continue
		}

		key := uuid.New().String()
		if err := p.Produce(string(jsonData), key); err != nil {
			log.Printf("Failed to produce message: %v", err)
		} else {
			log.Printf("Successfully produced message with key: %s, order_uid: %s", key, orderUID)
		}
	}
}
