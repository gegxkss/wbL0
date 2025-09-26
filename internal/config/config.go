package config

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DB_URL")

	if dsn == "" {
		dsn = "host=localhost user=gegxkss password=postgres dbname=wbl0_db port=5432 sslmode=disable"
		log.Println("DB_URL not set, using default:", dsn)
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Successfully connected to database")
}

func RestoreFromDB() {

}
