package crawler

import (
	"crypto/md5"
	"encoding/hex"
	"reflect"
	"strings"
	"testing"
)

const kHttp = "http://"
const kHttps = "https://"

func ShortArray[T any]() []T {
	return make([]T, 0, 10)
}

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

func toCsv(raw string) string {
	if strings.ContainsAny(raw, ",\"\n") {
		return "\"" + strings.ReplaceAll(raw, "\"", "\"\"") + "\""
	}
	return raw
}
