package lrw

type ConfigurationKey struct {
	Key        string
	Validation string
	Default    string
}

var (
	AuthVerifyTokenIp              = ConfigurationKey{Key: "AUTH_VERIFY_TOKEN_IP", Validation: IsBooleanTag, Default: "false"}
	AuthJWTKeyConfig               = ConfigurationKey{Key: "AUTH_JWT_KEY_CONFIG", Validation: "required"}
	AuthCookie                     = ConfigurationKey{Key: "AUTH_COOKIE", Validation: "required", Default: "lrw"}
	AuthCustomAuthenticationHeader = ConfigurationKey{Key: "AUTH_CUSTOM_AUTHENTICATION_HEADER", Validation: "required", Default: "Authorization"}
	AuthTokenTime                  = ConfigurationKey{Key: "AUTH_TOKEN_TIME", Validation: "gte=1", Default: "15"}
	AuthTokenAudience              = ConfigurationKey{Key: "AUTH_TOKEN_AUDIENCE", Validation: "required", Default: "LRW"}
	AuthDomain                     = ConfigurationKey{Key: "AUTH_DOMAIN", Validation: "required", Default: "localhost"}
	AuthCookieSecure               = ConfigurationKey{Key: "AUTH_COOKIE_SECURE", Validation: IsBooleanTag, Default: "false"}
	AuthCookieHttpOnly             = ConfigurationKey{Key: "AUTH_COOKIE_HTTP_ONLY", Validation: IsBooleanTag, Default: "true"}
	AuthBCryptCost                 = ConfigurationKey{Key: "AUTH_BCRYPT_COST", Validation: "gte=4,lte=31", Default: "10"}
)

var (
	FiltersFilterEmptyOrigin = ConfigurationKey{Key: "FILTERS_FILTER_EMPTY_ORIGIN", Validation: IsBooleanTag, Default: "false"}
)

var (
	ServiceAddress             = ConfigurationKey{Key: "SERVICE_ADDRESS", Validation: "required", Default: ":8000"}
	ServiceNetwork             = ConfigurationKey{Key: "SERVICE_NETWORK", Validation: "required", Default: "tcp4"}
	ServiceExposeInternalError = ConfigurationKey{Key: "SERVICE_EXPOSE_INTERNAL_ERROR", Validation: IsBooleanTag, Default: "true"}
	ServiceOriginalStatus      = ConfigurationKey{Key: "SERVICE_ORIGINAL_STATUS", Validation: IsBooleanTag, Default: "false"}
	ServiceAuthentication      = ConfigurationKey{Key: "SERVICE_AUTHENTICATION", Validation: IsBooleanTag, Default: "true"}
	ServiceName                = ConfigurationKey{Key: "SERVICE_NAME", Validation: "required", Default: "LRW"}
	ServicePath                = ConfigurationKey{"SERVICE_PATH", "required", "/api/v1"}
)

var (
	CorsAllowOrigins           = ConfigurationKey{Key: "CORS_ALLOW_ORIGINS", Validation: IsJsonStringArrayTag, Default: `[]`}
	CorsAllowMethods           = ConfigurationKey{Key: "CORS_ALLOW_METHODS", Validation: IsJsonStringArrayTag, Default: `["GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"]`}
	CorsAllowHeaders           = ConfigurationKey{Key: "CORS_ALLOW_HEADERS", Validation: IsJsonStringArrayTag, Default: `["Origin", "Content-Length", "Content-Type", "X-Request-Width", "Accept", "Authorization"]`}
	CorsExposeHeaders          = ConfigurationKey{Key: "CORS_EXPOSE_HEADERS", Validation: IsJsonStringArrayTag, Default: `[]`}
	CorsMaxAge                 = ConfigurationKey{Key: "CORS_MAX_AGE", Validation: "gte=0,lte=24", Default: "12"}
	CorsAllowAllOrigins        = ConfigurationKey{Key: "CORS_ALLOW_ALL_ORIGINS", Validation: IsBooleanTag, Default: "true"}
	CorsAllowCredentials       = ConfigurationKey{Key: "CORS_ALLOW_CREDENTIALS", Validation: IsBooleanTag, Default: "true"}
	CorsAllowWildcard          = ConfigurationKey{Key: "CORS_ALLOW_WILDCARD", Validation: IsBooleanTag, Default: "false"}
	CorsAllowBrowserExtensions = ConfigurationKey{Key: "CORS_ALLOW_BROWSER_EXTENSIONS", Validation: IsBooleanTag, Default: "false"}
	CorsAllowWebSockets        = ConfigurationKey{Key: "CORS_ALLOW_WEBSOCKETS", Validation: IsBooleanTag, Default: "false"}
	CorsAllowFiles             = ConfigurationKey{Key: "CORS_ALLOW_FILES", Validation: IsBooleanTag, Default: "false"}
)

var (
	GinMode = ConfigurationKey{Key: "GIN_MODE", Validation: IsGinModeTag, Default: "debug"}
)

var (
	DatabaseDSN                 = ConfigurationKey{Key: "DATABASE_DSN", Validation: "required", Default: "%s:%s@tcp(localhost:3306)/%s?charset=utf8&parseTime=True&loc=Local"}
	DatabaseName                = ConfigurationKey{Key: "DATABASE_NAME", Validation: "required", Default: "lrw"}
	DatabaseUser                = ConfigurationKey{Key: "DATABASE_USER", Validation: "required", Default: "lrw"}
	DatabasePassword            = ConfigurationKey{Key: "DATABASE_PASSWORD", Validation: "required", Default: "lrw"}
	DatabaseAttemptReconnectTry = ConfigurationKey{Key: "DATABASE_ATTEMPT_RECONNECT_TRY", Validation: "gte=0", Default: "0"}
	DatabaseDelayReconnectTry   = ConfigurationKey{Key: "DATABASE_DELAY_RECONNECT_TRY", Validation: "gte=1", Default: "5"}
	DatabaseMaxConnections      = ConfigurationKey{Key: "DATABASE_MAX_CONNECTIONS", Validation: "gte=1", Default: "2"}
	DatabaseMaxIdleConnections  = ConfigurationKey{Key: "DATABASE_MAX_IDLE_CONNECTIONS", Validation: "gte=1", Default: "2"}
	DatabaseLoggerMode          = ConfigurationKey{Key: "DATABASE_LOGGER_MODE", Validation: "gte=1,lte=4", Default: "4"}
	DatabaseMigration           = ConfigurationKey{Key: "DATABASE_MIGRATION", Validation: "is-boolean", Default: "true"}
)
