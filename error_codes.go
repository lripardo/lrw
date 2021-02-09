package lrw

type ErrorCode uint

const (
	errorClaimsInvalid                 ErrorCode = 1
	errorTokenInvalid                  ErrorCode = 2
	errorAuthUserQuery                 ErrorCode = 3
	errorLoginUserQuery                ErrorCode = 4
	errorTokenSignedString             ErrorCode = 5
	errorQueryExistentUser             ErrorCode = 6
	errorHashUserPasswordRegister      ErrorCode = 7
	errorCreateUserRegister            ErrorCode = 8
	errorAuthorizeIpFromBlacklistLogin ErrorCode = 9
)
