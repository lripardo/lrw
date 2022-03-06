package auth

import (
	"crypto/sha512"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const DefaultUserRole = "user"

type User struct {
	ID                 uint
	Email              string
	FirstName          string
	LastName           string
	Password           string
	Role               string
	LastLogin          *time.Time
	VerifiedOn         *time.Time
	LastChangePassword time.Time
	IgnoreTokenBefore  time.Time
}

func (u *User) IsPassword(clearPassword string) bool {
	hashSha := sha512.New().Sum([]byte(clearPassword))
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), hashSha); err == nil {
		return true
	}
	return false
}

func (u *User) HashPassword(clearPassword string, cost int) error {
	hashSha := sha512.New().Sum([]byte(clearPassword))
	hash, err := bcrypt.GenerateFromPassword(hashSha, cost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	u.LastChangePassword = time.Now()
	return nil
}

func NewUser(email, firstName, lastName, clearPassword string, cost int) (*User, error) {
	user := &User{
		Email:             email,
		FirstName:         firstName,
		LastName:          lastName,
		Password:          clearPassword,
		Role:              DefaultUserRole,
		IgnoreTokenBefore: time.Now(),
	}
	if err := user.HashPassword(clearPassword, cost); err != nil {
		return nil, err
	}
	return user, nil
}
