package mocks

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/plamen-v/tic-tac-toe-models/models"
	"github.com/stretchr/testify/mock"
)

type MockPlayerRepository struct {
	mock.Mock
}

func (m *MockPlayerRepository) Get(ctx context.Context, id uuid.UUID) (*models.Player, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Player), args.Error(1)
}

func (m *MockPlayerRepository) GetByLogin(ctx context.Context, login string) (*models.Player, error) {
	args := m.Called(ctx, login)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Player), args.Error(1)
}

func (m *MockPlayerRepository) UpdateStats(ctx context.Context, player *models.Player) error {
	args := m.Called(ctx, player)
	return args.Error(0)
}

func (m *MockPlayerRepository) GetRanking(ctx context.Context, page int, pageSize int) ([]*models.Player, int, int, int, error) {
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

type MockGameRepository struct {
	mock.Mock
}

func (m *MockGameRepository) Get(ctx context.Context, id uuid.UUID) (*models.Game, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Game), args.Error(1)
}

func (m *MockGameRepository) Create(ctx context.Context, game *models.Game) (uuid.UUID, error) {
	args := m.Called(ctx, game)
	if args.Get(0) == nil {
		return uuid.Nil, args.Error(1)
	}
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockGameRepository) Update(ctx context.Context, game *models.Game) error {
	args := m.Called(ctx, game)
	return args.Error(0)
}

type MockRoomRepository struct {
	mock.Mock
}

func (m *MockRoomRepository) Get(ctx context.Context, id uuid.UUID, lock bool) (*models.Room, error) {
	args := m.Called(ctx, id, lock)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Room), args.Error(1)
}

func (m *MockRoomRepository) GetByPlayerID(ctx context.Context, playerID uuid.UUID) (*models.Room, error) {
	args := m.Called(ctx, playerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Room), args.Error(1)
}

func (m *MockRoomRepository) GetList(ctx context.Context, phase models.RoomPhase, pPageSize, pPage int) ([]*models.Room, int, int, int, error) {
	args := m.Called(ctx, phase)
	rooms, okPlayers := args.Get(0).([]*models.Room)
	pageSize, okPageSize := args.Get(1).(int)
	page, okPage := args.Get(2).(int)
	total, okTotal := args.Get(3).(int)

	if rooms == nil || !okPlayers || !okPageSize || !okPage || !okTotal {
		return nil, 0, 0, 0, args.Error(4)
	}

	return rooms, pageSize, page, total, args.Error(4)
}

func (m *MockRoomRepository) Create(ctx context.Context, room *models.Room) (uuid.UUID, error) {
	args := m.Called(ctx, room)
	if args.Get(0) == nil {
		return uuid.Nil, args.Error(1)
	}
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockRoomRepository) Update(ctx context.Context, room *models.Room) error {
	args := m.Called(ctx, room)
	return args.Error(0)
}

func (m *MockRoomRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
