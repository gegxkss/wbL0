package models

import "testing"

func TestPaymentStruct(t *testing.T) {
	p := Payment{
		Transaction:  "tx",
		RequestID:    "req",
		Currency:     "RUB",
		Provider:     "bank",
		Amount:       1000,
		PaymentDt:    123456,
		Bank:         "Sber",
		DeliveryCost: 100,
		GoodsTotal:   900,
		CustomFee:    10,
	}
	if p.Amount != 1000 {
		t.Error("Amount not set")
	}
	if p.Currency != "RUB" {
		t.Error("Currency not set")
	}
	if p.Bank != "Sber" {
		t.Error("Bank not set")
	}
}
