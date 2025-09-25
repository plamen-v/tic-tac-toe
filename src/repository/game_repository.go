package repository

import (
	"context"
	"database/sql"

	"github.com/gofrs/uuid"
	"github.com/plamen-v/tic-tac-toe-models/models"
)

type GameRepository interface {
	Get(context.Context, uuid.UUID) (*models.Game, error)
	Create(context.Context, *models.Game) (uuid.UUID, error)
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

func (r *gameRepositoryImpl) Get(ctx context.Context, id uuid.UUID) (*models.Game, error) {

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

	row := r.db.QueryRowContext(ctx, sqlStr, id)

	var winnerID uuid.NullUUID
	game := &models.Game{}
	err := row.Scan(
		&game.ID, &game.Host.ID, &game.Host.Mark,
		&game.Guest.ID, &game.Guest.Mark, &game.CurrentPlayerID,
		&game.Board, &winnerID, &game.Phase)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.NewNotFoundErrorf("game with id equal to %d not exist", id)
		} else {
			return nil, models.NewGenericErrorWithCause("record scan error", err)
		}
	}

	if winnerID.Valid {
		game.WinnerID = &winnerID.UUID
	}

	return game, nil
}

func (r *gameRepositoryImpl) Create(ctx context.Context, game *models.Game) (uuid.UUID, error) {
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

	var id uuid.UUID
	err := r.db.QueryRowContext(ctx, sqlStr, game.Host.ID, game.Host.Mark, game.Guest.ID, game.Guest.Mark, game.CurrentPlayerID, game.Phase).Scan(&id)

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

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return models.NewGenericError("no game was updated")
	}

	return err
}
