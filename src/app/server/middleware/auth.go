package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/plamen-v/tic-tac-toe/src/services/auth"
)

const KEY_PLAYER_ID string = "KEY_PLAYER_ID"

func Authentication(authService auth.AuthenticationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader(auth.AUTHORIZATION_HEADER)
		tokenString = strings.TrimPrefix(tokenString, auth.AUTHORIZATION_HEADER_PREFIX)
		jwtToken, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		if claims, ok := jwtToken.Claims.(*auth.ExtendedClaims); ok {
			c.Set(KEY_PLAYER_ID, claims.PlayerID)
		} else {
			_ = c.Error(fmt.Errorf("Invalid token"))
			c.Abort()
		}

		c.Next()
	}
}
