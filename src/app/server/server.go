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

type apiServer struct {
	config *config.AppConfiguration
	logger logger.LoggerService
	//ginEngine   *gin.Engine
	server                *http.Server
	authenticationService auth.AuthenticationService
	gameEngineService     engine.GameEngineService
}

func NewAPI(config *config.AppConfiguration, logger logger.LoggerService, authenticationService auth.AuthenticationService, gameEngineService engine.GameEngineService) APIServer {
	return &apiServer{
		config:                config,
		logger:                logger,
		authenticationService: authenticationService,
		gameEngineService:     gameEngineService,
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
	setServerMode(s.config.AppMode)
	engine := gin.Default()
	s.setEndpoints(engine)

	address := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	s.server = &http.Server{
		Addr:    address,
		Handler: engine.Handler(),
	}
}

func (s *apiServer) setEndpoints(engine *gin.Engine) {
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

	game.GET("/rooms", nil)
	game.POST("/rooms", handlers.CreateNewRoomHandler(s.gameEngineService))
	game.GET("rooms/:id", nil)
	game.DELETE("/rooms/:id", nil)

	game.POST("rooms/:id/host/ready", nil)
	game.DELETE("rooms/:id/host/ready", nil)

	game.POST("rooms/:id/guest", nil)
	game.POST("rooms/:id/guest/ready", nil)
	game.DELETE("rooms/:id/guest", nil)

	game.GET("rooms/:id/games/:id", nil)
	game.PUT("rooms/:id/games/:id/board", nil)
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
