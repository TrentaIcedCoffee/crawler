package crawler

import (
	"testing"
)

func TestShortArrayReturnsEmptyArrayWithDesiredSmallCap(t *testing.T) {
	arr := shortArray[int]()
	expectEqualInTest(t, len(arr), 0)
	expectEqualInTest(t, cap(arr), 10)
}

func TestToCsvReturnsRawWhenNoSpecialChar(t *testing.T) {
	expectEqualInTest(t, toCsvRow("abc"), "abc")
}

func TestToCsvQuotesComma(t *testing.T) {
	expectEqualInTest(t, toCsvRow("a,bc"), "\"a,bc\"")
}

func TestToCsvQuotesNewline(t *testing.T) {
	expectEqualInTest(t, toCsvRow("a\nbc"), "\"a\nbc\"")
}

func TestToCsvDoubleQuotes(t *testing.T) {
	expectEqualInTest(t, toCsvRow("a\"bc"), "\"a\"\"bc\"")
}
