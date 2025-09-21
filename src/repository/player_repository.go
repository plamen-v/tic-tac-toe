package repository

//todo! error msg
import (
	"context"
	"database/sql"

	"github.com/plamen-v/tic-tac-toe-models/models"
)

type PlayerRepository interface {
	Get(context.Context, int64) (*models.Player, error)
	GetByLogin(context.Context, string) (*models.Player, error)
	UpdateStats(context.Context, *models.Player) error
}

func NewPlayerRepository(db Querier) PlayerRepository {
	return &playerRepositoryImpl{
		db: db,
	}
}

type playerRepositoryImpl struct {
	db Querier
}

func (r *playerRepositoryImpl) Get(ctx context.Context, id int64) (*models.Player, error) {
	sqlStr := `
		SELECT p.id, p.login, p.password, p.nickname, ps.wins, ps.losses, ps.draws
		FROM players AS p
		LEFT JOIN players_stats ps ON ps.player_id = p.id
		WHERE p.id = $1
		`

	var row *sql.Row
	row = r.db.QueryRowContext(ctx, sqlStr, id)
	player := &models.Player{}

	err := row.Scan(&player.ID, &player.Login, &player.Password, &player.Nickname, &player.Stats.Wins, &player.Stats.Losses, &player.Stats.Draws)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.NewNotFoundError("player not exist")
		} else {
			return nil, models.NewGenericErrorWithCause("player not selected", err) //todo!
		}
	}

	return player, nil
}

func (r *playerRepositoryImpl) GetByLogin(ctx context.Context, login string) (*models.Player, error) {
	sqlStr := `
		SELECT p.id, p.login, p.password, p.nickname, ps.wins, ps.losses, ps.draws
		FROM players AS p
		LEFT JOIN players_stats ps ON ps.player_id = p.id
		WHERE p.login = $1
		`
	player := &models.Player{}
	var row *sql.Row

	row = r.db.QueryRowContext(ctx, sqlStr, login)

	err := row.Scan(&player.ID, &player.Login, &player.Password, &player.Nickname, &player.Stats.Wins, &player.Stats.Losses, &player.Stats.Draws)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.NewNotFoundError("player not exist")
		} else {
			return nil, models.NewGenericErrorWithCause("player not selected", err) //todo!
		}
	}

	return player, nil
}

func (r *playerRepositoryImpl) UpdateStats(ctx context.Context, player *models.Player) error {
	sqlStr := `
		UPDATE players_stats
		SET wins    = $2,
    		losses  = $3,
    		draws   = $4
		WHERE player_id    = $1`

	result, err := r.db.ExecContext(ctx, sqlStr, player.ID, player.Stats.Wins, player.Stats.Losses, player.Stats.Draws)

	if err != nil {
		return models.NewGenericErrorWithCause("player not updated", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.NewGenericErrorWithCause("player not updated", err)
	}
	if rowsAffected != 1 {
		return models.NewGenericErrorWithCause("player not updated", err)
	}

	return err
}
