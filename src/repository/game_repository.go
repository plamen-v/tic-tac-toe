package repository

import (
	"context"
	"database/sql"

	"github.com/plamen-v/tic-tac-toe-models/models"
)

type GameRepository interface {
	Get(context.Context, int64) (*models.Game, error)
	Create(context.Context, *models.Game) (int64, error)
	Update(context.Context, *models.Game) error
}

// todo! error msg
func NewGameRepository(db Querier) GameRepository {
	return &gameRepositoryImpl{
		db: db,
	}
}

type gameRepositoryImpl struct {
	db Querier
}

func (r *gameRepositoryImpl) Get(ctx context.Context, id int64) (*models.Game, error) {

	sqlStr := `
		SELECT 
			g.id, 
			g.host_id, 
			g.host_mark, 
			g.guest_id, 
			g.guest_mark, 
			g.current_player_id, 
			g.board, 
			g.winner_id, 
			g.phase			
		FROM games AS g
		WHERE g.id = $1`

	var row *sql.Row
	row = r.db.QueryRowContext(ctx, sqlStr, id)

	var winnerID sql.NullInt64
	game := &models.Game{}
	err := row.Scan(
		&game.ID, &game.HostID, &game.HostMark,
		&game.GuestID, &game.GuestMark, &game.CurrentPlayerID,
		&game.Board, &winnerID, &game.Phase)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.NewNotFoundErrorf("game with id equal to %d not exist", id)
		} else {
			return nil, models.NewGenericErrorWithCause("game not selected", err)
		}
	}

	if winnerID.Valid {
		game.WinnerID = &winnerID.Int64
	}

	return game, nil
}

func (r *gameRepositoryImpl) Create(ctx context.Context, game *models.Game) (int64, error) {
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

	var id int64
	err := r.db.QueryRowContext(ctx, sqlStr, game.HostID, game.HostMark, game.GuestID, game.GuestMark, game.CurrentPlayerID, game.Phase).Scan(&id)

	if err != nil {
		err = models.NewGenericErrorWithCause("game insert failed", err)
	}

	return id, err
}

func (r *gameRepositoryImpl) Update(ctx context.Context, game *models.Game) error {
	sqlStr := `
		UPDATE games
		SET current_player_id = $2,
			board             = $3,
			phase 			  = $4, 
			winner_id 		  = $5
		WHERE id     		  = $1`

	result, err := r.db.ExecContext(ctx, sqlStr, game.ID, game.CurrentPlayerID, game.Board, game.Phase, game.WinnerID)

	if err != nil {
		return models.NewGenericErrorWithCause("game update failed", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.NewGenericErrorWithCause("could not get rows affected", err)
	}
	if rowsAffected == 0 {
		return models.NewGenericErrorWithCause("no game was deleted", err)
	}

	return err
}
