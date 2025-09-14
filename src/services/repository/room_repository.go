package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/plamen-v/tic-tac-toe/src/models"
)

type RoomRepository interface {
	GetByID(int64, *sql.Tx) (*models.Room, error)
	GetList(filter *models.RoomFilter) ([]*models.Room, error)
	Create(*models.Room, *sql.Tx) (int64, error)
	Update(*models.Room, *sql.Tx) (*models.Room, error)
	Delete(int, *sql.Tx) (int, error)
}

func NewRoomRepository(db *sql.DB) RoomRepository {
	return &roomRepository{
		db: db,
	}
}

type roomRepository struct {
	db *sql.DB
}

func (r *roomRepository) GetByID(id int64, tx *sql.Tx) (*models.Room, error) {
	return &models.Room{}, nil
}

func (r *roomRepository) GetList(roomfilter *models.RoomFilter) ([]*models.Room, error) {
	var sqlStrBuilder strings.Builder
	sqlStrBuilder.WriteString(`
		SELECT 
			r.id, 
			r.host_id, 
			r.host_ready, 
			r.guest_id, 
			r.guest_ready, 
			r.game_id, 
			r.title, 
			r.description, 
			r.status
		FROM rooms AS r
		INNER JOIN players AS ph ON ph.id = r.host_id
		LEFT JOIN players AS pg ON pg.id = r.guest_id
		WHERE 1=1
		%s`)
	index := 1
	args := make([]any, 1)

	sqlStrBuilder.WriteString(" AND r.status = $")
	sqlStrBuilder.WriteString(strconv.Itoa(index))
	index += 1
	args = append(args, roomfilter.Status)

	if len(roomfilter.Title) > 0 {
		sqlStrBuilder.WriteString(" AND r.title LIKE $")
		sqlStrBuilder.WriteString(strconv.Itoa(index))
		index += 1
		args = append(args, fmt.Sprintf("%%%s%%", roomfilter.Title))
	}
	if len(roomfilter.Description) > 0 {
		sqlStrBuilder.WriteString(" AND r.description LIKE $")
		sqlStrBuilder.WriteString(strconv.Itoa(index))
		index += 1
		args = append(args, fmt.Sprintf("%%%s%%", roomfilter.Description))
	}
	if len(roomfilter.Host) > 0 {
		sqlStrBuilder.WriteString(" AND ph.nickname LIKE $")
		sqlStrBuilder.WriteString(strconv.Itoa(index))
		index += 1
		args = append(args, fmt.Sprintf("%%%s%%", roomfilter.Host))
	}

	rows, err := r.db.Query(sqlStrBuilder.String(), args...)
	if err != nil {
		return make([]*models.Room, 1), err
	}
	defer rows.Close() // ensure rows are closed when done

	for rows.Next() {
		var id int
		var name string

		// scan the columns into variables
		err := rows.Scan(&id, &name)
		if err != nil {
			// handle error while scanning this row
			log.Fatal(err)
		}

		// process the row
		fmt.Printf("User: %d, %s\n", id, name)
	}

	// check for errors during iteration
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return make([]*models.Room, 1), nil
}

func (r *roomRepository) Create(room *models.Room, tx *sql.Tx) (int64, error) {

	sqlStr := `
		INSERT INTO rooms(host_id, host_ready, guest_id, guest_ready, game_id, prev_game_id, title, description, status)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	var (
		err error
		id  int64
	)

	if tx != nil {
		err = tx.QueryRow(sqlStr,
			room.Host.ID, room.Host.IsReady, room.Guest.ID, room.Guest.IsReady,
			room.GameID, room.PrevGameID, room.Title, room.Description, room.Status).Scan(&id)
	} else {
		err = tx.QueryRow(sqlStr,
			room.Host.ID, room.Host.IsReady, room.Guest.ID, room.Guest.IsReady,
			room.GameID, room.PrevGameID, room.Title, room.Description, room.Status).Scan(&id)
	}

	return id, err
}

func (r *roomRepository) Update(room *models.Room, tx *sql.Tx) (*models.Room, error) {
	return &models.Room{}, nil
}

func (r *roomRepository) Delete(id int, tx *sql.Tx) (int, error) {
	return 0, nil
}
