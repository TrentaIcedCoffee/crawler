package crawler

import (
	"testing"
)

func TestShortArrayReturnsEmptyArrayWithDesiredSmallCap(t *testing.T) {
	arr := ShortArray[int]()
	ExpectEqualInTest(t, len(arr), 0)
	ExpectEqualInTest(t, cap(arr), 10)
}

func TestMd5ReturnsHashEqualComputed(t *testing.T) {
	ExpectEqualInTest(t, Md5("abc"), "900150983cd24fb0d6963f7d28e17f72")
}

func TestToCsvReturnsRawWhenNoSpecialChar(t *testing.T) {
	ExpectEqualInTest(t, ToCsvRow("abc"), "abc")
}

func TestToCsvQuotesComma(t *testing.T) {
	ExpectEqualInTest(t, ToCsvRow("a,bc"), "\"a,bc\"")
}

func TestToCsvQuotesNewline(t *testing.T) {
	ExpectEqualInTest(t, ToCsvRow("a\nbc"), "\"a\nbc\"")
}

func TestToCsvDoubleQuotes(t *testing.T) {
	ExpectEqualInTest(t, ToCsvRow("a\"bc"), "\"a\"\"bc\"")
}
