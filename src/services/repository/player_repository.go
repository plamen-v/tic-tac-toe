package repository

import (
	"database/sql"
	"errors"

	"github.com/plamen-v/tic-tac-toe/src/models"
)

type PlayerRepository interface {
	GetByID(int) (*models.Player, error)
	GetByLogin(login string) (*models.Player, error)
	Update(*models.Player) (*models.Player, error)
}

func NewPlayerRepository(db *sql.DB) PlayerRepository {
	return &playerRepository{
		db: db,
	}
}

type playerRepository struct {
	db *sql.DB
}

func (r *playerRepository) GetByID(id int) (*models.Player, error) {
	var roomID sql.NullInt64
	var gameID sql.NullInt64
	player := &models.Player{}
	row := r.db.QueryRow(`
		SELECT p.id, p.login, p.password, p.nickname, p.room_id, p.game_id, p.wins, p.losses, p.draws
		FROM players AS p
		WHERE p.id = $1
	`, id)

	err := row.Scan(&player.ID, &player.Login, &player.Password, &player.Nickname, &roomID, &gameID, &player.Wins, &player.Losses, &player.Draws)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			//todo!fmt.Println("No user found with that ID")
		} else {
			//todo!fmt.Println("QueryRow error:", err)
		}
		return &models.Player{}, err
	}

	if roomID.Valid {
		player.RoomID = &roomID.Int64
	}

	if gameID.Valid {
		player.GameID = &gameID.Int64
	}

	return player, nil
}

func (r *playerRepository) GetByLogin(login string) (*models.Player, error) {
	var roomID sql.NullInt64
	var gameID sql.NullInt64
	player := &models.Player{}
	row := r.db.QueryRow(`
		SELECT p.id, p.login, p.password, p.nickname, p.room_id, p.game_id, p.wins, p.losses, p.draws
		FROM players AS p
		WHERE p.login = $1
	`, login)

	err := row.Scan(&player.ID, &player.Login, &player.Password, &player.Nickname, &player.RoomID, &player.GameID, &player.Wins, &player.Losses, &player.Draws)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			//todo!fmt.Println("No user found with that ID")
		} else {
			//todo!fmt.Println("QueryRow error:", err)
		}
		return &models.Player{}, err
	}

	if roomID.Valid {
		player.RoomID = &roomID.Int64
	}

	if gameID.Valid {
		player.GameID = &gameID.Int64
	}

	return player, nil
}

func (r *playerRepository) Update(player *models.Player) (*models.Player, error) {
	return &models.Player{}, nil
}
