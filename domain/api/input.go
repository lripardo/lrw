package api

import "github.com/go-playground/validator/v10"

type InputValidator interface {
	Read(ctx Context, input interface{}, validate *validator.Validate) *Response
}

type InputValidationFactory interface {
	InputValidator(t string) InputValidator
}
