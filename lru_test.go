package lru

import (
	"fmt"
	"testing"
	"time"
)

type value struct {
	data string
}

func (v value) Len() int64 {
	return int64(len(v.data))
}

func printList(t *testing.T, lru *LRU) {
	list := lru.Traversal()
	for i := range list {
		t.Logf("%#v", list[i])
	}
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
	printList(t, lru)
}

func TestGet(t *testing.T) {
	lru := NewLRU(100, 1)
	node1 := lru.newNode(123, value{"234"})
	node2 := lru.newNode(1231, value{"2342"})
	lru.add(node1)
	lru.add(node2)
	printList(t, lru)
	time.Sleep(2 * time.Second)
	t.Log(lru.Access(node1))
	t.Log(lru.Access(node2))
	printList(t, lru)
}

func TestEliminate(t *testing.T) {
	lru := &LRU{MaxSize: 100}
	for i := 0; i < 10; i++ {
		lru.AddNewNode(fmt.Sprintf("key%d", i), value{"1234567890"})
	}
	printList(t, lru)
	t.Log(lru.CurSize())
	lru.AddNewNode("key10", value{"hello"})
	printList(t, lru)
	t.Log(lru.CurSize())
}
