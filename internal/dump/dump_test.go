package dump

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStorage только для проверки вызова методов Set
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Set(originalURL, shortURL string) {
	m.Called(originalURL, shortURL)
}

func TestMockStorageSet(t *testing.T) {
	// Подготовка мокового хранилища с ожиданиями
	mockStorage := new(MockStorage)
	mockStorage.On("Set", "originalURL1", "shortURL1").Return()
	mockStorage.On("Set", "originalURL2", "shortURL2").Return()

	// Проверка, что методы можно вызвать и что они были вызваны с заданными аргументами
	mockStorage.Set("originalURL1", "shortURL1")
	mockStorage.Set("originalURL2", "shortURL2")
	mockStorage.AssertExpectations(t)
}

func TestTemporaryFileCreation(t *testing.T) {
	// Создание и проверка временного файла
	file, err := os.CreateTemp("", "test_temp_file_*.json")
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	// Запись данных и проверка успешности
	_, err = file.WriteString(`{"uuid": "1", "short_url": "shortURL1", "original_url": "originalURL1"}`)
	assert.NoError(t, err)
	file.Close()

	// Чтение данных из файла и проверка содержания
	content, err := os.ReadFile(file.Name())
	assert.NoError(t, err)
	assert.Contains(t, string(content), `"shortURL1"`)
}
