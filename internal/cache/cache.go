package cache

import (
	"fmt"
	"log"
	"sync"

	"github.com/gegxkss/wbL0/internal/models"
)

type Cache struct {
	orders  map[string]*models.Order
	size    int
	maxSize int
	mu      sync.RWMutex
}

const (
	cacheMaxSize = 500
)

func NewCache() *Cache {
	return &Cache{
		orders:  make(map[string]*models.Order),
		size:    0,
		maxSize: cacheMaxSize,
	}
}

func (c *Cache) Set(orderUID string, order *models.Order) error {
	if order == nil {
		return fmt.Errorf("can not be nil")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.size >= c.maxSize {
		c.deleteLast()
	}

	c.orders[orderUID] = order
	c.size++

	return nil
}

func (c *Cache) deleteLast() {
	for orderUID := range c.orders {
		delete(c.orders, orderUID)
		c.size--
		log.Printf("Removed order from cache: %s", orderUID)
		break
	}
}

func (c *Cache) Get(orderUID string) (*models.Order, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, exists := c.orders[orderUID]

	if !exists {
		return nil, fmt.Errorf("order not found: %s", orderUID)
	}

	return order, nil
}
