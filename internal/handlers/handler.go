package handlers

import (
	"encoding/json"
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
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Order ID is required",
			})
			return
		}

		orderUID := pathParts[2]
		getOrder(w, orderUID, cache, db)
	})
}

func getOrder(w http.ResponseWriter, orderUID string, cache *cache.Cache, db *gorm.DB) {
	order, err := cache.Get(orderUID)
	if err == nil {
		json.NewEncoder(w).Encode(order)
		return
	}
	var dbOrder models.Order
	result := db.Where("order_uid = ?", orderUID).First(&dbOrder)

	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Order not found",
			})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Database error",
			})
		}
		return
	}
	cache.Set(orderUID, &dbOrder)
	json.NewEncoder(w).Encode(dbOrder)
}
