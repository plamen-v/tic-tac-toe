package repository

import (
	"database/sql"

	"github.com/plamen-v/tic-tac-toe/src/models"
)

type RoomRepository interface {
	Get(int) (*models.Room, error)
	GetAll() ([]*models.Room, error)
	Create(*models.Room) (int64, error)
	Update(*models.Room) (*models.Room, error)
	Delete(int) (int, error)
}

func NewRoomRepository(db *sql.DB) RoomRepository {
	return &roomRepository{
		db: db,
	}
}

type roomRepository struct {
	db *sql.DB
}

func (r *roomRepository) Get(id int) (*models.Room, error) {
	return &models.Room{}, nil
}

func (r *roomRepository) GetAll() ([]*models.Room, error) {
	return make([]*models.Room, 1), nil
}

func (r *roomRepository) Create(room *models.Room) (int64, error) {

	result, err := r.db.Exec(`
		INSERT INTO rooms(host_id, host_ready, title, description, status)
		VALUES($1, $2, $3, $4, $5)`,
		room.Host.ID, true, room.Title, room.Description)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *roomRepository) Update(room *models.Room) (*models.Room, error) {
	return &models.Room{}, nil
}

func (r *roomRepository) Delete(id int) (int, error) {
	return 0, nil
}
