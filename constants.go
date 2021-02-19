package lrw

const (
	passwordUserMaxLength  = 100
	passwordUserMinLength  = 6
	nameUserMaxLength      = 255
	maxEmailLength         = 320
	maxAttemptReconnectTry = 100
	minAttemptReconnectTry = 1
	minDelayReconnectTry   = 1
	maxDelayReconnectTry   = 3600
	minMaxConnections      = 1
	maxMaxConnections      = 100
	minMaxIdleConnections  = 1
	maxMaxIdleConnections  = 100
)

const (
	RoleUser    = "user"
	RoleDefault = RoleUser
)
