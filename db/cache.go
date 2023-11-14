package db

import (
	"sync"
	"time"
)

type cachedItem struct {
	value      any
	expiration *time.Time
}

type cache struct {
	items map[string]cachedItem
	mutex *sync.RWMutex
}

func newCache() *cache {
	cache := &cache{
		items: make(map[string]cachedItem),
		mutex: &sync.RWMutex{},
	}

	go func() {
		for {
			time.Sleep(1 * time.Second)

			cache.mutex.Lock()
			for key, item := range cache.items {
				if item.expiration != nil && item.expiration.Before(time.Now()) {
					delete(cache.items, key)
				}
			}

			cache.mutex.Unlock()
		}
	}()

	return cache
}

func (c *cache) Get(key string) (any, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, ok := c.items[key]
	if !ok {
		return nil, false
	}

	if item.expiration != nil && item.expiration.Before(time.Now()) {
		return nil, false
	}

	return item.value, true
}

func (c *cache) GetOnce(key string) (any, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, ok := c.items[key]
	if !ok {
		return nil, false
	}

	delete(c.items, key)

	return item.value, true
}

func (c *cache) Set(key string, value any, exp *time.Time) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items[key] = cachedItem{
		value:      value,
		expiration: exp,
	}
}

func (c *cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.items, key)
}

func (c *cache) HasKey(key string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	_, ok := c.items[key]
	return ok
}
