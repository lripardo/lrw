package lrw

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http/httputil"
)

type Response uint8

const Next Response = 0

type Handler func(*gin.Context) Response

func (h Handler) Gin() func(*gin.Context) {
	return func(ginContext *gin.Context) {
		h(ginContext)
	}
}

var globalFilter Handler = func(ginContext *gin.Context) Response {
	if ginContext.GetHeader("Origin") == "" && !Configs.GetBool("allowEmptyOrigin") {
		return ResponseNotAuthorized(ginContext)
	}
	return Next
}

func read(params *StartServiceParameters) Handler {
	return func(ginContext *gin.Context) Response {
		jsonResponse := GetStartAppConfigFromGinContext(ginContext)
		var version Config
		if err := DB.Where("name = ?", "version").First(&version).Error; err != nil {
			return ResponseInternalError(err, errorGetVersionApp)(ginContext)
		}
		jsonResponse["version"] = version.Data
		if params.AuthReadResponse != nil {
			jr, err := params.AuthReadResponse(jsonResponse)
			if err != nil {
				return ResponseInternalError(err, errorCustomAuthReadResponse)(ginContext)
			}
			jsonResponse = jr
		}
		return ResponseOkWithData(jsonResponse)(ginContext)
	}
}

func ResponseInternalError(err error, code ErrorCode) Handler {
	return func(ginContext *gin.Context) Response {
		log.Println(err)
		c := uint(code)
		return response(500, fmt.Sprintf("Erro interno. Código: %d", code), nil, &c, ginContext)
	}
}

func ResponseOkWithData(data interface{}) Handler {
	return func(ginContext *gin.Context) Response {
		return response(200, "Ok", data, nil, ginContext)
	}
}

func ResponseCustom(status uint, message string) Handler {
	return func(ginContext *gin.Context) Response {
		return response(status, message, nil, nil, ginContext)
	}
}

func ResponseInvalid(message string) Handler {
	return func(ginContext *gin.Context) Response {
		return response(400, message, nil, nil, ginContext)
	}
}

func ResponseOk(ginContext *gin.Context) Response {
	return response(200, "Ok", nil, nil, ginContext)
}

func ResponseInvalidJsonInput(ginContext *gin.Context) Response {
	return response(400, "invalid json input", nil, nil, ginContext)
}

func ResponseNotAuthorized(ginContext *gin.Context) Response {
	return response(401, "", nil, nil, ginContext)
}

func ResponseForbidden(ginContext *gin.Context) Response {
	return response(403, "", nil, nil, ginContext)
}

func ResponseNotFound(ginContext *gin.Context) Response {
	return response(404, "resource not found", nil, nil, ginContext)
}

func response(status uint, message string, data interface{}, code *uint, ginContext *gin.Context) Response {
	if status != 200 || Configs.GetBool("logOkStatus") {
		userInterface, exists := ginContext.Get("user")
		var userId *uint64
		if exists {
			user, ok := userInterface.(User)
			if ok {
				userId = &user.ID
			}
		}
		var claimIp *string
		c := ginContext.GetString("claim_ip")
		if len(c) > 0 {
			claimIp = &c
		}
		var origin *string
		o := ginContext.Request.Header.Get("Origin")
		oLength := len(o)
		if oLength > 0 {
			if oLength > 255 {
				originReduced := o[:255]
				origin = &originReduced
			} else {
				origin = &o
			}
		}
		if status == 403 || status == 401 {
			if Configs.GetBool("printDeniedRequestDump") {
				dumpRequest, err := httputil.DumpRequest(ginContext.Request, true)
				if err == nil {
					log.Println(string(dumpRequest))
				}
			}
		}
		logModel := Log{
			Status:        status,
			Path:          ginContext.Request.RequestURI,
			IP:            ginContext.ClientIP(),
			ContentLength: ginContext.Request.ContentLength,
			Origin:        origin,
			ErrorCode:     code,
			Method:        ginContext.Request.Method,
			ClaimIP:       claimIp,
			User:          userId,
		}
		if err := DB.Create(&logModel).Error; err != nil {
			log.Println(err)
		}
	}
	//Security Headers
	ginContext.Header("Cache-Control", "no-store")
	ginContext.Header("X-Content-Type-Options", "nosniff")
	ginContext.Header("X-Frame-Options", "DENY")
	ginContext.Header("Referrer-Policy", "no-referrer")
	//
	ginContext.AbortWithStatusJSON(200, gin.H{"message": message, "status": status, "data": data})
	return Next
}
