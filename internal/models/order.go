package models

import (
	"time"
)

type Order struct {
	OrderUID          string    `gorm:"primaryKey" json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerId        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	ShardKey          string    `gorm:"column:shardkey" json:"shardkey"`
	SmId              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`

	Delivery Delivery `gorm:"foreignKey:OrderUID" json:"delivery"`
	Payment  Payment  `gorm:"foreignKey:OrderUID" json:"payment"`
	Items    []Items  `gorm:"foreignKey:OrderUID" json:"items"`
}
