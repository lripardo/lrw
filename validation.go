package lrw

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
)

const (
	IsJsonStringArrayTag = "is-json-string-array"
	IsBooleanTag         = "is-boolean"
	IsGinModeTag         = "is-gin-mode"
)

func NewValidator() *validator.Validate {
	validate := validator.New()
	if err := validate.RegisterValidation(IsGinModeTag, IsGinMode); err != nil {
		log.Panicln(err)
	}
	if err := validate.RegisterValidation(IsBooleanTag, IsBoolean); err != nil {
		log.Panicln(err)
	}
	if err := validate.RegisterValidation(IsJsonStringArrayTag, IsJSONStringArray); err != nil {
		log.Panicln(err)
	}
	return validate
}

func IsBoolean(fieldLevel validator.FieldLevel) bool {
	data := fieldLevel.Field().String()
	return data == "true" || data == "false"
}

func IsJSONStringArray(fieldLevel validator.FieldLevel) bool {
	var array []string
	if err := json.Unmarshal([]byte(fieldLevel.Field().String()), &array); err != nil {
		return false
	}
	for _, item := range array {
		if item == "" {
			return false
		}
	}
	return true
}

func IsGinMode(fieldLevel validator.FieldLevel) bool {
	ginMode := fieldLevel.Field().String()
	return ginMode == gin.DebugMode || ginMode == gin.ReleaseMode || ginMode == gin.TestMode
}
