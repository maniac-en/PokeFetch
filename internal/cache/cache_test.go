package cache

import (
	"testing"
	"time"
)

func TestCache_Add(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		interval time.Duration
		// Named input parameters for target function.
		key string
		val []byte
	}{
		{
			"Add to cache-1",
			2 * time.Second,
			"https://example.com",
			[]byte("pikapika"),
		},
		{
			"Add to cache-2",
			2 * time.Second,
			"https://example.com/newpath",
			[]byte("pikachu"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCache(tt.interval)
			c.Add(tt.key, tt.val)
			c.mu.RLock()
			defer c.mu.RUnlock()
			if entry, ok := c.entries[tt.key]; !ok || string(entry.val) != string(tt.val) {
				t.Errorf("%s\ntest entry %s was not added to the cache", tt.name, tt.key)
			}
		})
	}
}

func TestCache_removeExpired(t *testing.T) {
	expiryTime := 5 * time.Millisecond
	cache := NewCache(expiryTime)
	key := "https://example.com"
	cache.Add(key, []byte("test"))

	if _, ok := cache.Get(key); !ok {
		t.Errorf("expected to find %s", key)
	}
	time.Sleep(expiryTime * 2)
	if _, ok := cache.Get(key); !ok {
		t.Errorf("expected to not find %s", key)
	}
}
