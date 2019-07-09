package cache

import (
	"math"
	"sync"
)

type Cache struct {
	data *sync.Map
	meta *Policy
}

func NewCache(size uint64) *Cache {
	return &Cache{
		data: &sync.Map{},
		meta: NewPolicy(size),
	}
}

func (c *Cache) Get(key string) interface{} {
	// record access
	c.meta.Get(key)
	// return value
	value, _ := c.data.Load(key)
	return value
}

func (c *Cache) Set(key string, value interface{}) string {
	var victim string
	if victim = c.meta.Add(key); victim != "" {
		// delete victim from hash map
		c.data.Delete(victim)
	}
	// add new item
	c.data.Store(key, value)
	return victim
}

type Policy struct {
	sync.Mutex
	data map[string]uint8
	size uint64
}

func NewPolicy(size uint64) *Policy {
	return &Policy{
		data: make(map[string]uint8, size),
		size: size,
	}
}

func (p *Policy) Get(key string) {
	p.Lock()
	defer p.Unlock()
	p.data[key]++
}

func (p *Policy) Add(key string) (victim string) {
	p.Lock()
	defer p.Unlock()
	// do nothing if the item is already in the policy
	if _, exists := p.data[key]; exists {
		return
	}
	// check if eviction is needed
	if uint64(len(p.data)) >= p.size {
		// evict
		i, min := 0, uint8(math.MaxUint8)
		for k, v := range p.data {
			if v < min {
				victim, min = k, v
			}
			// stop at sample size
			if i++; i == 5 {
				break
			}
		}
		delete(p.data, victim)
	}
	// add item to policy
	p.data[key]++
	return
}
