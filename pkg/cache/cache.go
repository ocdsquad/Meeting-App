package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// Cache is a wrapper around the go-cache instance
type Cache struct {
	instance *cache.Cache
}

// MyCache is an example of how to initialize a cache
var MyCache = NewCache(5*time.Minute, 10*time.Minute)

// NewCache initializes a new cache
func NewCache(defaultExpiration, cleanupInterval time.Duration) *Cache {
	return &Cache{
		instance: cache.New(defaultExpiration, cleanupInterval),
	}
}

// Set stores a value in the cache with a TTL (time-to-live)
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.instance.Set(key, value, ttl)
}

// Get retrieves a value from the cache
func (c *Cache) Get(key string) (interface{}, bool) {
	return c.instance.Get(key)
}

// Delete removes a key from the cache
func (c *Cache) Delete(key string) {
	c.instance.Delete(key)
}
