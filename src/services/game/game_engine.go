package game

import (
	"github.com/plamen-v/tic-tac-toe/src/models"
	"github.com/plamen-v/tic-tac-toe/src/services/repository"
)

type GameEngine interface {
	CreateRoom(*models.Room) (int64, error)
}

type gameEngine struct {
	playerRepo repository.PlayerRepository
	roomRepo   repository.RoomRepository
	gameRepo   repository.GameRepository
}

func NewGameEngine(
	playerRepo repository.PlayerRepository,
	roomRepo repository.RoomRepository,
	gameRepo repository.GameRepository) GameEngine {
	return &gameEngine{
		playerRepo: playerRepo,
		roomRepo:   roomRepo,
		gameRepo:   gameRepo,
	}
}

func (g *gameEngine) GetOpenRooms() (models.Room, error) {
	return models.Room{}, nil
}

func (g *gameEngine) CreateRoom(room *models.Room) (int64, error) {
	return 0, nil
}

func (g *gameEngine) GetRoom(id int) (models.Room, error) {
	return models.Room{}, nil
}

func (g *gameEngine) DeleteRoom(id int) (models.Room, error) {
	return models.Room{}, nil
}

func (g *gameEngine) SetRoomHostReady(id int) (models.Room, error) {
	return models.Room{}, nil
}

func (g *gameEngine) AddRoomGuest(id int) (models.Room, error) {
	return models.Room{}, nil
}

func (g *gameEngine) SetRoomGuestReady(id int) (models.Room, error) {
	return models.Room{}, nil
}

func (g *gameEngine) RemoveRoomGuest(id int) (models.Room, error) {
	return models.Room{}, nil
}

func (g *gameEngine) GetGame(id int) (models.Room, error) {
	return models.Room{}, nil
}

func (g *gameEngine) SetGameBoard(id int) (models.Room, error) {
	return models.Room{}, nil
}
