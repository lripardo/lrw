package auth

import (
	"github.com/golang-jwt/jwt"
	"github.com/lripardo/lrw/domain/api"
	"time"
)

var (
	ResetPasswordKey          = api.NewKey("AUTH_RESET_PASSWORD_KEY", api.KeyValid, "")
	ResetPasswordTemplateName = api.NewKey("AUTH_RESET_PASSWORD_TEMPLATE_NAME", "required", "ResetPasswordTemplate")
	ResetPasswordAudience     = api.NewKey("AUTH_RESET_PASSWORD_AUDIENCE", "required", "auth.reset.password.audience")
	ResetPasswordIssuer       = api.NewKey("AUTH_RESET_PASSWORD_ISSUER", "required", "auth.reset.password.issuer")
	ResetPasswordExpires      = api.NewKey("AUTH_RESET_PASSWORD_EXPIRES", "gte=1", "5")
	ResetPasswordRoute        = api.NewKey("AUTH_RESET_PASSWORD_ROUTE", "required", "http://localhost:8080/api/v1/auth/reset-password")
	ResetPasswordFrom         = api.NewKey("AUTH_RESET_PASSWORD_FROM", "email", "reset_password@lripardo.github.com")
)

const (
	// ResetParam is the name of http param for token
	ResetParam = "tk"
)

type ResetPassword struct {
	route    string
	from     string
	template string
	audience string
	issuer   string
	expires  int
	key      []byte
	sender   api.EmailSender
}

func (r *ResetPassword) AllowChangePassword(user *User) bool {
	now := time.Now()
	return user.LastChangePassword.Add(time.Duration(r.expires+1) * time.Minute).Before(now)
}

func (r *ResetPassword) ResetPassword(token string) (*jwt.StandardClaims, error) {
	claims, err := Claims(token, r.key)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func (r *ResetPassword) Start(user *User) {
	go sendEmail(&SendEmailOptions{
		Sender:   r.sender,
		User:     user,
		Key:      r.key,
		Param:    ResetParam,
		Audience: r.audience,
		Issuer:   r.issuer,
		Route:    r.route,
		From:     r.from,
		Template: r.template,
		Expires:  r.expires,
	})
}

func NewResetPassword(configuration api.Configuration, sender api.EmailSender) *ResetPassword {
	key := configuration.String(ResetPasswordKey)
	if key == "" {
		api.D("reset password key was not found on configuration, a new one will be generated")
		key = api.RandomKey()
	}
	audience := configuration.String(ResetPasswordAudience)
	issuer := configuration.String(ResetPasswordIssuer)
	expires := configuration.Int(ResetPasswordExpires)
	template := configuration.String(ResetPasswordTemplateName)
	from := configuration.String(ResetPasswordFrom)
	route := configuration.String(ResetPasswordRoute)
	return &ResetPassword{
		template: template,
		from:     from,
		route:    route,
		expires:  expires,
		audience: audience,
		issuer:   issuer,
		key:      []byte(key),
		sender:   sender,
	}
}
