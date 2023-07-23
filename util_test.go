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

func TestJoinUrlJoinsAbsoluteUrl(t *testing.T) {
	url, err := joinUrl("http://a.com", "https://b.com")
	if err != nil {
		t.FailNow()
	}
	expectEqual(t, url, "https://b.com")
}

func TestJoinUrlJoinsRelativeUrl(t *testing.T) {
	url, err := joinUrl("https://a.com/b", "../c")
	if err != nil {
		t.FailNow()
	}
	expectEqual(t, url, "https://a.com/c")
}

func TestJoinUrlJoinsRelativeUrlBackwardOnly(t *testing.T) {
	url, err := joinUrl("https://a.com/b", "../")
	if err != nil {
		t.FailNow()
	}
	expectEqual(t, url, "https://a.com")
}

func TestJoinUrlErrorsWhenRelativeUrlInvalid(t *testing.T) {
	url, err := joinUrl("https://a.com/b", "../../c")
	if err == nil {
		t.Fatalf("Expect error but got %s", url)
	}
}

func TestToCsvReturnsRawWhenNoSpecialChar(t *testing.T) {
	expectEqual(t, toCsv("abc"), "abc")
}

func TestToCsvQuotesComma(t *testing.T) {
	expectEqual(t, toCsv("a,bc"), "\"a,bc\"")
}

func TestToCsvQuotesNewline(t *testing.T) {
	expectEqual(t, toCsv("a\nbc"), "\"a\nbc\"")
}

func TestToCsvDoubleQuotes(t *testing.T) {
	expectEqual(t, toCsv("a\"bc"), "\"a\"\"bc\"")
}
