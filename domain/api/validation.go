package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strings"
)

const (
	HeaderArray     = "is-header-array"
	JSONStringArray = "is-json-string-array"
	Boolean         = "is-boolean"
	GinMode         = "is-gin-mode"
	KeyValid        = "key-valid"
	PasswordDefault = "password-default"
)

const (
	MinSecretKeyLength = 128
	MinPasswordLength  = 6
)

type Validator struct {
	Tag           string
	ValidatorFunc validator.Func
}

func IsBoolean() *Validator {
	return &Validator{
		Tag: Boolean,
		ValidatorFunc: func(fieldLevel validator.FieldLevel) bool {
			data := fieldLevel.Field().String()
			return data == "true" || data == "false"
		},
	}
}

func IsJSONStringArray() *Validator {
	return &Validator{
		Tag: JSONStringArray,
		ValidatorFunc: func(fieldLevel validator.FieldLevel) bool {
			var array []string
			if err := json.Unmarshal([]byte(fieldLevel.Field().String()), &array); err != nil {
				return false
			}
			if array == nil {
				return false
			}
			for _, item := range array {
				if item == "" {
					return false
				}
			}
			return true
		},
	}
}

func IsHeaderArray() *Validator {
	return &Validator{
		Tag: HeaderArray,
		ValidatorFunc: func(fieldLevel validator.FieldLevel) bool {
			var array []string
			if err := json.Unmarshal([]byte(fieldLevel.Field().String()), &array); err != nil {
				return false
			}
			if array == nil {
				return false
			}
			for _, item := range array {
				if item == "" {
					return false
				}
				keyValue := strings.Split(item, ": ")
				if len(keyValue) != 2 {
					return false
				}
				if keyValue[0] == "" || keyValue[1] == "" {
					return false
				}
			}
			return true
		},
	}
}

func IsGinMode() *Validator {
	return &Validator{
		Tag: GinMode,
		ValidatorFunc: func(fieldLevel validator.FieldLevel) bool {
			ginMode := fieldLevel.Field().String()
			return ginMode == gin.DebugMode || ginMode == gin.ReleaseMode || ginMode == gin.TestMode
		},
	}
}

func IsKeyValid() *Validator {
	return &Validator{
		Tag: KeyValid,
		ValidatorFunc: func(fieldLevel validator.FieldLevel) bool {
			keyLength := len(fieldLevel.Field().String())
			return keyLength == 0 || keyLength >= MinSecretKeyLength
		},
	}
}

func IsPasswordDefault() *Validator {
	return &Validator{
		Tag: PasswordDefault,
		ValidatorFunc: func(fieldLevel validator.FieldLevel) bool {
			password := fieldLevel.Field().String()
			return len(password) >= MinPasswordLength
		},
	}
}

func NewValidator(validators ...*Validator) *validator.Validate {
	validate := validator.New()
	for _, v := range validators {
		if err := validate.RegisterValidation(v.Tag, v.ValidatorFunc); err != nil {
			Fatal(err)
		}
	}
	return validate
}
