package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/plamen-v/tic-tac-toe-models/models"
	"github.com/plamen-v/tic-tac-toe/src/app/server/middleware"
	"github.com/plamen-v/tic-tac-toe/src/services/engine"
)

func CreateRoomHandler(gameEngineService engine.GameEngineService) func(*gin.Context) {
	return func(c *gin.Context) {
		var request models.CreateRoomRequest
		var err error
		if err = c.BindJSON(&request); err != nil {
			_ = c.Error(models.NewValidationError("bad request"))
			return
		}

		playerID := c.GetInt64(middleware.KEY_PLAYER_ID)
		room := &models.Room{
			Host: models.RoomParticipant{
				ID:       playerID,
				Continue: true,
			},
			Title:       request.Title,
			Description: request.Description,
		}

		roomID, err := gameEngineService.CreateRoom(c.Request.Context(), room)
		if err != nil {
			_ = c.Error(err)
			return
		}

		response := models.Response{
			StatusCode: http.StatusCreated,
			Payload:    gin.H{"id": roomID},
		}

		c.JSON(response.StatusCode, response)
	}
}

func GetOpenRoomsHandler(gameEngineService engine.GameEngineService) func(*gin.Context) {
	return func(c *gin.Context) {
		var request models.GetOpenRoomsRequest
		var err error
		if err = c.BindJSON(&request); err != nil {
			_ = c.Error(models.NewValidationError("bad request"))
			return
		}
		rooms, err := gameEngineService.GetOpenRooms(c.Request.Context(), request.Host, request.Title, request.Description, request.Phase)
		if err != nil {
			_ = c.Error(err)
			return
		}

		response := models.Response{
			StatusCode: http.StatusOK,
			Payload:    gin.H{"rooms": rooms},
		}

		c.JSON(response.StatusCode, response)
	}
}

func PlayerJoinRoomHandler(gameEngineService engine.GameEngineService) func(*gin.Context) {
	return func(c *gin.Context) {

		pRoomID := c.Param("roomId")
		roomID, err := strconv.ParseInt(pRoomID, 10, 64)
		if err != nil {
			_ = c.Error(models.NewValidationErrorf("Invalid room id '%s'", pRoomID))
			return
		}

		playerID := c.GetInt64(middleware.KEY_PLAYER_ID)

		err = gameEngineService.PlayerJoinRoom(c.Request.Context(), roomID, playerID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		response := models.Response{
			StatusCode: http.StatusOK,
			Payload:    gin.H{},
		}

		c.JSON(response.StatusCode, response)
	}
}

func PlayerLeaveRoomHandler(gameEngineService engine.GameEngineService) func(*gin.Context) {
	return func(c *gin.Context) {
		pRoomID := c.Param("roomId")
		roomID, err := strconv.ParseInt(pRoomID, 10, 64)
		if err != nil {
			_ = c.Error(models.NewValidationErrorf("Invalid id '%s'", pRoomID))
			return
		}

		playerID := c.GetInt64(middleware.KEY_PLAYER_ID)

		err = gameEngineService.PlayerLeaveRoom(c.Request.Context(), roomID, playerID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		response := models.Response{
			StatusCode: http.StatusOK,
			Payload:    gin.H{},
		}

		c.JSON(response.StatusCode, response)
	}
}

func CreateGameHandler(gameEngineService engine.GameEngineService) func(*gin.Context) {
	return func(c *gin.Context) {

		pRoomID := c.Param("roomId")
		roomID, err := strconv.ParseInt(pRoomID, 10, 64)
		if err != nil {
			_ = c.Error(models.NewValidationErrorf("Invalid id '%s'", pRoomID))
			return
		}

		playerID := c.GetInt64(middleware.KEY_PLAYER_ID)

		gameID, err := gameEngineService.CreateGame(c.Request.Context(), roomID, playerID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		status := http.StatusAccepted
		payload := gin.H{}
		if gameID != 0 {
			status = http.StatusCreated
			payload = gin.H{"id": gameID}
		}

		response := models.Response{
			StatusCode: status,
			Payload:    payload,
		}

		c.JSON(response.StatusCode, response)
	}
}

func GetGameStateHandler(gameEngineService engine.GameEngineService) func(*gin.Context) {
	return func(c *gin.Context) {

		pRoomID := c.Param("roomId")
		roomID, err := strconv.ParseInt(pRoomID, 10, 64)
		if err != nil {
			_ = c.Error(models.NewValidationErrorf("Invalid id '%s'", pRoomID))
			return
		}

		playerID := c.GetInt64(middleware.KEY_PLAYER_ID)

		game, err := gameEngineService.GetGameState(c.Request.Context(), roomID, playerID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		response := models.Response{
			StatusCode: http.StatusOK,
			Payload:    gin.H{"game": game},
		}

		c.JSON(response.StatusCode, response)
	}
}

func MakeMoveHandler(gameEngineService engine.GameEngineService) func(*gin.Context) {
	return func(c *gin.Context) {

		pRoomID := c.Param("roomId")
		roomID, err := strconv.ParseInt(pRoomID, 10, 64)
		if err != nil {
			_ = c.Error(models.NewValidationErrorf("Invalid id '%s'", pRoomID))
			return
		}

		pPosition := c.Param("position")
		position, err := strconv.Atoi(pPosition)
		if err != nil {
			_ = c.Error(models.NewValidationErrorf("Invalid position '%s'", pPosition))
			return
		}

		playerID := c.GetInt64(middleware.KEY_PLAYER_ID)

		err = gameEngineService.PlayerMakeMove(c.Request.Context(), roomID, playerID, position)
		if err != nil {
			_ = c.Error(err)
			return
		}

		response := models.Response{
			StatusCode: http.StatusOK,
			Payload:    gin.H{},
		}

		c.JSON(response.StatusCode, response)
	}
}
