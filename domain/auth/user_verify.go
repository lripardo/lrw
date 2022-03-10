package auth

import (
	"github.com/golang-jwt/jwt"
	"github.com/lripardo/lrw/domain/api"
)

var (
	VerifyKey               = api.NewKey("AUTH_VERIFY_KEY", api.KeyValid, "")
	VerifyExpires           = api.NewKey("AUTH_VERIFY_EXPIRES", "gte=1", "30")
	VerifyRoute             = api.NewKey("AUTH_VERIFY_ROUTE", "required", "http://localhost:8080/api/v1/auth/verify")
	VerifyEmailTemplateName = api.NewKey("AUTH_VERIFY_EMAIL_TEMPLATE_NAME", "required", "VerifyEmailTemplate")
	VerifyAudience          = api.NewKey("AUTH_VERIFY_AUDIENCE", "required", "auth.verify.audience")
	VerifyIssuer            = api.NewKey("AUTH_VERIFY_ISSUER", "required", "auth.verify.issuer")
	VerifyFrom              = api.NewKey("AUTH_VERIFY_FROM", "email", "verify@lripardo.github.com")
)

const (
	VerifyParam = "tk"
)

type UserVerify struct {
	key      []byte
	sender   api.EmailSender
	expires  int
	audience string
	issuer   string
	route    string
	template string
	from     string
}

func (u *UserVerify) Verify(token string) (*jwt.StandardClaims, error) {
	claims, err := Claims(token, u.key)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func (u *UserVerify) Start(user *User) {
	if user.VerifiedOn != nil {
		api.D("user is already verified")
		return
	}
	go sendEmail(&SendEmailOptions{
		Sender:   u.sender,
		User:     user,
		Key:      u.key,
		Param:    VerifyParam,
		Audience: u.audience,
		Issuer:   u.issuer,
		Route:    u.route,
		From:     u.from,
		Template: u.template,
		Expires:  u.expires,
	})
}

func NewUserVerify(configuration api.Configuration, sender api.EmailSender) *UserVerify {
	key := configuration.String(VerifyKey)
	if key == "" {
		api.D("verify key was not found on configuration, a new one will be generated")
		key = api.RandomKey()
	}
	expires := configuration.Int(VerifyExpires)
	audience := configuration.String(VerifyAudience)
	issuer := configuration.String(VerifyIssuer)
	route := configuration.String(VerifyRoute)
	template := configuration.String(VerifyEmailTemplateName)
	from := configuration.String(VerifyFrom)
	return &UserVerify{
		expires:  expires,
		audience: audience,
		issuer:   issuer,
		route:    route,
		template: template,
		from:     from,
		key:      []byte(key),
		sender:   sender,
	}
}
