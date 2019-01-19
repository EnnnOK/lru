package lru

import (
	"time"
)

// Node lru data node
type Node struct {
	// Key key of node
	Key interface{}
	// length length of value
	length int64
	// Value value of node
	Value interface{}
	// Extra extra field of node. todo
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

// LRU manager of lru node with circular double linked list
type LRU struct {
	// MaxSize all value max size of lru(bytes).
	MaxSize int64
	// TTL time to live(second)
	TTL int64

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
	DeleteNodeCallBack func(key interface{})

	// EliminateLength eliminate length of lru double linked list,
	// if this field is nil, we calculate eliminate length with the
	// length of new node. you can define EliminateLength function
	// like this:
	// lru := &LRU{MaxSize: 100, ...}
	// lru.EliminateLength = func() int64 {
	// 		// closure with lru.MaxSize
	// 		return lru.MaxSize/10
	// }
	// so when eliminate happened, we release 1/10 of lru.MaxSize
	// space.
	EliminateLength func() int64

	// SetValue set node value, if SetValue is nil, store value in memory.
	// if length of value is big like elephant, we suggest you define
	// SetValue and store value in disk, for example:
	// lru := &LRU{...}
	// lru.SetValue = func(key, value interface) error {
	// 		fd, _ := os.Open(key.(string))
	// 		fd.Write(value.([]byte))
	//		return nil
	// }
	SetValue func(key, value interface{}) error

	// GetValue like SetValue, define the function of get value.
	GetValue func(interface{}) (interface{}, error)

	// curSize current all value size of lru(bytes).
	curSize int64

	// circular double linked list header.
	header *Node
}

// NewLRU return a new LRU instance, set time-to-live of lru node if ttl is greater than 0.
func NewLRU(maxSize, ttl int64) *LRU {
	lru := &LRU{MaxSize: maxSize}
	if ttl > 0 {
		lru.TTL = ttl
	}
	return lru
}

// NewLRUWithCallback return a new LRU instance with delete-node-callback.
func NewLRUWithCallback(ttl int64, callback func(interface{})) *LRU {
	lru := &LRU{DeleteNodeCallBack: callback}
	if ttl > 0 {
		lru.TTL = ttl
	}
	return lru
}

// getInterfaceLength get length of interface
func getInterfaceLength(i interface{}) int64 {
	// todo
	return 0
}

func (lru *LRU) newNode(key, value interface{}, extra ...interface{}) *Node {
	node := &Node{
		Key:    key,
		length: getInterfaceLength(value),
		Value:  value,
		Extra:  extra,
	}
	if lru.TTL > 0 {
		node.expire = time.Now().Unix() + lru.TTL
	}
	return node
}

func (lru *LRU) add(node *Node) {
	if lru.header == nil {
		node.previous = node
		node.next = node
		lru.header = node
	} else {
		// just one node
		if lru.header.next == lru.header {
			lru.header.previous = node
			lru.header.next = node
			node.previous = lru.header
			node.next = lru.header
		} else {
			lru.header.previous.next = node
			node.previous = lru.header.previous
			node.next = lru.header
			lru.header.previous = node
			lru.header = node
		}
	}
}

// NewNode return nil if lru.SetValue is nil or lru.SetValue return nil
func (lru *LRU) AddNewNode(key, value interface{}, extra ...interface{}) (*Node, error) {
	diff := lru.curSize + getInterfaceLength(value) - lru.MaxSize
	if diff > 0 {
		lru.eliminate(diff)
	}

	if lru.SetValue != nil {
		if err := lru.SetValue(key, value); err != nil {
			return nil, err
		}
	}

	node := lru.newNode(key, value, extra...)
	lru.add(node)
	return node, nil
}

// RemoveToHead move node to lru double linked list head
func (lru *LRU) moveToHead(node *Node) {
	if node != lru.header {
		node.previous.next = node.next
		node.next.previous = node.previous
		lru.add(node)
	}
}

func (lru *LRU) Access(node *Node) *Node {
	now := time.Now().Unix()
	if node.expire < now {
		lru.Delete(node)
		return nil
	}
	if lru.GetValue != nil {
		var err error
		node.Value, err = lru.GetValue(node.Key)
		if err != nil {
			lru.Delete(node)
			return nil
		}
	}

	lru.moveToHead(node)
	node.AccessTime = now
	node.AccessCount++
	return node
}

// eliminate eliminate old node
func (lru *LRU) eliminate(length int64) {
	for lru.header != nil && length > 0 {
		node := lru.header.previous
		length -= int64(node.length)
		lru.Delete(node)
	}
}

// Delete delete node from lru double linked list,
// node MUST not nil and is REAL node in lru list.
// Remove node reference to avoid escape GC.
func (lru *LRU) Delete(node *Node) {
	if lru.header.next == lru.header {
		lru.header.previous = nil
		lru.header.next = nil
		lru.header = nil
	} else {
		node.previous.next = node.next
		node.next.previous = node.previous
		node.previous = nil
		node.next = nil
	}

	// delete node callback
	if lru.DeleteNodeCallBack != nil {
		lru.DeleteNodeCallBack(node.Key)
	}
}
