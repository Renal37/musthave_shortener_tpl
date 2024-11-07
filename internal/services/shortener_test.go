package services_test

import (
	"errors"
	"testing"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStore - мок для интерфейса Store
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

func (m *MockStore) Get(shortID, originalURL string) (string, error) {
	args := m.Called(shortID, originalURL)
	return args.String(0), args.Error(1)
}

func (m *MockStore) GetFull(userID, BaseURL string) ([]map[string]string, error) {
	args := m.Called(userID, BaseURL)
	return args.Get(0).([]map[string]string), args.Error(1)
}

func (m *MockStore) DeleteURLs(userID, shortURL string, updateChan chan<- string) error {
	args := m.Called(userID, shortURL, updateChan)
	return args.Error(0)
}

// MockRepository - мок для интерфейса Repository
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

// Тест для метода Set (позитивный сценарий)
func TestShortenerService_Set(t *testing.T) {
	mockRepo := new(MockRepository) // Создаем мок для Repository
	mockStore := new(MockStore)     // Создаем мок для Store

	service := services.NewShortenerService("http://localhost", mockRepo, mockStore, true)

	mockStore.On("Create", "https://example.com", mock.Anything, "user1").Return(nil)

	shortURL, err := service.Set("user1", "https://example.com")

	assert.NoError(t, err)
	assert.Contains(t, shortURL, "http://localhost/")
	mockStore.AssertCalled(t, "Create", "https://example.com", mock.Anything, "user1")
}

// Тест для метода Set с ошибкой
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

// Тест для метода Get через кэш
func TestShortenerService_Get_FromCache(t *testing.T) {
	mockRepo := new(MockRepository) // Создаем мок для Repository
	mockStore := new(MockStore)     // Создаем мок для Store

	service := services.NewShortenerService("http://localhost", mockRepo, mockStore, false)

	mockRepo.On("Get", "short123").Return("https://example.com", true)

	originalURL, err := service.Get("short123")

	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", originalURL)
	mockRepo.AssertCalled(t, "Get", "short123")
}

// Тест для метода Get, если URL не найден в кэше
func TestShortenerService_Get_NotFound(t *testing.T) {
	mockRepo := new(MockRepository) // Создаем мок для Repository
	mockStore := new(MockStore)     // Создаем мок для Store

	service := services.NewShortenerService("http://localhost", mockRepo, mockStore, false)

	mockRepo.On("Get", "short123").Return("", false)

	originalURL, err := service.Get("short123")

	assert.Error(t, err)
	assert.Empty(t, originalURL)
	mockRepo.AssertCalled(t, "Get", "short123")
}

// Тест для метода Get с ошибкой при извлечении из базы данных
func TestShortenerService_Get_Error(t *testing.T) {
	mockRepo := new(MockRepository) // Создаем мок для Repository
	mockStore := new(MockStore)     // Создаем мок для Store

	service := services.NewShortenerService("http://localhost", mockRepo, mockStore, true)

	mockStore.On("Get", "short123", "").Return("", errors.New("db error"))

	originalURL, err := service.Get("short123")

	assert.Error(t, err)
	assert.Empty(t, originalURL)
	mockStore.AssertCalled(t, "Get", "short123", "")
}

// Тест для метода Ping (позитивный сценарий)
func TestShortenerService_Ping(t *testing.T) {
	mockRepo := new(MockRepository)
	mockStore := new(MockStore)

	service := services.NewShortenerService("http://localhost", mockRepo, mockStore, true)

	mockStore.On("PingStore").Return(nil)

	err := service.Ping()

	assert.NoError(t, err)
	mockStore.AssertCalled(t, "PingStore")
}

// Тест для метода Ping с ошибкой
func TestShortenerService_Ping_Error(t *testing.T) {
	mockRepo := new(MockRepository)
	mockStore := new(MockStore)

	service := services.NewShortenerService("http://localhost", mockRepo, mockStore, true)

	mockStore.On("PingStore").Return(errors.New("ping error"))

	err := service.Ping()

	assert.Error(t, err)
	assert.Equal(t, "pinging db-store: ping error", err.Error())
	mockStore.AssertCalled(t, "PingStore")
}

// Тест для метода CreateRep
func TestShortenerService_CreateRep(t *testing.T) {
	mockRepo := new(MockRepository)
	mockStore := new(MockStore)

	service := services.NewShortenerService("http://localhost", mockRepo, mockStore, true)

	mockStore.On("Create", "https://example.com", "short123", "user1").Return(nil)

	err := service.CreateRep("https://example.com", "short123", "user1")

	assert.NoError(t, err)
	mockStore.AssertCalled(t, "Create", "https://example.com", "short123", "user1")
}

// Тест для метода GetFullRep
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

// Тест для метода DeleteURLsRep
func TestShortenerService_DeleteURLsRep(t *testing.T) {
	mockRepo := new(MockRepository)
	mockStore := new(MockStore)

	service := services.NewShortenerService("http://localhost", mockRepo, mockStore, true)

	updateChan := make(chan string, 1)
	mockStore.On("DeleteURLs", "user1", "short123", updateChan).Return(nil)

	err := service.DeleteURLsRep("user1", []string{"short123"})

	assert.NoError(t, err)
	mockStore.AssertCalled(t, "DeleteURLs", "user1", "short123", updateChan)
}

// Тест для метода DeleteURLsRep с ошибкой
func TestShortenerService_DeleteURLsRep_Error(t *testing.T) {
	mockRepo := new(MockRepository)
	mockStore := new(MockStore)

	service := services.NewShortenerService("http://localhost", mockRepo, mockStore, true)

	updateChan := make(chan string, 1)
	mockStore.On("DeleteURLs", "user1", "short123", updateChan).Return(errors.New("delete error"))

	err := service.DeleteURLsRep("user1", []string{"short123"})

	assert.Error(t, err)
	mockStore.AssertCalled(t, "DeleteURLs", "user1", "short123", updateChan)
}
