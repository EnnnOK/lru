package main

import (
	"fmt"
	"lru"
)

type MemCache struct {
	m map[interface{}]*lru.Node

	lru *lru.LRU
}

func (mc *MemCache) Set(key, value interface{}) error {
	node, err := mc.lru.AddNewNode(key, value)
	if err != nil {
		return err
	}
	mc.m[key] = node
	return nil
}

func (mc *MemCache) Get(key interface{}) interface{} {
	if node, ok := mc.m[key]; ok {
		return mc.lru.Access(node).Value
	}
	return nil
}

func NewMemCache(maxSize, ttl int64) *MemCache {
	return &MemCache{
		m:   make(map[interface{}]*lru.Node),
		lru: lru.NewLRU(maxSize, ttl),
	}
}

func main() {
	memCache := NewMemCache(100, 10)
	memCache.Set("key1", "value1")
	fmt.Println(memCache.Get("key1").(string))
}
