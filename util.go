package crawler

import (
	"crypto/md5"
	"encoding/hex"
	"reflect"
	"strings"
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
