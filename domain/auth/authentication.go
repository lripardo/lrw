package auth

import (
	"github.com/golang-jwt/jwt"
	"github.com/lripardo/lrw/domain/api"
	"net/http"
	"strings"
	"time"
)

const (
	TokenContentPrefix = "Bearer "
)

var (
	JWTKey                     = api.NewKey("AUTH_JWT_KEY", api.KeyValid, "")
	Cookie                     = api.NewKey("AUTH_COOKIE", "required", "tk")
	CookiePath                 = api.NewKey("AUTH_COOKIE_PATH", "required", "/api/v1")
	CookieDomain               = api.NewKey("AUTH_COOKIE_DOMAIN", "required", "localhost")
	CookieSecure               = api.NewKey("AUTH_COOKIE_SECURE", api.Boolean, "false")
	CookieHttpOnly             = api.NewKey("AUTH_COOKIE_HTTP_ONLY", api.Boolean, "true")
	CustomAuthenticationHeader = api.NewKey("AUTH_CUSTOM_AUTHENTICATION_HEADER", "required", "Authorization")
	Audience                   = api.NewKey("AUTH_AUDIENCE", "required", "auth.audience")
	Issuer                     = api.NewKey("AUTH_ISSUER", "required", "auth.issuer")
	Expires                    = api.NewKey("AUTH_EXPIRES", "gte=1", "525600")
	SameSite                   = api.NewKey("AUTH_SAME_SITE", "gte=1,lte=4", "3")
)

type Authentication struct {
	jwtKey                []byte
	userContextRepository UserContextRepository
	cookie                string
	cookiePath            string
	cookieDomain          string
	cookieSecure          bool
	cookieHttpOnly        bool
	audience              string
	issuer                string
	expires               int
	header                string
	sameSite              http.SameSite
}

func (a *Authentication) DeleteUserContext(email string) error {
	if err := a.userContextRepository.DeleteUserContext(email); err != nil {
		return err
	}
	return nil
}

func (a *Authentication) SignOutUser(context api.Context, user UserContext) error {
	context.SetSameSite(a.sameSite)
	context.SetCookie(a.cookie, "", -1, a.cookiePath, a.cookieDomain, a.cookieSecure, a.cookieHttpOnly)
	return a.DeleteUserContext(user.Email)
}

func (a *Authentication) SignUser(context api.Context, user *User, onCookie bool) (map[string]interface{}, error) {
	claims := CreateClaims(user.Email, a.audience, a.issuer, a.expires)
	token, err := Sign(claims, a.jwtKey)
	if err != nil {
		return nil, err
	}
	if onCookie {
		context.SetSameSite(a.sameSite)
		context.SetCookie(a.cookie, token, a.expires, a.cookiePath, a.cookieDomain, a.cookieSecure, a.cookieHttpOnly)
	}
	data := map[string]interface{}{
		"token":   token,
		"header":  a.header,
		"expires": time.Unix(claims.ExpiresAt, 0),
	}
	return data, nil
}

func (a *Authentication) tokenStringFromCookieOrCustomHeader(context api.Context) string {
	if tokenCookie, err := context.Cookie(a.cookie); err == nil && tokenCookie != "" {
		return tokenCookie
	}
	if tokenHeader := context.GetHeader(a.header); tokenHeader != "" {
		return strings.Replace(tokenHeader, TokenContentPrefix, "", 1)
	}
	return ""
}

func (a *Authentication) userFromRepository(claims *jwt.StandardClaims) (*UserContext, error) {
	expires := time.Unix(claims.ExpiresAt, 0)
	user, err := a.userContextRepository.ReadUserContext(claims.Subject, expires)
	if err != nil {
		return nil, err
	}
	if user == nil {
		api.D("user email from claim was not found on database")
		return nil, nil
	}
	if time.Unix(claims.IssuedAt, 0).Before(user.IgnoreTokenBefore) {
		api.D("the token was issued before the allowed date time")
		return nil, nil
	}
	return user, nil
}

func (a *Authentication) Authenticate(context api.Context) *api.Response {
	tokenString := a.tokenStringFromCookieOrCustomHeader(context)
	claims, err := Claims(tokenString, a.jwtKey)
	if err != nil {
		return api.ResponseInternalError(err)
	}
	if claims == nil {
		return api.ResponseUnauthorized()
	}
	user, err := a.userFromRepository(claims)
	if err != nil {
		return api.ResponseInternalError(err)
	}
	if user == nil {
		return api.ResponseUnauthorized()
	}
	SetUserContext(context, *user)
	return nil
}

func NewAuthenticationService(configuration api.Configuration, repository UserContextRepository) *Authentication {
	key := configuration.String(JWTKey)
	if key == "" {
		api.D("auth key was not found on configuration, a new one will be generated")
		key = api.RandomKey()
	}

	cookie := configuration.String(Cookie)
	cookiePath := configuration.String(CookiePath)
	cookieDomain := configuration.String(CookieDomain)
	cookieSecure := configuration.Bool(CookieSecure)
	cookieHttpOnly := configuration.Bool(CookieHttpOnly)
	audience := configuration.String(Audience)
	issuer := configuration.String(Issuer)
	expires := configuration.Int(Expires)
	header := configuration.String(CustomAuthenticationHeader)
	sameSite := configuration.Int(SameSite)

	return &Authentication{
		cookie:                cookie,
		cookiePath:            cookiePath,
		cookieDomain:          cookieDomain,
		cookieHttpOnly:        cookieHttpOnly,
		cookieSecure:          cookieSecure,
		audience:              audience,
		issuer:                issuer,
		expires:               expires,
		header:                header,
		jwtKey:                []byte(key),
		userContextRepository: repository,
		sameSite:              http.SameSite(sameSite),
	}
}
