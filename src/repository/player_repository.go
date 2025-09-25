package repository

//todo! error msg
import (
	"context"
	"database/sql"

	"github.com/gofrs/uuid"
	"github.com/plamen-v/tic-tac-toe-models/models"
)

type PlayerRepository interface {
	Get(context.Context, uuid.UUID) (*models.Player, error)
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

func (r *playerRepositoryImpl) Get(ctx context.Context, id uuid.UUID) (*models.Player, error) {
	sqlStr := `
		SELECT p.id, p.login, p.password, p.nickname, ps.wins, ps.losses, ps.draws, COALESCE(hr.id, gr.id) AS room_id
		FROM players AS p
		LEFT JOIN players_stats ps ON ps.player_id = p.id
		LEFT JOIN rooms hr ON p.id = hr.host_id 
		LEFT JOIN rooms gr ON p.id = gr.guest_id
		WHERE p.id = $1
		`

	row := r.db.QueryRowContext(ctx, sqlStr, id)
	player := &models.Player{}

	var sqlRoomID uuid.NullUUID
	err := row.Scan(&player.ID, &player.Login, &player.Password, &player.Nickname, &player.Stats.Wins, &player.Stats.Losses, &player.Stats.Draws, &sqlRoomID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.NewNotFoundError("player not exist")
		} else {
			return nil, models.NewGenericErrorWithCause("record scan error", err)
		}
	}

	if sqlRoomID.Valid {
		player.RoomID = &sqlRoomID.UUID
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
	row := r.db.QueryRowContext(ctx, sqlStr, login)

	err := row.Scan(&player.ID, &player.Login, &player.Password, &player.Nickname, &player.Stats.Wins, &player.Stats.Losses, &player.Stats.Draws)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.NewNotFoundError("player not exist")
		} else {
			return nil, models.NewGenericErrorWithCause("record scan error", err)
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
		WHERE player_id = $1`

	result, err := r.db.ExecContext(ctx, sqlStr, player.ID, player.Stats.Wins, player.Stats.Losses, player.Stats.Draws)

	if err != nil {
		return models.NewGenericErrorWithCause("player not updated", err)
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return models.NewGenericError("no player was updated")
	}

	return err
}
