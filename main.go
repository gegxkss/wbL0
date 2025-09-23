package main

import (
	"flag"
	"log"

	"github.com/gegxkss/wbL0/internal/cache"
	"github.com/gegxkss/wbL0/internal/config"
	"github.com/gegxkss/wbL0/migrations"
)

var migrate = flag.Bool("m", false, "Run database migrations")

func main() {
	flag.Parse()

	if *migrate {
		migrations.Migration()
		log.Println("Migrations completed. Exiting.")
		return
	}
	config.ConnectDB()
	log.Println("Database connection successful!")
	log.Println("Add your Kafka and HTTP logic here later.")

	cache := cache.NewCache()
	err := cache.RestoreFromDB(config.DB)
	if err != nil {
		log.Printf("Failed to restore cache from DB: %v", err)
	} else {
		log.Printf("Cache restored from DB: %d orders loaded", cache.Size())
	}
}
