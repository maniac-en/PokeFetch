// Package cache is responsible to cache the API results, and other things
// around here, so the overall REPL's experience feels performant/snappier
package cache

import (
	"fmt"
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries map[string]cacheEntry
	mu      *sync.RWMutex
	ttl     time.Duration
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{time.Now(), val}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entries[key]
	return entry.val, ok
}

func (c *Cache) removeExpired() {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for key, val := range c.entries {
		if val.createdAt.Compare(time.Now()) >= 0 {
			fmt.Println("Removing the cached entry:", key)
			delete(c.entries, key)
		}
	}
}

func NewCache(interval time.Duration) Cache {
	// create a new cache with a configurable interval
	// and purge it from cache when interval passes
	cache := Cache{
		entries: make(map[string]cacheEntry),
		mu:      &sync.RWMutex{},
		ttl:     interval,
	}
	go func() {
		ticker := time.NewTicker(cache.ttl)
		defer ticker.Stop()
		for range ticker.C {
			cache.removeExpired()
		}
	}()
	return cache
}
