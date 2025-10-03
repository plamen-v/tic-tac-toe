package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/plamen-v/tic-tac-toe/src/app/server/handlers"
	"github.com/plamen-v/tic-tac-toe/src/app/server/middleware"
	"github.com/plamen-v/tic-tac-toe/src/config"
	"github.com/plamen-v/tic-tac-toe/src/services/auth"
	"github.com/plamen-v/tic-tac-toe/src/services/engine"
	"github.com/plamen-v/tic-tac-toe/src/services/logger"
)

type APIServer interface {
	Start() error
	Stop(context.Context) error
}

type apiServerImpl struct {
	config *config.AppConfiguration
	logger logger.LoggerService
	//ginEngine   *gin.Engine
	server                *http.Server
	authenticationService auth.AuthenticationService
	gameEngineService     engine.GameEngineService
}

func NewAPI(config *config.AppConfiguration, logger logger.LoggerService, authenticationService auth.AuthenticationService, gameEngineService engine.GameEngineService) APIServer {
	return &apiServerImpl{
		config:                config,
		logger:                logger,
		authenticationService: authenticationService,
		gameEngineService:     gameEngineService,
	}
}

func (s *apiServerImpl) Start() error {
	s.initialize()
	return s.server.ListenAndServe()
}

func (s *apiServerImpl) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *apiServerImpl) initialize() {
	setServerMode(s.config.AppMode)
	engine := gin.Default()
	s.setEndpoints(engine)

	address := fmt.Sprintf(":%d", s.config.Server.Port)
	s.server = &http.Server{
		Addr:    address,
		Handler: engine.Handler(),
	}
}

func (s *apiServerImpl) setEndpoints(engine *gin.Engine) {
	engine.Use(
		middleware.Logger(s.logger),
		middleware.ErrorHandler(),
		gin.RecoveryWithWriter(gin.DefaultErrorWriter,
			func(c *gin.Context, err any) {
				c.Error(errors.New("panic"))
				c.Abort()
			},
		),
	)

	api := engine.Group("/api")
	api.POST("/login", handlers.LoginHandler(s.authenticationService))

	game := api.Group("/")
	game.Use(middleware.Authentication(s.authenticationService))

	game.GET("/room", handlers.GetRoomHandler(s.gameEngineService))
	game.GET("/rooms", handlers.GetOpenRoomsHandler(s.gameEngineService))
	game.POST("/rooms", handlers.CreateRoomHandler(s.gameEngineService))
	game.POST("rooms/:roomId/player", handlers.PlayerJoinRoomHandler(s.gameEngineService))
	game.DELETE("rooms/:roomId/player", handlers.PlayerLeaveRoomHandler(s.gameEngineService))
	game.POST("rooms/:roomId/game", handlers.CreateGameHandler(s.gameEngineService))
	game.GET("rooms/:roomId/game/", handlers.GetGameStateHandler(s.gameEngineService))
	game.POST("rooms/:roomId/game/board/:position", handlers.MakeMoveHandler(s.gameEngineService))
	game.GET("ranking", handlers.GetRankingHandler(s.gameEngineService))
}

func setServerMode(mode config.AppMode) {
	switch mode {
	case config.ProductionAppMode:
		gin.SetMode(gin.ReleaseMode)
	case config.DevelopmentAppMode:
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.TestMode)
		gin.DefaultErrorWriter = io.Discard
	}
}
