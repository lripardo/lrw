package configuration

import (
	"github.com/lripardo/lrw/domain/api"
)

//MemoryConfiguration is for mocking configs for tests
type MemoryConfiguration struct {
	cache map[string]interface{}
}

func (c *MemoryConfiguration) AddWithValue(key api.Key, value interface{}) {
	c.cache[key.Name()] = value
}

func (c *MemoryConfiguration) Add(key api.Key) {
	c.AddWithValue(key, key.Value())
}

func (c *MemoryConfiguration) String(key api.Key) string {
	return c.cache[key.Name()].(string)
}

func (c *MemoryConfiguration) Uint(key api.Key) uint {
	return c.cache[key.Name()].(uint)
}

func (c *MemoryConfiguration) Bool(key api.Key) bool {
	return c.cache[key.Name()].(bool)
}

func (c *MemoryConfiguration) Strings(key api.Key) []string {
	return c.cache[key.Name()].([]string)
}

func (c *MemoryConfiguration) Int64(key api.Key) int64 {
	return c.cache[key.Name()].(int64)
}

func (c *MemoryConfiguration) Int(key api.Key) int {
	return int(c.Int64(key))
}

func NewMemoryConfiguration() api.Configuration {
	api.D("getting memory configuration")
	return &MemoryConfiguration{
		cache: make(map[string]interface{}),
	}
}
