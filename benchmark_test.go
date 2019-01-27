package lru

import "testing"

// go test -run=^^$ -bench=^BenchmarkLRUAdd$ -benchmem
func BenchmarkLRUAdd(b *testing.B) {
	lru := &LRU{MaxSize: 100}
	for i := 0; i < b.N; i++ {
		lru.AddNewNode("key", value{"test"})
	}
}
