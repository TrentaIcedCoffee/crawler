package crawler

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

const kHttp = "http://"
const kHttps = "https://"

func hash(input string) string {
	hash := md5.New()
	hash.Write([]byte(input))
	return hex.EncodeToString(hash.Sum(nil))
}

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
	} else if strings.HasPrefix(url, "./") {
		return domain + url[1:], nil
	} else if strings.HasPrefix(url, "../") {
		parent_end, child_start := len(domain)-1, 0
		for child_start+3 <= len(url) && url[child_start:child_start+3] == "../" {
			child_start += 3
			for parent_end > 0 && domain[parent_end] != '/' {
				parent_end--
			}
			if parent_end == 0 || parent_end == len(kHttps)-1 && strings.HasPrefix(domain, kHttps) || parent_end == len(kHttp)-1 && strings.HasPrefix(domain, kHttp) {
				return "", errors.New(fmt.Sprintf("joinUrl cannot slide backward on domain to join %s, %s", domain, url))
			}
			parent_end--
		}
		if child_start < len(url) {
			return domain[:parent_end+1] + "/" + url[child_start:], nil
		} else {
			return domain[:parent_end+1], nil
		}

	}
	return "", errors.New(fmt.Sprintf("Cannot join domain %s with url %s", domain, url))
}

func toCsv(raw string) string {
	if strings.ContainsAny(raw, ",\"\n") {
		return "\"" + strings.ReplaceAll(raw, "\"", "\"\"") + "\""
	}
	return raw
}
