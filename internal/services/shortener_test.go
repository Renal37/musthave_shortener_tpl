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

// MockStore - мок для интерфейса Store
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

// Тест для метода Set
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

// Тест для метода Set с ошибкой
func TestShortenerService_Set_Error(t *testing.T) {
	mockRepo := new(MockRepository)
	mockStore := new(MockStore)

	service := services.NewShortenerService("http://localhost", mockRepo, mockStore, true)

	mockStore.On("Create", "https://example.com", mock.Anything, "user1").Return(errors.New("ошибка базы данных"))

	shortURL, err := service.Set("user1", "https://example.com")

	assert.Error(t, err)
	assert.Empty(t, shortURL)
	mockStore.AssertCalled(t, "Create", "https://example.com", mock.Anything, "user1")
}

// Тест для метода Get из кэша
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

// Тест для метода Get из кэша, если ссылка не найдена
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

// Тест для метода GetExistURL, если ошибка отсутствует
func TestShortenerService_GetExistURL_NoError(t *testing.T) {
	mockRepo := new(MockRepository)
	mockStore := new(MockStore)

	service := services.NewShortenerService("http://localhost", mockRepo, mockStore, true)

	shortURL, err := service.GetExistURL("https://example.com", nil)

	assert.NoError(t, err)
	assert.Empty(t, shortURL) // Должен вернуть пустую строку, если ошибка уникальности отсутствует
}

// Тест для метода GetExistURL, если есть ошибка уникальности
func TestShortenerService_GetExistURL_Error(t *testing.T) {
	mockRepo := new(MockRepository)
	mockStore := new(MockStore)

	service := services.NewShortenerService("http://localhost", mockRepo, mockStore, true)

	pgErr := &pgconn.PgError{Code: pgerrcode.UniqueViolation}
	mockStore.On("Get", "", "https://example.com").Return("short123", nil)

	shortURL, err := service.GetExistURL("https://example.com", pgErr)

	assert.NoError(t, err)
	assert.Equal(t, "http://localhost/short123", shortURL)
	mockStore.AssertCalled(t, "Get", "", "https://example.com")
}

// Тест для метода Ping
func TestShortenerService_Ping(t *testing.T) {
	mockRepo := new(MockRepository)
	mockStore := new(MockStore)

	service := services.NewShortenerService("http://localhost", mockRepo, mockStore, true)

	mockStore.On("PingStore").Return(nil)

	err := service.Ping()

	assert.NoError(t, err)
	mockStore.AssertCalled(t, "PingStore")
}

// Тест для метода GetURLCount
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

// Тест для метода GetUserCount
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

func TestShortenerService_GetExistURL_UniqueViolation(t *testing.T) {
	mockStore := new(MockStore)
	service := services.NewShortenerService("http://short.url", new(MockRepository), mockStore, true)
	pgErr := &pgconn.PgError{Code: pgerrcode.UniqueViolation}
	mockStore.On("Get", "", "http://original.url").Return("shortID", nil)

	shortURL, err := service.GetExistURL("http://original.url", pgErr)
	assert.NoError(t, err)
	assert.Equal(t, "http://short.url/shortID", shortURL)
}

func TestShortenerService_GetURLCount_NoDB(t *testing.T) {
	service := services.NewShortenerService("http://short.url", new(MockRepository), new(MockStore), false)
	count, err := service.GetURLCount()
	assert.Error(t, err)
	assert.Equal(t, 0, count)
}

func TestShortenerService_GetUserCount_NoDB(t *testing.T) {
	service := services.NewShortenerService("http://short.url", new(MockRepository), new(MockStore), false)
	count, err := service.GetUserCount()
	assert.Error(t, err)
	assert.Equal(t, 0, count)
}

func TestShortenerService_GetFullRep(t *testing.T) {
	mockStore := new(MockStore)
	service := services.NewShortenerService("http://example.com", nil, mockStore, true)

	mockStore.On("GetFull", "user123", "http://example.com").Return([]map[string]string{
		{"shortURL": "short123", "originalURL": "http://original.com"},
	}, nil)

	urls, err := service.GetFullRep("user123")
	assert.NoError(t, err)
	assert.Len(t, urls, 1)
	assert.Equal(t, "short123", urls[0]["shortURL"])
}

func TestShortenerService_DeleteURLsRep(t *testing.T) {
	mockStore := new(MockStore)
	service := services.NewShortenerService("http://example.com", nil, mockStore, true)

	mockStore.On("DeleteURLs", "user123", "short123", mock.Anything).Return(nil)

	err := service.DeleteURLsRep("user123", []string{"short123"})
	assert.NoError(t, err)
}

func TestRandSeq(t *testing.T) {
	seq1 := services.RandSeq()
	seq2 := services.RandSeq()
	if seq1 == seq2 {
		t.Errorf("Expected unique random values, but got %s and %s", seq1, seq2)
	}
}
