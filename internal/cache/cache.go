package cache

import (
	"fmt"
	"log"
	"sync"
)

type Cache struct {
	items   map[string]interface{} // ИЗМЕНИЛ: interface{} вместо *models.Order
	size    int
	maxSize int
	mu      sync.RWMutex
}

const (
	cacheMaxSize = 500
)

func NewCache() *Cache {
	return &Cache{
		items:   make(map[string]interface{}),
		size:    0,
		maxSize: cacheMaxSize,
	}
}

func (c *Cache) Set(orderUID string, value interface{}) error { // ИЗМЕНИЛ: interface{} вместо *models.Order
	if value == nil {
		return fmt.Errorf("value can not be nil")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.size >= c.maxSize {
		c.deleteLast()
	}

	c.items[orderUID] = value
	c.size++

	return nil
}

func (c *Cache) deleteLast() {
	for orderUID := range c.items {
		delete(c.items, orderUID)
		c.size--
		log.Printf("Removed order from cache: %s", orderUID)
		break
	}
}

func (c *Cache) Get(orderUID string) (interface{}, bool) { // ИЗМЕНИЛ: возвращает interface{} и bool
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[orderUID]
	return item, exists
}
