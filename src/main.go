package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/plamen-v/tic-tac-toe/src/app"
	"github.com/plamen-v/tic-tac-toe/src/config"
	"github.com/plamen-v/tic-tac-toe/src/services/auth"
	"github.com/plamen-v/tic-tac-toe/src/services/engine"
	"github.com/plamen-v/tic-tac-toe/src/services/logger"
	"github.com/plamen-v/tic-tac-toe/src/services/repository"
)

func main() {
	flag.Parse()

	config, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	logger, err := logger.NewLoggerService(config.AppMode, config.LogLevel)
	if err != nil {
		panic(err)
	}
	defer func() {
		err = logger.Sync()
		if err != nil {
			log.Println(err)
		}
	}()

	dbConnection, err := repository.OpenDatabaseConnection(
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.Database)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := dbConnection.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	repo := repository.NewRepository(
		dbConnection,
		repository.NewPlayerRepository(dbConnection),
		repository.NewRoomRepository(dbConnection),
		repository.NewGameRepository(dbConnection),
	)
	app := app.NewApplication(
		config,
		logger,
		auth.NewAuthenticationService(config, repo),
		engine.NewGameEngineService(repo))

	go func() {
		if err = app.Start(); err != nil {
			panic(err)
		}
	}()

	//Gracefully shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	err = app.Stop(ctx)
	if err != nil {
		panic(err)
	}
}
