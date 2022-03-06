package auth

type UserRepository interface {
	UserExists(email string) (bool, error)
	CreateUser(user *User) error
	MakeVerified(email string) error
	UpdateTimestamp(email, field string) error
	UpdatePassword(user *User) error
	ReadUser(email string, fields ...string) (*User, error)
}
