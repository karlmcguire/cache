package cache

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestCache(t *testing.T) {
	p := NewPolicy(4)

	fmt.Println(p.Add("1"))
	fmt.Println(p.Add("2"))
	fmt.Println(p.Add("3"))
	fmt.Println(p.Add("4"))

	p.Get("2")
	p.Get("2")
	p.Get("2")
	p.Get("1")
	p.Get("1")
	p.Get("3")

	fmt.Println(p.Add("5"))

	spew.Dump(p)
}
