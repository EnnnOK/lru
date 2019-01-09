package lru

import "time"

// Node lru data node
type Node struct {
	// Key key of node
	Key interface{}
	// Value value of node
	Value interface{}
	// Extra extra field of node
	Extra interface{}
	// AccessTime timestamp of access time
	AccessTime int64
	// AccessCount access count
	AccessCount int64
	// expire timestamp of expire if this field greater than 0,
	// and if node is expired, delete from lru linked list.
	expire int64

	// double linked list
	previous *Node
	next     *Node
}

type LRU struct {
	// ttl time to live(second)
	ttl int64

	// ExpireCallBack if this field is not nil, when a node is expired,
	// delete it from lru linked list and callback this function.For example,
	// in http cache, we often use map+lru-list to cache some static files,
	// if we just delete data from lru list but not delete metadata in map,
	// we may get wrong result.So we can write a callback function like:
	// var globalMap = make(map[interface{}]*Node)
	// func DeleteMetadata(key interface{}) {
	// 		delete(globalMap, key)
	// }
	// lru := NewLRUWithCallback(300, DeleteMetadata)
	expireCallBack func(key interface{})

	header *Node
	tail   *Node
}

// NewLRU return a new LRU instance, set time-to-live of lru node if ttl is greater than 0.
func NewLRU(ttl int64) *LRU {
	lru := &LRU{}
	if ttl > 0 {
		lru.ttl = ttl
	}
	return lru
}

// NewLRUWithCallback return a new LRU instance with expire callback.
func NewLRUWithCallback(ttl int64, callback func(interface{})) *LRU {
	lru := &LRU{expireCallBack: callback}
	if ttl > 0 {
		lru.ttl = ttl
	}
	return lru
}

func (lru *LRU) NewNode(key, value interface{}, extra ...interface{}) *Node {
	node := &Node{
		Key:   key,
		Value: value,
	}
	if lru.ttl > 0 {
		node.expire = time.Now().Unix() + lru.ttl
	}

	if len(extra) > 0 {
		node.Extra = extra[0]
	}
	return node
}

func (lru *LRU) Add(node *Node) {
	if lru.header == nil {
		lru.header = node
		lru.tail = node
	} else {
		node.next = lru.header
		lru.header.previous = node
		lru.header = node
	}
}

func (lru *LRU) removeToHead(node *Node) {
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

func (lru *LRU) Delete(node *Node) {
	// todo
}
