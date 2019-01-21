package main

import (
	"fmt"

	"github.com/nuczzz/lru"
)

type value struct {
	data string
}

func (v value) Len() int64 {
	return int64(len(v.data))
}

func main() {
	l := lru.NewLRU(100, 40)
	node, _ := l.AddNewNode("key", value{"value"})
	fmt.Printf("%#v\n", node)
	newNode, _ := l.Access(node)
	fmt.Printf("%#v\n", newNode)
}
