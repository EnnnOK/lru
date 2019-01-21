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
	callback := func(interface{}) {}
	lru := NewLRUWithCallback(0, callback)
	t.Logf("%#v", lru)
}

func (lru *LRU) printLRUList(t *testing.T) {
	temp := lru.header
	for temp != nil && temp.next != lru.header {
		t.Logf("%#v", temp)
		temp = temp.next
	}
}

func TestNewNode(t *testing.T) {
	lru := NewLRU(100, 100)
	lru.AddNewNode(123, value{"234"})
	lru.AddNewNode(1231, value{"2342"})
	lru.printLRUList(t)
}

func TestGet(t *testing.T) {
	lru := NewLRU(100, 1)
	node1 := lru.newNode(123, value{"234"})
	node2 := lru.newNode(1231, value{"2342"})
	lru.add(node1)
	lru.add(node2)
	lru.printLRUList(t)
	time.Sleep(2 * time.Second)
	t.Log(lru.Access(node1))
	t.Log(lru.Access(node2))
}
