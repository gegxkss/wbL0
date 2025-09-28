package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gegxkss/wbL0/internal/cache"
)

// Мок для базы данных

type mockDB struct{}

func (m *mockDB) Preload(string) *mockDB               { return m }
func (m *mockDB) Where(string, ...interface{}) *mockDB { return m }
func (m *mockDB) First(dest interface{}) error {
	return nil // всегда "находит" заказ
}

func TestSetupRoutes_OrderNotFound(t *testing.T) {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./front"))
	mux.Handle("/", fs)
	mux.HandleFunc("/order/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 3 || pathParts[2] == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Order ID is required"})
			return
		}
	})

	req := httptest.NewRequest("GET", "/order/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestSetupRoutes_OrderFoundInCache(t *testing.T) {
	c := cache.NewCache()
	orderUID := "test-uid"
	c.Set(orderUID, map[string]string{"order_uid": orderUID})
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./front"))
	mux.Handle("/", fs)
	mux.HandleFunc("/order/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 3 || pathParts[2] == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Order ID is required"})
			return
		}
		orderUID := pathParts[2]
		if cached, found := c.Get(orderUID); found {
			json.NewEncoder(w).Encode(cached)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Order not found"})
	})

	req := httptest.NewRequest("GET", "/order/"+orderUID, nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}
