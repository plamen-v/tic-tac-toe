package errors

import "fmt"

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func NewNotFoundError(msg string) error {
	return &NotFoundError{
		Message: msg}
}

func NewNotFoundErrorf(format string, arg ...any) error {
	msg := fmt.Sprintf(format, arg...)
	return &NotFoundError{
		Message: msg}
}

func IsNotFoundError(err error) bool {
	if _, ok := err.(*NotFoundError); ok {
		return true
	}
	return false
}
