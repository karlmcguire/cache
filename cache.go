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
	if victimKey, victimCount = s.candidate(); victimKey != "" {
		s.evict(victimKey)
	}
	s.data[key]++
	return
}

func (s *Segment) evict(key string) {
	delete(s.data, key)
}

func (s *Segment) candidate() (string, uint8) {
	if uint64(len(s.data)) != s.size {
		return "", 0
	}
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
	// get the window victim and count if the window is full
	windowKey, windowCount := p.segs[0].Add(key)
	if windowKey == "" {
		// window has room, so nothing else needs to be done
		return
	}
	// get the main eviction candidate key and count if the main segment is full
	mainKey, mainCount := p.segs[1].candidate()
	if mainKey == "" {
		// main has room, so just move window victim to there
		goto move
	}
	// compare the window victim with the main candidate, also note that window
	// victims are preferred (>=) over main candidates, as we can assume that
	// window victims have been used more recently than the main candidate
	if windowCount >= mainCount {
		// main candidate lost to the window victim, so actually evict the main
		// candidate
		victim = mainKey
		p.segs[1].evict(mainKey)
		// main now has room for one more, so move the window victim to there
		goto move
	} else {
		// window victim lost to the main candidate, and the window victim has
		// already been evicted from the window, so nothing else needs to be
		// done
		victim = windowKey
	}
	return
move:
	// move moves the window key-count pair to the main segment
	p.segs[1].data[windowKey] = windowCount
	return
}
