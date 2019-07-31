package cache

import (
	"sync"
	"sync/atomic"
	"time"
)

type (
	Cache struct {
		sync.Mutex
		data *sync.Map
		sets chan interface{}
		used uint64
	}
	item struct {
		value interface{}
		size  uint64
		when  uint64
		hits  uint64
	}
)

func NewCache(size uint64) *Cache {
	cache := &Cache{
		data: &sync.Map{},
		sets: make(chan interface{}, 64),
	}
	go cache.Evictor(size)
	return cache
}

func (c *Cache) Get(key interface{}) (interface{}, bool) {
	if i := c.get(key); i != nil {
		i.hits++
		return i.value, true
	}
	return nil, false
}

func (c *Cache) get(key interface{}) *item {
	if i, _ := c.data.Load(key); i != nil {
		return i.(*item)
	}
	return nil
}

func (c *Cache) Set(key, value interface{}, size uint64) {
	// TODO: this doesn't handle value / cost updates
	if _, had := c.data.LoadOrStore(key, &item{
		value: value,
		size:  size,
		when:  uint64(time.Now().Unix()),
		hits:  1,
	}); had {
		atomic.AddUint64(&c.used, size)
		c.sets <- key
	}
}

func (c *Cache) Evictor(size uint64) {
	for key := range c.sets {
		_ = key
	}
}

func (i *item) priority() float64 {
	pass := float64(uint64(time.Now().Unix()) - i.when)
	if pass == 0 {
		pass = 1.0
	}
	hits := float64(i.hits)
	cost := 1 / float64(i.size)
	return (hits / pass) * cost
}
