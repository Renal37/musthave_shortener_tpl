package services

import (
	"github.com/stretchr/testify/mock"
)

// Мок для интерфейса Store
type MockStore struct {
	mock.Mock
}

func (m *MockStore) PingStore() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStore) Create(originalURL, shortURL, UserID string) error {
	args := m.Called(originalURL, shortURL, UserID)
	return args.Error(0)
}

func (m *MockStore) Get(shortIrl string, originalURL string) (string, error) {
	args := m.Called(shortIrl, originalURL)
	return args.String(0), args.Error(1)
}

func (m *MockStore) GetFull(userID string, BaseURL string) ([]map[string]string, error) {
	args := m.Called(userID, BaseURL)
	return args.Get(0).([]map[string]string), args.Error(1)
}

func (m *MockStore) DeleteURLs(userID string, shortURL string, updateChan chan<- string) error {
	args := m.Called(userID, shortURL, updateChan)
	return args.Error(0)
}
