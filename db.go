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
	DatabaseDialect             = "mysql"
	DatabaseUrl                 = "DATABASE_URL"
	DatabaseUser                = "DATABASE_USER"
	DatabasePassword            = "DATABASE_PASSWORD"
	DatabaseAttemptReconnectTry = "DATABASE_ATTEMPT_RECONNECT_TRY"
	DatabaseDelayReconnectTry   = "DATABASE_DELAY_RECONNECT_TRY"
	DatabaseMaxConnections      = "DATABASE_MAX_CONNECTIONS"
	DatabaseMaxIdleConnections  = "DATABASE_MAX_IDLE_CONNECTIONS"
	DatabaseMakeMigration       = "DATABASE_MAKE_MIGRATION"
	DatabaseDebugMode           = "DATABASE_DEBUG_MODE"
)

var DB *gorm.DB

func environmentVarNotSetMessage(environmentVariable string) string {
	return fmt.Sprintf("%s is not set or has a invalid format", environmentVariable)
}

func startDatabase(params *StartServiceParameters) {
	url := os.Getenv(DatabaseUrl)
	if len(url) == 0 {
		log.Fatal(environmentVarNotSetMessage(DatabaseUrl))
	}
	user := os.Getenv(DatabaseUser)
	if len(user) == 0 {
		log.Fatal(environmentVarNotSetMessage(DatabaseUser))
	}
	password := os.Getenv(DatabasePassword)
	attemptReconnectTryString := os.Getenv(DatabaseAttemptReconnectTry)
	attemptReconnectTry, err := strconv.Atoi(attemptReconnectTryString)
	if err != nil {
		log.Fatal(environmentVarNotSetMessage(DatabaseAttemptReconnectTry), err)
	}
	delayReconnectTryString := os.Getenv(DatabaseDelayReconnectTry)
	delayReconnectTry, err := strconv.Atoi(delayReconnectTryString)
	if err != nil {
		log.Fatal(environmentVarNotSetMessage(DatabaseDelayReconnectTry), err)
	}
	maxConnectionsString := os.Getenv(DatabaseMaxConnections)
	maxConnections, err := strconv.Atoi(maxConnectionsString)
	if err != nil {
		log.Fatal(environmentVarNotSetMessage(DatabaseMaxConnections), err)
	}
	maxIdleConnectionsString := os.Getenv(DatabaseMaxIdleConnections)
	maxIdleConnections, err := strconv.Atoi(maxIdleConnectionsString)
	if err != nil {
		log.Fatal(environmentVarNotSetMessage(DatabaseMaxIdleConnections), err)
	}
	makeMigrationString := os.Getenv(DatabaseMakeMigration)
	makeMigration, err := strconv.ParseBool(makeMigrationString)
	if err != nil {
		log.Fatal(environmentVarNotSetMessage(DatabaseMakeMigration), err)
	}
	for i := 0; i <= attemptReconnectTry; i++ {
		if i > 0 {
			log.Println(fmt.Sprintf("attempt %d from %d: trying reconnect to database in %d seconds...", i, attemptReconnectTry, delayReconnectTry), err)
			time.Sleep(time.Duration(delayReconnectTry) * time.Second)
		}
		DB, err = gorm.Open(DatabaseDialect, fmt.Sprintf(url, user, password))
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
	debugModeString := os.Getenv(DatabaseDebugMode)
	debugMode, err := strconv.ParseBool(debugModeString)
	if err != nil {
		log.Fatal(environmentVarNotSetMessage(DatabaseDebugMode), err)
	}
	DB.LogMode(debugMode)
	DB.DB().SetMaxOpenConns(maxConnections)
	DB.DB().SetMaxIdleConns(maxIdleConnections)
	if makeMigration {
		DB.AutoMigrate(getModelsMigrations()...)
		if params.ModelsMigration != nil {
			DB.AutoMigrate(params.ModelsMigration...)
		}
		if params.SetForeignKeys != nil {
			params.SetForeignKeys(DB)
		}
		setForeignKeys(DB)
	}
}
