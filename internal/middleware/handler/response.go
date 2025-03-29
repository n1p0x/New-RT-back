package handler

import "net/http"

type Response struct {
	StatusCode int
	Data       interface{}
	Err        error
}

func NewSuccessResponse(statusCode int, data interface{}) *Response {
	return &Response{
		StatusCode: statusCode,
		Data:       data,
	}
}

func NewErrorResponse(statusCode int, detail string) *Response {
	return &Response{
		StatusCode: statusCode,
		Err: &ErrorResponse{
			Detail: detail,
		},
	}
}

func NewUnprocessableErrorResponse(err error) *Response {
	return &Response{
		StatusCode: http.StatusUnprocessableEntity,
		Err:        err,
	}
}

func NewInternalErrorResponse(err error) *Response {
	return &Response{
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}
