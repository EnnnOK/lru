package lru

import (
	"testing"
	"time"
)

type value struct {
	data string
}

func (v value) Len() int64 {
	return int64(len(v.data))
}

func TestNewLRU(t *testing.T) {
	lru := NewLRU(100, 500)
	t.Logf("%#v", lru)
}

func TestNewLRUWithCallback(t *testing.T) {
	callback := func(interface{}) error { return nil }
	lru := NewLRUWithCallback(0, callback)
	t.Logf("%#v", lru)
}

func TestNewNode(t *testing.T) {
	lru := NewLRU(100, 100)
	lru.AddNewNode(123, value{"234"})
	lru.AddNewNode(1231, value{"2342"})
	t.Log(lru.Traversal())
}

func TestGet(t *testing.T) {
	lru := NewLRU(100, 1)
	node1 := lru.newNode(123, value{"234"})
	node2 := lru.newNode(1231, value{"2342"})
	lru.add(node1)
	lru.add(node2)
	t.Logf("%#v", lru.Traversal())
	time.Sleep(2 * time.Second)
	t.Log(lru.Access(node1))
	t.Log(lru.Access(node2))
	t.Logf("%#v", lru.Traversal())
}
