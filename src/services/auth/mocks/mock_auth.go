package mocks

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/plamen-v/tic-tac-toe-models/models"
	"github.com/stretchr/testify/mock"
)

type MockAuthenticationService struct {
	mock.Mock
}

func (m *MockAuthenticationService) ValidateToken(token string) (*jwt.Token, error) {
	args := m.Called(token)
	if token, ok := args.Get(0).(*jwt.Token); ok {
		return token, nil
	}
	return nil, args.Error(1)
}

func (m *MockAuthenticationService) Authenticate(ctx context.Context, login string, password string) (*models.Player, string, error) {
	args := m.Called(ctx, login, password)
	if player, ok := args.Get(0).(*models.Player); ok {
		return player, args.Get(1).(string), nil
	}
	return nil, "", args.Error(2)
}
