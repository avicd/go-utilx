package bufx

import (
	"sync"
)

type Cache[K comparable, V any] interface {
	Get(key K) (V, bool)
	Put(key K, val V)
	Remove(key K) bool
	Len() int
	Clear()
}

type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

type SimpleCache[K comparable, V any] struct {
	Size  int
	mutex sync.RWMutex
	store map[K]*Entry[K, V]
	keys  []K
}

func (it *SimpleCache[K, V]) Get(key K) (V, bool) {
	var ret V
	var exist bool
	it.mutex.RLock()
	if entry, ok := it.store[key]; ok {
		ret = entry.Value
		exist = true
	}
	it.mutex.RUnlock()
	return ret, exist
}

func (it *SimpleCache[K, V]) Put(key K, val V) {
	it.mutex.Lock()
	if it.store == nil {
		it.store = map[K]*Entry[K, V]{}
	}
	entry := it.store[key]
	if entry == nil {
		if it.Size > 0 && len(it.store)+1 > it.Size {
			delete(it.store, it.keys[0])
			it.keys = it.keys[1:]
		}
		entry = &Entry[K, V]{Key: key}
	}
	entry.Value = val
	it.store[key] = entry
	it.keys = append(it.keys, key)
	it.mutex.Unlock()
}

func (it *SimpleCache[K, V]) Remove(key K) bool {
	var ret bool
	it.mutex.Lock()
	if _, ok := it.store[key]; ok {
		delete(it.store, key)
		ret = true
	}
	it.mutex.Unlock()
	return ret
}

func (it *SimpleCache[K, V]) Len() int {
	return len(it.store)
}

func (it *SimpleCache[K, V]) Clear() {
	it.mutex.Lock()
	it.store = map[K]*Entry[K, V]{}
	it.mutex.Unlock()
}
