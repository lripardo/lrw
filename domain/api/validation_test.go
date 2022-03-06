package api_test

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lripardo/lrw/domain/api"
	"testing"
)

func validateFields(t *testing.T, validate *validator.Validate, tag string, valid, invalid []string) {
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

func TestIsBoolean(t *testing.T) {
	validatorBoolean := api.IsBoolean()
	validate := api.NewValidator(validatorBoolean)
	valid := []string{
		"true",
		"false",
	}
	invalid := []string{
		"1",
		"True",
		"False",
		"",
	}
	validateFields(t, validate, validatorBoolean.Tag, valid, invalid)
}

func TestIsGinMode(t *testing.T) {
	validatorGinMode := api.IsGinMode()
	validate := api.NewValidator(validatorGinMode)
	valid := []string{
		gin.ReleaseMode,
		gin.DebugMode,
		gin.TestMode,
	}
	invalid := []string{
		"something else",
		"Release",
		"Debug",
		"Test",
		"",
	}
	validateFields(t, validate, validatorGinMode.Tag, valid, invalid)
}

func TestIsJSONStringArray(t *testing.T) {
	validatorJsonStringArray := api.IsJSONStringArray()
	validate := api.NewValidator(validatorJsonStringArray)
	valid := []string{
		`[]`,
		`["a"]`,
		`["a", "b"]`,
		`["a", "b", "c"]`,
	}
	invalid := []string{
		`[`,
		`[""]`,
		`["a", "", "c"]`,
		"",
	}
	validateFields(t, validate, validatorJsonStringArray.Tag, valid, invalid)
}

func TestIsKeyValid(t *testing.T) {
	validatorKeyValid := api.IsKeyValid()
	validate := api.NewValidator(validatorKeyValid)
	valid := []string{
		"",
		"0cc175b9c0f1b6a831c399e2697726610cc175b9c0f1b6a831c399e2697726610cc175b9c0f1b6a831c399e2697726610cc175b9c0f1b6a831c399e269772661",  // 128
		"0cc175b9c0f1b6a831c399e2697726610cc175b9c0f1b6a831c399e2697726610cc175b9c0f1b6a831c399e2697726610cc175b9c0f1b6a831c399e2697726612", //129
	}
	invalid := []string{
		"1",
		"0cc175b9c0f1b6a831c399e2697726610cc175b9c0f1b6a831c399e2697726610cc175b9c0f1b6a831c399e2697726610cc175b9c0f1b6a831c399e26977266", // 127
	}
	validateFields(t, validate, validatorKeyValid.Tag, valid, invalid)
}

func TestIsPasswordDefault(t *testing.T) {
	validatorPasswordDefault := api.IsPasswordDefault()
	validate := api.NewValidator(validatorPasswordDefault)
	valid := []string{
		"123456",
		"1234567",
	}
	invalid := []string{
		"",
		"1",
	}
	validateFields(t, validate, validatorPasswordDefault.Tag, valid, invalid)
}
