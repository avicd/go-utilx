package bufx

import (
	"github.com/avicd/go-utilx/datax"
	"sync"
)

type LruCache[K comparable, V any] struct {
	Size        int
	store       map[K]*datax.LinkedNode[Entry[K, V]]
	first       *datax.LinkedNode[Entry[K, V]]
	last        *datax.LinkedNode[Entry[K, V]]
	mutex       sync.RWMutex
	firstLocker sync.Mutex
}

func (it *LruCache[K, V]) Get(key K) (V, bool) {
	var ret V
	var exist bool
	it.mutex.RLock()
	if node, ok := it.store[key]; ok {
		it.asFirst(node)
		ret = node.Item.Value
		exist = true
	}
	it.mutex.RUnlock()
	return ret, exist
}

func (it *LruCache[K, V]) asFirst(node *datax.LinkedNode[Entry[K, V]]) {
	it.firstLocker.Lock()
	if it.first != node {
		if it.first == nil {
			node.Remove()
			it.first = node
			it.last = node
		} else {
			if it.last == node {
				it.last = node.Pre
			}
			node.Remove()
			it.first.Insert(node)
			it.first = node
		}
	}
	it.firstLocker.Unlock()
}

func (it *LruCache[K, V]) removeLast() {
	if it.last != nil {
		delete(it.store, it.last.Item.Key)
		if it.last.Pre != nil {
			it.last = it.last.Pre
			it.last.Next = nil
		} else {
			it.last = nil
		}
	}
}

func (it *LruCache[K, V]) Put(key K, val V) {
	it.mutex.Lock()
	node := it.store[key]
	if node == nil {
		if it.Size > 0 && len(it.store)+1 > it.Size {
			it.removeLast()
		}
		node = &datax.LinkedNode[Entry[K, V]]{}
	}
	if it.store == nil {
		it.store = map[K]*datax.LinkedNode[Entry[K, V]]{}
	}
	node.Item = Entry[K, V]{Key: key, Value: val}
	it.store[key] = node
	it.asFirst(node)
	it.mutex.Unlock()
}

func (it *LruCache[K, V]) Remove(key K) bool {
	it.mutex.Lock()
	if node, ok := it.store[key]; ok {
		delete(it.store, node.Item.Key)
		if it.first == node {
			it.first = node.Next
		}
		if it.last == node {
			it.last = node.Pre
		}
		node.Remove()
	}
	it.mutex.Unlock()
	return false
}

func (it *LruCache[K, V]) Len() int {
	return len(it.store)
}

func (it *LruCache[K, V]) Clear() {
	it.mutex.Lock()
	it.store = map[K]*datax.LinkedNode[Entry[K, V]]{}
	it.first = nil
	it.last = nil
	it.mutex.Unlock()
}
