package crawler

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

const kHttp = "http://"
const kHttps = "https://"

func expectEqual(t *testing.T, actual interface{}, expect interface{}) {
	if !reflect.DeepEqual(actual, expect) {
		t.Fatalf("Expects %v, had %v", expect, actual)
	}
}

func getDomain(url string) (string, error) {
	var idx int
	if strings.HasPrefix(url, kHttp) {
		idx = len(kHttp)
	} else if strings.HasPrefix(url, kHttps) {
		idx = len(kHttps)
	} else {
		return "", errors.New(fmt.Sprintf("Did not find protocal in url, %s", url))
	}

	end := idx
	for ; end < len(url); end++ {
		if url[end] == '/' {
			break
		}
	}

	return url[:end], nil
}

func joinUrl(domain string, url string) (string, error) {
	if strings.HasPrefix(url, kHttp) || strings.HasPrefix(url, kHttps) {
		return url, nil
	} else if strings.HasPrefix(url, "/") {
		return domain + url, nil
	} else if strings.HasPrefix(url, "#") {
		return domain + url, nil
	}
	return "", errors.New(fmt.Sprintf("Cannot join domain %s with url %s", domain, url))
}
