package lrw

type ErrorCode uint

const (
	ErrorClaimsInvalid                 ErrorCode = 1
	ErrorTokenInvalid                  ErrorCode = 2
	ErrorAuthUserQuery                 ErrorCode = 3
	ErrorLoginUserQuery                ErrorCode = 4
	ErrorTokenSignedString             ErrorCode = 5
	ErrorQueryExistentUser             ErrorCode = 6
	ErrorHashUserPasswordRegister      ErrorCode = 7
	ErrorCreateUserRegister            ErrorCode = 8
	ErrorAuthorizeIpFromBlacklistLogin ErrorCode = 9
)
