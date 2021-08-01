package lrw

type Context interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	Cookie(name string) (string, error)
	GetHeader(key string) string
	ClientIP() string
	Param(key string) string
	ShouldBindJSON(obj interface{}) error
	SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool)
}
