package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/plamen-v/tic-tac-toe-models/models"
	"github.com/plamen-v/tic-tac-toe/src/services/auth"
)

func LoginHandler(authService auth.AuthenticationService) func(*gin.Context) {
	return func(c *gin.Context) {
		var loginRequest models.LoginRequest
		var err error
		if err = c.BindJSON(&loginRequest); err != nil {
			_ = c.Error(models.NewValidationError("bad request"))
			return
		}

		player, tokenStr, err := authService.Authenticate(c.Request.Context(), loginRequest.Login, loginRequest.Password)
		if err != nil {
			_ = c.Error(err)
			return
		}

		response := models.LoginResponse{
			Player: player,
			Token:  tokenStr,
		}

		c.JSON(http.StatusOK, response)
	}
}
