package cache

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestCache(t *testing.T) {
	c := NewCache(8)
	spew.Dump(c)
}
