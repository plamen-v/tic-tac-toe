package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/plamen-v/tic-tac-toe-models/models/errors"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			var statusCode int
			var errorCode, errorMessage string
			switch coreError := err.(type) {
			case *errors.NotFoundError:
				statusCode = http.StatusNotFound
				errorCode = string(errors.NotFoundErrorCode)
				errorMessage = coreError.Error()
			case *errors.ValidationError:
				statusCode = http.StatusBadRequest
				errorCode = string(errors.BadRequestErrorCode)
				errorMessage = coreError.Error()
			case *errors.AuthorizationError:
				statusCode = http.StatusForbidden
				errorCode = string(errors.UnauthorizedErrorCode)
				errorMessage = errors.AuthorizationErrorMessage
			case *errors.GenericError:
				statusCode = http.StatusInternalServerError
				errorCode = string(coreError.Code)
				errorMessage = errors.InternalServerErrorMessage
			default:
				statusCode = http.StatusInternalServerError
				errorCode = string(errors.InternalServerErrorErrorCode)
				errorMessage = errors.InternalServerErrorMessage
			}

			c.JSON(statusCode, errors.NewErrorMessage(errorCode, errorMessage))
		}
	}
}
