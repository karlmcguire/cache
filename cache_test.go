package cache

import (
	"testing"
)

func TestCache(t *testing.T) {
	c := NewCache(16)
	c.Set("1", 1, 1)
	if value, has := c.Get("1"); !has || value.(int) != 1 {
		t.Fatal("set/get error")
	}
}

func BenchmarkCacheGetOne(b *testing.B) {
	c := NewCache(16)
	c.Set("1", 1, 1)
	b.SetBytes(1)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Get("1")
		}
	})
}

func BenchmarkCacheSetOne(b *testing.B) {
	c := NewCache(16)
	b.SetBytes(1)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for i := uint64(0); pb.Next(); i++ {
			c.Set("1", 1, i)
		}
	})
}
