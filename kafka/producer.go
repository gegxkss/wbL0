package kafka

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	flushTimeout = 30 * time.Second
)

var errUnknownType = errors.New("unknown event type")

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(address []string) (*Producer, error) {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(address...),
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		BatchSize:    100,
		MaxAttempts:  3,
		RequiredAcks: kafka.RequireOne,
	}

	return &Producer{writer: writer}, nil
}

func (p *Producer) Produce(message, topic, key string) error {
	kafkaMsg := kafka.Message{
		Topic: topic,
		Value: []byte(message),
		Key:   []byte(key),
		Time:  time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), flushTimeout)
	defer cancel()

	err := p.writer.WriteMessages(ctx, kafkaMsg)
	if err != nil {
		return fmt.Errorf("error producing message: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
