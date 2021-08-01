package gorm

import (
	"errors"
	"github.com/lripardo/lrw"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

type Gorm struct {
	Configuration  lrw.Configuration
	DB             *gorm.DB
	Implementation DBImplementation
}

func (g *Gorm) FindUser(id uint64) (*lrw.User, error) {
	var user lrw.User
	if err := g.DB.Model(g.Implementation.UserModel()).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (g *Gorm) FindUserByEmail(email string) (*lrw.User, error) {
	var user lrw.User
	if err := g.DB.Model(g.Implementation.UserModel()).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (g *Gorm) CreateUser(user *lrw.User) error {
	if err := g.DB.Model(g.Implementation.UserModel()).Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (g *Gorm) ImplementationName() string {
	return g.Implementation.Name()
}

func (g *Gorm) AutoMigrate(models ...interface{}) error {
	if g.Configuration.Bool(lrw.DatabaseMigration) {
		if err := g.DB.AutoMigrate(models...); err != nil {
			return err
		}
	}
	return nil
}

func (g *Gorm) MigrateAuthentication() error {
	if err := g.AutoMigrate(g.Implementation.UserModel()); err != nil {
		return err
	}
	return nil
}

func (g *Gorm) Start() error {
	url := g.Configuration.String(lrw.DatabaseDSN)
	user := g.Configuration.String(lrw.DatabaseUser)
	password := g.Configuration.String(lrw.DatabasePassword)
	name := g.Configuration.String(lrw.DatabaseName)
	attemptReconnectTry := g.Configuration.Int(lrw.DatabaseAttemptReconnectTry)
	delayReconnectTry := g.Configuration.Int(lrw.DatabaseDelayReconnectTry)

	gormDialect := g.Implementation.DBDialect(url, user, password, name)
	gormConfig := &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.LogLevel(g.Configuration.Int(lrw.DatabaseLoggerMode))),
		SkipDefaultTransaction: true,
	}

	db, err := gorm.Open(gormDialect, gormConfig)

	if err != nil {
		count := 0
		for count <= attemptReconnectTry {
			log.Println("trying connection to database...", err)
			time.Sleep(time.Duration(delayReconnectTry) * time.Second)
			db, err = gorm.Open(gormDialect, gormConfig)
			if err == nil {
				log.Println("database connection successful")
				break
			}
			if attemptReconnectTry != 0 {
				count++
			}
		}
	}

	if db == nil {
		return errors.New("database connection fail")
	}

	genericDb, err := db.DB()
	if err != nil {
		return err
	}

	genericDb.SetMaxOpenConns(g.Configuration.Int(lrw.DatabaseMaxConnections))
	genericDb.SetMaxIdleConns(g.Configuration.Int(lrw.DatabaseMaxIdleConnections))

	g.DB = db

	return nil
}

func NewGormDB(configuration lrw.Configuration, implementation DBImplementation) lrw.DB {
	return &Gorm{
		Configuration:  configuration,
		Implementation: implementation,
	}
}
