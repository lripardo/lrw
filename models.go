package lrw

import (
	"time"
)

type Model struct {
	ID        uint64     `gorm:"type:serial;primary_key"`
	CreatedAt time.Time  `gorm:"not null"`
	UpdatedAt time.Time  `gorm:"not null"`
	DeletedAt *time.Time `sql:"index"`
}

type User struct {
	Model
	Email          string `gorm:"type:varchar(320);not null"`
	Password       string `gorm:"type:char(60);not null"`
	Name           string `gorm:"type:varchar(255);not null"`
	Role           string `gorm:"type:varchar(11);not null"`
	TokenTimestamp *time.Time
}

type Config struct {
	Model
	Name string `gorm:"type:varchar(45);unique_index;not null"`
	Data string `gorm:"type:text;not null"`
}

type Log struct {
	Model
	Status        uint    `gorm:"type:smallint;not null"`
	Path          string  `gorm:"type:varchar(255);not null"`
	IP            string  `gorm:"type:varchar(45);not null"`
	Method        string  `gorm:"type:varchar(10);not null"`
	ContentLength int64   `gorm:"type:bigint;not null"`
	Origin        *string `gorm:"type:varchar(255)"`
	User          *uint64 `gorm:"type:bigint unsigned"`
	ErrorCode     *uint   `gorm:"type:int unsigned"`
	ClaimIP       *string `gorm:"type:varchar(45)"`
}

func (User) TableName() string {
	return "users"
}

func (Config) TableName() string {
	return "configs"
}

func (Log) TableName() string {
	return "logs"
}
