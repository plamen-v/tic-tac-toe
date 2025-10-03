package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/plamen-v/tic-tac-toe/src/app"
	"github.com/plamen-v/tic-tac-toe/src/config"
	"github.com/plamen-v/tic-tac-toe/src/repository"
	"github.com/plamen-v/tic-tac-toe/src/services/auth"
	"github.com/plamen-v/tic-tac-toe/src/services/engine"
	"github.com/plamen-v/tic-tac-toe/src/services/logger"
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

	source := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password, config.Database.Database)
	db, err := sql.Open("postgres", source)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	app := app.NewApplication(
		config,
		logger,
		auth.NewAuthenticationService(config, db),
		engine.NewGameEngineService(db,
			repository.NewPlayerRepository,
			repository.NewGameRepository,
			repository.NewRoomRepository,
		))

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
