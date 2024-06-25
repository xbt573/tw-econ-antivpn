package env

import (
	"testing"
)

func TestGet(t *testing.T) {
	t.Setenv("FOO", "bar")

	got := Get("FOO")
	want := "bar"

	if got != want {
		t.Errorf("got %s, wanted %s", got, want)
	}
}

func TestGetDefault(t *testing.T) {
	got := GetDefault("FOO", "foo")
	want := "foo"

	if got != want {
		t.Errorf("got %s, wanted %s", got, want)
	}
}
