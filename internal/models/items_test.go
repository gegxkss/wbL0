package models

import "testing"

func TestItemsStruct(t *testing.T) {
	item := Items{
		ChrtId: 1,
		Tracknumber: "track",
		Price: 100,
		Rid: "rid",
		Name: "item",
		Sale: 10,
		Size: "M",
		TotalPrice: 90,
		NmId: 123,
		Brand: "brand",
		Status: 1,
	}
	if item.Name != "item" {
		t.Error("Name not set")
	}
	if item.Price != 100 {
		t.Error("Price not set")
	}
	if item.Brand != "brand" {
		t.Error("Brand not set")
	}
}
