package repository

//todo! error msg
import (
	"database/sql"
	"fmt"

	"github.com/plamen-v/tic-tac-toe-models/models"
	"github.com/plamen-v/tic-tac-toe-models/models/errors"
)

type PlayerRepository interface {
	Get(int64, *sql.Tx) (*models.Player, error)
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

func (r *playerRepository) Get(id int64, tx *sql.Tx) (*models.Player, error) {
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
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundErrorf("player with id equal to %d not exist", id)
		} else {
			return nil, errors.NewGenericErrorWithCause("player not selected", err)
		}
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
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundErrorf("player with id equal to %d not exist", player.ID)
		} else {
			return nil, errors.NewGenericErrorWithCause("player not selected", err)
		}
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
		SET room_id = $2,
    		game_id = $3,
    		wins    = $4,
    		losses  = $5,
    		draws   = $6
		WHERE id    = $1`

	var (
		result sql.Result
		err    error
	)
	if tx != nil {
		result, err = tx.Exec(sqlStr, player.ID, player.RoomID, player.GameID, player.Wins, player.Losses, player.Draws)
	} else {
		result, err = r.db.Exec(sqlStr, player.ID, player.RoomID, player.GameID, player.Wins, player.Losses, player.Draws)
	}

	if err != nil {
		return errors.NewGenericErrorWithCause("player not updated", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewGenericErrorWithCause("player not updated", err)
	}
	if rowsAffected != 1 {
		return errors.NewGenericErrorWithCause("player not updated", err)
	}

	return err
}
