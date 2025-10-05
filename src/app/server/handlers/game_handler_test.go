package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	. "github.com/onsi/ginkgo/v2"
	"github.com/plamen-v/tic-tac-toe-models/models"
	"github.com/plamen-v/tic-tac-toe/src/app/server/handlers"
	"github.com/plamen-v/tic-tac-toe/src/app/server/middleware"
	"github.com/plamen-v/tic-tac-toe/src/services/engine/mocks"
	"github.com/stretchr/testify/mock"

	. "github.com/onsi/gomega"
)

var _ = Describe("GameHandler", func() {
	var (
		mockGameEngineService *mocks.MockGameEngineService
		router                *gin.Engine
	)

	BeforeEach(func() {
		mockGameEngineService = new(mocks.MockGameEngineService)
		gin.SetMode(gin.TestMode)
		router = gin.Default()
		router.Use(middleware.ErrorHandler())
	})

	Context("GetRoom", func() {
		It("should return 400 if playerID missing from context", func() {
			request, err := http.NewRequest("GET", "/room", nil)
			Expect(err).To(BeNil())
			handler := handlers.GetRoomHandler(mockGameEngineService)
			router.GET("/room", handler)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 500 if server error occurs", func() {
			request, err := http.NewRequest("GET", "/room", nil)
			Expect(err).To(BeNil())

			mockGameEngineService.On("GetRoom", mock.Anything, mock.Anything).Return(nil, models.NewGenericError("server error"))
			handler := handlers.GetRoomHandler(mockGameEngineService)
			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			router.Use(insertPlayerIDInContextMiddleware(playerID))
			router.GET("/room", handler)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should return 200 if request is OK", func() {
			request, err := http.NewRequest("GET", "/room", nil)
			Expect(err).To(BeNil())
			mockGameEngineService.On("GetRoom", mock.Anything, mock.Anything).Return(&models.Room{}, nil)
			handler := handlers.GetRoomHandler(mockGameEngineService)
			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			router.Use(insertPlayerIDInContextMiddleware(playerID))
			router.GET("/room", handler)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusOK))
		})
	})

	Context("GetOpenRooms", func() {
		It("should return 200 if request is OK", func() {
			request, err := http.NewRequest("GET", "/rooms", nil)
			Expect(err).To(BeNil())
			response := httptest.NewRecorder()
			mockGameEngineService.On("GetOpenRooms", mock.Anything, mock.Anything).Return([]*models.Room{}, nil)
			handler := handlers.GetOpenRoomsHandler(mockGameEngineService)
			router.GET("/rooms", handler)
			router.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusOK))
		})

		It("should return 500 if internal error occurs", func() {
			request, err := http.NewRequest("GET", "/rooms", nil)
			Expect(err).To(BeNil())
			handler := handlers.GetOpenRoomsHandler(mockGameEngineService)
			router.GET("/rooms", handler)
			mockGameEngineService.On("GetOpenRooms", mock.Anything).Return(nil, models.NewGenericError("server error"))
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	Context("CreateRoomHandler", func() {
		It("should return 400 if request is invalid", func() {
			request, err := http.NewRequest("POST", "/rooms", nil)
			Expect(err).To(BeNil())
			request.Header.Set("Content-Type", "application/json")
			handler := handlers.CreateRoomHandler(mockGameEngineService)
			router.POST("/rooms", handler)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 400 if playerID missing from context", func() {
			createRoomRequest := models.CreateRoomRequest{
				Title:       "title",
				Description: "description",
			}
			requestBody, err := json.Marshal(createRoomRequest)
			Expect(err).To(BeNil())
			request, err := http.NewRequest("POST", "/rooms", bytes.NewBuffer(requestBody))
			Expect(err).To(BeNil())
			request.Header.Set("Content-Type", "application/json")
			handler := handlers.CreateRoomHandler(mockGameEngineService)
			router.POST("/rooms", handler)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 201", func() {
			createRoomRequest := models.CreateRoomRequest{
				Title:       "title",
				Description: "description",
			}
			requestBody, err := json.Marshal(createRoomRequest)
			Expect(err).To(BeNil())
			request, err := http.NewRequest("POST", "/rooms", bytes.NewBuffer(requestBody))
			Expect(err).To(BeNil())
			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			router.Use(insertPlayerIDInContextMiddleware(playerID))
			handler := handlers.CreateRoomHandler(mockGameEngineService)
			router.POST("/rooms", handler)
			mockGameEngineService.On("CreateRoom", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(uuid.Nil, nil)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusCreated))
		})

		It("should return 500 if server error occurs", func() {
			createRoomRequest := models.CreateRoomRequest{
				Title:       "title",
				Description: "description",
			}
			requestBody, err := json.Marshal(createRoomRequest)
			Expect(err).To(BeNil())
			request, err := http.NewRequest("POST", "/rooms", bytes.NewBuffer(requestBody))
			Expect(err).To(BeNil())
			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			router.Use(insertPlayerIDInContextMiddleware(playerID))
			handler := handlers.CreateRoomHandler(mockGameEngineService)
			router.POST("/rooms", handler)
			mockGameEngineService.On("CreateRoom", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, models.NewGenericError("server error"))
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusInternalServerError))
		})

	})

	Context("PlayerJoinRoomHandler", func() {
		It("should return 400 if roomId param is invalid", func() {
			invalidRoomID := "invalid-room-id"
			request, err := http.NewRequest("POST", fmt.Sprintf("/rooms/%s/player", invalidRoomID), nil)
			Expect(err).To(BeNil())
			request.Header.Set("Content-Type", "application/json")
			handler := handlers.PlayerJoinRoomHandler(mockGameEngineService)
			router.POST("/rooms/:roomId/player", handler)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 400 if playerID missing from context", func() {
			validRoomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			request, err := http.NewRequest("POST", fmt.Sprintf("/rooms/%s/player", validRoomID.String()), nil)
			Expect(err).To(BeNil())
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			handler := handlers.PlayerJoinRoomHandler(mockGameEngineService)
			router.POST("/rooms/:roomId/player", handler)
			router.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 200 if request is OK", func() {
			validRoomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			request, err := http.NewRequest("POST", fmt.Sprintf("/rooms/%s/player", validRoomID.String()), nil)
			Expect(err).To(BeNil())
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			handler := handlers.PlayerJoinRoomHandler(mockGameEngineService)
			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			router.Use(insertPlayerIDInContextMiddleware(playerID))
			router.POST("/rooms/:roomId/player", handler)
			mockGameEngineService.On("PlayerJoinRoom", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			router.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusOK))
		})

		It("should return 500 if server error occurs", func() {
			validRoomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			request, err := http.NewRequest("POST", fmt.Sprintf("/rooms/%s/player", validRoomID.String()), nil)
			Expect(err).To(BeNil())
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			router.Use(insertPlayerIDInContextMiddleware(playerID))
			handler := handlers.PlayerJoinRoomHandler(mockGameEngineService)
			router.POST("/rooms/:roomId/player", handler)
			mockGameEngineService.On("PlayerJoinRoom", mock.Anything, mock.Anything, mock.Anything).Return(models.NewGenericError("server error"))
			router.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	Context("PlayerLeaveRoomHandler", func() {
		It("should return 400 if roomId param is invalid", func() {
			invalidRoomID := "invalid-room-id"
			request, err := http.NewRequest("DELETE", fmt.Sprintf("/rooms/%s/player", invalidRoomID), nil)
			Expect(err).To(BeNil())
			request.Header.Set("Content-Type", "application/json")
			handler := handlers.PlayerLeaveRoomHandler(mockGameEngineService)
			router.DELETE("/rooms/:roomId/player", handler)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 400 if playerID missing from context", func() {
			validRoomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			request, err := http.NewRequest("DELETE", fmt.Sprintf("/rooms/%s/player", validRoomID), nil)
			Expect(err).To(BeNil())
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			handler := handlers.PlayerLeaveRoomHandler(mockGameEngineService)
			router.DELETE("/rooms/:roomId/player", handler)
			router.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 200 if request is OK", func() {
			validRoomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			request, err := http.NewRequest("DELETE", fmt.Sprintf("/rooms/%s/player", validRoomID), nil)
			Expect(err).To(BeNil())
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			router.Use(insertPlayerIDInContextMiddleware(playerID))
			handler := handlers.PlayerLeaveRoomHandler(mockGameEngineService)
			router.DELETE("/rooms/:roomId/player", handler)
			mockGameEngineService.On("PlayerLeaveRoom", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			router.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusOK))
		})

		It("should return 500 if server error occurs", func() {
			validRoomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			request, err := http.NewRequest("DELETE", fmt.Sprintf("/rooms/%s/player", validRoomID), nil)
			Expect(err).To(BeNil())
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			handler := handlers.PlayerLeaveRoomHandler(mockGameEngineService)
			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			router.Use(insertPlayerIDInContextMiddleware(playerID))
			router.DELETE("/rooms/:roomId/player", handler)
			mockGameEngineService.On("PlayerLeaveRoom", mock.Anything, mock.Anything, mock.Anything).Return(models.NewGenericError("server error"))
			router.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	Context("CreateGameHandler", func() {
		It("should return 400 if roomId param is invalid", func() {
			invalidRoomID := "invalid-room-id"
			request, err := http.NewRequest("POST", fmt.Sprintf("/rooms/%s/game", invalidRoomID), nil)
			Expect(err).To(BeNil())
			request.Header.Set("Content-Type", "application/json")
			handler := handlers.CreateGameHandler(mockGameEngineService)
			router.POST("/rooms/:roomId/game", handler)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			Expect(response.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 400 if playerID missing from context", func() {
			validRoomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			request, err := http.NewRequest("POST", fmt.Sprintf("/rooms/%s/game", validRoomID.String()), nil)
			Expect(err).To(BeNil())
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			handler := handlers.CreateGameHandler(mockGameEngineService)
			router.POST("/rooms/:roomId/game", handler)
			router.ServeHTTP(response, request)
			Expect(response.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 202 if no game is created and gameId is  set to default value", func() {
			validRoomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			request, err := http.NewRequest("POST", fmt.Sprintf("/rooms/%s/game", validRoomID.String()), nil)
			Expect(err).To(BeNil())
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			handler := handlers.CreateGameHandler(mockGameEngineService)
			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			router.Use(insertPlayerIDInContextMiddleware(playerID))
			router.POST("/rooms/:roomId/game", handler)
			mockGameEngineService.On("CreateGame", mock.Anything, mock.Anything, mock.Anything).Return(uuid.Nil, nil)
			router.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusAccepted))
		})

		It("should return 201 if game is created", func() {
			validRoomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			request, err := http.NewRequest("POST", fmt.Sprintf("/rooms/%s/game", validRoomID.String()), nil)
			Expect(err).To(BeNil())
			request.Header.Set("Content-Type", "application/json")
			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			router.Use(insertPlayerIDInContextMiddleware(playerID))
			handler := handlers.CreateGameHandler(mockGameEngineService)
			router.POST("/rooms/:roomId/game", handler)
			gameID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			mockGameEngineService.On("CreateGame", mock.Anything, mock.Anything, mock.Anything).Return(gameID, nil)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			Expect(response.Code).To(Equal(http.StatusCreated))
		})

		It("should return 500 if server error occurs", func() {
			validRoomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			request, err := http.NewRequest("POST", fmt.Sprintf("/rooms/%s/game", validRoomID.String()), nil)
			Expect(err).To(BeNil())
			request.Header.Set("Content-Type", "application/json")
			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			router.Use(insertPlayerIDInContextMiddleware(playerID))
			handler := handlers.CreateGameHandler(mockGameEngineService)
			router.POST("/rooms/:roomId/game", handler)
			mockGameEngineService.On("CreateGame", mock.Anything, mock.Anything, mock.Anything).Return(uuid.Nil, models.NewGenericError("server error"))
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	Context("GetGameStateHandler", func() {
		It("should return 400 if roomId param is invalid", func() {
			invalidRoomID := "invalid-room-id"
			request, err := http.NewRequest("GET", fmt.Sprintf("/rooms/%s/game", invalidRoomID), nil)
			Expect(err).To(BeNil())
			request.Header.Set("Content-Type", "application/json")
			handler := handlers.GetGameStateHandler(mockGameEngineService)
			router.GET("/rooms/:roomId/game", handler)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			Expect(response.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 400 if playerID missing from context", func() {
			validRoomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			request, err := http.NewRequest("GET", fmt.Sprintf("/rooms/%s/game", validRoomID.String()), nil)
			Expect(err).To(BeNil())
			request.Header.Set("Content-Type", "application/json")
			handler := handlers.GetGameStateHandler(mockGameEngineService)
			router.GET("/rooms/:roomId/game", handler)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			Expect(response.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 200 if request is OK", func() {
			validRoomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			request, err := http.NewRequest("GET", fmt.Sprintf("/rooms/%s/game", validRoomID.String()), nil)
			Expect(err).To(BeNil())
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			handler := handlers.GetGameStateHandler(mockGameEngineService)
			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			router.Use(insertPlayerIDInContextMiddleware(playerID))
			router.GET("/rooms/:roomId/game", handler)
			mockGameEngineService.On("GetGameState", mock.Anything, mock.Anything, mock.Anything).Return(&models.Game{}, nil)
			router.ServeHTTP(response, request)
			Expect(response.Code).To(Equal(http.StatusOK))
		})

		It("should return 500 if server error occurs", func() {
			validRoomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			request, err := http.NewRequest("GET", fmt.Sprintf("/rooms/%s/game", validRoomID.String()), nil)
			Expect(err).To(BeNil())
			request.Header.Set("Content-Type", "application/json")
			handler := handlers.GetGameStateHandler(mockGameEngineService)
			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			router.Use(insertPlayerIDInContextMiddleware(playerID))
			router.GET("/rooms/:roomId/game", handler)
			mockGameEngineService.On("GetGameState", mock.Anything, mock.Anything, mock.Anything).Return(nil, models.NewGenericError("server error"))
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			Expect(response.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	Context("MakeMoveHandler", func() {
		It("should return 400 if roomId param is invalid", func() {
			invalidRoomID := "invalid-room-id"
			invalidPosition := "invalid-position"
			request, err := http.NewRequest("POST", fmt.Sprintf("/test/%s/game/board/%s", invalidRoomID, invalidPosition), nil)
			Expect(err).To(BeNil())
			request.Header.Set("Content-Type", "application/json")
			handler := handlers.MakeMoveHandler(mockGameEngineService)
			router.POST("/test/:roomId/game/board/:position", handler)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			Expect(response.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 400 if position param is invalid", func() {
			validRoomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			invalidPosition := "invalid-position"
			request, err := http.NewRequest("POST", fmt.Sprintf("/test/%s/game/board/%s", validRoomID, invalidPosition), nil)
			Expect(err).To(BeNil())
			request.Header.Set("Content-Type", "application/json")
			handler := handlers.MakeMoveHandler(mockGameEngineService)
			router.POST("/test/:roomId/game/board/:position", handler)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			Expect(response.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 400 if playerID missing from context", func() {
			validRoomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			validPosition := 4
			request, err := http.NewRequest("POST", fmt.Sprintf("/test/%s/game/board/%d", validRoomID, validPosition), nil)
			Expect(err).To(BeNil())
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			handler := handlers.MakeMoveHandler(mockGameEngineService)
			router.POST("/test/:roomId/game/board/:position", handler)
			router.ServeHTTP(response, request)
			Expect(response.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 200 if request is OK", func() {
			validRoomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			validPosition := 4
			request, err := http.NewRequest("POST", fmt.Sprintf("/test/%s/game/board/%d", validRoomID, validPosition), nil)
			Expect(err).To(BeNil())
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			handler := handlers.MakeMoveHandler(mockGameEngineService)
			validPlayerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			router.Use(insertPlayerIDInContextMiddleware(validPlayerID))
			router.POST("/test/:roomId/game/board/:position", handler)
			mockGameEngineService.On("PlayerMakeMove", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			router.ServeHTTP(response, request)
			Expect(response.Code).To(Equal(http.StatusOK))
		})

		It("should return 500 if server error occurs", func() {
			validRoomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			validPosition := 4
			request, err := http.NewRequest("POST", fmt.Sprintf("/test/%s/game/board/%d", validRoomID, validPosition), nil)
			Expect(err).To(BeNil())
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			handler := handlers.MakeMoveHandler(mockGameEngineService)
			validPlayerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			router.Use(insertPlayerIDInContextMiddleware(validPlayerID))
			router.POST("/test/:roomId/game/board/:position", handler)
			mockGameEngineService.On("PlayerMakeMove", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(models.NewGenericError("server error"))
			router.ServeHTTP(response, request)
			Expect(response.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	Context("GetRankingHandler", func() {

		It("should return 200 if request is OK", func() {
			rankingRequest := models.RankingRequest{
				PageInfo: models.PageInfo{
					Page:     1,
					PageSize: 10,
				},
			}
			requestBody, err := json.Marshal(rankingRequest)
			Expect(err).To(BeNil())
			request, err := http.NewRequest("GET", "/ranking", bytes.NewBuffer(requestBody))
			Expect(err).To(BeNil())
			request.Header.Set("Content-Type", "application/json")
			handler := handlers.GetRankingHandler(mockGameEngineService)
			router.GET("/ranking", handler)
			mockGameEngineService.On("GetRanking", mock.Anything, mock.Anything, mock.Anything).Return([]*models.Player{}, 1, 1, 1, nil)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			Expect(response.Code).To(Equal(http.StatusOK))
		})

		It("should return 500 if server error occurs", func() {
			rankingRequest := models.RankingRequest{
				PageInfo: models.PageInfo{
					Page:     1,
					PageSize: 10,
				},
			}
			requestBody, err := json.Marshal(rankingRequest)
			Expect(err).To(BeNil())
			request, err := http.NewRequest("GET", "/ranking", bytes.NewBuffer(requestBody))
			Expect(err).To(BeNil())
			request.Header.Set("Content-Type", "application/json")
			handler := handlers.GetRankingHandler(mockGameEngineService)
			router.GET("/ranking", handler)
			mockGameEngineService.On("GetRanking", mock.Anything, mock.Anything, mock.Anything).Return(nil, 0, 0, 0, models.NewGenericError("server error"))
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			Expect(response.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})

func insertPlayerIDInContextMiddleware(id uuid.UUID) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(middleware.KEY_PLAYER_ID,
			uuid.NullUUID{UUID: id,
				Valid: true})
		c.Next()
	}
}
