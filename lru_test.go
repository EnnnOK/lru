package lru

import (
	"testing"
	"time"
)

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
	for temp != nil {
		t.Logf("%#v", temp)
		temp = temp.next
	}
}

func TestNewNode(t *testing.T) {
	lru := NewLRU(100, 100)
	lru.AddNewNode(123, 234)
	lru.AddNewNode(1231, 2342)
	lru.printLRUList(t)
}

func TestGet(t *testing.T) {
	lru := NewLRU(100, 1)
	node1 := newNode(123, 234)
	node2 := newNode(1231, 2342)
	lru.add(node1)
	lru.add(node2)
	lru.printLRUList(t)
	time.Sleep(2 * time.Second)
	t.Log(lru.Access(node1))
	t.Log(lru.Access(node2))
}
