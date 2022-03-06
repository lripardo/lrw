package auth

import (
	"gorm.io/gorm"
	"time"
)

type UserDTO struct {
	gorm.Model
	Email              string     `gorm:"type:varchar(255);not null;unique"`
	FirstName          string     `gorm:"type:varchar(255);not null"`
	LastName           string     `gorm:"type:varchar(255);not null"`
	Password           string     `gorm:"type:char(60);not null"`
	Role               string     `gorm:"type:varchar(45);not null"`
	LastLogin          *time.Time `gorm:"type:datetime"`
	VerifiedOn         *time.Time `gorm:"type:datetime"`
	LastChangePassword time.Time  `gorm:"type:datetime;not null"`
	IgnoreTokenBefore  time.Time  `gorm:"type:datetime;not null"`
}

func (u *UserDTO) TableName() string {
	return "users"
}
