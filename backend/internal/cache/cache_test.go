package cache

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryCache_SetAndGet(t *testing.T) {
	c := NewInMemoryCache()
	c.Set("key1", "value1", 1*time.Hour)

	val, ok := c.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "value1", val)
}

func TestInMemoryCache_Miss(t *testing.T) {
	c := NewInMemoryCache()

	val, ok := c.Get("nonexistent")
	assert.False(t, ok)
	assert.Equal(t, "", val)
}

func TestInMemoryCache_TTLExpires(t *testing.T) {
	c := NewInMemoryCache()
	c.Set("key1", "value1", 1*time.Millisecond)

	time.Sleep(5 * time.Millisecond)

	val, ok := c.Get("key1")
	assert.False(t, ok)
	assert.Equal(t, "", val)
}

func TestInMemoryCache_Overwrite(t *testing.T) {
	c := NewInMemoryCache()
	c.Set("key1", "old", 1*time.Hour)
	c.Set("key1", "new", 1*time.Hour)

	val, ok := c.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "new", val)
}

func TestInMemoryCache_ConcurrentAccess(t *testing.T) {
	c := NewInMemoryCache()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(2)
		key := "key"
		go func() {
			defer wg.Done()
			c.Set(key, "value", 1*time.Hour)
		}()
		go func() {
			defer wg.Done()
			c.Get(key)
		}()
	}

	wg.Wait()
}

func TestDiffCacheKey(t *testing.T) {
	key1 := DiffCacheKey("diff content A")
	key2 := DiffCacheKey("diff content B")
	key1Again := DiffCacheKey("diff content A")

	assert.NotEqual(t, key1, key2)
	assert.Equal(t, key1, key1Again)
	assert.Len(t, key1, 64) // SHA-256 hex = 64 chars
}
