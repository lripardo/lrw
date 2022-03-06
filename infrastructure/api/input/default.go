package input

import (
	"github.com/go-playground/validator/v10"
	"github.com/lripardo/lrw/domain/api"
)

type Default struct{}

func (d *Default) Read(ctx api.Context, input interface{}, validate *validator.Validate) *api.Response {
	if err := ctx.ShouldBindJSON(input); err != nil {
		return api.ResponseInvalidJSONDataFormat()
	}
	err := validate.Struct(input)
	if err != nil {
		return api.ResponseInvalidInput()
	}
	return nil
}

func NewDefaultInputValidator() api.InputValidator {
	api.D("getting default input validation implementation")
	return &Default{}
}
