package api_test

import (
	"github.com/lripardo/lrw/domain/api"
	"testing"
)

func TestNewKey(t *testing.T) {
	name := "my name"
	tag := "my tag"
	value := "my value"
	key := api.NewKey(name, tag, value)
	if key.Name() != name || key.Validation() != tag || key.Value() != value {
		t.Fatal("key vars needs to be equals to property")
	}
}
