# lru
LRU algorithm written in golang, NOT thread-safe.

# Example
```
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
	l.AddNewNode("key1", value{"value1"})
	l.AddNewNode("key2", value{"value2"})
	list := l.Traversal()
	for i := range list {
		fmt.Printf("%#v\n", list[i])
	}
}
```