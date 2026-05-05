package helper

import (
	"reflect"
	"testing"
)

func AssertEqual(t *testing.T, expected any, actual any, message string) {
	t.Helper()

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("%s: expected=%#v actual=%#v", message, expected, actual)
	}
}
