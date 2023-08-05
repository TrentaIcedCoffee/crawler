package crawler

import (
	"testing"
)

func TestShortArrayReturnsEmptyArrayWithDesiredSmallCap(t *testing.T) {
	arr := ShortArray[int]()
	expectEqual(t, len(arr), 0)
	expectEqual(t, cap(arr), 10)
}

func TestHash(t *testing.T) {
	expectEqual(t, hash("abc"), "900150983cd24fb0d6963f7d28e17f72")
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
