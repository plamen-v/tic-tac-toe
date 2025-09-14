package errors

type ErrorCode string

const (
	UnauthorizedErrorCode        ErrorCode = "UNAUTHORIZED"
	InternalServerErrorErrorCode ErrorCode = "INTERNAL_SERVER_ERROR"
	NotFoundErrorCode            ErrorCode = "NOT_FOUND"
	InvalidInputErrorCode        ErrorCode = "INVALID_INPUT"
)

const InternalServerErrorMessage string = "Internal server error"

type ErrorMessage struct {
	Code    string `json:"error_code"`
	Message string `json:"error_message"`
}

func NewErrorMessage(code string, message string) ErrorMessage {
	return ErrorMessage{
		Code:    code,
		Message: message,
	}
}
