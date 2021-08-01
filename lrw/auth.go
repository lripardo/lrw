package lrw

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/lripardo/lrw"
	"time"
)

const UserInfoKey = "userInfoConfig"

type UserClaims struct {
	ID uint64 `json:"id"`
	IP string `json:"ip"`
	jwt.StandardClaims
}

type UserInfo struct {
	ID    uint64 `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

type UserConfig struct {
	User    UserInfo `json:"user"`
	Expires int64    `json:"expires"`
	ClaimIP string   `json:"claim_ip"`
}

func GetUserConfig(context lrw.Context) *UserConfig {
	config, exists := context.Get(UserInfoKey)
	if !exists {
		return nil
	}
	appStartupConfig, ok := config.(UserConfig)
	if !ok {
		return nil
	}
	return &appStartupConfig
}

func GetUserInfo(context lrw.Context) *UserInfo {
	appStartupConfig := GetUserConfig(context)
	if appStartupConfig == nil {
		return nil
	}
	return &appStartupConfig.User
}

func SetUserConfig(context lrw.Context, user *lrw.User, userClaims *UserClaims) {
	context.Set(UserInfoKey, UserConfig{
		User: UserInfo{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Role:  user.Role,
		},
		Expires: userClaims.ExpiresAt,
		ClaimIP: userClaims.IP,
	})
}

func GetTokenStringFromCookieOrCustomHeader(context lrw.Context, config lrw.Configuration) string {
	if tokenCookie, err := context.Cookie(config.String(lrw.AuthCookie)); err == nil && tokenCookie != "" {
		return tokenCookie
	}
	if tokenHeader := context.GetHeader(config.String(lrw.AuthCustomAuthenticationHeader)); tokenHeader != "" {
		return tokenHeader
	}
	return ""
}

func GetUserClaims(context lrw.Context, publicKey *rsa.PublicKey, config lrw.Configuration) (*UserClaims, *lrw.Response) {
	tokenString := GetTokenStringFromCookieOrCustomHeader(context, config)
	if tokenString == "" {
		return nil, lrw.ResponseUnauthorized()
	}
	tokenObject, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(*jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil || tokenObject == nil {
		return nil, lrw.ResponseUnauthorized()
	}
	if !tokenObject.Valid {
		if _, ok := err.(*jwt.ValidationError); ok {
			return nil, lrw.ResponseUnauthorized()
		} else {
			return nil, lrw.ResponseInternalError(err)
		}
	}
	userClaims, ok := tokenObject.Claims.(*UserClaims)
	if !ok {
		return nil, lrw.ResponseInternalError(errors.New("invalid custom claims conversion"))
	}
	if len(userClaims.IP) == 0 || userClaims.ID == 0 {
		return nil, lrw.ResponseUnauthorized()
	}
	if config.Bool(lrw.AuthVerifyTokenIp) {
		if userClaims.IP != context.ClientIP() {
			return nil, lrw.ResponseUnauthorized()
		}
	}
	return userClaims, nil
}

func GetUserAndValidateClaims(userClaims *UserClaims, db lrw.DB) (*lrw.User, *lrw.Response) {
	user, err := db.FindUser(userClaims.ID)
	if err != nil {
		return nil, lrw.ResponseInternalError(err)
	}
	if user == nil {
		return nil, lrw.ResponseUnauthorized()
	}
	if time.Unix(userClaims.IssuedAt, 0).Before(user.TokenTimestamp) {
		return nil, lrw.ResponseUnauthorized()
	}
	return user, nil
}

func NewAuthKey(config lrw.Configuration) (*rsa.PrivateKey, error) {
	key := config.String(lrw.AuthJWTKeyConfig)
	if key == "" {
		privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			return nil, err
		}
		return privateKey, nil
	}
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	pkcs1PrivateKey, err := x509.ParsePKCS1PrivateKey(keyBytes)
	if err != nil {
		return nil, err
	}
	return pkcs1PrivateKey, nil
}
