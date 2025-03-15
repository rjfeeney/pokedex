package pokecache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	CreatedAt time.Time
	Val       []byte
}

type Cache struct {
	Entries  map[string]CacheEntry
	Mutex    sync.Mutex
	Interval time.Duration
}
