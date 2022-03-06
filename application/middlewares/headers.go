package middlewares

import (
	"fmt"
	"github.com/lripardo/lrw/domain/api"
	"strings"
)

var (
	AppAllowEmptyOrigin = api.NewKey("MIDDLEWARES_APP_ALLOW_EMPTY_ORIGIN", api.Boolean, "true")
	AppServerName       = api.NewKey("MIDDLEWARES_APP_SERVER_NAME", "required", "LRW")
	AppCustomHeaders    = api.NewKey("MIDDLEWARES_APP_CUSTOM_HEADERS", api.HeaderArray, DefaultCustomHeaders)
)

type HeadersApp struct {
	allowEmptyOrigin bool
	serverName       string
	version          string
	headers          []*Header
}

type Header struct {
	Name  string
	Value string
}

func (a *HeadersApp) Version(ctx api.Context) *api.Response {
	for _, h := range a.headers {
		ctx.Header(h.Name, h.Value)
	}
	ctx.Header("Server", fmt.Sprintf("%s/%s", a.serverName, a.version))
	return nil
}

func (a *HeadersApp) Origin(ctx api.Context) *api.Response {
	if !a.allowEmptyOrigin && ctx.GetHeader("Origin") == "" {
		api.D("origin is empty, returning unauthorized...")
		return api.ResponseUnauthorized()
	}
	return nil
}

func (a *HeadersApp) Routes() []api.Route {
	root := api.NewRootRoute("")

	headers := root.Append(api.Route{
		Handlers: []api.Handler{a.Version, a.Origin},
	})

	return []api.Route{headers}
}

func NewHeadersApp(configuration api.Configuration, version string) api.App {
	allowEmptyOrigin := configuration.Bool(AppAllowEmptyOrigin)
	serverName := configuration.String(AppServerName)

	headersStr := configuration.Strings(AppCustomHeaders)
	headers := make([]*Header, 0)
	for _, item := range headersStr {
		keyValue := strings.Split(item, ": ")
		headers = append(headers, &Header{
			Name:  keyValue[0],
			Value: keyValue[1],
		})
	}

	return &HeadersApp{
		allowEmptyOrigin: allowEmptyOrigin,
		version:          version,
		serverName:       serverName,
		headers:          headers,
	}
}
