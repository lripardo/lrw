package gin

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lripardo/lrw"
	"log"
	"net/http"
	"time"
)

type Gin struct {
	Routes        gin.IRoutes
	Configuration lrw.Configuration
	ServerParams  *lrw.ServerParams
	Engine        *gin.Engine
}

func transformHandlersInGinHandlers(serverParams *lrw.ServerParams, handlers ...lrw.Handler) []gin.HandlerFunc {
	ginHandlers := make([]gin.HandlerFunc, 0)
	for _, handler := range handlers {
		ginHandlers = append(ginHandlers, func(context *gin.Context) {
			response := handler(context)

			if response != nil {
				if response.Error != nil {
					log.Println(response.Error)

					if serverParams.ExposeInternalError {
						response.Message = response.Error.Error()
					}
				}

				if response.Message == "" {
					response.Message = http.StatusText(response.Status)
				}

				context.Header("Cache-Control", "no-store")
				context.Header("X-Content-Type-Options", "nosniff")
				context.Header("X-Frame-Options", "DENY")
				context.Header("Referrer-Policy", "no-referrer")

				code := http.StatusOK
				if serverParams.OriginalStatus {
					code = response.Status
				}

				context.AbortWithStatusJSON(code, response)
			}
		})
	}
	return ginHandlers
}

func (g *Gin) RegisterHandlers(method, path string, handlers ...lrw.Handler) {
	if len(handlers) == 0 {
		log.Panicf("no handlers associated to %s and method %s", path, method)
	}
	ginHandlers := transformHandlersInGinHandlers(g.ServerParams, handlers...)
	g.Routes.Handle(method, path, ginHandlers...)
}

func (g *Gin) Start(serverParams *lrw.ServerParams, globalFilters ...lrw.Handler) error {
	corsConfig := cors.Config{
		AllowAllOrigins:        g.Configuration.Bool(lrw.CorsAllowAllOrigins),
		AllowOrigins:           g.Configuration.Strings(lrw.CorsAllowOrigins),
		AllowMethods:           g.Configuration.Strings(lrw.CorsAllowMethods),
		AllowHeaders:           g.Configuration.Strings(lrw.CorsAllowHeaders),
		AllowCredentials:       g.Configuration.Bool(lrw.CorsAllowCredentials),
		ExposeHeaders:          g.Configuration.Strings(lrw.CorsExposeHeaders),
		AllowWildcard:          g.Configuration.Bool(lrw.CorsAllowWildcard),
		AllowBrowserExtensions: g.Configuration.Bool(lrw.CorsAllowBrowserExtensions),
		AllowWebSockets:        g.Configuration.Bool(lrw.CorsAllowWebSockets),
		AllowFiles:             g.Configuration.Bool(lrw.CorsAllowFiles),
		MaxAge:                 time.Duration(g.Configuration.Int64(lrw.CorsMaxAge)) * time.Hour,
	}

	gin.SetMode(g.Configuration.String(lrw.GinMode))

	ginEngine := gin.Default()
	ginEngine.Use(cors.New(corsConfig))

	routes := ginEngine.Group(serverParams.Path)
	ginGlobalFilters := transformHandlersInGinHandlers(serverParams, globalFilters...)

	routes.Use(ginGlobalFilters...)

	g.Routes = routes
	g.ServerParams = serverParams
	g.Engine = ginEngine

	return nil
}

func (g *Gin) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	g.Engine.ServeHTTP(responseWriter, request)
}

func NewGinServer(configuration lrw.Configuration) lrw.Server {
	return &Gin{
		Configuration: configuration,
	}
}
