package lrw

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

type InfoUser struct {
	ID                  uint64 `json:"id"`
	Name                string `json:"name"`
	Role                string `json:"role"`
	Email               string `json:"email"`
	HasToChangePassword bool   `json:"has_to_change_password"`
}

func GetUserFromGinContext(ginContext *gin.Context) *User {
	userContextInterface, exists := ginContext.Get("user")
	if !exists {
		log.Fatal("gin context user not exists")
	}
	userContext, ok := userContextInterface.(User)
	if !ok {
		log.Fatal("error parsing user gin context")
	}
	return &userContext
}

func GetUint64ParamFromGinContext(param string, ginContext *gin.Context) (uint64, error) {
	idString := ginContext.Params.ByName(param)
	if idString == "" {
		return 0, errors.New("param id invalid")
	}
	id, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		return 0, err
	}
	if id == 0 {
		return 0, errors.New("id zero not allowed")
	}
	return id, nil
}

func stringLen(s string) int {
	return len([]rune(s))
}
