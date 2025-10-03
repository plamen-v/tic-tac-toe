package mocks

import (
	"github.com/plamen-v/tic-tac-toe/src/services/logger"
	"github.com/stretchr/testify/mock"
)

type MockLoggerService struct {
	mock.Mock
}

func (m *MockLoggerService) Info(msg string, fields ...logger.Field) {
	m.Called(msg, fields)
}
func (m *MockLoggerService) Debug(msg string, fields ...logger.Field) {
	m.Called(msg, fields)
}
func (m *MockLoggerService) Error(msg string, fields ...logger.Field) {
	m.Called(msg, fields)
}

func (m *MockLoggerService) Sync() error {
	args := m.Called()
	return args.Error(0)
}
