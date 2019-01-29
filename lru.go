package lru

import (
	"time"
)

type Value interface {
	Len() int64
}

// Node lru data node
type Node struct {
	// Key key of node
	Key interface{}
	// Length length of value
	Length int64
	// Value value of node
	Value Value
	// Extra extra field of node.
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
	// curSize current all value size of lru(bytes).
	curSize int64
	// TTL time to live(second)
	TTL int64

	// AddNodeCallBack callback after lru add node
	AddNodeCallBack func(node *Node)

	// deleteNodeCallBack if this field is not nil, when a node is deleted,
	// from lru linked list and callback this function.For example, in http
	// cache, we often use map+lru-list to cache some static files, if we
	// just delete data from lru list but not delete metadata in map, we may
	// get wrong result. So we can write a callback function like:
	// var globalMap = make(map[interface{}]*Node)
	// func DeleteMetadata(key interface{}) error {
	// 		delete(globalMap, key)
	//		return nil
	// }
	// lru := NewLRUWithCallback(300, DeleteMetadata)
	DeleteNodeCallBack func(key interface{}) error

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
	// space rather than the oldest value's length.
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
	GetValue func(interface{}) (Value, error)

	// header header of circular double linked list.
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

func (lru *LRU) CurSize() int64 {
	return lru.curSize
}

// NewLRUWithCallback return a new LRU instance with delete-node-callback.
func NewLRUWithCallback(ttl int64, callback func(interface{}) error) *LRU {
	lru := &LRU{DeleteNodeCallBack: callback}
	if ttl > 0 {
		lru.TTL = ttl
	}
	return lru
}

func (lru *LRU) newNode(key interface{}, value Value, extra ...interface{}) *Node {
	node := &Node{
		Key:   key,
		Value: value,
	}
	if len(extra) > 0 {
		node.Extra = extra[0]
	}
	if value != nil {
		node.Length = value.Len()
	}
	if lru.TTL > 0 {
		node.expire = time.Now().Unix() + lru.TTL
	}
	return node
}

func (lru *LRU) Add(node *Node) {
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
			lru.header = node
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
func (lru *LRU) AddNewNode(key interface{}, value Value, extra ...interface{}) error {
	diff := lru.curSize - lru.MaxSize
	if value != nil {
		diff += value.Len()
	}
	if diff > 0 {
		var err error
		if lru.EliminateLength != nil {
			err = lru.eliminate(lru.EliminateLength())
		} else {
			err = lru.eliminate(diff)
		}
		if err != nil {
			return err
		}
	}

	if lru.SetValue != nil {
		if err := lru.SetValue(key, value); err != nil {
			return err
		}
	}

	node := lru.newNode(key, value, extra...)
	lru.Add(node)
	if lru.AddNodeCallBack != nil {
		lru.AddNodeCallBack(node)
	}
	lru.curSize += node.Length
	return nil
}

// RemoveToHead move node to lru double linked list head
func (lru *LRU) moveToHead(node *Node) {
	if node != lru.header {
		node.previous.next = node.next
		node.next.previous = node.previous
		lru.Add(node)
	}
}

func (lru *LRU) Access(node *Node) (*Node, error) {
	now := time.Now().Unix()
	if lru.TTL > 0 {
		if node.expire < now {
			if err := lru.delete(node); err != nil {
				return nil, err
			}
			return nil, nil
		}
	}
	if lru.GetValue != nil {
		var err error
		node.Value, err = lru.GetValue(node.Key)
		if err != nil {
			return nil, err
		}
	}

	lru.moveToHead(node)
	node.AccessTime = now
	node.AccessCount++
	return node, nil
}

// eliminate eliminate old node
func (lru *LRU) eliminate(length int64) error {
	for lru.header != nil && length > 0 {
		node := lru.header.previous
		length -= node.Length
		if err := lru.delete(node); err != nil {
			return err
		}
	}
	return nil
}

func (lru *LRU) replace(node *Node, value Value, extra ...interface{}) error {
	if value != nil {
		node.Length = value.Len()
	} else {
		node.Length = 0
	}

	if lru.SetValue != nil {
		if err := lru.SetValue(node.Key, value); err != nil {
			return err
		}
	} else {
		node.Value = value
	}

	if len(extra) > 0 {
		node.Extra = extra[0]
	}
	node.AccessCount = 0
	node.AccessTime = 0
	if lru.TTL > 0 {
		node.expire = time.Now().Unix() + lru.TTL
	}
	return nil
}

func (lru *LRU) Replace(node *Node, value Value, extra ...interface{}) error {
	lru.moveToHead(node)
	length := node.Length
	diff := lru.curSize - node.Length - lru.MaxSize
	if value != nil {
		diff += value.Len()
		length -= value.Len()
	}
	if diff > 0 {
		var err error
		if lru.EliminateLength != nil {
			err = lru.eliminate(lru.EliminateLength())
		} else {
			err = lru.eliminate(diff)
		}
		if err != nil {
			return err
		}
	}
	lru.curSize -= length
	return lru.replace(node, value, extra...)
}

// delete delete node from lru double linked list,
// node MUST not nil and is REAL node in lru list.
// Remove node reference to avoid escape GC. After
// removed linked node, DeleteNodeCallBack will be
// executed.
func (lru *LRU) delete(node *Node) error {
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

	if lru.DeleteNodeCallBack != nil {
		// delete node callback
		return lru.DeleteNodeCallBack(node.Key)
	}
	return nil
}

// Traversal return lru node list
func (lru *LRU) Traversal() []*Node {
	var list []*Node
	node := lru.header
	for node != nil {
		list = append(list, node)
		if node.next == lru.header {
			break
		}
		node = node.next
	}
	return list
}
