package app

import (
	"os"
	"testing"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/stretchr/testify/assert"
)

// createTempFile создает временный файл для тестов
func createTempFile(t *testing.T, content string) string {
	tmpFile, err := os.CreateTemp("", "testfile")
	assert.NoError(t, err)

	_, err = tmpFile.Write([]byte(content))
	assert.NoError(t, err)

	err = tmpFile.Close()
	assert.NoError(t, err)

	return tmpFile.Name()
}

// removeTempFile удаляет временный файл после теста
func removeTempFile(t *testing.T, filePath string) {
	err := os.Remove(filePath)
	assert.NoError(t, err)
}

// TestApp_Start_Success тестирует успешный запуск приложения
func TestApp_Start_Success(t *testing.T) {
	// Создаем временный файл для данных
	tempFilePath := createTempFile(t, "test data")
	defer removeTempFile(t, tempFilePath)

	// Настраиваем хранилище и конфигурацию
	mockStorage := &storage.Storage{}
	mockConfig := &config.Config{
		DBPath:     "test_db_path.db", // Используем временную базу данных
		FilePath:   tempFilePath,
		ServerAddr: "localhost:8080",
		BaseURL:    "http://localhost",
		LogLevel:   "info",
	}

	// Создаем приложение и запускаем его
	app := NewApp(mockStorage, mockConfig)
	err := app.Start()

	// Проверяем, что приложение запустилось без ошибок
	assert.NoError(t, err)

	// Проверяем, что данные из файла загружены в хранилище
	// Можно сделать assert для проверки данных в mockStorage
}

// TestApp_Start_InitDBError тестирует ошибку инициализации базы данных
func TestApp_Start_InitDBError(t *testing.T) {
	// Настраиваем хранилище и конфигурацию с неверным путем к базе данных
	mockStorage := &storage.Storage{}
	mockConfig := &config.Config{
		DBPath:     "", // Пустой путь для ошибки
		FilePath:   "test_file_path",
		ServerAddr: "localhost:8080",
		BaseURL:    "http://localhost",
		LogLevel:   "info",
	}

	// Создаем приложение и запускаем его
	app := NewApp(mockStorage, mockConfig)
	err := app.Start()

	// Проверяем, что приложение не вернуло ошибку (только логируется ошибка)
	assert.NoError(t, err)
}

// Тест для метода Stop
func TestApp_Stop(t *testing.T) {
	// Настройка конфигурации и хранилища для теста
	mockStorage := &storage.Storage{}
	mockConfig := &config.Config{
		DBPath:   "test_db_path",
		FilePath: "test_file_path",
	}

	// Создаем экземпляр App
	app := NewApp(mockStorage, mockConfig)

	// Вызываем метод Stop и проверяем, что он не вызывает паники
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Stop вызвал панику: %v", r)
		}
	}()

	app.Stop()
}
