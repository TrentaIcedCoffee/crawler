package crawler

import (
	"testing"
)

func TestShortArrayReturnsEmptyArrayWithDesiredSmallCap(t *testing.T) {
	arr := shortArray[int]()
	expectEqualInTest(t, len(arr), 0)
	expectEqualInTest(t, cap(arr), 10)
}

func TestGetHostOfUrl(t *testing.T) {
	host, err := getHost("https://www.example.com/some/path?query=123")
	if err != nil {
		t.Fatal(err)
	}
	expectEqualInTest(t, host, "www.example.com")

	host, err = getHost("https://example.com/a")
	if err != nil {
		t.Fatal(err)
	}
	expectEqualInTest(t, host, "example.com")
}

func TestSameHost(t *testing.T) {
	is_same_domain, err := isSameHost("https://www.example.com", "https://www.example.com/a/b/c")
	if err != nil {
		t.Fatal(err)
	}
	expectEqualInTest(t, is_same_domain, true)

	is_same_domain, err = isSameHost("http://www.example.com", "http://www.another.com/a/b/c")
	if err != nil {
		t.Fatal(err)
	}
	expectEqualInTest(t, is_same_domain, false)
}
