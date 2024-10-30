package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

func (m *MockStore) PingStore() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStore) Create(originalURL, shortURL, userID string) error {
	args := m.Called(originalURL, shortURL, userID)
	return args.Error(0)
}

func (m *MockStore) Get(shortURL string, originalURL string) (string, error) {
	args := m.Called(shortURL, originalURL)
	return args.String(0), args.Error(1)
}

func TestShortenerService_Set(t *testing.T) {
	mockStorage := new(MockStorage)
	mockStore := new(MockStore)
	mockService := NewShortenerService("http://localhost:8080", mockStorage, mockStore, false)

	mockStorage.On("Set", mock.Anything, "originalURL1").Return()
	shortURL, err := mockService.Set("userID", "originalURL1")

	assert.NoError(t, err)
	assert.Contains(t, shortURL, "http://localhost:8080/")
}

func TestShortenerService_Get(t *testing.T) {
	mockStorage := new(MockStorage)	
	mockStore := new(MockStore)
	mockService := NewShortenerService("http://localhost:8080", mockStorage, mockStore, true)

	mockStore.On("Get", "shortURL1", "").Return("originalURL1", nil)

	originalURL, err := mockService.Get("shortURL1")
	assert.NoError(t, err)
	assert.Equal(t, "originalURL1", originalURL)
}
