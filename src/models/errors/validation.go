package errors

import "fmt"

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func NewValidationError(msg string) error {
	return &ValidationError{
		Message: msg}
}

func NewValidationErrorf(format string, arg ...any) error {
	msg := fmt.Sprintf(format, arg...)
	return &ValidationError{
		Message: msg}
}

func IsValidationErrorf(err error) bool {
	if _, ok := err.(*ValidationError); ok {
		return true
	}
	return false
}
