package cache

import (
	"sync"
	"time"
)

type Cache struct {
	items  sync.Map
	ticker *time.Ticker
	done   chan bool
}

type cacheItem struct {
	data   any
	expiry time.Time
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		ticker: time.NewTicker(interval),
		done:   make(chan bool),
	}
	go c.startCleanup()
	return c
}

func (c *Cache) Add(key string, data any, expiry time.Duration) {
	item := &cacheItem{
		data:   data,
		expiry: time.Now().Add(expiry),
	}
	c.items.Store(key, item)
}

func (c *Cache) Get(key string) (any, bool) {
	if item, ok := c.items.Load(key); ok {
		cacheItem := item.(*cacheItem)
		if time.Now().Before(cacheItem.expiry) {
			return cacheItem.data, true
		}
		c.items.Delete(key)
	}
	return nil, false
}

func (c *Cache) Delete(key string) {
	c.items.Delete(key)
}

func (c *Cache) Clear() {
	c.items.Range(func(key, _ any) bool {
		c.items.Delete(key)
		return true
	})
}

func (c *Cache) startCleanup() {
	for {
		select {
		case <-c.ticker.C:
			c.items.Range(func(key, value any) bool {
				item := value.(*cacheItem)
				if time.Now().After(item.expiry) {
					c.items.Delete(key)
				}
				return true
			})
		case <-c.done:
			c.ticker.Stop()
			return
		}
	}
}

func (c *Cache) Stop() {
	close(c.done)
}
