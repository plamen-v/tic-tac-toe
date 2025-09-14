package app

//TODO! singletone
import (
	"context"

	_ "github.com/lib/pq"
	"github.com/plamen-v/tic-tac-toe/src/app/server"
	"github.com/plamen-v/tic-tac-toe/src/config"
	"github.com/plamen-v/tic-tac-toe/src/services/auth"
	"github.com/plamen-v/tic-tac-toe/src/services/engine"
	"github.com/plamen-v/tic-tac-toe/src/services/logger"
)

type Application interface {
	Start() error
	Stop(context.Context) error
}

type application struct {
	config                *config.AppConfiguration
	logger                logger.LoggerService
	server                server.APIServer
	authenticationService auth.AuthenticationService
	gameEngineService     engine.GameEngineService
}

// todo! singletone
func NewApplication(
	configuration *config.AppConfiguration,
	logger logger.LoggerService,
	authenticationService auth.AuthenticationService,
	gameEngineService engine.GameEngineService) Application {
	return &application{
		config:                configuration,
		logger:                logger,
		authenticationService: authenticationService,
		gameEngineService:     gameEngineService,
	}
}

func (a *application) Start() error {
	if err := a.initialize(); err != nil {
		return err
	}
	a.logger.Info("OK Bro TODO! HERE")
	return a.server.Start()
}

func (a *application) Stop(ctx context.Context) error {
	err := a.finalize(ctx)
	a.logger.Info("END Bro TODO! HERE")
	a.logger.Sync() //todo!
	return err
}

func (a *application) initialize() error {
	a.server = server.NewAPI(a.config, a.logger, a.authenticationService, a.gameEngineService)
	return nil
}

func (a *application) finalize(ctx context.Context) error {
	return a.server.Stop(ctx)
}
