package auth

import "time"

type UserContextRepository interface {
	ReadUserContext(email string, expires time.Time) (*UserContext, error)
	DeleteUserContext(email string) error
	CreateUserContext(user *UserContext) error
}
