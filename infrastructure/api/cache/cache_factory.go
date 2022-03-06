package cache

import "github.com/lripardo/lrw/domain/api"

const (
	MemoryCacheType = "memory"
	RedisCacheType  = "redis"
)

var (
	Type = api.NewKey("CACHE_TYPE", "required", MemoryCacheType)
)

func NewCache(configuration api.Configuration, db int) api.Cache {
	cacheType := configuration.String(Type)

	if cacheType == RedisCacheType {
		return NewRedisCache(configuration, db)
	}

	return NewMemoryCache()
}
