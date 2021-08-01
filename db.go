package lrw

type DB interface {
	AutoMigrate(models ...interface{}) error
	MigrateAuthentication() error
	ImplementationName() string
	FindUser(id uint64) (*User, error)
	FindUserByEmail(email string) (*User, error)
	CreateUser(user *User) error
	Start() error
}
