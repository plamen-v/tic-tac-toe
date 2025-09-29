package middleware

import (
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
			switch coreError := err.(type) {
			case *models.NotFoundError:
				statusCode = http.StatusNotFound
				errorCode = string(models.NotFoundErrorCode)
				errorMessage = coreError.Error()
			case *models.ValidationError:
				statusCode = http.StatusBadRequest
				errorCode = string(models.BadRequestErrorCode)
				errorMessage = coreError.Error()
			case *models.AuthorizationError:
				statusCode = http.StatusUnauthorized
				errorCode = string(models.UnauthorizedErrorCode)
				errorMessage = models.AuthorizationErrorMessage
			case *models.GenericError:
				statusCode = http.StatusInternalServerError
				errorCode = string(coreError.Code)
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
