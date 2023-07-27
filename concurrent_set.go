package crawler

import (
	"sync"
)

type concurrentSet struct {
	mutex sync.RWMutex
	data  map[string]interface{}
}

func newConcurrentSet() *concurrentSet {
	return &concurrentSet{
		data: make(map[string]interface{}),
	}
}

func (this *concurrentSet) has(key string) bool {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	_, exists := this.data[key]
	return exists
}

func (this *concurrentSet) add(key string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.data[key] = nil
}

func (this *concurrentSet) size() int {
	return len(this.data)
}
