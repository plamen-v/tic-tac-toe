package errors

import "fmt"

type AuthorizationError struct {
	Message string
}

func (e *AuthorizationError) Error() string {
	return e.Message
}

func NewAuthorizationError(msg string) error {
	return &AuthorizationError{
		Message: msg}
}

func NewAuthorizationErrorf(format string, arg ...any) error {
	msg := fmt.Sprintf(format, arg...)
	return &AuthorizationError{
		Message: msg}
}

func IsAuthorizationError(err error) bool {
	if _, ok := err.(*AuthorizationError); ok {
		return true
	}
	return false
}
