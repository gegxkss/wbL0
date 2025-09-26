// internal/kafka/consumer.go
package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gegxkss/wbL0/internal/cache"
	"github.com/gegxkss/wbL0/internal/models"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

type Consumer struct {
	reader   *kafka.Reader
	db       *gorm.DB
	stopChan chan struct{}
	cache    *cache.Cache
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewConsumer(address []string, topic, groupID string, db *gorm.DB, cache *cache.Cache) (*Consumer, error) {
	log.Printf("Connecting to Kafka brokers: %v", address)
	config := kafka.ReaderConfig{
		Brokers:  address,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		MaxWait:  5 * time.Second,
	}

	reader := kafka.NewReader(config)

	ctx, cancel := context.WithCancel(context.Background())

	return &Consumer{
		reader:   reader,
		db:       db,
		stopChan: make(chan struct{}),
		cache:    cache,
		ctx:      ctx,
		cancel:   cancel,
	}, nil
}

func (c *Consumer) Start() {
	log.Println("Starting Kafka consumer")
	defer log.Println("Kafka consumer stopped")

	for {
		select {
		case <-c.stopChan:
			return
		default:
			// Используем контекст с таймаутом для чтения
			ctx, cancel := context.WithTimeout(c.ctx, 20*time.Second)
			msg, err := c.reader.ReadMessage(ctx)
			cancel()

			if err != nil {
				if err == context.DeadlineExceeded {
					continue // Таймаут - нормально, продолжаем
				}
				if strings.Contains(err.Error(), "context canceled") {
					log.Println("Consumer context canceled, stopping")
					return
				}
				log.Printf("Consumer error: %v", err)
				continue
			}

			// Проверка пустого сообщения
			if len(msg.Value) == 0 {
				log.Println("We got empty message, continue")
				continue
			}

			if err := c.processMessage(msg.Value); err != nil {
				log.Printf("Failed to process message: %v", err)
				continue
			}

			log.Printf("Message processed successfully for topic %s, partition %d, offset %d",
				msg.Topic, msg.Partition, msg.Offset)
		}
	}
}

// Структура для парсинга JSON из Kafka
type KafkaOrder struct {
	OrderUID          string          `json:"order_uid"`
	TrackNumber       string          `json:"track_number"`
	Entry             string          `json:"entry"`
	Delivery          models.Delivery `json:"delivery"`
	Payment           models.Payment  `json:"payment"`
	Items             []models.Items  `json:"items"`
	Locale            string          `json:"locale"`
	InternalSignature string          `json:"internal_signature"`
	CustomerId        string          `json:"customer_id"`
	DeliveryService   string          `json:"delivery_service"`
	ShardKey          string          `json:"shardkey"`
	SmId              int             `json:"sm_id"`
	DateCreated       string          `json:"date_created"`
	OofShard          string          `json:"oof_shard"`
}

func (c *Consumer) processMessage(data []byte) error {
	// Парсим JSON в временную структуру
	var kafkaOrder KafkaOrder
	if err := json.Unmarshal(data, &kafkaOrder); err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}

	// Валидация обязательных полей
	if kafkaOrder.OrderUID == "" {
		return fmt.Errorf("invalid order: order_uid is empty")
	}

	// Начинаем транзакцию
	tx := c.db.Begin()

	// 1. Сохраняем основной заказ
	order := models.Order{
		OrderUID:          kafkaOrder.OrderUID,
		TrackNumber:       kafkaOrder.TrackNumber,
		Entry:             kafkaOrder.Entry,
		Locale:            kafkaOrder.Locale,
		InternalSignature: kafkaOrder.InternalSignature,
		CustomerId:        kafkaOrder.CustomerId,
		DeliveryService:   kafkaOrder.DeliveryService,
		ShardKey:          kafkaOrder.ShardKey,
		SmId:              kafkaOrder.SmId,
		OofShard:          kafkaOrder.OofShard,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("create order failed: %w", err)
	}

	// 2. Сохраняем Delivery (устанавливаем OrderUID)
	kafkaOrder.Delivery.OrderUID = kafkaOrder.OrderUID
	if err := tx.Create(&kafkaOrder.Delivery).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("create delivery failed: %w", err)
	}

	// 3. Сохраняем Payment (устанавливаем OrderUID)
	kafkaOrder.Payment.OrderUID = kafkaOrder.OrderUID
	if err := tx.Create(&kafkaOrder.Payment).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("create payment failed: %w", err)
	}

	// 4. Сохраняем Items (устанавливаем OrderUID для каждого)
	for i := range kafkaOrder.Items {
		kafkaOrder.Items[i].OrderUID = kafkaOrder.OrderUID
		if err := tx.Create(&kafkaOrder.Items[i]).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("create item failed: %w", err)
		}
	}

	// Коммитим транзакцию
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("commit transaction failed: %w", err)
	}

	// Добавляем в кэш (только основную информацию о заказе)
	if err := c.cache.Set(order.OrderUID, &order); err != nil {
		log.Printf("Warning: failed to add order to cache: %v", err)
	}

	log.Printf("Successfully saved order %s with all related data", order.OrderUID)
	return nil
}

func (c *Consumer) Stop() {
	close(c.stopChan)
	c.cancel()

	if err := c.reader.Close(); err != nil {
		log.Printf("Error closing consumer: %v", err)
	}
}
