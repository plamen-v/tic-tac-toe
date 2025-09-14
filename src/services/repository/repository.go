package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	DatabaseDriver = "postgres"
)

type Repository interface {
	Begin() (*sql.Tx, error)
	Players() PlayerRepository
	Rooms() RoomRepository
	Games() GameRepository
}

func OpenDatabaseConnection(host string, port int, user string, password string, database string) (*sql.DB, error) {
	source := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, database)
	db, err := sql.Open(DatabaseDriver, source)
	return db, err
}

func NewRepository(db *sql.DB, playerRepository PlayerRepository,
	roomRepository RoomRepository,
	gameRepository GameRepository) Repository {
	return &repository{
		db:               db,
		playerRepository: playerRepository,
		roomRepository:   roomRepository,
		gameRepository:   gameRepository,
	}
}

type repository struct {
	db               *sql.DB
	playerRepository PlayerRepository
	roomRepository   RoomRepository
	gameRepository   GameRepository
}

func (r *repository) Begin() (*sql.Tx, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err //TODO!
	}
	return tx, nil
}

func (r *repository) Players() PlayerRepository {
	return r.playerRepository
}

func (r *repository) Rooms() RoomRepository {
	return r.roomRepository
}

func (r *repository) Games() GameRepository {
	return r.gameRepository
}
