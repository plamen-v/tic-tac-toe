package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/plamen-v/tic-tac-toe/src/app"
	"github.com/plamen-v/tic-tac-toe/src/config"
	"github.com/plamen-v/tic-tac-toe/src/services/auth"
	"github.com/plamen-v/tic-tac-toe/src/services/game"
	"github.com/plamen-v/tic-tac-toe/src/services/repository"
)

const (
	DATABASE_DRIVER = "postgres"
)

func main() {
	flag.Parse()

	config := config.GetConfig()
	//todo! log init:

	source := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password, config.Database.Database)
	database, err := sql.Open(DATABASE_DRIVER, source)
	if err != nil {
		//TODO! log
	}
	//todo! single repo
	playerRepo := repository.NewPlayerRepository(database)
	roomRepo := repository.NewRoomRepository(database)
	gameRepo := repository.NewGameRepository(database)

	app := app.NewApplication(config,
		auth.NewAuthenticationService(config, playerRepo),
		game.NewGameEngine(playerRepo, roomRepo, gameRepo))

	go func() {
		if err = app.Start(); err != nil {
			//todo!panic(err)
		}
	}()

	//Gracefully shutdown with a timeout.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	//todo! use timeout
	if database.Close(); err != nil {
		//todo! log.Println("todo:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = app.Stop(ctx)
	if err != nil {
		//todo!
	}
}
