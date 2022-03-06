package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/lripardo/lrw/domain/api"
	"time"
)

var (
	RedisUrl      = api.NewKey("REDIS_URL", "required", "%s:%d")
	RedisHost     = api.NewKey("REDIS_HOST", "required", "localhost")
	RedisPort     = api.NewKey("REDIS_PORT", "gte=1,lte=65535", "6379")
	RedisPassword = api.NewKey("REDIS_PASSWORD", "required", "redis")
)

type RedisCache struct {
	client *redis.Client
}

func (r *RedisCache) Set(key string, v interface{}, expires time.Duration) error {
	if _, err := r.client.Set(context.Background(), key, v, expires).Result(); err != nil {
		return err
	}
	return nil
}

func (r *RedisCache) Del(key string) error {
	if _, err := r.client.Del(context.Background(), key).Result(); err != nil {
		return err
	}
	return nil
}

func (r *RedisCache) Get(key string, v interface{}) error {
	if err := r.client.Get(context.Background(), key).Scan(v); err != nil && err != redis.Nil {
		return err
	}
	return nil
}

func NewRedisCache(configuration api.Configuration, db int) api.Cache {
	url := configuration.String(RedisUrl)
	host := configuration.String(RedisHost)
	port := configuration.Uint(RedisPort)
	password := configuration.String(RedisPassword)
	cacheUrl := fmt.Sprintf(url, host, port)

	client := redis.NewClient(&redis.Options{
		Addr:     cacheUrl,
		Password: password,
		DB:       db,
	})

	api.D("getting redis cache")
	return &RedisCache{
		client: client,
	}
}
