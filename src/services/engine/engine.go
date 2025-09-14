package engine

import (
	"errors"

	"github.com/plamen-v/tic-tac-toe/src/models"
	"github.com/plamen-v/tic-tac-toe/src/services/repository"
)

type GameEngineService interface {
	CreateNewRoom(*models.Room) (int64, error)
	GetOpenRooms(roomFilter *models.RoomFilter) ([]*models.Room, error)
}

type gameEngineService struct {
	repo repository.Repository
}

func NewGameEngineService(repo repository.Repository) GameEngineService {
	return &gameEngineService{
		repo: repo,
	}
}

func (g *gameEngineService) GetOpenRooms(roomFilter *models.RoomFilter) ([]*models.Room, error) {
	roomFilter.Status = models.ROOM_OPEN
	rooms, err := g.repo.Rooms().GetList(roomFilter)
	if err != nil {
		return make([]*models.Room, 1), err
	}

	return rooms, nil
}

func (g *gameEngineService) CreateNewRoom(room *models.Room) (id int64, err error) {
	tx, err := g.repo.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	host, err := g.repo.Players().GetByID(*room.Host.ID, tx) //todo!
	if err != nil {
		return
	}

	err = g.validateCreateNewRoom(host, room)
	if err != nil {
		return
	}

	id, err = g.repo.Rooms().Create(room, tx)
	if err != nil {
		return
	}

	host.RoomID = &id
	err = g.repo.Players().Update(host, tx)
	if err != nil {
		return
	}

	return
}

func (g *gameEngineService) validateCreateNewRoom(host *models.Player, room *models.Room) error {
	if host.RoomID != nil {
		return errors.New("")
	}

	if len(room.Title) == 0 {
		return errors.New("")
	}

	return nil
}

func (g *gameEngineService) GetRoom(id int) (*models.Room, error) {
	return &models.Room{}, nil
}

func (g *gameEngineService) DeleteRoom(id int) (models.Room, error) {
	return models.Room{}, nil
}

func (g *gameEngineService) SetRoomHostReady(id int) (models.Room, error) {
	return models.Room{}, nil
}

func (g *gameEngineService) AddRoomGuest(id int) (models.Room, error) {
	return models.Room{}, nil
}

func (g *gameEngineService) SetRoomGuestReady(id int) (models.Room, error) {
	return models.Room{}, nil
}

func (g *gameEngineService) RemoveRoomGuest(id int) (models.Room, error) {
	return models.Room{}, nil
}

func (g *gameEngineService) GetGame(id int) (models.Room, error) {
	return models.Room{}, nil
}

func (g *gameEngineService) SetGameBoard(id int) (models.Room, error) {
	return models.Room{}, nil
}
