package repository

import (
	"database/sql"

	"github.com/plamen-v/tic-tac-toe/src/models"
)

type GameRepository interface {
	Get(int) (*models.Game, error)
	Create(*models.Game) (*models.Game, error)
	Update(*models.Game) (*models.Game, error)
	Delete(int) (int, error)
}

func NewGameRepository(db *sql.DB) GameRepository {
	return &gameRepository{
		db: db,
	}
}

type gameRepository struct {
	db *sql.DB
}

func (r *gameRepository) Get(id int) (*models.Game, error) {
	return &models.Game{}, nil
}

func (r *gameRepository) Create(room *models.Game) (*models.Game, error) {
	return &models.Game{}, nil
}

func (r *gameRepository) Update(room *models.Game) (*models.Game, error) {
	return &models.Game{}, nil
}

func (r *gameRepository) Delete(id int) (int, error) {
	return 0, nil
}
