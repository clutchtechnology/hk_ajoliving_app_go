package tools

import "errors"

var (
	ErrNotFound       = errors.New("resource not found")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrInvalidInput   = errors.New("invalid input")
	ErrAlreadyExists  = errors.New("resource already exists")
	ErrInternalServer = errors.New("internal server error")
	ErrNotImplemented = errors.New("not implemented")
)

// BusinessError 自定义业务错误
type BusinessError struct {
	Code    int
	Message string
	Err     error
}

func (e *BusinessError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// NewBusinessError 创建业务错误
func NewBusinessError(code int, message string) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: message,
	}
}

// NewBusinessErrorWithErr 创建带底层错误的业务错误
func NewBusinessErrorWithErr(code int, message string, err error) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
