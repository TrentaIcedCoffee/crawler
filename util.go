package crawler

import (
	"reflect"
	"strings"
	"sync"
	"testing"
)

func shortArray[T any]() []T {
	return make([]T, 0, 10)
}

func expectEqualInTest(t *testing.T, actual interface{}, expect interface{}) {
	if !reflect.DeepEqual(actual, expect) {
		t.Fatalf("Expects %v, had %v", expect, actual)
	}
}

func toCsvRow(raw string) string {
	if strings.ContainsAny(raw, ",\"\n") {
		return "\"" + strings.ReplaceAll(raw, "\"", "\"\"") + "\""
	}
	return raw
}

// ConcurrentSet

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

// Returns true if added a new entry.
func (this *concurrentSet) add(key string) bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	_, exists := this.data[key]
	if exists {
		return false
	} else {
		this.data[key] = nil
		return true
	}
}

func (this *concurrentSet) size() int {
	return len(this.data)
}
