package lrw

import (
	"errors"
	"github.com/lripardo/lrw"
	"strconv"
)

func GetUint64ParamFromContext(param string, context lrw.Context) (uint64, error) {
	idString := context.Param(param)
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

func Roles(roles ...string) lrw.Handler {
	return func(context lrw.Context) *lrw.Response {
		if user := GetUserInfo(context); user != nil {
			for _, role := range roles {
				if user.Role == role {
					return nil
				}
			}
		}
		return lrw.ResponseForbidden()
	}
}
