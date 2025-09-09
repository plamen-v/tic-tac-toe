package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/plamen-v/tic-tac-toe/src/app/server/middleware"
	"github.com/plamen-v/tic-tac-toe/src/models"
	"github.com/plamen-v/tic-tac-toe/src/services/game"
)

func CreateRoomHandler(game game.GameEngine) func(*gin.Context) {
	return func(c *gin.Context) {
		var createRoomRequest models.CreateRoomRequest
		var err error
		// Parse and bind JSON to struct
		if err = c.BindJSON(&createRoomRequest); err != nil {
			c.Error(err)
			return
		}

		id, exists := c.Get(middleware.KEY_PLAYER_ID) //

		room := &models.Room{
			Host: models.RoomParticipant{
				ID: id,
			},
			Title:       createRoomRequest.Title,
			Description: createRoomRequest.Description,
		}
		roomID, err := game.CreateRoom(room)
		if err != nil {
			c.Error(err)
			return
		}

		c.Header("Location", fmt.Sprintf("/rooms/%d", roomID))
		c.JSON(http.StatusOK, roomID)

	}
}
