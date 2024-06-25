package parse

import (
	"testing"
	"reflect"
)

func TestArrayBasic(t *testing.T) {
	got := GetArray("foo,bar")
	want := []string { "foo", "bar" }

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func TestArrayNoDelim(t *testing.T) {
	got := GetArray("foo")
	want := []string { "foo" }

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, wanted %v", got, want)
	}
}
