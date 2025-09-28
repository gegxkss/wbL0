package models

import "testing"

func TestDeliveryStruct(t *testing.T) {
	d := Delivery{
		Name:    "Ivan",
		Phone:   "123",
		Zip:     "456",
		City:    "Moscow",
		Address: "Lenina 1",
		Region:  "RU",
		Email:   "ivan@example.com",
	}
	if d.Name != "Ivan" {
		t.Error("Name not set")
	}
	if d.City != "Moscow" {
		t.Error("City not set")
	}
	if d.Email != "ivan@example.com" {
		t.Error("Email not set")
	}
}
