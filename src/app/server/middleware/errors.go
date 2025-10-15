package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/plamen-v/tic-tac-toe-models/models"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		var validationError *models.ValidationError
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			var statusCode int
			var errorCode models.ErrorCode
			var errorMessage string
			switch {
			case errors.As(err, new(*models.NotFoundError)):
				statusCode = http.StatusNotFound
				errorCode = models.NotFoundErrorCode
				errorMessage = err.Error()
			case errors.As(err, &validationError):
				statusCode = http.StatusBadRequest
				errorCode = validationError.Code()
				errorMessage = err.Error()
			case errors.As(err, new(*models.AuthorizationError)):
				statusCode = http.StatusUnauthorized
				errorCode = models.UnauthorizedErrorCode
				errorMessage = models.AuthorizationErrorMessage
			case errors.As(err, new(*models.GenericError)):
				statusCode = http.StatusInternalServerError
				errorCode = models.InternalServerErrorErrorCode
				errorMessage = models.InternalServerErrorMessage
			default:
				statusCode = http.StatusInternalServerError
				errorCode = models.InternalServerErrorErrorCode
				errorMessage = models.InternalServerErrorMessage
			}

			response := models.ErrorResponse{
				Code:    string(errorCode),
				Message: errorMessage,
			}

			c.JSON(statusCode, response)
		}
	}
}
