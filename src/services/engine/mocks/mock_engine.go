package mocks

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/plamen-v/tic-tac-toe-models/models"
	"github.com/stretchr/testify/mock"
)

type MockGameEngineService struct {
	mock.Mock
}

func (m *MockGameEngineService) GetRoom(ctx context.Context, playerID uuid.UUID) (*models.Room, error) {
	args := m.Called(ctx, playerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Room), args.Error(1)
}
func (m *MockGameEngineService) GetOpenRooms(ctx context.Context, pPageSize int, pPage int) ([]*models.Room, int, int, int, error) {
	args := m.Called(ctx)

	rooms, okRooms := args.Get(0).([]*models.Room)
	pageSize, okPageSize := args.Get(1).(int)
	page, okPage := args.Get(2).(int)
	total, okTotal := args.Get(3).(int)

	if rooms == nil || !okRooms || !okPageSize || !okPage || !okTotal {
		return nil, 0, 0, 0, args.Error(4)
	}

	return rooms, pageSize, page, total, args.Error(4)
}
func (m *MockGameEngineService) CreateRoom(ctx context.Context, playerID uuid.UUID, title string, description string) (uuid.UUID, error) {
	args := m.Called(ctx, playerID, title, description)
	if args.Get(0) == nil {
		return uuid.Nil, args.Error(1)
	}
	return args.Get(0).(uuid.UUID), args.Error(1)
}
func (m *MockGameEngineService) PlayerJoinRoom(ctx context.Context, roomID uuid.UUID, playerID uuid.UUID) error {
	args := m.Called(ctx, roomID, playerID)
	return args.Error(0)
}
func (m *MockGameEngineService) PlayerLeaveRoom(ctx context.Context, roomID uuid.UUID, playerID uuid.UUID) error {
	args := m.Called(ctx, roomID, playerID)
	return args.Error(0)
}
func (m *MockGameEngineService) CreateGame(ctx context.Context, roomID uuid.UUID, playerID uuid.UUID) (uuid.UUID, error) {
	args := m.Called(ctx, roomID, playerID)
	if args.Get(0) == nil {
		return uuid.Nil, args.Error(1)
	}
	return args.Get(0).(uuid.UUID), args.Error(1)
}
func (m *MockGameEngineService) GetGameState(ctx context.Context, roomID uuid.UUID, playerID uuid.UUID) (*models.Game, error) {
	args := m.Called(ctx, roomID, playerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Game), args.Error(1)
}
func (m *MockGameEngineService) PlayerMakeMove(ctx context.Context, roomID uuid.UUID, playerID uuid.UUID, position int) error {
	args := m.Called(ctx, roomID, playerID, position)
	return args.Error(0)
}
func (m *MockGameEngineService) GetRanking(ctx context.Context, pageSize int, page int) ([]*models.Player, int, int, int, error) {
	args := m.Called(ctx)

	players, okPlayers := args.Get(0).([]*models.Player)
	pageSize, okPageSize := args.Get(1).(int)
	page, okPage := args.Get(2).(int)
	total, okTotal := args.Get(3).(int)

	if players == nil || !okPlayers || !okPageSize || !okPage || !okTotal {
		return nil, 0, 0, 0, args.Error(4)
	}

	return players, pageSize, page, total, args.Error(4)
}
