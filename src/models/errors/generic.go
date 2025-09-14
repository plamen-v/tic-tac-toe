package errors

import "fmt"

type GenericError struct {
	Message string
	Code    ErrorCode
	Cause   error
}

func (e *GenericError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Cause.Error())
	}
	return e.Message
}

func NewGenericError(msg string) error {
	return &GenericError{
		Message: msg}
}

func NewGenericErrorWithCause(msg string, cause error) error {
	return &GenericError{
		Message: msg,
		Cause:   cause,
	}
}

func NewGenericErrorf(format string, arg ...any) error {
	msg := fmt.Sprintf(format, arg...)
	return &GenericError{
		Message: msg}
}

func IsGenericError(err error) bool {
	if _, ok := err.(*GenericError); ok {
		return true
	}
	return false
}
