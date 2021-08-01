package lrw

import (
	"time"
)

const ImplementationDBMemory = "DB_MEMORY_IMPLEMENTATION"

type MemoryDB struct {
	Users map[uint64]map[string]string
}

func getUserFromUserMap(id uint64, userMap map[string]string) (*User, error) {
	t, err := time.Parse(time.RFC3339, userMap["tokenTimeStamp"])
	if err != nil {
		return nil, err
	}
	user := User{
		ID:             id,
		Email:          userMap["email"],
		Password:       userMap["password"],
		Name:           userMap["name"],
		Role:           userMap["role"],
		TokenTimestamp: t,
	}
	return &user, nil
}

func (db *MemoryDB) AutoMigrate(...interface{}) error {
	return nil
}

func (db *MemoryDB) ImplementationName() string {
	return ImplementationDBMemory
}

func (db *MemoryDB) Start() error {
	return nil
}

func (db *MemoryDB) MigrateAuthentication() error {
	return nil
}

func (db *MemoryDB) FindUser(id uint64) (*User, error) {
	userMap := db.Users[id]
	if userMap != nil {
		return getUserFromUserMap(id, userMap)
	}
	return nil, nil
}

func (db *MemoryDB) FindUserByEmail(email string) (*User, error) {
	for id, userMap := range db.Users {
		if userMap["email"] == email {
			return getUserFromUserMap(id, userMap)
		}
	}
	return nil, nil
}

func (db *MemoryDB) CreateUser(user *User) error {
	if db.Users == nil {
		db.Users = make(map[uint64]map[string]string)
	}

	db.Users[user.ID] = map[string]string{
		"name":           user.Name,
		"email":          user.Email,
		"role":           user.Role,
		"password":       user.Password,
		"tokenTimeStamp": user.TokenTimestamp.Format(time.RFC3339),
	}

	return nil
}
