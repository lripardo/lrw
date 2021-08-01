package lrw

import (
	"crypto/rsa"
	"github.com/lripardo/lrw"
)

func AuthenticateFilter(db lrw.DB, publicKey *rsa.PublicKey, configuration lrw.Configuration) lrw.Handler {
	return func(context lrw.Context) *lrw.Response {
		userClaims, response := GetUserClaims(context, publicKey, configuration)
		if response != nil {
			return response
		}
		user, response := GetUserAndValidateClaims(userClaims, db)
		if response != nil {
			return response
		}
		SetUserConfig(context, user, userClaims)
		return nil
	}
}

func OriginFilter(configuration lrw.Configuration) lrw.Handler {
	return func(context lrw.Context) *lrw.Response {
		if configuration.Bool(lrw.FiltersFilterEmptyOrigin) && context.GetHeader("Origin") == "" {
			return lrw.ResponseUnauthorized()
		}
		return nil
	}
}
