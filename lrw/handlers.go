package lrw

import (
	"crypto/rsa"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/lripardo/lrw"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type LoginUser struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Cookie   bool   `json:"cookie"`
}

type RegisterUser struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name"`
}

func ReadUserConfig() lrw.Handler {
	return func(context lrw.Context) *lrw.Response {
		userConfig := GetUserConfig(context)
		if userConfig == nil {
			return lrw.ResponseInternalError(errors.New("user config not found"))
		}
		return lrw.ResponseOkWithData(userConfig)
	}
}

func Login(db lrw.DB, privateKey *rsa.PrivateKey, configuration lrw.Configuration) lrw.Handler {
	return func(context lrw.Context) *lrw.Response {
		var loginUser LoginUser
		if err := context.ShouldBindJSON(&loginUser); err != nil {
			return lrw.ResponseInvalidJsonInput()
		}
		user, err := db.FindUserByEmail(loginUser.Email)
		if err != nil {
			return lrw.ResponseInternalError(err)
		}
		if user == nil {
			return lrw.ResponseNotFound()
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password)); err != nil {
			return lrw.ResponseUnauthorized()
		}
		tokenTime := configuration.Int(lrw.AuthTokenTime)
		userClaims := UserClaims{
			ID: user.ID,
			IP: context.ClientIP(),
			StandardClaims: jwt.StandardClaims{
				IssuedAt:  time.Now().Unix(),
				Audience:  configuration.String(lrw.AuthTokenAudience),
				Subject:   user.Email,
				ExpiresAt: time.Now().Add(time.Duration(tokenTime) * time.Minute).Unix(),
				Issuer:    configuration.String(lrw.ServiceName)}}
		userConfig := UserConfig{
			User: UserInfo{
				ID:    user.ID,
				Email: user.Email,
				Name:  user.Name,
				Role:  user.Role,
			},
			Expires: userClaims.ExpiresAt,
			ClaimIP: userClaims.IP,
		}
		tokenObject := jwt.NewWithClaims(jwt.SigningMethodRS512, userClaims)
		tokenString, err := tokenObject.SignedString(privateKey)
		if err != nil {
			return lrw.ResponseInternalError(err)
		}
		if loginUser.Cookie {
			context.SetCookie(
				configuration.String(lrw.AuthCookie),
				tokenString,
				tokenTime,
				configuration.String(lrw.ServicePath),
				configuration.String(lrw.AuthDomain),
				configuration.Bool(lrw.AuthCookieSecure),
				configuration.Bool(lrw.AuthCookieHttpOnly))
		}
		dataAuth := map[string]interface{}{
			"token":  tokenString,
			"header": configuration.String(lrw.AuthCustomAuthenticationHeader),
		}
		data := map[string]interface{}{
			"config": userConfig,
			"auth":   dataAuth,
		}
		return lrw.ResponseOkWithData(data)
	}
}

func Logout(configuration lrw.Configuration) lrw.Handler {
	return func(context lrw.Context) *lrw.Response {
		context.SetCookie(
			configuration.String(lrw.AuthCookie),
			"",
			-1,
			configuration.String(lrw.ServicePath),
			configuration.String(lrw.AuthDomain),
			configuration.Bool(lrw.AuthCookieSecure),
			configuration.Bool(lrw.AuthCookieHttpOnly))
		return lrw.ResponseOk()
	}
}

func Register(db lrw.DB, configuration lrw.Configuration) lrw.Handler {
	return func(context lrw.Context) *lrw.Response {
		var registerUser RegisterUser
		if err := context.ShouldBindJSON(&registerUser); err != nil {
			return lrw.ResponseInvalidJsonInput()
		}
		userExisting, err := db.FindUserByEmail(registerUser.Email)
		if err != nil {
			return lrw.ResponseInternalError(err)
		}
		if userExisting != nil {
			return lrw.ResponseCustom(409, "user already exists")
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(registerUser.Password), configuration.Int(lrw.AuthBCryptCost))
		if err != nil {
			return lrw.ResponseInternalError(err)
		}
		user := lrw.User{
			Email:          registerUser.Email,
			Password:       string(hash),
			Name:           registerUser.Name,
			TokenTimestamp: time.Now(),
		}
		if err := db.CreateUser(&user); err != nil {
			return lrw.ResponseInternalError(err)
		}
		return lrw.ResponseOk()
	}
}
