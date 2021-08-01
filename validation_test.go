package lrw

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"testing"
)

func validateFields(validate *validator.Validate, tag string, valid []string, invalid []string, t *testing.T) {
	for _, validValue := range valid {
		errs := validate.Var(validValue, tag)
		if errs != nil {
			t.Fatalf("%s should be valid on %s", validValue, tag)
		}
	}
	for _, invalidValue := range invalid {
		errs := validate.Var(invalidValue, tag)
		if errs == nil {
			t.Fatalf("%s should be invalid on %s", invalidValue, tag)
		}
	}
}

func TestNewValidator(t *testing.T) {
	validate := NewValidator()
	if validate == nil {
		t.Fatal("validate must be != nil")
	}
}

func TestIsBoolean(t *testing.T) {
	validate := NewValidator()
	var valid = []string{
		"true",
		"false",
	}
	var invalid = []string{
		"1",
		"True",
		"False",
		"",
	}
	validateFields(validate, IsBooleanTag, valid, invalid, t)
}

func TestIsGinMode(t *testing.T) {
	validate := NewValidator()
	var valid = []string{
		gin.ReleaseMode,
		gin.DebugMode,
		gin.TestMode,
	}
	var invalid = []string{
		"something else",
		"Release",
		"Debug",
		"Test",
		"",
	}
	validateFields(validate, IsGinModeTag, valid, invalid, t)
}

func TestIsJSONStringArray(t *testing.T) {
	validate := NewValidator()
	var valid = []string{
		`[]`,
		`["a"]`,
		`["a", "b"]`,
		`["a", "b", "c"]`,
	}
	var invalid = []string{
		`[`,
		`[""]`,
		`["a", "", "c"]`,
		"",
	}
	validateFields(validate, IsJsonStringArrayTag, valid, invalid, t)
}
