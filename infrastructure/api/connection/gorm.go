package connection

import (
	"fmt"
	"github.com/lripardo/lrw/domain/api"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	MySQLGormYpe   = "mysql"
	SQLiteGormType = "sqlite"
)

var (
	GormDriverType             = api.NewKey("GORM_DRIVER_TYPE", "required", SQLiteGormType)
	GormLoggerMode             = api.NewKey("GORM_LOGGER_MODE", "gte=1,lte=4", "4")
	GormLoggerDiscard          = api.NewKey("GORM_LOGGER_DISCARD", api.Boolean, "false")
	GormSkipDefaultTransaction = api.NewKey("GORM_SKIP_DEFAULT_TRANSACTION", api.Boolean, "true")
	GormMigrate                = api.NewKey("GORM_MIGRATE", api.Boolean, "true")

	SQLiteDatabase = api.NewKey("SQLITE_DATABASE", "required", "lrw")

	MySQLUrl      = api.NewKey("MYSQL_URL", "required", "%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local")
	MySQLHost     = api.NewKey("MYSQL_HOST", "required", "localhost")
	MySQLPort     = api.NewKey("MYSQL_PORT", "gte=1,lte=65535", "3306")
	MySQLDatabase = api.NewKey("MYSQL_DATABASE", "required", "lrw")
	MySQLUser     = api.NewKey("MYSQL_USER", "required", "lrw")
	MySQLPassword = api.NewKey("MYSQL_PASSWORD", "required", "lrw")

	MySQLMaxIdleConnections = api.NewKey("MYSQL_MAX_IDLE_CONNECTIONS", "gte=0", "2")
	MySQLMaxOpenConnections = api.NewKey("MYSQL_MAX_OPEN_CONNECTIONS", "gte=0", "2")
)

func newGormConfig(configuration api.Configuration) *gorm.Config {
	discard := configuration.Bool(GormLoggerDiscard)
	loggerMode := configuration.Int(GormLoggerMode)
	skipDefaultTransaction := configuration.Bool(GormSkipDefaultTransaction)

	gormLogger := logger.Default
	if discard {
		gormLogger = logger.Discard
	}

	gormLogger = gormLogger.LogMode(logger.LogLevel(loggerMode))

	return &gorm.Config{
		Logger:                 gormLogger,
		SkipDefaultTransaction: skipDefaultTransaction,
	}
}

func newSQLiteGormDb(configuration api.Configuration, config *gorm.Config) (*gorm.DB, error) {
	database := configuration.String(SQLiteDatabase)
	api.D("getting database from sqlite implementation")
	return gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", database)), config)
}

func newMySQLGormDb(configuration api.Configuration, config *gorm.Config) (*gorm.DB, error) {
	url := configuration.String(MySQLUrl)
	user := configuration.String(MySQLUser)
	password := configuration.String(MySQLPassword)
	host := configuration.String(MySQLHost)
	port := configuration.Uint(MySQLPort)
	database := configuration.String(MySQLDatabase)
	databaseUrl := fmt.Sprintf(url, user, password, host, port, database)

	db, err := gorm.Open(mysql.Open(databaseUrl), config)
	if err != nil {
		return nil, err
	}

	mysqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	maxIdleConnections := configuration.Int(MySQLMaxIdleConnections)
	maxOpenConnections := configuration.Int(MySQLMaxOpenConnections)

	mysqlDB.SetMaxIdleConns(maxIdleConnections)
	mysqlDB.SetMaxOpenConns(maxOpenConnections)

	api.D("getting database from mysql implementation")
	return db, nil
}

func NewGormDB(configuration api.Configuration) (*gorm.DB, error) {
	config := newGormConfig(configuration)

	gormType := configuration.String(GormDriverType)
	if gormType == MySQLGormYpe {
		return newMySQLGormDb(configuration, config)
	}
	return newSQLiteGormDb(configuration, config)
}
