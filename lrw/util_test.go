package lrw

import (
	"github.com/lripardo/lrw"
	"testing"
)

func TestGetUint64ParamFromContext(t *testing.T) {
	context := &lrw.ContextMock{
		Params: map[string]string{
			"id": "0",
		},
	}
	_, err := GetUint64ParamFromContext("id", context)
	if err == nil {
		t.Fatal("id param with zero value must return error")
	}
	context.Params["id"] = ""
	_, err = GetUint64ParamFromContext("id", context)
	if err == nil {
		t.Fatal("id param with empty value must return error")
	}
	context.Params["id"] = "-1"
	_, err = GetUint64ParamFromContext("id", context)
	if err == nil {
		t.Fatal("id param with negative values must return error")
	}
	context.Params["id"] = "test"
	_, err = GetUint64ParamFromContext("id", context)
	if err == nil {
		t.Fatal("id param with non numeric values must return error")
	}
	context.Params["id"] = "1"
	_, err = GetUint64ParamFromContext("id", context)
	if err != nil {
		t.Fatal("id param should accept numbers greater than zero")
	}
}

func TestRoles(t *testing.T) {
	context := &lrw.ContextMock{Data: map[string]interface{}{
		UserInfoKey: UserConfig{
			User: UserInfo{
				Role: "user",
			},
		},
	}}
	response := Roles("admin")(context)
	if response == nil {
		t.Fatal("user without role shouldn't return nil")
	}
	response = Roles("user")(context)
	if response != nil {
		t.Fatal("user with role should return nil")
	}
}
