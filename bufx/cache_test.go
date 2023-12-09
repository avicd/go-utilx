package bufx

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestSimpleCache_Put(t *testing.T) {
	size := 1000
	cache := &SimpleCache[int, string]{Size: size}
	count := 0
	var keys []int
	cond := sync.NewCond(&sync.Mutex{})
	mutex := sync.Mutex{}
	testTotal := 5000
	for i := 0; i < 100; i++ {
		go func(id int) {
			for j := 0; j < testTotal/100; j++ {
				key := id*100 + j
				cache.Put(key, fmt.Sprintf("key:%d", key))
				mutex.Lock()
				count++
				if count > testTotal-size {
					keys = append(keys, key)
				}
				mutex.Unlock()
				if count == testTotal {
					cond.Signal()
				}
			}
		}(i)
	}
	cond.L.Lock()
	cond.Wait()
	assert.Equal(t, len(keys), cache.Len())
	for _, id := range keys {
		go func(key int) {
			val, ok := cache.Get(key)
			assert.Equal(t, true, ok)
			assert.Equal(t, fmt.Sprintf("key:%d", key), val)
		}(id)
	}
}
