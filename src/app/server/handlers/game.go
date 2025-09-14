package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/plamen-v/tic-tac-toe/src/app/server/middleware"
	"github.com/plamen-v/tic-tac-toe/src/models"
	"github.com/plamen-v/tic-tac-toe/src/models/requests"
	"github.com/plamen-v/tic-tac-toe/src/services/engine"
)

func CreateNewRoomHandler(gameEngineService engine.GameEngineService) func(*gin.Context) {
	return func(c *gin.Context) {
		var createRoomRequest requests.CreateRoomRequest
		var err error
		// Parse and bind JSON to struct
		if err = c.BindJSON(&createRoomRequest); err != nil {
			c.Error(err)
			return
		}

		//todo!id, exists := c.Get(middleware.KEY_PLAYER_ID) //
		//token guarantees that KEY_PLAYER_ID is set
		hostID := c.GetInt64(middleware.KEY_PLAYER_ID)
		room := &models.Room{
			Host: models.RoomParticipant{
				ID: &hostID, //todo!
			},
			Title:       createRoomRequest.Title,
			Description: createRoomRequest.Description,
		}
		roomID, err := gameEngineService.CreateNewRoom(room)
		if err != nil {
			c.Error(err)
			return
		}

		c.Header("Location", fmt.Sprintf("/rooms/%d", roomID))
		c.JSON(http.StatusOK, roomID)

	}
}
