package dump

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
)

type MockStorage struct {
	mock.Mock
	storage.Storage
}

func (m *MockStorage) Set(originalURL, shortURL string) {
	m.Called(originalURL, shortURL)
}

func TestFillFromStorage(t *testing.T) {
	// Подготовка мокового хранилища
	mockStorage := new(MockStorage)
	mockStorage.On("Set", "originalURL1", "shortURL1").Return()
	mockStorage.On("Set", "originalURL2", "shortURL2").Return()

	// Создаем временный файл
	file, _ := os.CreateTemp("", "test_fill_from_storage_*.json")
	defer os.Remove(file.Name())

	// Записываем данные в JSON формате
	file.WriteString(`{"uuid": "1", "short_url": "shortURL1", "original_url": "originalURL1"}
	{"uuid": "2", "short_url": "shortURL2", "original_url": "originalURL2"}`)
	file.Close()

	// Тестируем FillFromStorage
	err := FillFromStorage(mockStorage, file.Name())
	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

func TestSet(t *testing.T) {
	// Подготовка мокового хранилища
	mockStorage := &storage.Storage{URLs: map[string]string{"shortURL1": "originalURL1"}}

	// Создаем временный файл
	file, _ := os.CreateTemp("", "test_set_*.json")
	defer os.Remove(file.Name())

	// Тестируем Set
	err := Set(mockStorage, file.Name())
	assert.NoError(t, err)

	// Проверяем содержимое файла
	content, _ := os.ReadFile(file.Name())
	assert.Contains(t, string(content), "shortURL1")
}
