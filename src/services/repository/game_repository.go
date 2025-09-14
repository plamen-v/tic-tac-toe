package repository

import (
	"database/sql"

	"github.com/plamen-v/tic-tac-toe/src/models"
)

type GameRepository interface {
	Get(int, *sql.Tx) (*models.Game, error)
	Create(*models.Game, *sql.Tx) (*models.Game, error)
	Update(*models.Game, *sql.Tx) (*models.Game, error)
	Delete(int, *sql.Tx) (int, error)
}

func NewGameRepository(db *sql.DB) GameRepository {
	return &gameRepository{
		db: db,
	}
}

type gameRepository struct {
	db *sql.DB
}

func (r *gameRepository) Get(id int, tx *sql.Tx) (*models.Game, error) {
	return &models.Game{}, nil
}

func (r *gameRepository) Create(room *models.Game, tx *sql.Tx) (*models.Game, error) {
	return &models.Game{}, nil
}

func (r *gameRepository) Update(room *models.Game, tx *sql.Tx) (*models.Game, error) {
	return &models.Game{}, nil
}

func (r *gameRepository) Delete(id int, tx *sql.Tx) (int, error) {
	return 0, nil
}
