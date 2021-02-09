package lrw

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"strconv"
	"strings"
)

type configHelper uint8

var Configs configHelper

var (
	configStorage map[string]string
	jwtKey        *rsa.PrivateKey
)

type mapConfig struct {
	Key       string
	Value     string
	Validator func(string) bool
}

func validateGinMode(ginMode string) bool {
	return ginMode == gin.DebugMode || ginMode == gin.ReleaseMode || ginMode == gin.TestMode
}

func validateCommaArrayString(commaArrayString string) bool {
	if len(commaArrayString) > 0 {
		headers := strings.Split(commaArrayString, ",")
		for _, h := range headers {
			if len(h) == 0 {
				return false
			}
		}
		return len(headers) != 0
	}
	return false
}

func validatePath(path string) bool {
	if len(path) > 1 {
		if path[0] == '/' {
			return len(path[1:]) > 0
		}
	}
	return false
}

func validateTokenTime(tokenTime string) bool {
	time, err := strconv.ParseInt(tokenTime, 10, 64)
	return err == nil && time > 0
}

func validateStringNotEmpty(str string) bool {
	return len(str) != 0
}

func validateNotZeroUint64(uint64Value string) bool {
	parsedNumber, err := strconv.ParseUint(uint64Value, 10, 64)
	return err == nil && parsedNumber > 0
}

func validateBoolean(booleanValue string) bool {
	return booleanValue == "true" || booleanValue == "false"
}

func startConfig() {
	jwtAlias := "jwtKey"
	configMapper := []mapConfig{
		{Key: "cookie", Value: "session", Validator: validateStringNotEmpty},
		{Key: "ginMode", Value: "debug", Validator: validateGinMode},
		{Key: "allowHeaders", Value: "Origin,X-Request-Width,Content-Type,Accept", Validator: validateCommaArrayString},
		{Key: "allowOrigins", Value: "http://localhost:8080", Validator: validateCommaArrayString},
		{Key: "customAuthenticationHeader", Value: "Authorization", Validator: validateStringNotEmpty},
		{Key: "path", Value: "/api/v1", Validator: validatePath},
		{Key: "domain", Value: "localhost", Validator: validateStringNotEmpty},
		{Key: "tokenTime", Value: "31540000000", Validator: validateTokenTime},
		{Key: "tokenAudience", Value: "RP_WEB_LIB", Validator: validateStringNotEmpty},
		{Key: "tokenIssuer", Value: "RP_WEB_LIB", Validator: validateStringNotEmpty},
		{Key: "verifyTokenIp", Value: "true", Validator: validateBoolean},
		{Key: "bruteForceCountAttemptsByIp", Value: "3", Validator: validateNotZeroUint64},
		{Key: "bruteForceTimeMinutesAttemptsByIp", Value: "5", Validator: validateNotZeroUint64},
		{Key: "logOkStatus", Value: "false", Validator: validateBoolean},
		{Key: "printDeniedRequestDump", Value: "true", Validator: validateBoolean},
		{Key: jwtAlias, Value: "", Validator: nil},
	}
	for _, c := range configMapper {
		configDb := Config{}
		if err := DB.Where("name = ?", c.Key).First(&configDb).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				if c.Key == jwtAlias {
					jwtKey, err = rsa.GenerateKey(rand.Reader, 4096)
					if err != nil {
						log.Fatal("cannot generate RSA key", err)
					}
					c.Value = base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(jwtKey))
				}
				configNewDb := Config{Name: c.Key, Data: c.Value}
				if err := DB.Create(&configNewDb).Error; err != nil {
					log.Fatal(fmt.Sprintf("error when create default config %s", c.Key), err)
				}
				configDb = configNewDb
			} else {
				log.Fatal("error in query config on database", err)
			}
		}
		if c.Validator != nil {
			if !c.Validator(configDb.Data) {
				log.Fatal(fmt.Sprintf("%s has a invalid config", c.Key))
			}
		}
		if configStorage == nil {
			configStorage = make(map[string]string)
		}
		configStorage[c.Key] = configDb.Data
	}
	keyBytes, err := base64.StdEncoding.DecodeString(configStorage[jwtAlias])
	if err != nil {
		log.Fatal("cannot decode key from base64")
	}
	jwtKey, err = x509.ParsePKCS1PrivateKey(keyBytes)
	if err != nil {
		log.Fatal("cannot parse RSA JWT key", err)
	}
}

func (config *configHelper) GetString(param string) string {
	return configStorage[param]
}

func getSignKey() *rsa.PrivateKey {
	return jwtKey
}

func getVerifyKey() *rsa.PublicKey {
	return &jwtKey.PublicKey
}

func (config *configHelper) GetUint64(param string) uint64 {
	n, _ := strconv.ParseUint(configStorage[param], 10, 64)
	return n
}

func (config *configHelper) GetBool(param string) bool {
	b, _ := strconv.ParseBool(configStorage[param])
	return b
}

func (config *configHelper) GetInt64(param string) int64 {
	n, _ := strconv.ParseInt(configStorage[param], 10, 64)
	return n
}

func (config *configHelper) GetFullConfig() map[string]string {
	return configStorage
}
