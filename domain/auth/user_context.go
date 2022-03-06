package auth

import (
	"encoding/json"
	"github.com/lripardo/lrw/domain/api"
	"time"
)

type UserContext struct {
	Email             string    `json:"email"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Role              string    `json:"role"`
	Expires           time.Time `json:"expires"`
	IgnoreTokenBefore time.Time `json:"ignore_token_before"`
}

const UserContextKey = "user"

func GetUserContext(context api.Context) UserContext {
	if user, exists := context.Get(UserContextKey); exists {
		if userContext, ok := user.(UserContext); ok {
			return userContext
		}
	}
	panic("user not found on context")
}

func SetUserContext(context api.Context, userContext UserContext) {
	context.Set(UserContextKey, userContext)
}

func (u *UserContext) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *UserContext) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
