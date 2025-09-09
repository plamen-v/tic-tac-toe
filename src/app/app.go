package app

//TODO! singletone
import (
	"context"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/plamen-v/tic-tac-toe/src/app/server"
	"github.com/plamen-v/tic-tac-toe/src/config"
	"github.com/plamen-v/tic-tac-toe/src/services/auth"
	"github.com/plamen-v/tic-tac-toe/src/services/game"
	"github.com/plamen-v/tic-tac-toe/src/services/repository"
)

type Application interface {
	Start() error
	Stop(context.Context) error
}

type application struct {
	database       *sql.DB
	config         *config.AppConfiguration
	server         server.APIServer
	authentication auth.AuthenticationService
	playerRepo     repository.PlayerRepository
	gameRepo       repository.GameRepository
	roomRepo       repository.RoomRepository
	gameEngine     game.GameEngine
}

func NewApplication(configuration *config.AppConfiguration,
	auth auth.AuthenticationService,
	gameEngine game.GameEngine) Application {
	return &application{
		config:         configuration,
		authentication: auth,
		gameEngine:     gameEngine,
	}
}

func (a *application) Start() error {
	if err := a.initialize(); err != nil {
		return err
	}

	if err := a.server.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}

	return nil
}

func (a *application) Stop(ctx context.Context) error {
	err := a.finalize(ctx)
	return err
}

func (a *application) initialize() error {
	a.server = server.NewAPI(a.config, a.authentication, a.playerRepo, a.roomRepo)
	return nil
}

func (a *application) finalize(ctx context.Context) error {
	if err := a.server.Stop(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}
	log.Println("Server exiting")

	return nil
}
