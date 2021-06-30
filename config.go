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

type MapConfig struct {
	Key       string
	Value     string
	Validator func(string) bool
}

func ValidateGinMode(ginMode string) bool {
	return ginMode == gin.DebugMode || ginMode == gin.ReleaseMode || ginMode == gin.TestMode
}

func ValidateCommaArrayString(commaArrayString string) bool {
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

func ValidatePath(path string) bool {
	if len(path) > 0 {
		return path[0] == '/'
	}
	return false
}

func ValidateNotZeroInt64(tokenTime string) bool {
	time, err := strconv.ParseInt(tokenTime, 10, 64)
	return err == nil && time > 0
}

func ValidateStringNotEmpty(str string) bool {
	return len(str) != 0
}

func ValidateNotZeroUint64(uint64Value string) bool {
	parsedNumber, err := strconv.ParseUint(uint64Value, 10, 64)
	return err == nil && parsedNumber > 0
}

func ValidateBoolean(booleanValue string) bool {
	return booleanValue == "true" || booleanValue == "false"
}

func startConfig(params *StartServiceParameters) {
	jwtAlias := "jwtKey"
	configMapper := []MapConfig{
		{Key: "ginMode", Value: "debug", Validator: ValidateGinMode},
		{Key: "allowHeaders", Value: "Origin,X-Request-Width,Content-Type,Accept,Authorization", Validator: ValidateCommaArrayString},
		{Key: "allowOrigins", Value: "http://localhost:8080", Validator: ValidateCommaArrayString},
		{Key: "path", Value: "/api/v1", Validator: ValidatePath},
		{Key: "logOkStatus", Value: "false", Validator: ValidateBoolean},
		{Key: "printDeniedRequestDump", Value: "false", Validator: ValidateBoolean},
		{Key: "allowEmptyOrigin", Value: "true", Validator: ValidateBoolean},
	}
	if params.AuthFramework {
		configMapperExtra := []MapConfig{
			{Key: "cookie", Value: "session", Validator: ValidateStringNotEmpty},
			{Key: "customAuthenticationHeader", Value: "Authorization", Validator: ValidateStringNotEmpty},
			{Key: "domain", Value: "localhost", Validator: ValidateStringNotEmpty},
			{Key: "tokenTime", Value: "31540000000", Validator: ValidateNotZeroInt64},
			{Key: "tokenAudience", Value: "RP_WEB_LIB", Validator: ValidateStringNotEmpty},
			{Key: "tokenIssuer", Value: "RP_WEB_LIB", Validator: ValidateStringNotEmpty},
			{Key: "verifyTokenIp", Value: "false", Validator: ValidateBoolean},
			{Key: "bruteForceCountAttemptsByIp", Value: "3", Validator: ValidateNotZeroUint64},
			{Key: "bruteForceTimeMinutesAttemptsByIp", Value: "5", Validator: ValidateNotZeroUint64},
			{Key: jwtAlias, Value: "", Validator: nil},
		}
		configMapper = append(configMapper, configMapperExtra...)
	}
	if params.ExtraConfigs != nil {
		configMapper = append(configMapper, params.ExtraConfigs...)
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
	if params.AuthFramework {
		keyBytes, err := base64.StdEncoding.DecodeString(configStorage[jwtAlias])
		if err != nil {
			log.Fatal("cannot decode key from base64")
		}
		jwtKey, err = x509.ParsePKCS1PrivateKey(keyBytes)
		if err != nil {
			log.Fatal("cannot parse RSA JWT key", err)
		}
	}
	if params.StartConfig != nil {
		params.StartConfig()
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
