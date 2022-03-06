package api

import (
	"net/http"
)

func ResponseOk() *Response {
	return &Response{Status: http.StatusOK}
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

func ResponseInvalidInput() *Response {
	return ResponseBadRequest("invalid input data")
}

func ResponseInvalidJSONDataFormat() *Response {
	return ResponseBadRequest("invalid json data format")
}

func ResponseBadRequest(message string) *Response {
	return &Response{Status: http.StatusBadRequest, Message: message}
}

func ResponseConflict() *Response {
	return &Response{Status: http.StatusConflict}
}

func ResponseInternalError(err error) *Response {
	return &Response{Status: http.StatusInternalServerError, Error: err}
}

func ResponseOkWithData(data interface{}) *Response {
	return &Response{Status: http.StatusOK, Data: data}
}
