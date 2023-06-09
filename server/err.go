package server

import "fmt"

type ErrorResponse struct {
	Message string `json:"msg"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf(`{"msg": "%s"}`, e.Message)
}

func NewErrorResponse(msg string) *ErrorResponse {
	return &ErrorResponse{
		Message: msg,
	}
}

func (e *ErrorResponse) ByteResponse() []byte {
	return []byte(e.Error())
}

func ErrorToResponse(err error) []byte {
	return []byte(err.Error())
}
