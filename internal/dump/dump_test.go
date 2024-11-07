package dump

import (
	"os"
	"testing"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestFillFromStorage(t *testing.T) {
	// Создаем временный файл для теста
	tempFile, err := os.CreateTemp("", "test_fill_from_storage_*.json")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// Записываем тестовые данные в файл
	_, err = tempFile.WriteString(`{"uuid": "1", "short_url": "shortURL1", "original_url": "originalURL1"}
{"uuid": "2", "short_url": "shortURL2", "original_url": "originalURL2"}`)
	assert.NoError(t, err)
	tempFile.Close()

	// Инициализируем новое хранилище
	storageInstance := storage.NewStorage()

	// Вызываем FillFromStorage
	err = FillFromStorage(storageInstance, tempFile.Name())
	assert.NoError(t, err)

	// Проверяем, что данные корректно загружены в хранилище
	value, exists := storageInstance.Get("originalURL1")
	assert.True(t, exists)
	assert.Equal(t, "shortURL1", value)

	value, exists = storageInstance.Get("originalURL2")
	assert.True(t, exists)
	assert.Equal(t, "shortURL2", value)
}
