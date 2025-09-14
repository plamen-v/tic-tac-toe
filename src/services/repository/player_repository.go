package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/plamen-v/tic-tac-toe/src/models"
)

type PlayerRepository interface {
	GetByID(int64, *sql.Tx) (*models.Player, error)
	GetByLogin(string, *sql.Tx) (*models.Player, error)
	Update(*models.Player, *sql.Tx) error
}

func NewPlayerRepository(db *sql.DB) PlayerRepository {
	return &playerRepository{
		db: db,
	}
}

type playerRepository struct {
	db *sql.DB
}

func (r *playerRepository) GetByID(id int64, tx *sql.Tx) (*models.Player, error) {
	var roomID sql.NullInt64
	var gameID sql.NullInt64
	player := &models.Player{}

	forUpdate := ""
	if tx != nil {
		forUpdate = "FOR UPDATE"
	}
	sqlStr := fmt.Sprintf(`
		SELECT p.id, p.login, p.password, p.nickname, p.room_id, p.game_id, p.wins, p.losses, p.draws
		FROM players AS p
		WHERE p.id = $1
		%s`, forUpdate)

	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow(sqlStr, id)
	} else {
		row = r.db.QueryRow(sqlStr, id)
	}

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

func (r *playerRepository) GetByLogin(login string, tx *sql.Tx) (*models.Player, error) {
	var roomID sql.NullInt64
	var gameID sql.NullInt64
	player := &models.Player{}

	forUpdate := ""
	if tx != nil {
		forUpdate = "FOR UPDATE"
	}
	sqlStr := fmt.Sprintf(`
		SELECT p.id, p.login, p.password, p.nickname, p.room_id, p.game_id, p.wins, p.losses, p.draws
		FROM players AS p
		WHERE p.login = $1
		%s`, forUpdate)

	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow(sqlStr, login)
	} else {
		row = r.db.QueryRow(sqlStr, login)
	}

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

func (r *playerRepository) Update(player *models.Player, tx *sql.Tx) error {
	sqlStr := `
		UPDATE players
		SET nickname = $1,
    		room_id = $2,
    		game_id = $3,
    		wins = $4,
    		losses = $5,
    		draws = $6
		WHERE id = $7`

	var (
		result sql.Result
		err    error
	)
	if tx != nil {
		result, err = tx.Exec(sqlStr, player.Nickname, player.RoomID, player.GameID, player.Wins, player.Losses, player.Draws, player.ID)
	} else {
		result, err = r.db.Exec(sqlStr, player.Nickname, player.RoomID, player.GameID, player.Wins, player.Losses, player.Draws, player.ID)
	}

	_, err = result.RowsAffected()
	//TODO! check affected rows count

	return err

}
