package api

import "time"

type Cache interface {
	Set(key string, v interface{}, duration time.Duration) error
	Get(key string, v interface{}) error
	Del(key string) error
}
