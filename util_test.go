package crawler

import (
	"testing"
)

func TestGetDomainReturnsDomainExpectedly(t *testing.T) {
	domain, err := getDomain("http://abc.com/a/b/c")
	if err != nil {
		t.Fatalf("Error in getDomain")
	}
	expectEqual(t, domain, "http://abc.com")
}

func TestGetDomainReturnsDomainWhenInputOnlyDomain(t *testing.T) {
	domain, err := getDomain("https://abc.com")
	if err != nil {
		t.Fatalf("Error in getDomain")
	}
	expectEqual(t, domain, "https://abc.com")
}

func TestGetDomainErrorsWhenNoProtocal(t *testing.T) {
	_, err := getDomain("abc.com")
	if err == nil {
		t.Fatalf("Url without protocal should returns an error")
	}
}
