package pokecache

import (
	"time"
)

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		Entries:  make(map[string]CacheEntry),
		Interval: interval,
	}
	go c.ReapLoop()
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Entries[key] = CacheEntry{
		CreatedAt: time.Now(),
		Val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	val, ok := c.Entries[key]
	if ok {
		return val.Val, true
	} else {
		return nil, false
	}
}

func (c *Cache) ReapLoop() {
	ticker := time.NewTicker(c.Interval)
	for range ticker.C {
		c.Mutex.Lock()
		for key, entry := range c.Entries {
			age := time.Since(entry.CreatedAt)
			if age > c.Interval {
				delete(c.Entries, key)
			}
		}
		c.Mutex.Unlock()
	}
}
