package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStore — заглушка для интерфейса Store.
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

func (m *MockStore) Get(shortID string, originalURL string) (string, error) {
	args := m.Called(shortID, originalURL)
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

func (m *MockStore) GetURLCount() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockStore) GetUserCount() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

// MockRepository — заглушка для интерфейса Repository.
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Set(shortID string, originalURL string) {
	m.Called(shortID, originalURL)
}

func (m *MockRepository) Get(shortID string) (string, bool) {
	args := m.Called(shortID)
	return args.String(0), args.Bool(1)
}

func TestShortenerService_Set(t *testing.T) {
	mockStore := new(MockStore)
	mockRepo := new(MockRepository)
	service := NewShortenerService("http://localhost", mockRepo, mockStore, true)

	mockStore.On("Create", "https://example.com", mock.Anything, "user123").Return(nil)

	shortURL, err := service.Set("user123", "https://example.com")

	assert.NoError(t, err)
	assert.Contains(t, shortURL, "http://localhost")
	mockStore.AssertCalled(t, "Create", "https://example.com", mock.Anything, "user123")
}

func TestShortenerService_Get(t *testing.T) {
	mockStore := new(MockStore)
	mockRepo := new(MockRepository)
	service := NewShortenerService("http://localhost", mockRepo, mockStore, true)

	mockStore.On("Get", "short123", "").Return("https://example.com", nil)

	originalURL, err := service.Get("short123")

	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", originalURL)
	mockStore.AssertCalled(t, "Get", "short123", "")
}

func TestShortenerService_Ping(t *testing.T) {
	mockStore := new(MockStore)
	service := NewShortenerService("http://localhost", nil, mockStore, true)

	mockStore.On("PingStore").Return(nil)

	err := service.Ping()

	assert.NoError(t, err)
	mockStore.AssertCalled(t, "PingStore")
}
