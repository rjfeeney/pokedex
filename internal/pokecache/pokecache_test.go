package pokecache

import (
	"testing"
	"time"
)

func TestRequestAddAndGet(t *testing.T) {
	const interval = 5 * time.Second
	cache := NewCache(interval)
	key := "https://example.com/data"
	_, found := cache.Get(key)
	if found {
		t.Fatalf("unexpectedly found key: %s in the cache", key)
	}
	value := []byte("test-response")
	cache.Add(key, value)
	result, found := cache.Get(key)
	if !found {
		t.Fatalf("expected key: %s to be in the cache", key)
	}
	if string(result) != string(value) {
		t.Errorf("expected value: %s, got value: %s", value, result)
	}
}
