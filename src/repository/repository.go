package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/plamen-v/tic-tac-toe-models/models"
)

const (
	DatabaseDriver            = "postgres"
	NoRecordsAffectedErrorMsg = "no records affected"
)

type Querier interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type GameRepository interface {
	Get(context.Context, uuid.UUID) (*models.Game, error)
	Create(context.Context, *models.Game) (uuid.UUID, error)
	Update(context.Context, *models.Game) error
}

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
			return nil, models.NewNotFoundErrorf("game '%s' not exist", id.String())
		} else {
			return nil, models.NewGenericError(err.Error())
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
		err = models.NewGenericError(err.Error())
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
		return models.NewGenericError(err.Error())
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return models.NewGenericError(NoRecordsAffectedErrorMsg)
	}

	return err
}

type PlayerRepository interface {
	Get(context.Context, uuid.UUID) (*models.Player, error)
	GetByLogin(context.Context, string) (*models.Player, error)
	UpdateStats(context.Context, *models.Player) error
	GetRanking(context.Context, int, int) ([]*models.Player, int, int, int, error)
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
		SELECT p.id, p.login, p.password, p.nickname, ps.wins, ps.losses, ps.draws
		FROM players AS p
		LEFT JOIN players_stats ps ON ps.player_id = p.id
		WHERE p.id = $1
		`

	row := r.db.QueryRowContext(ctx, sqlStr, id)
	player := &models.Player{}

	err := row.Scan(&player.ID, &player.Login, &player.Password, &player.Nickname, &player.Stats.Wins, &player.Stats.Losses, &player.Stats.Draws)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.NewNotFoundErrorf("player '%s' not exist", id.String())
		} else {
			return nil, models.NewGenericError(err.Error())
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
	row := r.db.QueryRowContext(ctx, sqlStr, login)

	err := row.Scan(&player.ID, &player.Login, &player.Password, &player.Nickname, &player.Stats.Wins, &player.Stats.Losses, &player.Stats.Draws)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.NewNotFoundErrorf("player '%s' not exist", login)
		} else {
			return nil, models.NewGenericError(err.Error())
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
		return models.NewGenericError(err.Error())
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return models.NewGenericError(NoRecordsAffectedErrorMsg)
	}

	return err
}

func (r *playerRepositoryImpl) GetRanking(ctx context.Context, page int, pageSize int) ([]*models.Player, int, int, int, error) {
	sqlStr := `
		SELECT COUNT(*)
		FROM players AS p
		LEFT JOIN players_stats ps ON ps.player_id = p.id
		ORDER BY ps.wins DESC, ps.draws DESC, ps.losses ASC
		`

	totalCnt := 0
	row := r.db.QueryRowContext(ctx, sqlStr)
	err := row.Scan(&totalCnt)
	if err != nil {
		return nil, 0, 0, 0, models.NewGenericError(err.Error())
	}

	lastPage := totalCnt/pageSize + 1
	if page > lastPage {
		page = lastPage
	}

	limit := pageSize
	offset := (page - 1) * pageSize

	sqlStr = `
		SELECT p.id, p.nickname, ps.wins, ps.losses, ps.draws
		FROM players AS p
		LEFT JOIN players_stats ps ON ps.player_id = p.id
		ORDER BY ps.wins DESC, ps.draws DESC, ps.losses ASC
		LIMIT $1 OFFSET $2
		`

	rows, err := r.db.QueryContext(ctx, sqlStr, limit, offset)
	if err != nil {
		return nil, 0, 0, 0, models.NewGenericError(err.Error())
	}
	defer rows.Close()

	players := make([]*models.Player, 0)
	for rows.Next() {
		player := &models.Player{}
		err := rows.Scan(&player.ID, &player.Nickname, &player.Stats.Wins, &player.Stats.Losses, &player.Stats.Draws)
		if err != nil {
			return nil, 0, 0, 0, models.NewGenericError(err.Error())
		}
		players = append(players, player)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, 0, 0, models.NewGenericError(err.Error())
	}

	return players, pageSize, page, totalCnt, nil
}

type RoomRepository interface {
	Get(context.Context, uuid.UUID, bool) (*models.Room, error)
	GetByPlayerID(context.Context, uuid.UUID) (*models.Room, error)
	GetList(context.Context, models.RoomPhase) ([]*models.Room, error)
	Create(context.Context, *models.Room) (uuid.UUID, error)
	Update(context.Context, *models.Room) error
	Delete(context.Context, uuid.UUID) error
}

func NewRoomRepository(db Querier) RoomRepository {
	return &roomRepositoryImpl{
		db: db,
	}
}

type roomRepositoryImpl struct {
	db Querier
}

func (r *roomRepositoryImpl) Get(ctx context.Context, id uuid.UUID, lock bool) (*models.Room, error) {
	lockCmd := ""
	if lock {
		lockCmd = "FOR UPDATE OF r"
	}
	sqlStr := fmt.Sprintf(`
		SELECT 
			r.id,
			ph.id AS host_id, 
			ph.nickname AS host_nickname,
			r.host_request_new_game, 
			pg.id AS guest_id, 
			pg.nickname AS guest_nickname,
			r.guest_request_new_game,
			r.game_id, 
			r.title, 
			r.description, 
			r.phase
		FROM rooms AS r
		INNER JOIN players AS ph ON ph.id = r.host_id
		LEFT JOIN players AS pg ON pg.id = r.guest_id
		WHERE r.id = $1
		%s;`, lockCmd)

	row := r.db.QueryRowContext(ctx, sqlStr, id)

	var (
		sqlGuestID             uuid.NullUUID
		sqlGuestNickname       sql.NullString
		sqlGuestRequestNewGame sql.NullBool
		sqlGameID              uuid.NullUUID
		sqlDescription         sql.NullString
	)

	room := &models.Room{}
	err := row.Scan(
		&room.ID,
		&room.Host.ID,
		&room.Host.Nickname,
		&room.Host.RequestNewGame,
		&sqlGuestID,
		&sqlGuestNickname,
		&sqlGuestRequestNewGame,
		&sqlGameID,
		&room.Title,
		&sqlDescription,
		&room.Phase)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.NewNotFoundErrorf("room '%s' not exist", id.String())
		} else {
			return nil, models.NewGenericError(err.Error())
		}
	}

	if sqlGuestID.Valid {
		room.Guest = &models.RoomPlayer{ID: sqlGuestID.UUID}

		if sqlGuestNickname.Valid {
			room.Guest.Nickname = sqlGuestNickname.String
		}

		if sqlGuestRequestNewGame.Valid {
			room.Guest.RequestNewGame = sqlGuestRequestNewGame.Bool
		}
	}

	if sqlGameID.Valid {
		room.GameID = &sqlGameID.UUID
	}

	if sqlDescription.Valid {
		room.Description = sqlDescription.String
	}

	return room, nil
}

func (r *roomRepositoryImpl) GetByPlayerID(ctx context.Context, playerID uuid.UUID) (*models.Room, error) {
	sqlStr := `
		SELECT 
			ph.id AS host_id, 
			ph.nickname AS host_nickname,
			r.host_request_new_game, 
			pg.id AS guest_id, 
			pg.nickname AS guest_nickname,
			r.guest_request_new_game,
			r.game_id, 
			r.title, 
			r.description, 
			r.phase
		FROM rooms AS r
		INNER JOIN players AS ph ON ph.id = r.host_id
		LEFT JOIN players AS pg ON pg.id = r.guest_id
		WHERE (r.host_id = $1) OR (r.guest_id IS NOT NULL AND r.guest_id = $1)
		`
	row := r.db.QueryRowContext(ctx, sqlStr, playerID)

	var (
		sqlGuestID             uuid.NullUUID
		sqlGuestNickname       sql.NullString
		sqlGuestRequestNewGame sql.NullBool
		sqlGameID              uuid.NullUUID
		sqlDescription         sql.NullString
	)

	room := &models.Room{}
	err := row.Scan(
		&room.Host.ID,
		&room.Host.Nickname,
		&room.Host.RequestNewGame,
		&sqlGuestID,
		&sqlGuestNickname,
		&sqlGuestRequestNewGame,
		&sqlGameID,
		&room.Title,
		&sqlDescription,
		&room.Phase)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.NewNotFoundErrorf("plaier '%s' not participate in a room", playerID.String())
		} else {
			return nil, models.NewGenericError(err.Error())
		}
	}

	if sqlGuestID.Valid {
		room.Guest = &models.RoomPlayer{}

		if sqlGuestNickname.Valid {
			room.Guest.Nickname = sqlGuestNickname.String
		}

		if sqlGuestRequestNewGame.Valid {
			room.Guest.RequestNewGame = sqlGuestRequestNewGame.Bool
		}
	}

	if sqlGameID.Valid {
		room.GameID = &sqlGameID.UUID
	}

	if sqlDescription.Valid {
		room.Description = sqlDescription.String
	}

	return room, nil
}

func (r *roomRepositoryImpl) GetList(ctx context.Context, phase models.RoomPhase) ([]*models.Room, error) {
	sqlStr := `
		SELECT 
			ph.id AS host_id, 
			ph.nickname AS host_nickname,
			r.title, 
			r.description, 
			r.phase
		FROM rooms AS r
		INNER JOIN players AS ph ON ph.id = r.host_id
		WHERE (r.phase = $1)
		`
	rows, err := r.db.QueryContext(ctx, sqlStr, phase)
	if err != nil {
		return nil, models.NewGenericError(err.Error())
	}
	defer rows.Close()

	rooms := make([]*models.Room, 0)
	var sqlDescription sql.NullString
	for rows.Next() {
		room := &models.Room{}
		err := rows.Scan(&room.Host.ID, &room.Host.Nickname, &room.Title, &sqlDescription, &room.Phase)
		if err != nil {
			return nil, models.NewGenericError(err.Error())
		}

		if sqlDescription.Valid {
			room.Description = sqlDescription.String
		}

		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return nil, models.NewGenericError(err.Error())
	}

	return rooms, nil
}

func (r *roomRepositoryImpl) Create(ctx context.Context, room *models.Room) (uuid.UUID, error) {
	sqlStr := `
		INSERT INTO rooms(host_id, host_request_new_game, title, description, phase)
		VALUES($1, $2, $3, $4, $5)
		RETURNING id
		`
	var id uuid.UUID
	err := r.db.QueryRowContext(ctx, sqlStr, room.Host.ID, room.Host.RequestNewGame, room.Title, room.Description, room.Phase).Scan(&id)
	if err != nil {
		err = models.NewGenericError(err.Error())
	}

	return id, err
}

func (r *roomRepositoryImpl) Update(ctx context.Context, room *models.Room) error {
	sqlStr := `
		UPDATE rooms
		SET host_id       		   = $2,
			host_request_new_game  = $3,
			guest_id       		   = $4,
			guest_request_new_game = $5,
			game_id         	   = $6,
			phase 		           = $7 
		WHERE id     	           = $1
		`
	var (
		sqlGuestID             uuid.NullUUID
		sqlGuestRequestNewGame bool = false
	)

	if room.Guest != nil {
		sqlGuestID.UUID = room.Guest.ID
		sqlGuestID.Valid = true
		sqlGuestRequestNewGame = room.Guest.RequestNewGame
	}

	result, err := r.db.ExecContext(ctx, sqlStr, room.ID, room.Host.ID, room.Host.RequestNewGame, sqlGuestID, sqlGuestRequestNewGame, room.GameID, room.Phase)
	if err != nil {
		return models.NewGenericError(err.Error())
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return models.NewGenericError(NoRecordsAffectedErrorMsg)
	}

	return err
}

func (r *roomRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	sqlStr := `DELETE FROM rooms WHERE id = $1`

	result, err := r.db.ExecContext(ctx, sqlStr, id)
	if err != nil {
		return models.NewGenericError(err.Error())
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return models.NewGenericError(NoRecordsAffectedErrorMsg)
	}

	return err
}
