package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cache           map[string]cacheEntry
	cleanupInterval time.Duration
	mu              sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		cache:           make(map[string]cacheEntry),
		cleanupInterval: interval,
	}

	go c.readLoop()

	return c
}

func (c *Cache) Add(key string, val []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}

	return nil
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.cache[key]

	return entry.val, ok
}

func (c *Cache) readLoop() {
	//Check if each item in the cache has expired
	//Use time.Ticker to trigger the cleanup at the specified interval

	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		c.mu.Lock()
		defer c.mu.Unlock()
		for key, entry := range c.cache {
			if now.Sub(entry.createdAt) > c.cleanupInterval {
				delete(c.cache, key)
			}
		}
	}
}
