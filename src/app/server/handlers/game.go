package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"

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

		playerID, ok := GetPlayerIDFromContext(c, middleware.KEY_PLAYER_ID)
		if !ok {
			_ = c.Error(models.NewValidationError("Missing player_id claim"))
			return
		}

		room := &models.Room{
			Host: models.RoomPlayer{
				ID:             playerID,
				RequestNewGame: true,
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
			Payload:    roomID,
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
		rooms, err := gameEngineService.GetOpenRooms(c.Request.Context(), request.Keyword)
		if err != nil {
			_ = c.Error(err)
			return
		}

		response := models.Response{
			StatusCode: http.StatusOK,
			Payload:    rooms,
		}

		c.JSON(response.StatusCode, response)
	}
}

func PlayerJoinRoomHandler(gameEngineService engine.GameEngineService) func(*gin.Context) {
	return func(c *gin.Context) {

		pRoomID := c.Param("roomId")
		roomID, err := uuid.FromString(pRoomID)
		if err != nil {
			_ = c.Error(models.NewValidationErrorf("Invalid room id '%s'", pRoomID))
			return
		}

		playerID, ok := GetPlayerIDFromContext(c, middleware.KEY_PLAYER_ID)
		if !ok {
			_ = c.Error(models.NewValidationError("Missing player_id claim"))
			return
		}

		err = gameEngineService.PlayerJoinRoom(c.Request.Context(), roomID, playerID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		response := models.Response{
			StatusCode: http.StatusOK,
		}

		c.JSON(response.StatusCode, response)
	}
}

func PlayerLeaveRoomHandler(gameEngineService engine.GameEngineService) func(*gin.Context) {
	return func(c *gin.Context) {
		pRoomID := c.Param("roomId")
		roomID, err := uuid.FromString(pRoomID)
		if err != nil {
			_ = c.Error(models.NewValidationErrorf("Invalid id '%s'", pRoomID))
			return
		}

		playerID, ok := GetPlayerIDFromContext(c, middleware.KEY_PLAYER_ID)
		if !ok {
			_ = c.Error(models.NewValidationError("Missing player_id claim"))
			return
		}

		err = gameEngineService.PlayerLeaveRoom(c.Request.Context(), roomID, playerID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		response := models.Response{
			StatusCode: http.StatusOK,
		}

		c.JSON(response.StatusCode, response)
	}
}

func CreateGameHandler(gameEngineService engine.GameEngineService) func(*gin.Context) {
	return func(c *gin.Context) {
		pRoomID := c.Param("roomId")
		roomID, err := uuid.FromString(pRoomID)
		if err != nil {
			_ = c.Error(models.NewValidationErrorf("Invalid id '%s'", pRoomID))
			return
		}

		playerID, ok := GetPlayerIDFromContext(c, middleware.KEY_PLAYER_ID)
		if !ok {
			_ = c.Error(models.NewValidationError("Missing player_id claim"))
			return
		}

		gameID, err := gameEngineService.CreateGame(c.Request.Context(), roomID, playerID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		status := http.StatusAccepted
		var payload *uuid.UUID
		if gameID != uuid.Nil {
			status = http.StatusCreated
			payload = &gameID
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
		roomID, err := uuid.FromString(pRoomID)
		if err != nil {
			_ = c.Error(models.NewValidationErrorf("Invalid id '%s'", pRoomID))
			return
		}

		playerID, ok := GetPlayerIDFromContext(c, middleware.KEY_PLAYER_ID)
		if !ok {
			_ = c.Error(models.NewValidationError("Missing player_id claim"))
			return
		}

		game, err := gameEngineService.GetGameState(c.Request.Context(), roomID, playerID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		response := models.Response{
			StatusCode: http.StatusOK,
			Payload:    game,
		}

		c.JSON(response.StatusCode, response)
	}
}

func MakeMoveHandler(gameEngineService engine.GameEngineService) func(*gin.Context) {
	return func(c *gin.Context) {

		pRoomID := c.Param("roomId")
		roomID, err := uuid.FromString(pRoomID)
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

		playerID, ok := GetPlayerIDFromContext(c, middleware.KEY_PLAYER_ID)
		if !ok {
			_ = c.Error(models.NewValidationError("Missing player_id claim"))
			return
		}

		err = gameEngineService.PlayerMakeMove(c.Request.Context(), roomID, playerID, position)
		if err != nil {
			_ = c.Error(err)
			return
		}

		response := models.Response{
			StatusCode: http.StatusOK,
		}

		c.JSON(response.StatusCode, response)
	}
}

func GetPlayerIDFromContext(c *gin.Context, key string) (uuid.UUID, bool) {
	val, exists := c.Get(key)
	if !exists {
		return uuid.Nil, false
	}
	u, ok := val.(uuid.NullUUID)
	return u.UUID, ok
}
