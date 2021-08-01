package lrw

import (
	"time"
)

type User struct {
	ID             uint64
	Email          string
	Password       string
	Name           string
	Role           string
	TokenTimestamp time.Time
}
