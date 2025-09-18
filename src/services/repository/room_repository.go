package repository

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/plamen-v/tic-tac-toe-models/models"
	"github.com/plamen-v/tic-tac-toe-models/models/errors"
)

// todo! error msg
type RoomRepository interface {
	Get(int64, *sql.Tx) (*models.Room, error)
	GetList(string, string, string, models.RoomPhase) ([]*models.Room, error)
	Create(*models.Room, *sql.Tx) (int64, error)
	Update(*models.Room, *sql.Tx) error
	Delete(int64, *sql.Tx) error
}

func NewRoomRepository(db *sql.DB) RoomRepository {
	return &roomRepository{
		db: db,
	}
}

type roomRepository struct {
	db *sql.DB
}

func (r *roomRepository) Get(id int64, tx *sql.Tx) (*models.Room, error) {
	var (
		sqlGuestID       sql.NullInt64
		sqlGuestNickname sql.NullString
		// sqlGuestRoomID   sql.NullInt64
		// sqlGuestGameID   sql.NullInt64
		// sqlGuestWins     sql.NullInt64
		// sqlGuestLosses   sql.NullInt64
		// sqlGuestDraws    sql.NullInt64
		sqlGuestContinue sql.NullBool
		sqlGameID        sql.NullInt64
		sqlPrevGameID    sql.NullInt64
		sqlDescription   sql.NullString
	)
	room := &models.Room{}

	forUpdate := ""
	if tx != nil {
		forUpdate = "FOR UPDATE"
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
			r.prev_game_id,
			r.title, 
			r.description, 
			r.phase
		FROM rooms AS r
		INNER JOIN players AS ph ON ph.id = r.host_id
		LEFT JOIN players AS pg ON pg.id = r.guest_id
		WHERE r.id = $1
		%s`, forUpdate)

	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow(sqlStr, id)
	} else {
		row = r.db.QueryRow(sqlStr, id)
	}

	err := row.Scan(
		&room.ID,

		&room.Host.ID,
		&room.Host.Nickname,
		// &room.Host.RoomID,
		// &room.Host.GameID,
		// &room.Host.Wins,
		// &room.Host.Losses,
		// &room.Host.Draws,
		&room.Host.Continue,

		&sqlGuestID,
		&sqlGuestNickname,
		// &sqlGuestRoomID,
		// &sqlGuestGameID,
		// &sqlGuestWins,
		// &sqlGuestLosses,
		// &sqlGuestDraws,
		&sqlGuestContinue,

		&sqlGameID,
		&sqlPrevGameID,
		&room.Title,
		&sqlDescription,
		&room.Phase)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundErrorf("room with id equal to %d not exist", id)
		} else {
			return nil, errors.NewGenericErrorWithCause("room not selected", err)
		}
	}

	if sqlGuestID.Valid {
		room.Guest = &models.RoomParticipant{}

		if sqlGuestNickname.Valid {
			room.Guest.Nickname = sqlGuestNickname.String
		}

		// if sqlGuestRoomID.Valid {
		// 	room.Guest.RoomID = &sqlGuestRoomID.Int64
		// }

		// if sqlGuestGameID.Valid {
		// 	room.Guest.GameID = &sqlGuestGameID.Int64
		// }

		// if sqlGuestWins.Valid {
		// 	room.Guest.Wins = sqlGuestWins.Int64
		// }

		// if sqlGuestLosses.Valid {
		// 	room.Guest.Losses = sqlGuestLosses.Int64
		// }

		// if sqlGuestDraws.Valid {
		// 	room.Guest.Draws = sqlGuestDraws.Int64
		// }

		if sqlGuestContinue.Valid {
			room.Guest.Continue = sqlGuestContinue.Bool
		}
	}

	if sqlGameID.Valid {
		room.GameID = &sqlGameID.Int64
	}

	if sqlPrevGameID.Valid {
		room.PrevGameID = &sqlPrevGameID.Int64
	}

	if sqlDescription.Valid {
		room.Description = sqlDescription.String
	}

	return room, nil
}

func (r *roomRepository) GetList(host string, title string, description string, phase models.RoomPhase) ([]*models.Room, error) {
	var (
		rooms          = make([]*models.Room, 10)
		sqlStrBuilder  strings.Builder
		sqlDescription sql.NullString
	)

	sqlStrBuilder.WriteString(`
		SELECT 
			r.id, 
			ph.id AS host_id, 
			ph.nickname AS host_nickname,
			r.title, 
			r.description, 
			r.phase
		FROM rooms AS r
		INNER JOIN players AS ph ON ph.id = r.host_id
		WHERE 1=1
		%s`)
	index := 1
	args := make([]any, 1)

	sqlStrBuilder.WriteString(" AND r.phase = $")
	sqlStrBuilder.WriteString(strconv.Itoa(index))
	index += 1
	args = append(args, phase)

	if len(title) > 0 {
		sqlStrBuilder.WriteString(" AND r.title LIKE $")
		sqlStrBuilder.WriteString(strconv.Itoa(index))
		index += 1
		args = append(args, fmt.Sprintf("%%%s%%", title))
	}
	if len(description) > 0 {
		sqlStrBuilder.WriteString(" AND r.description LIKE $")
		sqlStrBuilder.WriteString(strconv.Itoa(index))
		index += 1
		args = append(args, fmt.Sprintf("%%%s%%", description))
	}
	if len(host) > 0 {
		sqlStrBuilder.WriteString(" AND ph.nickname LIKE $")
		sqlStrBuilder.WriteString(strconv.Itoa(index))
		index += 1
		args = append(args, fmt.Sprintf("%%%s%%", host))
	}

	rows, err := r.db.Query(sqlStrBuilder.String(), args...)
	if err != nil {
		return nil, errors.NewGenericErrorWithCause("rooms not selected", err)
	}
	defer rows.Close() // ensure rows are closed when done

	for rows.Next() {
		room := &models.Room{}
		err := rows.Scan(&room.ID, &room.Host.ID, &room.Host.Nickname, &room.Title, &sqlDescription, &room.Phase)
		if err != nil {
			return nil, errors.NewGenericErrorWithCause("rooms not selected", err)
		}

		if sqlDescription.Valid {
			room.Description = sqlDescription.String
		}

		rooms = append(rooms, room)
	}

	// check for errors during iteration
	if err = rows.Err(); err != nil {
		return nil, errors.NewGenericErrorWithCause("rooms not selected", err)
	}

	return rooms, nil
}

func (r *roomRepository) Create(room *models.Room, tx *sql.Tx) (int64, error) {
	sqlStr := `
		INSERT INTO rooms(host_id, host_continue, title, description, phase)
		VALUES($1, $2, $3, $4, $5)
		RETURNING id`

	var (
		err error
		id  int64
	)

	if tx != nil {
		err = tx.QueryRow(sqlStr, room.Host.ID, room.Host.Continue, room.Title, room.Description, room.Phase).Scan(&id)
	} else {
		err = tx.QueryRow(sqlStr, room.Host.ID, room.Host.Continue, room.Title, room.Description, room.Phase).Scan(&id)
	}

	if err != nil {
		err = errors.NewGenericErrorWithCause("room not created", err)
	}

	return id, err
}

func (r *roomRepository) Update(room *models.Room, tx *sql.Tx) error {
	sqlStr := `
		UPDATE rooms
		SET host_continue  = $2,
			guest_id       = $3,
			guest_continue = $4,
			game_id 	   = $5, 
			prev_game_id   = $6, 
			phase 		   = $7 
		WHERE id     	   = $1`

	var (
		result sql.Result
		err    error
	)
	if tx != nil {
		result, err = tx.Exec(sqlStr, room.Host.Continue, room.Guest.ID, room.GameID, room.PrevGameID, room.Phase)
	} else {
		result, err = r.db.Exec(sqlStr, room.Host.Continue, room.Guest.ID, room.GameID, room.PrevGameID, room.Phase)
	}

	if err != nil {
		return errors.NewGenericErrorWithCause("room not updated", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewGenericErrorWithCause("room not updated", err)
	}
	if rowsAffected != 1 {
		return errors.NewGenericErrorWithCause("room not updated", err)
	}

	return err
}

func (r *roomRepository) Delete(id int64, tx *sql.Tx) error {
	sqlStr := `DELETE FROM rooms WHERE id = $1`

	var (
		result sql.Result
		err    error
	)
	if tx != nil {
		result, err = tx.Exec(sqlStr, id)
	} else {
		result, err = r.db.Exec(sqlStr, id)
	}

	if err != nil {
		return errors.NewGenericErrorWithCause("room not deleted", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewGenericErrorWithCause("room not deleted", err)
	}
	if rowsAffected != 1 {
		return errors.NewGenericErrorWithCause("room not deleted", err)
	}

	return err
}
