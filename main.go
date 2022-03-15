package main

import (
	appAuth "github.com/lripardo/lrw/application/auth"
	"github.com/lripardo/lrw/application/middlewares"
	"github.com/lripardo/lrw/domain/api"
	domainAuth "github.com/lripardo/lrw/domain/auth"
	"github.com/lripardo/lrw/infrastructure/api/configuration"
	"github.com/lripardo/lrw/infrastructure/api/connection"
	"github.com/lripardo/lrw/infrastructure/api/email"
	"github.com/lripardo/lrw/infrastructure/api/input"
	"github.com/lripardo/lrw/infrastructure/api/server"
	infraAuth "github.com/lripardo/lrw/infrastructure/auth"
)

var VERSION = "Development"

func main() {
	//configuration instances
	environmentConfiguration := configuration.NewEnvironmentConfiguration()

	//singleton instances
	api.InitLogger(environmentConfiguration)

	//connections instances
	gormDB, err := connection.NewGormDB(environmentConfiguration)
	if err != nil {
		api.Fatal(err)
	}

	//Migrate tables
	if err := gormDB.AutoMigrate(
		&infraAuth.UserDTO{},
		//List of table models
	); err != nil {
		api.Fatal(err)
	}

	//repository instances
	userRepository, err := infraAuth.NewUserRepository(environmentConfiguration, gormDB.DB)
	if err != nil {
		api.Fatal(err)
	}

	//internal services instances
	emailService, err := email.NewEmailService(environmentConfiguration)
	if err != nil {
		api.Fatal(err)
	}
	authenticationService := domainAuth.NewAuthenticationService(environmentConfiguration, userRepository)
	resetPasswordService := domainAuth.NewResetPassword(environmentConfiguration, emailService)
	userVerifyService := domainAuth.NewUserVerify(environmentConfiguration, emailService)

	//factory instances
	inputValidatorFactory := input.NewInputValidatorFactory(environmentConfiguration)

	//apps instances
	authApp := appAuth.NewApp(
		environmentConfiguration,
		authenticationService,
		userRepository,
		userVerifyService,
		resetPasswordService,
		inputValidatorFactory,
	)
	headersApp := middlewares.NewHeadersApp(environmentConfiguration, VERSION)

	//server instances
	ginServer := server.NewGinServer(environmentConfiguration)

	//middlewares registrations
	ginServer.RegisterMiddlewares(headersApp)

	//apps registrations
	ginServer.RegisterApps(authApp)

	//start servers
	go ginServer.Start()

	api.I("version: " + VERSION)

	//block main thread
	select {}
}
