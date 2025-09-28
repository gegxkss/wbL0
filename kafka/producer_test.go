package kafka

import (
	"testing"
)

type mockWriter struct{}

func (m *mockWriter) WriteMessages(...interface{}) error { return nil }
func (m *mockWriter) Close() error                       { return nil }

func TestNewProducer(t *testing.T) {
	p, err := NewProducer([]string{"localhost:9091"})
	if err != nil {
		t.Fatalf("Producer creation failed: %v", err)
	}
	if p.writer == nil {
		t.Error("Writer is nil")
	}
}
