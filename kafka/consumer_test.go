package kafka

import (
	"testing"

	"github.com/gegxkss/wbL0/internal/cache"
	"gorm.io/gorm"
)

func TestNewConsumer(t *testing.T) {
	cache := cache.NewCache()
	var db *gorm.DB
	c, err := NewConsumer([]string{"localhost:9091"}, "order", "group", db, cache)
	if err != nil {
		t.Fatalf("Consumer creation failed: %v", err)
	}
	if c.reader == nil {
		t.Error("Reader is nil")
	}
	if c.cache != cache {
		t.Error("Cache not set")
	}
}
