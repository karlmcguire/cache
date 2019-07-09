package cache

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestCache(t *testing.T) {
	c := NewCache(4)
	fmt.Println(c.Set("1", 1))
	fmt.Println(c.Set("2", 1))
	fmt.Println(c.Set("3", 1))
	fmt.Println(c.Set("4", 1))
	fmt.Println(c.Set("5", 1))
	fmt.Println(c.Set("6", 1))

	c.Get("3")
	c.Get("3")

	spew.Dump(c.meta.data)
}

func BenchmarkCache(b *testing.B) {
	c := NewCache(4)
	b.SetBytes(1)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Get("1")
		}
	})
}
