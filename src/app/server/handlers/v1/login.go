package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/plamen-v/tic-tac-toe/src/models"
	"github.com/plamen-v/tic-tac-toe/src/services/auth"
)

func LoginHandler(authService auth.AuthenticationService) func(*gin.Context) {
	return func(c *gin.Context) {
		var loginRequest models.LoginRequest
		var err error
		// Parse and bind JSON to struct
		if err = c.BindJSON(&loginRequest); err != nil {
			c.Error(err)
			return
		}

		player, err := authService.AuthenticatePlayer(loginRequest.Login, loginRequest.Password)
		if err != nil {
			c.Error(err)
			return
		}

		tokenStr, err := authService.CreateToken(player)
		if err != nil {
			c.Error(err)
			return
		}

		c.Header(auth.AUTHORIZATION_HEADER, fmt.Sprintf("%s%s", auth.AUTHORIZATION_HEADER_PREFIX, tokenStr))
		c.JSON(http.StatusOK, player)
	}
}
