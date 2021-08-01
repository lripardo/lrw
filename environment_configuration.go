package lrw

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"log"
	"os"
	"strconv"
)

type EnvironmentConfiguration struct {
	Validate *validator.Validate
}

func getValidatedOrDefaultField(key ConfigurationKey, configuration *EnvironmentConfiguration) string {
	v := os.Getenv(key.Key)
	if v == "" {
		return key.Default
	}
	errs := configuration.Validate.Var(v, key.Validation)
	if errs != nil {
		return key.Default
	}
	return v
}

func showErrorAndThrowPanic(key ConfigurationKey, value string, typeOfKey string, err error) {
	log.Println(err)
	log.Panicf("cannot parse %s with value %s as %s", key.Key, value, typeOfKey)
}

func (c *EnvironmentConfiguration) Strings(key ConfigurationKey) []string {
	validatedField := getValidatedOrDefaultField(key, c)
	var arrayStrings []string
	if err := json.Unmarshal([]byte(validatedField), &arrayStrings); err != nil {
		showErrorAndThrowPanic(key, validatedField, "[]string", err)
	}
	return arrayStrings
}

func (c *EnvironmentConfiguration) String(key ConfigurationKey) string {
	return getValidatedOrDefaultField(key, c)
}

func (c *EnvironmentConfiguration) Uint(key ConfigurationKey) uint {
	validatedField := getValidatedOrDefaultField(key, c)
	fieldConverted, err := strconv.ParseUint(validatedField, 10, 32)
	if err != nil {
		showErrorAndThrowPanic(key, validatedField, "uint64", err)
	}
	return uint(fieldConverted)
}

func (c *EnvironmentConfiguration) Bool(key ConfigurationKey) bool {
	validatedField := getValidatedOrDefaultField(key, c)
	fieldConverted, err := strconv.ParseBool(validatedField)
	if err != nil {
		showErrorAndThrowPanic(key, validatedField, "bool", err)
	}
	return fieldConverted
}

func (c *EnvironmentConfiguration) Int(key ConfigurationKey) int {
	validatedField := getValidatedOrDefaultField(key, c)
	fieldConverted, err := strconv.ParseInt(validatedField, 10, 32)
	if err != nil {
		showErrorAndThrowPanic(key, validatedField, "int64", err)
	}
	return int(fieldConverted)
}

func (c *EnvironmentConfiguration) Int64(key ConfigurationKey) int64 {
	validatedField := getValidatedOrDefaultField(key, c)
	fieldConverted, err := strconv.ParseInt(validatedField, 10, 64)
	if err != nil {
		showErrorAndThrowPanic(key, validatedField, "int64", err)
	}
	return fieldConverted
}

func NewEnvironmentConfiguration(validate *validator.Validate) Configuration {
	return &EnvironmentConfiguration{Validate: validate}
}
