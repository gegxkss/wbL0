package models

import (
	"testing"
	"time"
)

func TestOrderStruct(t *testing.T) {
	order := Order{
		OrderUID:          "uid",
		TrackNumber:       "track",
		Entry:             "entry",
		Locale:            "ru",
		InternalSignature: "sig",
		CustomerId:        "cust",
		DeliveryService:   "service",
		ShardKey:          "key",
		SmId:              1,
		DateCreated:       time.Now(),
		OofShard:          "shard",
		Delivery:          Delivery{Name: "Ivan"},
		Payment:           Payment{Amount: 100},
		Items:             []Items{{Name: "item1"}},
	}
	if order.OrderUID != "uid" {
		t.Error("OrderUID not set")
	}
	if order.Delivery.Name != "Ivan" {
		t.Error("Delivery not set")
	}
	if order.Payment.Amount != 100 {
		t.Error("Payment not set")
	}
	if len(order.Items) != 1 || order.Items[0].Name != "item1" {
		t.Error("Items not set")
	}
}
