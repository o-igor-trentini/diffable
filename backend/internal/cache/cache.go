package cache

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

type Cache interface {
	Get(key string) (string, bool)
	Set(key, value string, ttl time.Duration)
}

type cacheEntry struct {
	value     string
	expiresAt time.Time
}

type InMemoryCache struct {
	data sync.Map
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{}
}

func (c *InMemoryCache) Get(key string) (string, bool) {
	val, ok := c.data.Load(key)
	if !ok {
		return "", false
	}
	entry := val.(cacheEntry)
	if time.Now().After(entry.expiresAt) {
		c.data.Delete(key)
		return "", false
	}
	return entry.value, true
}

func (c *InMemoryCache) Set(key, value string, ttl time.Duration) {
	c.data.Store(key, cacheEntry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	})
}

func DiffCacheKey(diff string) string {
	h := sha256.Sum256([]byte(diff))
	return fmt.Sprintf("%x", h)
}
