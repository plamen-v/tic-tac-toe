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

func GetRoomHandler(gameEngineService engine.GameEngineService) func(*gin.Context) {
	return func(c *gin.Context) {
		playerID, ok := getPlayerIDFromContext(c, middleware.KEY_PLAYER_ID)
		if !ok {
			_ = c.Error(models.NewValidationError("Missing player_id claim"))
			return
		}

		room, err := gameEngineService.GetRoom(c.Request.Context(), playerID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		response := models.RoomResponse{
			Room: room,
		}

		c.JSON(http.StatusOK, response)
	}
}

func GetOpenRoomsHandler(gameEngineService engine.GameEngineService) func(*gin.Context) {
	return func(c *gin.Context) {
		rooms, err := gameEngineService.GetOpenRooms(c.Request.Context())
		if err != nil {
			_ = c.Error(err)
			return
		}

		response := models.RoomListResponse{
			Rooms: rooms,
		}

		c.JSON(http.StatusOK, response)
	}
}

func CreateRoomHandler(gameEngineService engine.GameEngineService) func(*gin.Context) {
	return func(c *gin.Context) {
		var request models.CreateRoomRequest
		var err error
		if err = c.BindJSON(&request); err != nil {
			_ = c.Error(models.NewValidationError("bad request"))
			return
		}

		playerID, ok := getPlayerIDFromContext(c, middleware.KEY_PLAYER_ID)
		if !ok {
			_ = c.Error(models.NewValidationError("Missing player_id claim"))
			return
		}

		roomID, err := gameEngineService.CreateRoom(c.Request.Context(), playerID, request.Title, request.Description)
		if err != nil {
			_ = c.Error(err)
			return
		}

		response := models.CreateRoomResponse{
			RoomID: roomID,
		}

		c.JSON(http.StatusCreated, response)
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

		playerID, ok := getPlayerIDFromContext(c, middleware.KEY_PLAYER_ID)
		if !ok {
			_ = c.Error(models.NewValidationError("Missing player_id claim"))
			return
		}

		err = gameEngineService.PlayerJoinRoom(c.Request.Context(), roomID, playerID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.Status(http.StatusOK)
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

		playerID, ok := getPlayerIDFromContext(c, middleware.KEY_PLAYER_ID)
		if !ok {
			_ = c.Error(models.NewValidationError("Missing player_id claim"))
			return
		}

		err = gameEngineService.PlayerLeaveRoom(c.Request.Context(), roomID, playerID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.Status(http.StatusOK)
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

		playerID, ok := getPlayerIDFromContext(c, middleware.KEY_PLAYER_ID)
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
		if gameID != uuid.Nil {
			status = http.StatusCreated
		}

		c.Status(status)
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

		playerID, ok := getPlayerIDFromContext(c, middleware.KEY_PLAYER_ID)
		if !ok {
			_ = c.Error(models.NewValidationError("Missing player_id claim"))
			return
		}

		game, err := gameEngineService.GetGameState(c.Request.Context(), roomID, playerID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		response := models.GameResponse{
			Game: game,
		}

		c.JSON(http.StatusOK, response)
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

		playerID, ok := getPlayerIDFromContext(c, middleware.KEY_PLAYER_ID)
		if !ok {
			_ = c.Error(models.NewValidationError("Missing player_id claim"))
			return
		}

		err = gameEngineService.PlayerMakeMove(c.Request.Context(), roomID, playerID, position)
		if err != nil {
			_ = c.Error(err)
			return
		}
		c.Status(http.StatusOK)
	}
}

func GetRankingHandler(gameEngineService engine.GameEngineService) func(*gin.Context) {
	return func(c *gin.Context) {
		players, err := gameEngineService.GetRanking(c.Request.Context())
		if err != nil {
			_ = c.Error(err)
			return
		}

		response := models.RankResponse{
			Players: players,
		}

		c.JSON(http.StatusOK, response)
	}
}

func getPlayerIDFromContext(c *gin.Context, key string) (uuid.UUID, bool) {
	val, exists := c.Get(key)
	if !exists {
		return uuid.Nil, false
	}
	u, ok := val.(uuid.NullUUID)
	return u.UUID, ok
}
