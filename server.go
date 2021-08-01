package lrw

import "net/http"

type Server interface {
	RegisterHandlers(method string, path string, handlers ...Handler)
	Start(serverParams *ServerParams, globalFilters ...Handler) error
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type ServerParams struct {
	Path                string
	ExposeInternalError bool
	OriginalStatus      bool
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   error       `json:"-"`
}

type Handler func(Context) *Response

func ResponseInternalError(err error) *Response {
	return &Response{Status: http.StatusInternalServerError, Error: err}
}

func ResponseOkWithData(data interface{}) *Response {
	return &Response{Status: http.StatusOK, Data: data}
}

func ResponseCustom(status int, message string) *Response {
	return &Response{Status: status, Message: message}
}

func ResponseOk() *Response {
	return &Response{Status: http.StatusOK}
}

func ResponseInvalidJsonInput() *Response {
	return &Response{Status: http.StatusBadRequest, Message: "invalid json input"}
}

func ResponseUnauthorized() *Response {
	return &Response{Status: http.StatusUnauthorized}
}

func ResponseForbidden() *Response {
	return &Response{Status: http.StatusForbidden}
}

func ResponseNotFound() *Response {
	return &Response{Status: http.StatusNotFound}
}

func ResponseCustomWithData(status int, message string, data interface{}) *Response {
	return &Response{Status: status, Message: message, Data: data}
}

func ResponseBadRequest(message string) *Response {
	return &Response{Status: http.StatusBadRequest, Message: message}
}
