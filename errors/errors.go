package errors

import (
	"errors"
	"net/http"
)

type Error struct {
	error
	Code Code
}

func Wrap(e error, code Code) *Error {
	return &Error{e, code}
}

type Response struct {
	Error string `json:"error,omitempty"`
	Code  Code   `json:"code,omitempty"`
}

func (r *Response) HTTPStatus() int {
	return r.Code.HTTPStatus()
}

func ErrorResponse(err error) *Response {
	var code Code

	switch e := err.(type) {
	case *Error:
		code = e.Code
	case Error:
		code = e.Code
	default:
		code = CodeUnknown
	}

	return &Response{
		Error: err.Error(),
		Code:  code,
	}
}

type Code string

func (c Code) HTTPStatus() int {
	if s, ok := httpStatus[c]; ok {
		return s
	}

	return http.StatusInternalServerError
}

const (
	CodeUnknown          Code = "unknown"
	CodeResourceNotFound Code = "resource_not_found"
	CodeBadRequest       Code = "bad_request"
	CodeQuerySyntax      Code = "query_syntax"
	CodeUnauthorized     Code = "unauthorized"
	CodeForbidden        Code = "forbidden"
	CodeDeviceNotFound   Code = "device_not_found"
)

var httpStatus = map[Code]int{
	CodeUnknown:          http.StatusInternalServerError,
	CodeResourceNotFound: http.StatusNotFound,
	CodeBadRequest:       http.StatusBadRequest,
	CodeForbidden:        http.StatusForbidden,
	CodeUnauthorized:     http.StatusUnauthorized,
	CodeQuerySyntax:      http.StatusBadRequest,
	CodeDeviceNotFound:   http.StatusNotFound,
}

// Some predefined errors

var (
	ErrResourceNotFound = &Error{errors.New("Resource not found"), CodeResourceNotFound}
	ErrDeviceNotFound   = &Error{errors.New("Device not found"), CodeDeviceNotFound}
)
