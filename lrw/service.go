package lrw

import (
	"github.com/lripardo/lrw"
	"github.com/lripardo/lrw/gin"
	"github.com/lripardo/lrw/gorm"
	"net"
	"net/http"
)

type Service struct {
	Configuration      lrw.Configuration
	DB                 lrw.DB
	Server             lrw.Server
	AuthenticateFilter lrw.Handler
}

func (l *Service) Start() error {
	if err := l.DB.Start(); err != nil {
		return err
	}

	enabledAuthentication := l.Configuration.Bool(lrw.ServiceAuthentication)

	if enabledAuthentication {
		if err := l.DB.MigrateAuthentication(); err != nil {
			return err
		}
	}

	serverParams := &lrw.ServerParams{
		Path:                l.Configuration.String(lrw.ServicePath),
		ExposeInternalError: l.Configuration.Bool(lrw.ServiceExposeInternalError),
		OriginalStatus:      l.Configuration.Bool(lrw.ServiceOriginalStatus),
	}

	if err := l.Server.Start(serverParams, OriginFilter(l.Configuration)); err != nil {
		return err
	}

	privateKey, err := NewAuthKey(l.Configuration)
	if err != nil {
		return err
	}

	if enabledAuthentication {
		l.AuthenticateFilter = AuthenticateFilter(l.DB, &privateKey.PublicKey, l.Configuration)

		l.Server.RegisterHandlers(http.MethodGet, "", l.AuthenticateFilter, ReadUserConfig())
		l.Server.RegisterHandlers(http.MethodPost, "/auth/login", Login(l.DB, privateKey, l.Configuration))
		l.Server.RegisterHandlers(http.MethodPost, "/auth/logout", Logout(l.Configuration))
		l.Server.RegisterHandlers(http.MethodPost, "/auth/register", Register(l.DB, l.Configuration))
	}

	network := l.Configuration.String(lrw.ServiceNetwork)
	address := l.Configuration.String(lrw.ServiceAddress)

	server := &http.Server{Handler: l.Server}
	networkListener, err := net.Listen(network, address)
	if err != nil {
		return err
	}
	if err := server.Serve(networkListener); err != nil {
		return err
	}
	return nil
}

func NewEnvironmentGinGormMySQL() *Service {
	validate := lrw.NewValidator()
	configuration := lrw.NewEnvironmentConfiguration(validate)
	implementation := gorm.NewMySQLDBImplementation()
	db := gorm.NewGormDB(configuration, implementation)
	server := gin.NewGinServer(configuration)

	return &Service{
		Configuration: configuration,
		DB:            db,
		Server:        server,
	}
}
