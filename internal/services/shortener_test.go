package services_test

import (
	"errors"
	"testing"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/services"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStore - mock for the Store interface
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

// MockRepository - mock for the Repository interface
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

// Test for the Set method (positive scenario)
func TestShortenerService_Set(t *testing.T) {
	mockRepo := new(MockRepository)
	mockStore := new(MockStore)

	service := services.NewShortenerService("http://localhost", mockRepo, mockStore, true)

	mockStore.On("Create", "https://example.com", mock.Anything, "user1").Return(nil)

	shortURL, err := service.Set("user1", "https://example.com")

	assert.NoError(t, err)
	assert.Contains(t, shortURL, "http://localhost/")
	mockStore.AssertCalled(t, "Create", "https://example.com", mock.Anything, "user1")
}

// Test for the Set method with an error
func TestShortenerService_Set_Error(t *testing.T) {
	mockRepo := new(MockRepository)
	mockStore := new(MockStore)

	service := services.NewShortenerService("http://localhost", mockRepo, mockStore, true)

	mockStore.On("Create", "https://example.com", mock.Anything, "user1").Return(errors.New("database error"))

	shortURL, err := service.Set("user1", "https://example.com")

	assert.Error(t, err)
	assert.Empty(t, shortURL)
	mockStore.AssertCalled(t, "Create", "https://example.com", mock.Anything, "user1")
}

// Test for the Get method from cache
func TestShortenerService_Get_FromCache(t *testing.T) {
	mockRepo := new(MockRepository)
	mockStore := new(MockStore)

	service := services.NewShortenerService("http://localhost", mockRepo, mockStore, false)

	mockRepo.On("Get", "short123").Return("https://example.com", true)

	originalURL, err := service.Get("short123")

	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", originalURL)
	mockRepo.AssertCalled(t, "Get", "short123")
}

// Test for the Get method from cache, if the link is missing
func TestShortenerService_Get_FromCache_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	mockStore := new(MockStore)

	service := services.NewShortenerService("http://localhost", mockRepo, mockStore, false)

	mockRepo.On("Get", "short123").Return("", false)

	originalURL, err := service.Get("short123")

	assert.Error(t, err)
	assert.Empty(t, originalURL)
	mockRepo.AssertCalled(t, "Get", "short123")
}

// Test for the Ping method
func TestShortenerService_Ping(t *testing.T) {
	mockRepo := new(MockRepository)
	mockStore := new(MockStore)

	service := services.NewShortenerService("http://localhost", mockRepo, mockStore, true)

	mockStore.On("PingStore").Return(nil)

	err := service.Ping()

	assert.NoError(t, err)
	mockStore.AssertCalled(t, "PingStore")
}

// Test for the CreateRep method
func TestShortenerService_CreateRep(t *testing.T) {
	mockRepo := new(MockRepository)
	mockStore := new(MockStore)

	service := services.NewShortenerService("http://localhost", mockRepo, mockStore, true)

	mockStore.On("Create", "https://example.com", "short123", "user1").Return(nil)

	err := service.CreateRep("https://example.com", "short123", "user1")

	assert.NoError(t, err)
	mockStore.AssertCalled(t, "Create", "https://example.com", "short123", "user1")
}

// Test for the GetFullRep method
func TestShortenerService_GetFullRep(t *testing.T) {
	mockRepo := new(MockRepository)
	mockStore := new(MockStore)

	service := services.NewShortenerService("http://localhost", mockRepo, mockStore, true)

	expectedResult := []map[string]string{
		{"short_url": "http://localhost/short123", "original_url": "https://example.com"},
	}

	mockStore.On("GetFull", "user1", "http://localhost").Return(expectedResult, nil)

	result, err := service.GetFullRep("user1")

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockStore.AssertCalled(t, "GetFull", "user1", "http://localhost")
}

// Test for the GetExistURL method
func TestShortenerService_GetExistURL(t *testing.T) {
	mockRepo := new(MockRepository)
	mockStore := new(MockStore)

	service := services.NewShortenerService("http://localhost", mockRepo, mockStore, true)

	mockStore.On("Get", "https://example.com").Return("short123", nil)

	pgErr := &pgconn.PgError{Code: pgerrcode.UniqueViolation}
	shortURL, err := service.GetExistURL("https://example.com", pgErr)

	assert.NoError(t, err)
	assert.Equal(t, "http://localhost/short123", shortURL)
	mockStore.AssertCalled(t, "Get", "https://example.com")
}

// Test for the GetURLCount method
func TestShortenerService_GetURLCount(t *testing.T) {
	mockRepo := new(MockRepository)
	mockStore := new(MockStore)

	service := services.NewShortenerService("http://localhost", mockRepo, mockStore, true)

	mockStore.On("GetURLCount").Return(42, nil)

	count, err := service.GetURLCount()

	assert.NoError(t, err)
	assert.Equal(t, 42, count)
	mockStore.AssertCalled(t, "GetURLCount")
}
func TestShortenerService_GetUserCount(t *testing.T) {
	mockRepo := new(MockRepository)
	mockStore := new(MockStore)

	service := services.NewShortenerService("http://localhost", mockRepo, mockStore, true)

	mockStore.On("GetUserCount").Return(42, nil)

	count, err := service.GetUserCount()

	assert.NoError(t, err)
	assert.Equal(t, 42, count)
	mockStore.AssertCalled(t, "GetUserCount")
}
