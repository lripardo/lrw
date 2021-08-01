package gorm

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

const ImplementationDBGormMySQL = "DB_GORM_MYSQL_IMPLEMENTATION"

type Model struct {
	ID        uint64         `gorm:"type:serial;primary_key"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type User struct {
	Model
	Email          string    `gorm:"type:varchar(320);not null"`
	Password       string    `gorm:"type:char(60);not null"`
	Name           string    `gorm:"type:varchar(255);not null"`
	Role           string    `gorm:"type:varchar(11);not null"`
	TokenTimestamp time.Time `gorm:"type:datetime;not null"`
}

func (User) TableName() string {
	return "users"
}

type MySQLDBImplementation struct {
}

func (i *MySQLDBImplementation) DBDialect(url, user, password, name string) gorm.Dialector {
	return mysql.Open(fmt.Sprintf(url, user, password, name))
}

func (i *MySQLDBImplementation) UserModel() interface{} {
	return &User{}
}

func (i *MySQLDBImplementation) Name() string {
	return ImplementationDBGormMySQL
}

func NewMySQLDBImplementation() DBImplementation {
	return &MySQLDBImplementation{}
}
