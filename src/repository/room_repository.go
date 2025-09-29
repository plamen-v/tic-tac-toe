package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/plamen-v/tic-tac-toe-models/models"
)

// todo! error msg
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
			return nil, models.NewNotFoundError("room not exist")
		} else {
			return nil, models.NewGenericErrorWithCause("record scan error", err)
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
			return nil, models.NewNotFoundError("room not exist")
		} else {
			return nil, models.NewGenericErrorWithCause("record scan error", err)
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
		return nil, models.NewGenericErrorWithCause("query failed ", err)
	}
	defer rows.Close()

	rooms := make([]*models.Room, 0)
	var sqlDescription sql.NullString
	for rows.Next() {
		room := &models.Room{}
		err := rows.Scan(&room.Host.ID, &room.Host.Nickname, &room.Title, &sqlDescription, &room.Phase)
		if err != nil {
			return nil, models.NewGenericErrorWithCause("room scan failed", err)
		}

		if sqlDescription.Valid {
			room.Description = sqlDescription.String
		}

		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return nil, models.NewGenericErrorWithCause("rooms iteration error", err)
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
		err = models.NewGenericErrorWithCause("room insert failed", err)
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
		return models.NewGenericErrorWithCause("room update failed", err)
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return models.NewGenericError("no room was updated")
	}

	return err
}

func (r *roomRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	sqlStr := `DELETE FROM rooms WHERE id = $1`

	result, err := r.db.ExecContext(ctx, sqlStr, id)
	if err != nil {
		return models.NewGenericErrorWithCause("room deletion failed", err)
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return models.NewGenericError("no room was deleted")
	}

	return err
}
