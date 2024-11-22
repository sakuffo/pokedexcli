package pokecache

import (
	"sync"
	"time"

	"github.com/sakuffo/pokedexcli/internal/logger"
)

type Cache struct {
	cache           map[string]cacheEntry
	cleanupInterval time.Duration
	mu              sync.Mutex
	logger          *logger.Logger
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration, logger *logger.Logger) *Cache {
	c := &Cache{
		cache:           make(map[string]cacheEntry),
		cleanupInterval: interval,
		logger:          logger,
	}

	go c.readLoop()

	return c
}

func (c *Cache) Add(key string, val []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.logger.Debug("Adding item to cache: %s", key)
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

	if ok {
		c.logger.Debug("Cache hit: %s", key)
	} else {
		c.logger.Debug("Cache miss: %s", key)
	}

	return entry.val, ok
}

func (c *Cache) readLoop() {
	//Check if each item in the cache has expired
	//Use time.Ticker to trigger the cleanup at the specified interval

	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {

		c.mu.Lock()
		initialSize := len(c.cache)
		now := time.Now()

		for key, entry := range c.cache {
			if now.Sub(entry.createdAt) > c.cleanupInterval {
				c.logger.Debug("Removing expired item from cache: %s", key)
				delete(c.cache, key)
			}
		}

		itemsRemoved := initialSize - len(c.cache)
		if itemsRemoved > 0 {
			c.logger.Info("Removed %d items from cache", itemsRemoved)
		}

		c.mu.Unlock()
	}
}
