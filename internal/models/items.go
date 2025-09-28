package models

type Items struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"-"`
	OrderUID    string `gorm:"not null;column:order_uid" json:"-"`
	ChrtId      int    `json:"chrt_id"`
	Tracknumber string `gorm:"column:track_number" json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmId        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}
