package cache

import (
	"sync"
)

type Cache struct {
	sync.RWMutex

	data map[string][]byte
}

func New() *Cache {
	return &Cache{data: make(map[string][]byte)}
}

func (c *Cache) Get(k string) []byte {
	c.RLock()
	defer c.RUnlock()

	if v, exists := c.data[k]; exists {
		return v
	}

	return nil
}

func (c *Cache) Getter(k string) func() []byte {
	return func() []byte {
		c.RLock()
		defer c.RUnlock()

		return c.Get(k)
	}
}

func (c *Cache) Set(k string, v []byte) {
	c.Lock()
	defer c.Unlock()

	c.data[k] = v
}

func (c *Cache) Setter(k string) func([]byte) {
	return func(v []byte) {
		c.Lock()
		defer c.Unlock()

		c.data[k] = v
	}
}
