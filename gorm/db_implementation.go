package gorm

import "gorm.io/gorm"

type DBImplementation interface {
	DBDialect(url, user, password, name string) gorm.Dialector
	UserModel() interface{}
	Name() string
}
