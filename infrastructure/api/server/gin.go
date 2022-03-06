package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lripardo/lrw/domain/api"
	"net/http"
)

var (
	GinTrustedProxies = api.NewKey("GIN_TRUSTED_PROXIES", api.JSONStringArray, `["127.0.0.1"]`)
	GinEnableLog      = api.NewKey("GIN_ENABLE_LOG", api.Boolean, "true")
)

type Gin struct {
	engine *gin.Engine
}

func doResponse(ginContext *gin.Context, response *api.Response) {
	if response.Redirect != "" {
		ginContext.Redirect(response.Status, response.Redirect)
		ginContext.Abort()
		return
	}

	if response.Error != nil {
		if gin.IsDebugging() {
			response.Message = response.Error.Error()
		}
		api.E(ginContext, "Error", response.Error)
	}

	if response.Message == "" {
		response.Message = http.StatusText(response.Status)
	}

	ginContext.AbortWithStatusJSON(response.Status, response)
}

func ginHandlersWrap(handlers ...api.Handler) func(*gin.Context) {
	return func(ginContext *gin.Context) {
		for _, h := range handlers {
			response := h(ginContext)

			if response != nil {
				doResponse(ginContext, response)
				return
			}
		}
	}
}

func (g *Gin) Start() {
	if err := g.engine.Run(); err != nil {
		api.Fatal(err)
	}
}

func (g *Gin) route(route api.Route) {
	for _, m := range route.Methods {
		g.engine.Handle(m, route.Path, ginHandlersWrap(route.Handlers...))
	}
}

func (g *Gin) filter(filter api.Route) {
	g.engine.Use(ginHandlersWrap(filter.Handlers...))
}

func (g *Gin) RegisterMiddlewares(middlewares ...api.App) {
	for _, m := range middlewares {
		for _, r := range m.Routes() {
			g.filter(r)
		}
	}
}

func (g *Gin) RegisterApps(apps ...api.App) {
	for _, app := range apps {
		for _, r := range app.Routes() {
			g.route(r)
		}
	}
}

func NewGinServer(configuration api.Configuration) api.Server {
	corsConfig := NewCorsConfig(configuration)
	engine := gin.New()
	trustedProxies := configuration.Strings(GinTrustedProxies)
	if err := engine.SetTrustedProxies(trustedProxies); err != nil {
		api.Fatal(err)
	}
	if ginEnableLog := configuration.Bool(GinEnableLog); ginEnableLog {
		engine.Use(gin.Logger())
	}
	engine.Use(gin.Recovery(), cors.New(corsConfig))
	api.D("getting gin server implementation")
	return &Gin{engine: engine}
}
