package lrw

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	databaseDialect             = "mysql"
	databaseUrl                 = "DATABASE_URL"
	databaseName                = "DATABASE_NAME"
	databaseUser                = "DATABASE_USER"
	databasePassword            = "DATABASE_PASSWORD"
	databaseAttemptReconnectTry = "DATABASE_ATTEMPT_RECONNECT_TRY"
	databaseDelayReconnectTry   = "DATABASE_DELAY_RECONNECT_TRY"
	databaseMaxConnections      = "DATABASE_MAX_CONNECTIONS"
	databaseMaxIdleConnections  = "DATABASE_MAX_IDLE_CONNECTIONS"
	databaseMakeMigration       = "DATABASE_MAKE_MIGRATION"
	databaseDebugMode           = "DATABASE_DEBUG_MODE"
)

var DB *gorm.DB

func environmentVarNotSetMessage(environmentVariable string) string {
	return fmt.Sprintf("%s is not set or has a invalid format", environmentVariable)
}

func environmentVarBetweenMessage(environmentVariable string, min int, max int) string {
	return fmt.Sprintf("%s must be between %d and %d", environmentVariable, min, max)
}

func startDatabase(params *StartServiceParameters) {
	url := os.Getenv(databaseUrl)
	if len(url) == 0 {
		log.Fatal(environmentVarNotSetMessage(databaseUrl))
	}
	name := os.Getenv(databaseName)
	if len(name) == 0 {
		log.Fatal(environmentVarNotSetMessage(databaseName))
	}
	user := os.Getenv(databaseUser)
	if len(user) == 0 {
		log.Fatal(environmentVarNotSetMessage(databaseUser))
	}
	password := os.Getenv(databasePassword)
	attemptReconnectTryString := os.Getenv(databaseAttemptReconnectTry)
	attemptReconnectTry, err := strconv.Atoi(attemptReconnectTryString)
	if err != nil {
		log.Fatal(environmentVarNotSetMessage(databaseAttemptReconnectTry), err)
	}
	if attemptReconnectTry < minAttemptReconnectTry || attemptReconnectTry > maxAttemptReconnectTry {
		log.Fatal(environmentVarBetweenMessage(databaseAttemptReconnectTry, minAttemptReconnectTry, maxAttemptReconnectTry))
	}
	delayReconnectTryString := os.Getenv(databaseDelayReconnectTry)
	delayReconnectTry, err := strconv.Atoi(delayReconnectTryString)
	if err != nil {
		log.Fatal(environmentVarNotSetMessage(databaseDelayReconnectTry), err)
	}
	if delayReconnectTry < minDelayReconnectTry || delayReconnectTry > maxDelayReconnectTry {
		log.Fatal(environmentVarBetweenMessage(databaseDelayReconnectTry, minDelayReconnectTry, maxDelayReconnectTry))
	}
	maxConnectionsString := os.Getenv(databaseMaxConnections)
	maxConnections, err := strconv.Atoi(maxConnectionsString)
	if err != nil {
		log.Fatal(environmentVarNotSetMessage(databaseMaxConnections), err)
	}
	if maxConnections < minMaxConnections || maxConnections > maxMaxConnections {
		log.Fatal(environmentVarBetweenMessage(databaseMaxConnections, minMaxConnections, maxMaxConnections))
	}
	maxIdleConnectionsString := os.Getenv(databaseMaxIdleConnections)
	maxIdleConnections, err := strconv.Atoi(maxIdleConnectionsString)
	if err != nil {
		log.Fatal(environmentVarNotSetMessage(databaseMaxIdleConnections), err)
	}
	if maxIdleConnections < minMaxIdleConnections || maxIdleConnections > maxMaxIdleConnections {
		log.Fatal(environmentVarBetweenMessage(databaseMaxIdleConnections, minMaxIdleConnections, maxMaxIdleConnections))
	}
	makeMigrationString := os.Getenv(databaseMakeMigration)
	makeMigration, err := strconv.ParseBool(makeMigrationString)
	if err != nil {
		log.Fatal(environmentVarNotSetMessage(databaseMakeMigration), err)
	}
	for i := 0; i <= attemptReconnectTry; i++ {
		if i > 0 {
			log.Println(fmt.Sprintf("attempt %d from %d: trying reconnect to database in %d seconds...", i, attemptReconnectTry, delayReconnectTry), err)
			time.Sleep(time.Duration(delayReconnectTry) * time.Second)
		}
		DB, err = gorm.Open(databaseDialect, fmt.Sprintf(url, user, password, name))
		if err == nil {
			if i > 0 {
				log.Println("database connection successful")
			}
			break
		} else {
			DB = nil
		}
	}
	if DB == nil {
		log.Fatal("database connection fail")
	}
	debugModeString := os.Getenv(databaseDebugMode)
	debugMode, err := strconv.ParseBool(debugModeString)
	if err != nil {
		log.Fatal(environmentVarNotSetMessage(databaseDebugMode), err)
	}
	DB.LogMode(debugMode)
	DB.DB().SetMaxOpenConns(maxConnections)
	DB.DB().SetMaxIdleConns(maxIdleConnections)
	if makeMigration {
		if params.StartDefaultModels {
			DB.AutoMigrate(getModelsMigrations()...)
			setForeignKeys(DB)
		}
		if params.ModelsMigration != nil {
			DB.AutoMigrate(params.ModelsMigration...)
		}
		if params.SetForeignKeys != nil {
			params.SetForeignKeys(DB)
		}
	}
}
