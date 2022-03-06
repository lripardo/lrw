package configuration

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/lripardo/lrw/domain/api"
	"os"
	"strconv"
)

type EnvironmentConfiguration struct {
	validate *validator.Validate
}

func getEnvOrDefault(key api.Key) string {
	e := os.Getenv(key.Name())
	if e == "" {
		e = key.Value()
	}
	return e
}

func (c *EnvironmentConfiguration) String(key api.Key) string {
	e := getEnvOrDefault(key)
	if errs := c.validate.Var(e, key.Validation()); errs != nil {
		panic(errs)
	}
	return e
}

func (c *EnvironmentConfiguration) Uint(key api.Key) uint {
	e := getEnvOrDefault(key)
	n64, err := strconv.ParseUint(e, 10, 64)
	if err != nil {
		panic(err)
	}
	n := uint(n64)
	if errs := c.validate.Var(n, key.Validation()); errs != nil {
		panic(errs)
	}
	return n
}

func (c *EnvironmentConfiguration) Bool(key api.Key) bool {
	e := getEnvOrDefault(key)
	if errs := c.validate.Var(e, key.Validation()); errs != nil {
		panic(errs)
	}
	b, err := strconv.ParseBool(e)
	if err != nil {
		panic(err)
	}
	return b
}

func (c *EnvironmentConfiguration) Strings(key api.Key) []string {
	e := getEnvOrDefault(key)
	var array []string
	if err := json.Unmarshal([]byte(e), &array); err != nil {
		panic(err)
	}
	return array
}

func (c *EnvironmentConfiguration) Int64(key api.Key) int64 {
	e := getEnvOrDefault(key)
	n64, err := strconv.ParseInt(e, 10, 64)
	if err != nil {
		panic(err)
	}
	if errs := c.validate.Var(n64, key.Validation()); errs != nil {
		panic(errs)
	}
	return n64
}

func (c *EnvironmentConfiguration) Int(key api.Key) int {
	return int(c.Int64(key))
}

func NewEnvironmentConfiguration() api.Configuration {
	validate := api.NewValidator(
		api.IsBoolean(),
		api.IsGinMode(),
		api.IsJSONStringArray(),
		api.IsKeyValid(),
		api.IsHeaderArray(),
	)
	api.D("getting environment configuration")
	return &EnvironmentConfiguration{
		validate: validate,
	}
}
