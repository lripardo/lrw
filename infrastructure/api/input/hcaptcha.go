package input

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/go-resty/resty/v2"
	"github.com/lripardo/lrw/domain/api"
	"net/http"
	"reflect"
)

const (
	HCaptchaAPI           = "https://hcaptcha.com/siteverify"
	HCaptchaDummyHostname = "dummy-key-pass"
)

var (
	HCaptchaSecretKey = api.NewKey("H_CAPTCHA_SECRET_KEY", "required", "0x0000000000000000000000000000000000000000")
	HCaptchaSiteKey   = api.NewKey("H_CAPTCHA_SITE_KEY", "required", "10000000-ffff-ffff-ffff-000000000001")
)

type HCaptchaResponse struct {
	Success     bool     `json:"success"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	Credit      bool     `json:"credit"`
	ErrorCodes  []string `json:"error-codes"`
}

type HCaptcha struct {
	HCaptchaSecretKey string
	HCaptchaSiteKey   string
	Validator         api.InputValidator
	Client            *resty.Client
}

func (d *HCaptcha) validate(token string) *api.Response {
	resp, err := d.Client.R().SetFormData(map[string]string{
		"response": token,
		"secret":   d.HCaptchaSecretKey,
		"sitekey":  d.HCaptchaSiteKey,
	}).Post(HCaptchaAPI)
	if err != nil {
		return api.ResponseInternalError(err)
	}
	if resp.StatusCode() != http.StatusOK {
		api.D(fmt.Sprintf("the h captcha server status code was %d", resp.StatusCode()))
		return api.ResponseInternalError(errors.New("failed to verify the h captcha token"))
	}
	hCaptchaResponse := HCaptchaResponse{}
	if err := json.Unmarshal(resp.Body(), &hCaptchaResponse); err != nil {
		return api.ResponseInternalError(err)
	}
	if !hCaptchaResponse.Success {
		api.D("unsuccessful token validation", hCaptchaResponse)
		return api.ResponseUnauthorized()
	}
	if hCaptchaResponse.Hostname == HCaptchaDummyHostname {
		api.W("using dummy h captcha keys, please do not run this on environment production")
	}
	return nil
}

func (d *HCaptcha) token(input interface{}) string {
	e := reflect.ValueOf(input).Elem()
	for i := 0; i < e.NumField(); i++ {
		fieldName := e.Type().Field(i).Name
		fieldType := e.Type().Field(i).Type
		fieldValue := e.Field(i).Interface()
		if fieldName == "Token" && fieldType.String() == "string" {
			return fieldValue.(string)
		}
	}
	return ""
}

func (d *HCaptcha) Read(ctx api.Context, input interface{}, validate *validator.Validate) *api.Response {
	response := d.Validator.Read(ctx, input, validate)
	if response != nil {
		return response
	}
	token := d.token(input)
	if token == "" {
		api.D("token param is empty, returning unauthorized")
		return api.ResponseUnauthorized()
	}
	return d.validate(token)
}

func NewHCaptchaInputValidator(v api.InputValidator, configuration api.Configuration) api.InputValidator {
	api.D("getting h captcha input validation implementation")
	hCaptchaSecretKey := configuration.String(HCaptchaSecretKey)
	hCaptchaSiteKey := configuration.String(HCaptchaSiteKey)
	restyClient := resty.New()
	return &HCaptcha{
		HCaptchaSiteKey:   hCaptchaSiteKey,
		HCaptchaSecretKey: hCaptchaSecretKey,
		Client:            restyClient,
		Validator:         v,
	}
}
