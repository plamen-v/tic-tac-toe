package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/plamen-v/tic-tac-toe/src/services/auth"
)

const KEY_PLAYER_ID string = "KEY_PLAYER_ID"

func AuthenticationMiddleware(authService auth.AuthenticationService) func(*gin.Context) {
	return func(c *gin.Context) {

		tokenString := c.GetHeader(auth.AUTHORIZATION_HEADER)
		tokenString = strings.TrimPrefix(tokenString, auth.AUTHORIZATION_HEADER_PREFIX)

		jwtToken, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.Error(err)
			//TODO!
			c.Abort()
			return
		}
		//c.Set("token", jwtToken) //todo!
		playerID, _ := jwtToken.Claims.GetIssuer()
		c.Set(KEY_PLAYER_ID, playerID) //todo!
		c.Next()
	}
}
