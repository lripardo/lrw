package auth

import (
	"github.com/lripardo/lrw/domain/api"
	"github.com/lripardo/lrw/domain/auth"
	"github.com/lripardo/lrw/infrastructure/api/cache"
	"github.com/lripardo/lrw/infrastructure/api/connection"
	"gorm.io/gorm"
	"time"
)

const UserContextDB = 1

type UserRepository struct {
	cache api.Cache
	db    *gorm.DB
}

func (s *UserRepository) UserExists(email string) (bool, error) {
	var count int64
	if err := s.db.Model(&UserDTO{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count != 0, nil
}

func (s *UserRepository) CreateUser(user *auth.User) error {
	userDTO := UserDTO{
		Email:              user.Email,
		FirstName:          user.FirstName,
		LastName:           user.LastName,
		Password:           user.Password,
		Role:               user.Role,
		VerifiedOn:         user.VerifiedOn,
		LastChangePassword: user.LastChangePassword,
		IgnoreTokenBefore:  user.IgnoreTokenBefore,
	}
	if err := s.db.Create(&userDTO).Error; err != nil {
		return err
	}
	return nil
}

func (s *UserRepository) CreateUserContext(user *auth.UserContext) error {
	duration := user.Expires.Sub(time.Now())
	if err := s.cache.Set(user.Email, user, duration); err != nil {
		return err
	}
	return nil
}

func (s *UserRepository) DeleteUserContext(email string) error {
	if err := s.cache.Del(email); err != nil {
		return err
	}
	return nil
}

func (s *UserRepository) ReadUserContext(email string, expires time.Time) (*auth.UserContext, error) {
	userContext := auth.UserContext{}
	if err := s.cache.Get(email, &userContext); err != nil {
		return nil, err
	}
	if userContext.Email == "" {
		api.D("user was not found on cache, trying getting from database")
		if err := s.db.Model(&UserDTO{}).
			Select("email", "first_name", "last_name", "role", "ignore_token_before").
			Where("email = ?", email).Scan(&userContext).Error; err != nil {
			return nil, err
		}
		if userContext.Email == "" {
			api.D("user was not found on database either")
			return nil, nil
		}
		userContext.Expires = expires
		if err := s.CreateUserContext(&userContext); err != nil {
			return nil, err
		}
	}
	return &userContext, nil
}

func (s *UserRepository) MakeVerified(email string) error {
	now := time.Now()
	if err := s.db.Model(&UserDTO{}).Where("email = ? and verified_on is NULL", email).Update("verified_on", now).Error; err != nil {
		return err
	}
	return nil
}

func (s *UserRepository) UpdateTimestamp(email, field string) error {
	now := time.Now()
	if err := s.db.Model(&UserDTO{}).Where("email = ?", email).Update(field, now).Error; err != nil {
		return err
	}
	return nil
}

func (s *UserRepository) UpdatePassword(user *auth.User) error {
	updates := map[string]interface{}{
		"password":             user.Password,
		"last_change_password": user.LastChangePassword,
		"ignore_token_before":  user.LastChangePassword,
	}
	if err := s.db.Model(&UserDTO{}).Where("email = ?", user.Email).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}

func (s *UserRepository) ReadUser(email string, fields ...string) (*auth.User, error) {
	user := auth.User{}
	if err := s.db.Model(&UserDTO{}).Select("id", fields).Where("email = ?", email).Scan(&user).Error; err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, nil
	}
	return &user, nil
}

func NewUserRepository(configuration api.Configuration, db *gorm.DB) (*UserRepository, error) {
	if migrate := configuration.Bool(connection.GormMigrate); migrate {
		api.W("migration table is enabled, migrating on every startup can compromise the performance")
		if err := db.AutoMigrate(&UserDTO{}); err != nil {
			return nil, err
		}
	}
	c := cache.NewCache(configuration, UserContextDB)
	return &UserRepository{
		db:    db,
		cache: c,
	}, nil
}
