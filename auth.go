package lrw

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"time"
)

type userClaims struct {
	ID uint64 `json:"id"`
	IP string `json:"ip"`
	jwt.StandardClaims
}

type inputLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Cookie   bool   `json:"cookie"`
}

type registerUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func isProduction() bool {
	return Configs.GetString("ginMode") == gin.ReleaseMode
}

func getTokenStringFromCookieOrCustomHeader(ginContext *gin.Context) string {
	if tokenString, err := ginContext.Cookie(Configs.GetString("cookie")); err == nil && tokenString != "" {
		return tokenString
	}
	if tokenString := ginContext.GetHeader(Configs.GetString("customAuthenticationHeader")); tokenString != "" {
		return tokenString
	}
	return ""
}

var authenticate Handler = func(ginContext *gin.Context) Response {
	tokenString := getTokenStringFromCookieOrCustomHeader(ginContext)
	if tokenString == "" {
		return ResponseNotAuthorized(ginContext)
	}
	tokenObject, err := jwt.ParseWithClaims(tokenString, &userClaims{}, func(*jwt.Token) (interface{}, error) {
		return getVerifyKey(), nil
	})
	if err != nil || tokenObject == nil {
		return ResponseNotAuthorized(ginContext)
	}
	if !tokenObject.Valid {
		if _, ok := err.(*jwt.ValidationError); ok {
			return ResponseNotAuthorized(ginContext)
		} else {
			return ResponseInternalError(err, ErrorTokenInvalid)(ginContext)
		}
	}
	uc, ok := tokenObject.Claims.(*userClaims)
	if !ok {
		return ResponseInternalError(errors.New("invalid custom claims conversion"), ErrorClaimsInvalid)(ginContext)
	}
	if len(uc.IP) == 0 || uc.ID == 0 {
		return ResponseNotAuthorized(ginContext)
	}
	if Configs.GetBool("verifyTokenIp") {
		if uc.IP != ginContext.ClientIP() {
			return ResponseNotAuthorized(ginContext)
		}
	}
	var userModel User
	if err := DB.First(&userModel, uc.ID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return ResponseNotAuthorized(ginContext)
		} else {
			return ResponseInternalError(err, ErrorAuthUserQuery)(ginContext)
		}
	}
	if userModel.TokenTimestamp != nil {
		if time.Unix(uc.IssuedAt, 0).Before(*userModel.TokenTimestamp) {
			return ResponseNotAuthorized(ginContext)
		}
	}
	setStartAppConfigToGinContext(ginContext, userModel, uc.ExpiresAt, uc.IP)
	return Next
}

func authorizeIpFromBlacklistBruteForce(ginContext *gin.Context) (bool, error) {
	bruteForceCountAttemptsByIp := Configs.GetUint64("bruteForceCountAttemptsByIp")
	bruteForceTimeMinutesAttemptsByIp := Configs.GetUint64("bruteForceTimeMinutesAttemptsByIp")
	lastTimestamp := time.Now().Add(time.Duration(-int(bruteForceTimeMinutesAttemptsByIp)) * time.Minute)
	var attempts uint64
	if err := DB.Model(&Log{}).
		Where("status = ? and created_at >= ? and ip = ?", 401, lastTimestamp, ginContext.ClientIP()).
		Count(&attempts).Error; err != nil {
		return false, err
	}
	return attempts <= bruteForceCountAttemptsByIp, nil
}

func validEmail(email string) bool {
	emailLength := StringLen(email)
	if emailLength == 0 || emailLength > MaxEmailLength {
		return false
	}
	emailPattern := "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	isValidEmail, err := regexp.Match(emailPattern, []byte(email))
	if err != nil {
		return false
	}
	return isValidEmail
}

func getStartAppConfigFromGinContext(ginContext *gin.Context) gin.H {
	userContext := GetUserFromGinContext(ginContext)
	expires := ginContext.GetInt64("expires")
	claimIp := ginContext.GetString("claim_ip")
	userInfo := InfoUser{ID: userContext.ID, Name: userContext.Name, Role: userContext.Role, Email: userContext.Email}
	jsonResponse := gin.H{
		"user":     userInfo,
		"expires":  time.Unix(expires, 0),
		"claim_ip": claimIp,
	}
	return jsonResponse
}

func setStartAppConfigToGinContext(ginContext *gin.Context, user User, expires int64, claimIp string) {
	ginContext.Set("user", user)
	ginContext.Set("expires", expires)
	ginContext.Set("claim_ip", claimIp)
}

var login Handler = func(ginContext *gin.Context) Response {
	authorizeIp, err := authorizeIpFromBlacklistBruteForce(ginContext)
	if err != nil {
		return ResponseInternalError(err, ErrorAuthorizeIpFromBlacklistLogin)(ginContext)
	}
	if !authorizeIp {
		return ResponseCustom(429, "too many tries")(ginContext)
	}
	var il inputLogin
	if err := ginContext.ShouldBindJSON(&il); err != nil {
		return ResponseInvalidJsonInput(ginContext)
	}
	if !validEmail(il.Email) {
		return ResponseInvalid("invalid email")(ginContext)
	}
	//use rune
	passwordStringLength := StringLen(il.Password)
	if passwordStringLength > PasswordUserMaxLength || passwordStringLength < PasswordUserMinLength {
		return ResponseInvalid("invalid password")(ginContext)
	}
	var userModel User
	if err := DB.Where("email = ?", il.Email).First(&userModel).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return ResponseCustom(406, "user not found")(ginContext)
		} else {
			return ResponseInternalError(err, ErrorLoginUserQuery)(ginContext)
		}
	}
	if err := bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(il.Password)); err != nil {
		return ResponseNotAuthorized(ginContext)
	}
	tokenTime := Configs.GetInt64("tokenTime")
	uc := userClaims{
		userModel.ID,
		ginContext.ClientIP(),
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			Audience:  Configs.GetString("tokenAudience"),
			Subject:   userModel.Email,
			ExpiresAt: time.Now().Add(time.Duration(tokenTime) * time.Millisecond).Unix(),
			Issuer:    Configs.GetString("tokenIssuer")}}
	setStartAppConfigToGinContext(ginContext, userModel, uc.ExpiresAt, uc.IP)
	tokenObject := jwt.NewWithClaims(jwt.SigningMethodRS512, uc)
	tokenString, err := tokenObject.SignedString(getSignKey())
	if err != nil {
		return ResponseInternalError(err, ErrorTokenSignedString)(ginContext)
	}
	if il.Cookie {
		ginContext.SetCookie(
			Configs.GetString("cookie"),
			tokenString,
			int(tokenTime/1000),
			Configs.GetString("path"),
			Configs.GetString("domain"),
			isProduction(),
			true)
	}
	jsonResponseAuth := gin.H{
		"token":  tokenString,
		"header": Configs.GetString("customAuthenticationHeader"),
	}
	jsonResponseConfig := getStartAppConfigFromGinContext(ginContext)
	return ResponseOkWithData(gin.H{"config": jsonResponseConfig, "auth": jsonResponseAuth})(ginContext)
}

var logout Handler = func(ginContext *gin.Context) Response {
	ginContext.SetCookie(
		Configs.GetString("cookie"),
		"",
		-1,
		Configs.GetString("path"),
		Configs.GetString("domain"),
		isProduction(),
		true)
	return ResponseOk(ginContext)
}

var register Handler = func(ginContext *gin.Context) Response {
	var ru registerUser
	if err := ginContext.ShouldBindJSON(&ru); err != nil {
		return ResponseInvalidJsonInput(ginContext)
	}
	if !validEmail(ru.Email) {
		return ResponseInvalid("invalid email")(ginContext)
	}
	passwordLength := StringLen(ru.Password)
	if passwordLength < PasswordUserMinLength || passwordLength > PasswordUserMaxLength {
		return ResponseInvalid("invalid password")(ginContext)
	}
	nameLength := StringLen(ru.Name)
	if nameLength == 0 || nameLength > NameUserMaxLength {
		return ResponseInvalid("invalid name")(ginContext)
	}
	userExisting := &User{}
	if err := DB.Where("email = ?", ru.Email).First(userExisting).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			userExisting = nil
		} else {
			return ResponseInternalError(err, ErrorQueryExistentUser)(ginContext)
		}
	}
	if userExisting != nil {
		return ResponseCustom(409, "user already exists")(ginContext)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(ru.Password), bcrypt.DefaultCost)
	if err != nil {
		return ResponseInternalError(err, ErrorHashUserPasswordRegister)(ginContext)
	}
	userNewRegister := User{Email: ru.Email, Password: string(hash), Name: ru.Name, Role: RoleDefault}
	if err := DB.Create(&userNewRegister).Error; err != nil {
		return ResponseInternalError(err, ErrorCreateUserRegister)(ginContext)
	}
	return ResponseOk(ginContext)
}
