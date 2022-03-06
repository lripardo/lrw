package server

import (
	"github.com/gin-contrib/cors"
	"github.com/lripardo/lrw/domain/api"
	"time"
)

var (
	CorsAllowOrigins           = api.NewKey("CORS_ALLOW_ORIGINS", api.JSONStringArray, `[]`)
	CorsAllowMethods           = api.NewKey("CORS_ALLOW_METHODS", api.JSONStringArray, `["GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"]`)
	CorsAllowHeaders           = api.NewKey("CORS_ALLOW_HEADERS", api.JSONStringArray, `["Origin", "Content-Length", "Content-Type", "X-Request-Width", "Accept", "Authorization"]`)
	CorsExposeHeaders          = api.NewKey("CORS_EXPOSE_HEADERS", api.JSONStringArray, `[]`)
	CorsMaxAge                 = api.NewKey("CORS_MAX_AGE", "gte=0,lte=24", "12")
	CorsAllowAllOrigins        = api.NewKey("CORS_ALLOW_ALL_ORIGINS", api.Boolean, "true")
	CorsAllowCredentials       = api.NewKey("CORS_ALLOW_CREDENTIALS", api.Boolean, "true")
	CorsAllowWildcard          = api.NewKey("CORS_ALLOW_WILDCARD", api.Boolean, "false")
	CorsAllowBrowserExtensions = api.NewKey("CORS_ALLOW_BROWSER_EXTENSIONS", api.Boolean, "false")
	CorsAllowWebSockets        = api.NewKey("CORS_ALLOW_WEBSOCKETS", api.Boolean, "false")
	CorsAllowFiles             = api.NewKey("CORS_ALLOW_FILES", api.Boolean, "false")
)

func NewCorsConfig(configuration api.Configuration) cors.Config {
	maxAge := configuration.Int64(CorsMaxAge)
	api.D("getting cors configuration")
	return cors.Config{
		AllowAllOrigins:        configuration.Bool(CorsAllowAllOrigins),
		AllowOrigins:           configuration.Strings(CorsAllowOrigins),
		AllowMethods:           configuration.Strings(CorsAllowMethods),
		AllowHeaders:           configuration.Strings(CorsAllowHeaders),
		AllowCredentials:       configuration.Bool(CorsAllowCredentials),
		ExposeHeaders:          configuration.Strings(CorsExposeHeaders),
		AllowWildcard:          configuration.Bool(CorsAllowWildcard),
		AllowBrowserExtensions: configuration.Bool(CorsAllowBrowserExtensions),
		AllowWebSockets:        configuration.Bool(CorsAllowWebSockets),
		AllowFiles:             configuration.Bool(CorsAllowFiles),
		MaxAge:                 time.Duration(maxAge) * time.Hour,
	}
}
