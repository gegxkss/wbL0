package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gegxkss/wbL0/internal/cache"
	"github.com/gegxkss/wbL0/internal/models"
	"gorm.io/gorm"
)

func SetupRoutes(cache *cache.Cache, db *gorm.DB) {
	fs := http.FileServer(http.Dir("./front"))
	http.Handle("/", fs)
	http.HandleFunc("/order/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 3 || pathParts[2] == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Order ID is required"})
			return
		}

		orderUID := pathParts[2]
		getOrder(w, orderUID, cache, db)
	})
}

func getOrder(w http.ResponseWriter, orderUID string, cache *cache.Cache, db *gorm.DB) {
	// Пытаемся получить из кэша
	if cached, found := cache.Get(orderUID); found {
		log.Printf("Found order in cache: %s", orderUID)
		json.NewEncoder(w).Encode(cached)
		return
	}

	log.Printf("Order not found in cache, searching in DB: %s", orderUID)

	// Загружаем из БД с отношениями
	var order models.Order
	if err := db.Preload("Delivery").Preload("Payment").Preload("Items").
		Where("order_uid = ?", orderUID).First(&order).Error; err != nil {
		w.WriteHeader(http.StatusNotFound) // Просто устанавливаем статус
		json.NewEncoder(w).Encode(map[string]string{"error": "Order not found"})
		return
	}

	// Сохраняем в кэш
	cache.Set(orderUID, &order)
	json.NewEncoder(w).Encode(order)
}
