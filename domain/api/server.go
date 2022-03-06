package api

import "net/http"

type Context interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	Cookie(name string) (string, error)
	GetHeader(key string) string
	Header(key, value string)
	ShouldBindJSON(obj interface{}) error
	Query(key string) string
	SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool)
	SetSameSite(ss http.SameSite)
}

type Handler func(Context) *Response

type Route struct {
	Path     string
	Methods  []string
	Handlers []Handler
}

type App interface {
	Routes() []Route
}

type Server interface {
	Start()
	RegisterMiddlewares(middlewares ...App)
	RegisterApps(apps ...App)
}

type Response struct {
	Status   int         `json:"status"`
	Code     string      `json:"code"`
	Message  string      `json:"message"`
	Data     interface{} `json:"data"`
	Error    error       `json:"-"`
	Redirect string      `json:"-"`
}

func (m Route) Append(route Route) Route {
	if route.Path != "" {
		m.Path += "/" + route.Path
	}
	if len(route.Methods) > 0 {
		m.Methods = append(m.Methods, route.Methods...)
	}
	if len(route.Handlers) > 0 {
		m.Handlers = append(m.Handlers, route.Handlers...)
	}
	return m
}

func NewRootRoute(root string) Route {
	return Route{
		Path:     root,
		Methods:  make([]string, 0),
		Handlers: make([]Handler, 0),
	}
}
