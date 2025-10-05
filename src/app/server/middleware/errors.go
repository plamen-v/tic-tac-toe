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

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			var statusCode int
			var errorCode, errorMessage string
			switch {
			case errors.As(err, new(*models.NotFoundError)):
				statusCode = http.StatusNotFound
				errorCode = string(models.NotFoundErrorCode)
				errorMessage = err.Error()
			case errors.As(err, new(*models.ValidationError)):
				statusCode = http.StatusBadRequest
				errorCode = string(models.BadRequestErrorCode)
				errorMessage = err.Error()
			case errors.As(err, new(*models.AuthorizationError)):
				statusCode = http.StatusUnauthorized
				errorCode = string(models.UnauthorizedErrorCode)
				errorMessage = models.AuthorizationErrorMessage
			case errors.As(err, new(*models.GenericError)):
				statusCode = http.StatusInternalServerError
				errorCode = string(models.InternalServerErrorErrorCode)
				errorMessage = models.InternalServerErrorMessage
			default:
				statusCode = http.StatusInternalServerError
				errorCode = string(models.InternalServerErrorErrorCode)
				errorMessage = models.InternalServerErrorMessage
			}

			response := models.ErrorResponse{
				Code:    errorCode,
				Message: errorMessage,
			}

			c.JSON(statusCode, response)
		}
	}
}
