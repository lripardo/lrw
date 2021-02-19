package lrw

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

const (
	address = "ADDRESS"
)

type StartServiceParameters struct {
	ModelsMigration []interface{}
	SetForeignKeys  func(*gorm.DB)
	Routes          func(*gin.RouterGroup)
	Network         string
}

func DefaultStartServiceParams() *StartServiceParameters {
	return &StartServiceParameters{
		ModelsMigration: nil,
		SetForeignKeys:  nil,
		Routes:          nil,
		Network:         "tcp4",
	}
}

func StartService(params *StartServiceParameters) {
	a := os.Getenv(address)
	if len(a) == 0 {
		log.Fatal(environmentVarNotSetMessage(address))
	}
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
	rootRouterGroup.Use(globalFilter.Gin())
	rootRouterGroup.GET("", Authenticate.Gin(), read.Gin())
	authRouterGroup := rootRouterGroup.Group("auth")
	authRouterGroup.POST("login", login.Gin())
	authRouterGroup.POST("logout", logout.Gin())
	authRouterGroup.POST("register", register.Gin())
	if params.Routes != nil {
		params.Routes(rootRouterGroup)
	}
	server := &http.Server{Handler: ginEngine}
	networkListener, err := net.Listen(params.Network, a)
	if err != nil {
		log.Fatal(err)
	}
	if err := server.Serve(networkListener); err != nil {
		log.Fatal(err)
	}
}
