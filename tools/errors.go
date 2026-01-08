package tools

import (
	"errors"
	"fmt"
)

// 预定义错误
var (
	ErrNotFound       = errors.New("resource not found")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrInvalidInput   = errors.New("invalid input")
	ErrAlreadyExists  = errors.New("resource already exists")
	ErrInternalServer = errors.New("internal server error")
	ErrInvalidToken   = errors.New("invalid token")
	ErrExpiredToken   = errors.New("token expired")
)

// BusinessError 业务错误
type BusinessError struct {
	Code    int
	Message string
	Err     error
}

func (e *BusinessError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *BusinessError) Unwrap() error {
	return e.Err
}

// NewError 创建业务错误
func NewError(code int, message string) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: message,
	}
}

// WrapError 包装底层错误
func WrapError(code int, message string, err error) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
