package cache

import (
	"sync"
	"time"

	"github.com/dgraph-io/ristretto"
)

type Cache struct {
	data ristretto.Map
	meta *policy
}

func NewCache(size uint64) *Cache {
	return &Cache{
		data: ristretto.NewMap(),
		meta: newPolicy(size),
	}
}

func (c *Cache) Get(key string) interface{} {
	c.meta.hit(key)
	return c.data.Get(key)
}

func (c *Cache) Set(key string, value interface{}, cost uint64) []string {
	return nil
}

////////////////////////////////////////////////////////////////////////////////

type (
	policy struct {
		sync.Mutex
		data   map[string]policyItem
		sample *policySample
		size   uint64
		used   uint64
		ravg   uint64
	}
	policySample []*policyItem
	policyItem   struct {
		key     string
		hits    uint64
		cost    uint64
		created int64
		expire  int64
	}
)

func newPolicyItem(key string, cost uint64) policyItem {
	return policyItem{0, cost, time.Now().UnixNano(), -1}
}

func newPolicy(size uint64) *policy {
	return &policy{
		data:   make(map[string]policyItem, 0),
		sample: make(policySample, 0, 5),
		size:   size,
	}
}

func (p *policy) hit(key string) {
	p.Lock()
	defer p.Unlock()
	if item, exists := p.data[key]; exists {
		item.hits++
		p.data[key] = item
	}
}

func (p *policy) add(key string, cost uint64) []string {
	p.Lock()
	defer p.Unlock()
	// if already in the cache, just update the cost value
	if item, exists := p.data[key]; exists {
		item.cost = cost
		p.data[key] = item
		return nil
	}
	// if eviction is needed
	if p.used+cost > p.size {

	}

	// add to the policy with no eviction needed
	p.data[key] = newPolicyItem(key, cost)
	return nil
}
