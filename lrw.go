package lrw

import (
	"fmt"
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

var (
	VERSION = "Development"
)

type StartServiceParameters struct {
	StartConfig     func()
	ExtraConfigs    []MapConfig
	ModelsMigration []interface{}
	SetForeignKeys  func(*gorm.DB)
	Routes          func(*gin.RouterGroup)
	Network         string
}

func DefaultStartServiceParams() *StartServiceParameters {
	return &StartServiceParameters{
		StartConfig:     nil,
		ExtraConfigs:    nil,
		ModelsMigration: nil,
		SetForeignKeys:  nil,
		Routes:          nil,
		Network:         "tcp4",
	}
}

func StartService(params *StartServiceParameters) {
	log.Println(fmt.Sprintf("Started with version: %s", VERSION))
	a := os.Getenv(address)
	if len(a) == 0 {
		log.Fatal(environmentVarNotSetMessage(address))
	}
	startDatabase(params)
	startConfig(params)
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
