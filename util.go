package crawler

import (
	"errors"
	"fmt"
	net_url "net/url"
	"reflect"
	"sync"
	"testing"
	"unicode/utf8"
)

func isUtf8(strings []string) bool {
	for _, str := range strings {
		if !utf8.ValidString(str) {
			return false
		}
	}
	return true
}

func shortArray[T any]() []T {
	return make([]T, 0, 10)
}

func expectEqualInTest(t *testing.T, actual interface{}, expect interface{}) {
	if !reflect.DeepEqual(actual, expect) {
		t.Fatalf("Expects %v, had %v", expect, actual)
	}
}

func getHost(url string) (string, error) {
	url_object, err := net_url.Parse(url)
	if err != nil {
		return "", err
	}
	return url_object.Hostname(), nil
}

func isSameHost(parent_url string, child_url string) (bool, error) {
	parent_host, err := getHost(parent_url)
	if err != nil {
		return false, err
	}

	child_host, err := getHost(child_url)
	if err != nil {
		return false, err
	}

	if parent_host == "" {
		return false, errors.New(fmt.Sprintf("Host is empty for parent url %s", parent_url))
	} else if child_host == "" {
		return false, errors.New(fmt.Sprintf("Host is empty for child url %s", child_url))
	}

	return parent_host == child_host, nil
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
