package cache

import (
	"math"
	"sync"
)

/*
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

func (c *Cache) Set(key string, value interface{}) {
	if victim := c.meta.Add(key); victim != "" {
		// delete victim from hash map
		c.data.Delete(victim)
	}
	// add new item
	c.data.Store(key, value)
}
*/

type Policy struct {
	sync.Mutex
	segs [2]*Segment
}

type Segment struct {
	data map[string]uint8
	size uint64
}

func NewSegment(size uint64) *Segment {
	return &Segment{
		data: make(map[string]uint8, size),
		size: size,
	}
}

func (s *Segment) Add(key string) (victimKey string, victimCount uint8) {
	// check if eviction is needed
	if uint64(len(s.data)) == s.size {
		victimKey, victimCount = s.candidate()
		s.evict(victimKey)
	}
	s.data[key]++
	return
}

func (s *Segment) evict(key string) {
	delete(s.data, key)
}

func (s *Segment) candidate() (string, uint8) {
	i, minKey, minCount := 0, "", uint8(math.MaxUint8)
	for key, count := range s.data {
		if count < minCount {
			minKey, minCount = key, count
		}
		if i++; i == 5 {
			break
		}
	}
	return minKey, minCount
}

func NewPolicy(size uint64) *Policy {
	return &Policy{
		segs: [2]*Segment{
			// window
			NewSegment(uint64(math.Ceil(float64(size) * 0.01))),
			// main
			NewSegment(uint64(math.Floor(float64(size) * 0.99))),
		},
	}
}

func (p *Policy) Seg(key string) int {
	if p.segs[0].data[key] != 0 {
		return 0
	} else if p.segs[1].data[key] != 0 {
		return 1
	}
	return -1
}

func (p *Policy) Get(key string) {
	p.Lock()
	defer p.Unlock()
	if seg := p.Seg(key); seg != -1 {
		p.segs[seg].data[key]++
	}
}

func (p *Policy) Add(key string) (victim string) {
	p.Lock()
	defer p.Unlock()
	// do nothing if the item is already in the policy
	if p.Seg(key) != -1 {
		return
	}
	if windowKey, windowCount := p.segs[0].Add(key); windowKey != "" {
		if mainKey, mainCount := p.segs[1].candidate(); mainKey != "" {
			if windowCount > mainCount {
				p.segs[1].evict(mainKey)
			}
		}
		victim, _ = p.segs[1].Add(windowKey)
	}
	return
}
