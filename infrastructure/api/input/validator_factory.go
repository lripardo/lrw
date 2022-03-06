package input

import "github.com/lripardo/lrw/domain/api"

const (
	HCaptchaInputValidator = "h captcha"
)

// ValidatorFactory need to have configuration instance
type ValidatorFactory struct {
	Configuration api.Configuration
	HCaptcha      api.InputValidator
	Default       api.InputValidator
}

func (v *ValidatorFactory) InputValidator(t string) api.InputValidator {
	if t == HCaptchaInputValidator {
		if v.HCaptcha == nil {
			v.HCaptcha = NewHCaptchaInputValidator(v.Default, v.Configuration)
		}
		return v.HCaptcha
	}
	return v.Default
}

func NewInputValidatorFactory(configuration api.Configuration) api.InputValidationFactory {
	d := NewDefaultInputValidator()
	return &ValidatorFactory{
		Configuration: configuration,
		Default:       d,
	}
}
