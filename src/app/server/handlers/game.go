package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/plamen-v/tic-tac-toe-models/models"
	"github.com/plamen-v/tic-tac-toe-models/models/errors"
	"github.com/plamen-v/tic-tac-toe-models/models/requests"
	"github.com/plamen-v/tic-tac-toe/src/app/server/middleware"
	"github.com/plamen-v/tic-tac-toe/src/services/engine"
)

func CreateRoomHandler(gameEngineService engine.GameEngineService) func(*gin.Context) {
	return func(c *gin.Context) {
		var request requests.CreateRoomRequest
		var err error
		// Parse and bind JSON to struct
		if err = c.BindJSON(&request); err != nil {
			_ = c.Error(err)
			return
		}

		//token guarantees that KEY_PLAYER_ID is set
		hostID := c.GetInt64(middleware.KEY_PLAYER_ID)
		room := &models.Room{
			Host: models.RoomParticipant{
				ID:       hostID,
				Continue: true,
			},
			Title:       request.Title,
			Description: request.Description,
		}

		roomID, err := gameEngineService.CreateRoom(room)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.Header("Location", fmt.Sprintf("/rooms/%d", roomID))
		c.JSON(http.StatusCreated, roomID)
	}
}

func GetOpenRoomsHandler(gameEngineService engine.GameEngineService) func(*gin.Context) {
	return func(c *gin.Context) {
		var request requests.GetOpenRoomsRequest
		var err error
		// Parse and bind JSON to struct
		if err = c.BindJSON(&request); err != nil {
			_ = c.Error(err)
			return
		}
		rooms, err := gameEngineService.GetOpenRooms(request.Host, request.Title, request.Description, request.Phase)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, rooms)
	}
}

func GetRoomStateHandler(gameEngineService engine.GameEngineService) func(*gin.Context) {
	return func(c *gin.Context) {

		pRoomID := c.Param("roomId")
		roomID, err := strconv.ParseInt(pRoomID, 10, 64)
		if err != nil {
			c.Error(errors.NewValidationErrorf("Invalid id '%s'", pRoomID)) //todo!
			return
		}

		//token guarantees that KEY_PLAYER_ID is set
		playerID := c.GetInt64(middleware.KEY_PLAYER_ID)

		room, err := gameEngineService.GetRoomState(roomID, playerID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.Header("Location", fmt.Sprintf("/rooms/%d", room.ID))
		c.JSON(http.StatusOK, room)
	}
}

func RegisterGuestHandler(gameEngineService engine.GameEngineService) func(*gin.Context) {
	return func(c *gin.Context) {

		pRoomID := c.Param("roomId")
		roomID, err := strconv.ParseInt(pRoomID, 10, 64)
		if err != nil {
			c.Error(errors.NewValidationErrorf("Invalid id '%s'", pRoomID)) //todo!
			return
		}

		//token guarantees that KEY_PLAYER_ID is set
		guestID := c.GetInt64(middleware.KEY_PLAYER_ID)

		err = gameEngineService.RegisterGuest(roomID, guestID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.Header("Location", fmt.Sprintf("/rooms/%d", roomID))
		c.JSON(http.StatusOK, gin.H{}) //todo!
	}
}

func GuestLeaveHandler(gameEngineService engine.GameEngineService) func(*gin.Context) {
	return func(c *gin.Context) {

		pRoomID := c.Param("roomId")
		roomID, err := strconv.ParseInt(pRoomID, 10, 64)
		if err != nil {
			c.Error(errors.NewValidationErrorf("Invalid id '%s'", pRoomID)) //todo!
			return
		}

		//token guarantees that KEY_PLAYER_ID is set
		guestID := c.GetInt64(middleware.KEY_PLAYER_ID)

		err = gameEngineService.GuestLeave(roomID, guestID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.Header("Location", "/rooms/")
		c.JSON(http.StatusOK, gin.H{}) //todo!
	}
}

func HostLeaveHandler(gameEngineService engine.GameEngineService) func(*gin.Context) {
	return func(c *gin.Context) {

		pRoomID := c.Param("roomId")
		roomID, err := strconv.ParseInt(pRoomID, 10, 64)
		if err != nil {
			c.Error(errors.NewValidationErrorf("Invalid id '%s'", pRoomID)) //todo!
			return
		}

		//token guarantees that KEY_PLAYER_ID is set
		hostID := c.GetInt64(middleware.KEY_PLAYER_ID)

		err = gameEngineService.HostLeave(roomID, hostID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.Header("Location", "/rooms/")
		c.JSON(http.StatusOK, gin.H{}) //todo!
	}
}

func CreateGameHandler(gameEngineService engine.GameEngineService) func(*gin.Context) {
	return func(c *gin.Context) {

		pRoomID := c.Param("roomId")
		roomID, err := strconv.ParseInt(pRoomID, 10, 64)
		if err != nil {
			c.Error(errors.NewValidationErrorf("Invalid id '%s'", pRoomID)) //todo!
			return
		}

		//token guarantees that KEY_PLAYER_ID is set
		playerID := c.GetInt64(middleware.KEY_PLAYER_ID)

		gameID, err := gameEngineService.CreateGame(roomID, playerID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		status := http.StatusAccepted
		location := fmt.Sprintf("/rooms/%d/", roomID)
		if gameID != 0 {
			status = http.StatusCreated
			location = fmt.Sprintf("/rooms/%d/games/%d", roomID, gameID)
		}

		c.Header("Location", location)
		c.JSON(status, gin.H{}) //todo!
	}
}

func GetGameStateHandler(gameEngineService engine.GameEngineService) func(*gin.Context) {
	return func(c *gin.Context) {

		pRoomID := c.Param("roomId")
		roomID, err := strconv.ParseInt(pRoomID, 10, 64)
		if err != nil {
			c.Error(errors.NewValidationErrorf("Invalid id '%s'", pRoomID)) //todo!
			return
		}

		pGameID := c.Param("roomId")
		gameID, err := strconv.ParseInt(pGameID, 10, 64)
		if err != nil {
			c.Error(errors.NewValidationErrorf("Invalid id '%s'", pGameID)) //todo!
			return
		}

		//token guarantees that KEY_PLAYER_ID is set
		playerID := c.GetInt64(middleware.KEY_PLAYER_ID)

		game, err := gameEngineService.GetGameState(roomID, gameID, playerID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.Header("Location", fmt.Sprintf("/rooms/%d/games/%d", roomID, gameID))
		c.JSON(http.StatusOK, game)
	}
}

func SetMarkHandler(gameEngineService engine.GameEngineService) func(*gin.Context) {
	return func(c *gin.Context) {

		pRoomID := c.Param("roomId")
		roomID, err := strconv.ParseInt(pRoomID, 10, 64)
		if err != nil {
			c.Error(errors.NewValidationErrorf("Invalid id '%s'", pRoomID)) //todo!
			return
		}

		pGameID := c.Param("roomId")
		gameID, err := strconv.ParseInt(pGameID, 10, 64)
		if err != nil {
			c.Error(errors.NewValidationErrorf("Invalid id '%s'", pGameID)) //todo!
			return
		}

		pPosition := c.Param("position")
		position, err := strconv.Atoi(pPosition)
		if err != nil {
			c.Error(errors.NewValidationErrorf("Invalid position '%s'", pPosition)) //todo!
			return
		}
		//token guarantees that KEY_PLAYER_ID is set
		playerID := c.GetInt64(middleware.KEY_PLAYER_ID)

		err = gameEngineService.SetMark(roomID, gameID, playerID, position)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.Header("Location", fmt.Sprintf("/rooms/%d/games/%d", roomID, gameID))
		c.JSON(http.StatusOK, gin.H{})
	}
}
