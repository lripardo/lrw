package lrw

import "encoding/json"

type ContextMock struct {
	Data    map[string]interface{}
	Cookies map[string]string
	Headers map[string]string
	IP      string
	Params  map[string]string
	JSON    string
}

func (c *ContextMock) Get(key string) (interface{}, bool) {
	value, ok := c.Data[key]
	if !ok {
		return nil, false
	}
	return value, true
}

func (c *ContextMock) Set(key string, value interface{}) {
	if c.Data == nil {
		c.Data = make(map[string]interface{})
	}
	c.Data[key] = value
}

func (c *ContextMock) Cookie(name string) (string, error) {
	return c.Cookies[name], nil
}

func (c *ContextMock) GetHeader(key string) string {
	return c.Headers[key]
}

func (c *ContextMock) ClientIP() string {
	return c.IP
}

func (c *ContextMock) Param(key string) string {
	return c.Params[key]
}

func (c *ContextMock) ShouldBindJSON(value interface{}) error {
	if err := json.Unmarshal([]byte(c.JSON), value); err != nil {
		return err
	}
	return nil
}

func (c *ContextMock) SetCookie(string, string, int, string, string, bool, bool) {
	//DO NOTHING
}
