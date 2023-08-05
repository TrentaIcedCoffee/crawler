package crawler

import (
	"crypto/md5"
	"encoding/hex"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func Md5(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

func ShortArray[T any]() []T {
	return make([]T, 0, 10)
}

func ExpectEqualInTest(t *testing.T, actual interface{}, expect interface{}) {
	if !reflect.DeepEqual(actual, expect) {
		t.Fatalf("Expects %v, had %v", expect, actual)
	}
}

func ToCsvRow(raw string) string {
	if strings.ContainsAny(raw, ",\"\n") {
		return "\"" + strings.ReplaceAll(raw, "\"", "\"\"") + "\""
	}
	return raw
}

// ConcurrentSet

type ConcurrentSet struct {
	mutex sync.RWMutex
	data  map[string]interface{}
}

func NewConcurrentSet() *ConcurrentSet {
	return &ConcurrentSet{
		data: make(map[string]interface{}),
	}
}

func (this *ConcurrentSet) Has(key string) bool {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	_, exists := this.data[key]
	return exists
}

func (this *ConcurrentSet) Add(key string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.data[key] = nil
}

func (this *ConcurrentSet) Size() int {
	return len(this.data)
}
