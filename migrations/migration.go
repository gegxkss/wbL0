package migrations

import (
	"log"

	"github.com/gegxkss/wbL0/internal/config"
	"github.com/gegxkss/wbL0/internal/models"
)

func Migration() {
	config.ConnectDB()

	err := config.DB.AutoMigrate(
		&models.Order{},
		&models.Delivery{},
		&models.Payment{},
		&models.Items{},
	)

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration completed successfully")
}
