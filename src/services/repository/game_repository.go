package repository

import (
	"database/sql"
	"fmt"

	"github.com/plamen-v/tic-tac-toe-models/models"
	"github.com/plamen-v/tic-tac-toe-models/models/errors"
)

type GameRepository interface {
	Get(int64, *sql.Tx) (*models.Game, error)
	Create(*models.Game, *sql.Tx) (int64, error)
	Update(*models.Game, *sql.Tx) error
}

// todo! error msg
func NewGameRepository(db *sql.DB) GameRepository {
	return &gameRepository{
		db: db,
	}
}

type gameRepository struct {
	db *sql.DB
}

func (r *gameRepository) Get(id int64, tx *sql.Tx) (*models.Game, error) {
	var winnerID sql.NullInt64
	var loserID sql.NullInt64
	game := &models.Game{}

	forUpdate := ""
	if tx != nil {
		forUpdate = "FOR UPDATE"
	}
	sqlStr := fmt.Sprintf(`
		SELECT 
			g.id, 
			g.host_id, 
			g.host_mark, 
			g.guest_id, 
			g.guest_mark, 
			g.current_player_id, 
			g.board, 
			g.winner_id, 
			g.loser_id,
			g.phase			
		FROM games AS g
		WHERE g.id = $1
		%s`, forUpdate)

	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow(sqlStr, id)
	} else {
		row = r.db.QueryRow(sqlStr, id)
	}

	err := row.Scan(
		&game.ID, &game.HostID, &game.HostMark,
		&game.GuestID, &game.GuestMark, &game.CurrentPlayerID,
		&game.Board, &game.Phase, &winnerID, &loserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundErrorf("game with id equal to %d not exist", id)
		} else {
			return nil, errors.NewGenericErrorWithCause("game not selected", err)
		}
	}

	if winnerID.Valid {
		game.WinnerID = &winnerID.Int64
	}

	if loserID.Valid {
		game.LoserID = &loserID.Int64
	}

	return game, nil
}

func (r *gameRepository) Create(game *models.Game, tx *sql.Tx) (int64, error) {
	sqlStr := `
		INSERT INTO games(
			host_id, 
			host_mark, 
			guest_id, 
			guest_mark, 
			current_player_id, 
			phase)
		VALUES($1, $2, $3, $4, $5, $6)
		RETURNING id`

	var (
		err error
		id  int64
	)

	if tx != nil {
		err = tx.QueryRow(sqlStr, game.HostID, game.HostMark, game.GuestID, game.GuestMark, game.CurrentPlayerID, game.Phase).Scan(&id)
	} else {
		err = tx.QueryRow(sqlStr, game.HostID, game.HostMark, game.GuestID, game.GuestMark, game.CurrentPlayerID, game.Phase).Scan(&id)
	}

	if err != nil {
		err = errors.NewGenericErrorWithCause("game not created", err)
	}

	return id, err
}

func (r *gameRepository) Update(game *models.Game, tx *sql.Tx) error {
	sqlStr := `
		UPDATE games
		SET current_player_id = $2,
			board             = $3,
			status 			  = $4, 
			winner_id 		  = $5, 
			loser_id 		  = $6
		WHERE id     		  = $1`

	var (
		result sql.Result
		err    error
	)
	if tx != nil {
		result, err = tx.Exec(sqlStr, game.ID, game.Board, game.Phase, game.WinnerID, game.LoserID)
	} else {
		result, err = r.db.Exec(sqlStr, game.ID, game.Board, game.Phase, game.WinnerID, game.LoserID)
	}

	if err != nil {
		return errors.NewGenericErrorWithCause("game not updated", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewGenericErrorWithCause("game not updated", err)
	}
	if rowsAffected != 1 {
		return errors.NewGenericErrorWithCause("game not updated", err)
	}

	return err
}
