package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	handlerv1 "github.com/plamen-v/tic-tac-toe/src/app/server/handlers/v1"
	"github.com/plamen-v/tic-tac-toe/src/app/server/middleware"
	"github.com/plamen-v/tic-tac-toe/src/config"
	"github.com/plamen-v/tic-tac-toe/src/services/auth"
	"github.com/plamen-v/tic-tac-toe/src/services/game"
	"github.com/plamen-v/tic-tac-toe/src/services/repository"
)

type APIServer interface {
	Start() error
	Stop(context.Context) error
}

type apiServer struct {
	config *config.AppConfiguration
	//ginEngine   *gin.Engine
	server     *http.Server
	auth       auth.AuthenticationService
	game       game.GameEngine
	playerRepo repository.PlayerRepository
	roomRepo   repository.RoomRepository
}

func NewAPI(config *config.AppConfiguration, authService auth.AuthenticationService, playerRepo repository.PlayerRepository, roomRepo repository.RoomRepository) APIServer {
	return &apiServer{
		config:     config,
		auth:       authService,
		playerRepo: playerRepo,
		roomRepo:   roomRepo,
	}
}

func (s *apiServer) Start() error {
	s.initialize()
	return s.server.ListenAndServe()
}

func (s *apiServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *apiServer) initialize() {
	engine := gin.Default()
	s.setEndpoints(engine)

	address := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	s.server = &http.Server{
		Addr:    address,
		Handler: engine.Handler(),
	}
}

func (s *apiServer) setEndpoints(engine *gin.Engine) {
	api := engine.Group("/api")

	v1 := api.Group("/v1")

	v1.POST("/login", handlerv1.LoginHandler(s.auth))

	auth := v1.Group("/")
	auth.Use(middleware.AuthenticationMiddleware(s.auth))

	//List open rooms
	auth.GET("/rooms", nil)
	//Create room
	auth.POST("/rooms", handlerv1.CreateRoomHandler(s.game))
	//Get room state
	auth.GET("rooms/:id", nil)
	//Delete room
	auth.DELETE("/rooms/:id", nil)

	//Setting the room host as ready to continue the game
	auth.POST("rooms/:id/host/ready", nil)
	//Delete room
	auth.DELETE("rooms/:id/host/ready", nil)

	auth.POST("rooms/:id/guest", nil)
	auth.POST("rooms/:id/guest/ready", nil)
	auth.DELETE("rooms/:id/guest", nil)

	auth.GET("rooms/:id/games/:id", nil)
	auth.PUT("rooms/:id/games/:id/board", nil)
}
