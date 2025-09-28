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
		MinBytes: 10e3,
		MaxBytes: 10e6,
		MaxWait:  30 * time.Second,
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
			ctx, cancel := context.WithTimeout(c.ctx, 20*time.Second)
			msg, err := c.reader.ReadMessage(ctx)
			cancel()

			if err != nil {
				if err == context.DeadlineExceeded {
					continue
				}
				if strings.Contains(err.Error(), "context canceled") {
					log.Println("Consumer context canceled, stopping")
					return
				}
				log.Printf("Consumer error: %v", err)
				continue
			}

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

func (c *Consumer) processMessage(data []byte) error {
	var order models.Order
	if err := json.Unmarshal(data, &order); err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}

	log.Printf("Received order: %s, items count: %d", order.OrderUID, len(order.Items))

	if order.OrderUID == "" {
		return fmt.Errorf("invalid order: order_uid is empty")
	}

	tx := c.db.Begin()

	// Сохраняем заказ
	orderToSave := models.Order{
		OrderUID:          order.OrderUID,
		TrackNumber:       order.TrackNumber,
		Entry:             order.Entry,
		Locale:            order.Locale,
		InternalSignature: order.InternalSignature,
		CustomerId:        order.CustomerId,
		DeliveryService:   order.DeliveryService,
		ShardKey:          order.ShardKey,
		SmId:              order.SmId,
		DateCreated:       order.DateCreated,
		OofShard:          order.OofShard,
	}

	if err := tx.Create(&orderToSave).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("create order failed: %w", err)
	}

	// Сохраняем доставку
	delivery := order.Delivery
	delivery.ID = 0
	delivery.OrderUID = order.OrderUID
	if err := tx.Create(&delivery).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("create delivery failed: %w", err)
	}

	// Сохраняем товары
	for i := range order.Items {
		item := order.Items[i]
		item.ID = 0
		item.OrderUID = order.OrderUID
		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("create item failed: %w", err)
		}
	}

	// Сохраняем оплату
	payment := order.Payment
	payment.ID = 0
	payment.OrderUID = order.OrderUID
	if err := tx.Create(&payment).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("create payment failed: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("commit transaction failed: %w", err)
	}

	// Сохраняем в кэш оригинальный order
	log.Printf("Saving order to cache: %s", order.OrderUID)
	if err := c.cache.Set(order.OrderUID, &order); err != nil {
		log.Printf("Warning: failed to add order to cache: %v", err)
	}

	log.Printf("Successfully saved order %s", order.OrderUID)
	return nil
}

func (c *Consumer) Stop() {
	close(c.stopChan)
	c.cancel()
	c.reader.Close()
}
