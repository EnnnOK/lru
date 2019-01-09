package lru

type Node struct {
	data interface{}

	// double linked list
	previous *Node
	next     *Node
}

type LRU struct {
	header *Node
	tail   *Node
}

func NewLRU() *LRU {
	return &LRU{}
}

func NewNode(data interface{}) *Node {
	return &Node{data: data}
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

func (lru *LRU) RemoveToHead(node *Node) {
	// todo
}

func (lru *LRU) Delete(node *Node) {
	// todo
}

func (lru *LRU) Get(node *Node) interface{} {
	return node.data
}
