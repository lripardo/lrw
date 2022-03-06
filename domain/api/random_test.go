package api_test

import (
	"github.com/lripardo/lrw/domain/api"
	"testing"
)

func contains(s string, r int32) bool {
	for _, c := range s {
		if c == r {
			return true
		}
	}
	return false
}

func TestRandomKey(t *testing.T) {
	key := api.RandomKey()
	if len(key) != api.MinSecretKeyLength {
		t.Fatalf("generated key must have a length of %d", api.MinSecretKeyLength)
	}
	for _, l := range key {
		if !contains(key, l) {
			t.Fatalf("rune %d is not on alphabet key", l)
		}
	}
	key2 := api.RandomKey()
	if key == key2 {
		t.Fatal("generated keys needs to be random key")
	}
}
