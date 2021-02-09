package lrw

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"net"
	"net/http"
	"strings"
)

type StartServiceParameters struct {
	ModelsMigration []interface{}
	SetForeignKeys  func(*gorm.DB)
	Routes          func(*gin.RouterGroup)
	BindAddress     string
	Network         string
}

func DefaultStartServiceParams() *StartServiceParameters {
	return &StartServiceParameters{
		ModelsMigration: nil,
		SetForeignKeys:  nil,
		Routes:          nil,
		BindAddress:     "8000",
		Network:         "tcp4",
	}
}

func StartService(params *StartServiceParameters) {
	startDatabase(params)
	startConfig()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = strings.Split(Configs.GetString("allowHeaders"), ",")
	corsConfig.AllowOrigins = strings.Split(Configs.GetString("allowOrigins"), ",")
	gin.SetMode(Configs.GetString("ginMode"))
	ginEngine := gin.Default()
	ginEngine.Use(cors.New(corsConfig))
	rootRouterGroup := ginEngine.Group(Configs.GetString("path"))
	rootRouterGroup.GET("", Authenticate.Gin(), read.Gin())
	authRouterGroup := rootRouterGroup.Group("auth")
	authRouterGroup.POST("login", login.Gin())
	authRouterGroup.POST("logout", logout.Gin())
	authRouterGroup.POST("register", register.Gin())
	if params.Routes != nil {
		params.Routes(rootRouterGroup)
	}
	server := &http.Server{Handler: ginEngine}
	networkListener, err := net.Listen(params.Network, params.BindAddress)
	if err != nil {
		log.Fatal(err)
	}
	if err := server.Serve(networkListener); err != nil {
		log.Fatal(err)
	}
}
