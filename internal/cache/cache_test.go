package cache

import (
	"testing"
)

func TestNewCache(t *testing.T) {
	c := NewCache()
	if c == nil {
		t.Fatal("Cache is nil")
	}
	if c.maxSize != 1000 {
		t.Errorf("Expected maxSize 1000, got %d", c.maxSize)
	}
}

func TestSetAndGet(t *testing.T) {
	c := NewCache()
	orderUID := "test-uid"
	value := "test-value"
	if err := c.Set(orderUID, value); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	got, ok := c.Get(orderUID)
	if !ok {
		t.Fatalf("Get failed: not found")
	}
	if got != value {
		t.Errorf("Expected %v, got %v", value, got)
	}
}

func TestSetNilValue(t *testing.T) {
	c := NewCache()
	if err := c.Set("uid", nil); err == nil {
		t.Error("Expected error for nil value")
	}
}
