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
			ctx, cancel := context.WithTimeout(c.ctx, 5*time.Second)
			msg, err := c.reader.ReadMessage(ctx)
			cancel()

			if err != nil {
				if err == context.DeadlineExceeded {
					log.Println("Polling timeout, no new messages")
					continue
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
				log.Printf("Failed to process message %v", err)
				continue
			}

			log.Printf("Message processed successfully for topic %s, partition %d, offset %d\n",
				msg.Topic, msg.Partition, msg.Offset)
		}
	}
}

func (c *Consumer) processMessage(data []byte) (err error) {
	// Защита от паник
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic in process Message: %v, message: %q", r, string(data))
			err = fmt.Errorf("panic in processMessage: %v", r)
		}
	}()

	// Парсинг JSON
	var order models.Order
	if err := json.Unmarshal(data, &order); err != nil {
		return fmt.Errorf("unmarshal error for message %q: %w", string(data), err)
	}

	// Валидация ключевых полей
	if order.OrderUID == "" {
		return fmt.Errorf("invalid order: order_uid is empty, message: %q", string(data))
	}
	if order.Delivery.Name == "" || order.Delivery.Phone == "" {
		return fmt.Errorf("invalid order: delivery name or phone is empty, order_uid: %s", order.OrderUID)
	}
	if order.Payment.Transaction == "" {
		return fmt.Errorf("invalid order: payment transaction is empty, order_uid: %s", order.OrderUID)
	}

	// Retry-логика для сохранения в БД
	const maxRetries = 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		tx := c.db.Begin()
		// Откат при панике внутри транзакции
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				log.Printf("Panic during transaction for order %s, attempt %d: %v", order.OrderUID, attempt, r)
			}
		}()

		// Сохранение заказа
		if err := tx.Create(&order).Error; err != nil {
			tx.Rollback()
			if attempt == maxRetries {
				return fmt.Errorf("create order failed after %d attempts, order_uid: %s: %w", maxRetries, order.OrderUID, err)
			}
			log.Printf("Attempt %d failed for order %s: %v, retrying...", attempt, order.OrderUID, err)
			time.Sleep(time.Second * time.Duration(attempt))
			continue
		}

		// Коммит транзакции
		if err := tx.Commit().Error; err != nil {
			return fmt.Errorf("commit failed for order %s: %w", order.OrderUID, err)
		}

		c.cache.Set(order)
		log.Printf("Successfully saved order %s and added to cache\n", order.OrderUID)
		return nil
	}
	return fmt.Errorf("create order failed after %d attempts, order_uid: %s", maxRetries, order.OrderUID)
}

func (c *Consumer) Stop() {
	close(c.stopChan)
	c.cancel() // Отменяем контекст для остановки ReadMessage

	if err := c.reader.Close(); err != nil {
		log.Printf("Error closing consumer: %v", err)
	}
}
