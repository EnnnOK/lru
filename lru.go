package lru

import (
	"time"
)

// Node lru data node
type Node struct {
	// Key key of node
	Key interface{}
	// length length of value
	length int //todo
	// Value value of node
	Value interface{}
	// Extra extra field of node. todo: this field
	Extra interface{}
	// AccessTime timestamp of access time
	AccessTime int64
	// AccessCount access count.
	// todo: current version this field is not used to sort, it just count.
	AccessCount int64
	// expire timestamp of expire if this field greater than 0,
	// and if node is expired, delete from lru linked list.
	expire int64

	// double linked list
	previous *Node
	next     *Node
}

type LRU struct {
	// maxSize all value max size of lru
	maxSize int64
	// curSize current all value size of lru
	curSize int64

	// ttl time to live(second)
	ttl int64

	// deleteNodeCallBack if this field is not nil, when a node is deleted,
	// from lru linked list and callback this function.For example, in http
	// cache, we often use map+lru-list to cache some static files, if we
	// just delete data from lru list but not delete metadata in map, we may
	// get wrong result. So we can write a callback function like:
	// var globalMap = make(map[interface{}]*Node)
	// func DeleteMetadata(key interface{}) {
	// 		delete(globalMap, key)
	// }
	// lru := NewLRUWithCallback(300, DeleteMetadata)
	deleteNodeCallBack func(key interface{})

	// lru double linked list header and tail pointer.
	header *Node
	tail   *Node
}

// NewLRU return a new LRU instance, set time-to-live of lru node if ttl is greater than 0.
func NewLRU(maxSize, ttl int64) *LRU {
	lru := &LRU{maxSize: maxSize}
	if ttl > 0 {
		lru.ttl = ttl
	}
	return lru
}

// NewLRUWithCallback return a new LRU instance with delete-node-callback.
func NewLRUWithCallback(ttl int64, callback func(interface{})) *LRU {
	lru := &LRU{deleteNodeCallBack: callback}
	if ttl > 0 {
		lru.ttl = ttl
	}
	return lru
}

// getInterfaceLength get length of interface
func getInterfaceLength(i interface{}) int {
	// todo
	return 0
}

func newNode(key, value interface{}, extra ...interface{}) *Node {
	node := &Node{
		Key:    key,
		length: getInterfaceLength(value),
		Value:  value,
	}
	if len(extra) > 0 {
		node.Extra = extra[0]
	}
	return node
}

func (lru *LRU) add(node *Node) {
	if lru.ttl > 0 {
		node.expire = time.Now().Unix() + lru.ttl
	}

	if lru.header == nil {
		lru.header = node
		lru.tail = node
	} else {
		node.next = lru.header
		lru.header.previous = node
		lru.header = node
	}
}

func (lru *LRU) NewNode(key, value interface{}, extra ...interface{}) {
	lru.add(newNode(key, value, extra...))
}

// RemoveToHead move node to lru double linked list head
func (lru *LRU) moveToHead(node *Node) {
	if node != lru.header {
		if node == lru.tail {
			lru.tail = lru.tail.previous
			lru.tail.next = nil
		} else {
			node.next.previous = node.previous
			node.previous.next = node.next
		}
		node.previous = nil
		node.next = lru.header
		lru.header.previous = node
		lru.header = node
	}
}

func (lru *LRU) Get(node *Node) *Node {
	now := time.Now().Unix()
	if node.expire < now {
		lru.Delete(node)
		return nil
	}
	lru.moveToHead(node)
	node.AccessTime = now
	node.AccessCount++

	//todo: define get node info with interface
	return node
}

// todo
func (lru *LRU) eliminate(length int64) {
	for lru.tail != nil && length > 0 {
		node := lru.tail
		length -= int64(node.length)
		lru.Delete(node)
	}
}

// Delete delete node from lru double linked list,
// node MUST not nil and is REAL node in lru list.
func (lru *LRU) Delete(node *Node) {
	if lru.header == lru.tail {
		if lru.tail == node {
			lru.header = node.next
			lru.header.previous = nil
			node.next = nil
		} else if lru.header == node {
			lru.tail = node.previous
			lru.tail.next = nil
			node.previous = nil
		} else {
			node.previous.next = node.next
			node.next.previous = node.previous
			node.previous = nil
			node.next = nil
		}
	} else {
		// just one node in lru double linked list
		lru.header = nil
		lru.tail = nil
	}

	// delete node callback
	if lru.deleteNodeCallBack != nil {
		lru.deleteNodeCallBack(node.Key)
	}
}
