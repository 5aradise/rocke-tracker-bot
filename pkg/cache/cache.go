package cache

import (
	"sync"
	"time"
)

type Cache[K comparable, V any] struct {
	data sync.Map
	stop chan struct{}
}

func New[K comparable, V any](clearDelay time.Duration) *Cache[K, V] {
	c := &Cache[K, V]{}
	go c.run(clearDelay)
	return c
}

func (c *Cache[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	v, loaded := c.data.LoadOrStore(key, value)
	if loaded {
		value = v.(V)
	}
	return value, loaded
}

func (c *Cache[K, V]) Delete(key K) {
	c.data.Delete(key)
}

func (c *Cache[K, V]) Stop() {
	close(c.stop)
}

func (c *Cache[K, V]) run(clearDelay time.Duration) {
	update := time.Tick(clearDelay)
	for {
		select {
		case <-update:
			c.data.Clear()
		case <-c.stop:
			return
		}
	}
}
