package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/plamen-v/tic-tac-toe-models/models"
)

// todo! error msg
type RoomRepository interface {
	Get(context.Context, int64, bool) (*models.Room, error)
	GetByPlayerID(context.Context, int64) (*models.Room, error)
	GetList(context.Context, string, string, string, models.RoomPhase) ([]*models.Room, error)
	Create(context.Context, *models.Room) (int64, error)
	Update(context.Context, *models.Room) error
	Delete(context.Context, int64) error
}

func NewRoomRepository(db Querier) RoomRepository {
	return &roomRepositoryImpl{
		db: db,
	}
}

type roomRepositoryImpl struct {
	db Querier
}

func (r *roomRepositoryImpl) Get(ctx context.Context, id int64, lock bool) (*models.Room, error) {
	lockCmd := ""
	if lock {
		lockCmd = "FOR UPDATE OF r"
	}
	sqlStr := fmt.Sprintf(`
		SELECT 
			r.id, 
			ph.id AS host_id, 
			ph.nickname AS host_nickname,
			r.host_continue, 
			pg.id AS guest_id, 
			pg.nickname AS guest_nickname,
			r.guest_continue,
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
		sqlGuestID       sql.NullInt64
		sqlGuestNickname sql.NullString
		sqlGuestContinue sql.NullBool
		sqlGameID        sql.NullInt64
		sqlDescription   sql.NullString
	)

	room := &models.Room{}
	err := row.Scan(
		&room.ID,
		&room.Host.ID,
		&room.Host.Nickname,
		&room.Host.Continue,
		&sqlGuestID,
		&sqlGuestNickname,
		&sqlGuestContinue,
		&sqlGameID,
		&room.Title,
		&sqlDescription,
		&room.Phase)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.NewNotFoundError("room not exist")
		} else {
			return nil, models.NewGenericErrorWithCause("room select failed", err)
		}
	}

	if sqlGuestID.Valid {
		room.Guest = &models.RoomParticipant{ID: sqlGuestID.Int64}

		if sqlGuestNickname.Valid {
			room.Guest.Nickname = sqlGuestNickname.String
		}

		if sqlGuestContinue.Valid {
			room.Guest.Continue = sqlGuestContinue.Bool
		}
	}

	if sqlGameID.Valid {
		room.GameID = &sqlGameID.Int64
	}

	if sqlDescription.Valid {
		room.Description = sqlDescription.String
	}

	return room, nil
}

func (r *roomRepositoryImpl) GetByPlayerID(ctx context.Context, playerID int64) (*models.Room, error) {
	sqlStr := `
		SELECT 
			r.id, 
			ph.id AS host_id, 
			ph.nickname AS host_nickname,
			r.host_continue, 
			pg.id AS guest_id, 
			pg.nickname AS guest_nickname,
			r.guest_continue,
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
		sqlGuestID       sql.NullInt64
		sqlGuestNickname sql.NullString
		sqlGuestContinue sql.NullBool
		sqlGameID        sql.NullInt64
		sqlDescription   sql.NullString
	)

	room := &models.Room{}
	err := row.Scan(
		&room.ID,
		&room.Host.ID,
		&room.Host.Nickname,
		&room.Host.Continue,
		&sqlGuestID,
		&sqlGuestNickname,
		&sqlGuestContinue,
		&sqlGameID,
		&room.Title,
		&sqlDescription,
		&room.Phase)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.NewNotFoundError("room not exist")
		} else {
			return nil, models.NewGenericErrorWithCause("room select failed", err)
		}
	}

	if sqlGuestID.Valid {
		room.Guest = &models.RoomParticipant{}

		if sqlGuestNickname.Valid {
			room.Guest.Nickname = sqlGuestNickname.String
		}

		if sqlGuestContinue.Valid {
			room.Guest.Continue = sqlGuestContinue.Bool
		}
	}

	if sqlGameID.Valid {
		room.GameID = &sqlGameID.Int64
	}

	if sqlDescription.Valid {
		room.Description = sqlDescription.String
	}

	return room, nil
}

func (r *roomRepositoryImpl) GetList(ctx context.Context, host string, title string, description string, phase models.RoomPhase) ([]*models.Room, error) {
	sqlStr := `
		SELECT 
			r.id, 
			ph.id AS host_id, 
			ph.nickname AS host_nickname,
			r.title, 
			r.description, 
			r.phase
		FROM rooms AS r
		INNER JOIN players AS ph ON ph.id = r.host_id
		WHERE (r.phase = $1)
		AND ($2::text IS NULL OR r.title LIKE $2)
		AND ($3::text IS NULL OR r.description LIKE $3)
		AND ($4::text IS NULL OR  ph.nickname LIKE $4)
		`

	var titleArg sql.NullString
	if title != "" {
		titleArg.String = fmt.Sprintf("%%%s%%", title)
		titleArg.Valid = true
	}

	var descriptionArg sql.NullString
	if description != "" {
		descriptionArg.String = fmt.Sprintf("%%%s%%", description)
		descriptionArg.Valid = true
	}

	var hostArg sql.NullString
	if host != "" {
		hostArg.String = fmt.Sprintf("%%%s%%", host)
		hostArg.Valid = true
	}
	rows, err := r.db.QueryContext(ctx, sqlStr, phase, titleArg, descriptionArg, hostArg)
	if err != nil {
		return nil, models.NewGenericErrorWithCause("rooms select failed", err)
	}
	defer rows.Close()

	rooms := make([]*models.Room, 0)
	var sqlDescription sql.NullString
	for rows.Next() {
		room := &models.Room{}
		err := rows.Scan(&room.ID, &room.Host.ID, &room.Host.Nickname, &room.Title, &sqlDescription, &room.Phase)
		if err != nil {
			return nil, models.NewGenericErrorWithCause("room scan failed", err)
		}

		if sqlDescription.Valid {
			room.Description = sqlDescription.String
		}

		rooms = append(rooms, room)
	}

	// check for errors during iteration
	if err = rows.Err(); err != nil {
		return nil, models.NewGenericErrorWithCause("rooms iteration error", err)
	}

	return rooms, nil
}

func (r *roomRepositoryImpl) Create(ctx context.Context, room *models.Room) (int64, error) {
	sqlStr := `
		INSERT INTO rooms(host_id, host_continue, title, description, phase)
		VALUES($1, $2, $3, $4, $5)
		RETURNING id
		`

	var id int64
	err := r.db.QueryRowContext(ctx, sqlStr, room.Host.ID, room.Host.Continue, room.Title, room.Description, room.Phase).Scan(&id)
	if err != nil {
		err = models.NewGenericErrorWithCause("room insert failed", err)
	}

	return id, err
}

func (r *roomRepositoryImpl) Update(ctx context.Context, room *models.Room) error {
	sqlStr := `
		UPDATE rooms
		SET host_continue  = $2,
			guest_id       = $3,
			guest_continue = $4,
			game_id 	   = $5,
			phase 		   = $6 
		WHERE id     	   = $1
		`
	var (
		sqlGuestID       *int64
		sqlGuestContinue bool = false
	)

	if room.Guest != nil {
		sqlGuestID = &room.Guest.ID
		sqlGuestContinue = room.Guest.Continue
	}

	result, err := r.db.ExecContext(ctx, sqlStr, room.ID, room.Host.Continue, sqlGuestID, sqlGuestContinue, room.GameID, room.Phase)
	if err != nil {
		return models.NewGenericErrorWithCause("room update failed", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.NewGenericErrorWithCause("could not get rows affected", err)
	}
	if rowsAffected == 0 {
		return models.NewGenericError("no room was updated")
	}

	return err
}

func (r *roomRepositoryImpl) Delete(ctx context.Context, id int64) error {
	sqlStr := `DELETE FROM rooms WHERE id = $1`

	result, err := r.db.ExecContext(ctx, sqlStr, id)
	if err != nil {
		return models.NewGenericErrorWithCause("room deletion failed", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.NewGenericErrorWithCause("could not get rows affected", err)
	}
	if rowsAffected == 0 {
		return models.NewGenericError("no room was deleted")
	}

	return err
}
