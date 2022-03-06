package cache

import (
	"encoding/json"
	"github.com/lripardo/lrw/domain/api"
	"sync"
	"time"
)

type Item struct {
	Value   []byte
	Expires time.Time
}

type MemoryCache struct {
	cache sync.Map
}

func (m *MemoryCache) Set(key string, v interface{}, expires time.Duration) error {
	d, err := json.Marshal(v)
	if err != nil {
		return err
	}
	item := &Item{Value: d}
	if expires != 0 {
		item.Expires = time.Now().Add(expires)
	}
	m.cache.Store(key, item)
	return nil
}

func (m *MemoryCache) Del(key string) error {
	m.cache.Delete(key)
	return nil
}

func (m *MemoryCache) Get(key string, v interface{}) error {
	if i, ok := m.cache.Load(key); ok {
		if item, ok := i.(*Item); ok {
			if !item.Expires.IsZero() {
				if item.Expires.Before(time.Now()) {
					return m.Del(key)
				}
			}
			if err := json.Unmarshal(item.Value, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func NewMemoryCache() api.Cache {
	api.D("getting new memory cache")
	return &MemoryCache{}
}
