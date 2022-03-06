package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/lripardo/lrw/domain/api"
	"time"
)

func CreateClaims(email, audience, issuer string, expires int) *jwt.StandardClaims {
	return &jwt.StandardClaims{
		Audience:  audience,
		ExpiresAt: time.Now().Add(time.Duration(expires) * time.Minute).Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    issuer,
		Subject:   email,
	}
}

func Claims(tokenString string, key []byte) (*jwt.StandardClaims, error) {
	if tokenString == "" {
		api.D("empty token string, returning nil")
		return nil, nil
	}
	tokenObject, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return key, nil
		})
	if err != nil || tokenObject == nil {
		api.D("error", err)
		return nil, nil
	}
	if !tokenObject.Valid {
		if _, ok := err.(*jwt.ValidationError); ok {
			api.D("validation error", err)
			return nil, nil
		} else {
			return nil, err
		}
	}
	claims, ok := tokenObject.Claims.(*jwt.StandardClaims)
	if !ok {
		return nil, errors.New("claims must be of type jwt standard claims")
	}
	return claims, nil
}

func Sign(claims *jwt.StandardClaims, key []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	ss, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return ss, nil
}
