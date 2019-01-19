# lru
LRU algorithm written in golang, NOT thread-safe.

# Example
```
package main

import (
	"fmt"

	"github.com/nuczzz/lru"
)

func main() {
	l := lru.NewLRU(100, 40)
	node, _ := l.AddNewNode("key", "value")
	fmt.Printf("%#v\n", node)
	new := l.Access(node)
	fmt.Printf("%#v\n", new)
}
```