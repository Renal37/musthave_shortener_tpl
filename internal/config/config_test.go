package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	// Сохраняем текущее значение переменных окружения для восстановления после теста
	oldServerAddr := os.Getenv("SERVER_ADDRESS")
	oldBaseURL := os.Getenv("BASE_URL")
	oldLogLevel := os.Getenv("FLAG_LOG_LEVEL")
	oldFilePath := os.Getenv("FILE_STORAGE_PATH")
	oldDBPath := os.Getenv("db")

	// Очищаем переменные окружения перед тестом
	os.Clearenv()

	// Проверяем значения по умолчанию
	config := InitConfig()
	assert.Equal(t, "localhost:8080", config.ServerAddr)
	assert.Equal(t, "http://localhost:8080", config.BaseURL)
	assert.Equal(t, "info", config.LogLevel)
	assert.Equal(t, "short-url-db.json", config.FilePath)
	assert.Equal(t, "", config.DBPath)

	// Устанавливаем переменные окружения для следующего теста
	os.Setenv("SERVER_ADDRESS", "127.0.0.1:9090")
	os.Setenv("BASE_URL", "http://127.0.0.1:9090")
	os.Setenv("FLAG_LOG_LEVEL", "debug")
	os.Setenv("FILE_STORAGE_PATH", "custom-file.json")
	os.Setenv("db", "postgres://user:password@localhost:5432/mydb")

	// Проверяем, что значения переменных окружения теперь переопределяют значения по умолчанию
	config = InitConfig()
	assert.Equal(t, "127.0.0.1:9090", config.ServerAddr)
	assert.Equal(t, "http://127.0.0.1:9090", config.BaseURL)
	assert.Equal(t, "debug", config.LogLevel)
	assert.Equal(t, "custom-file.json", config.FilePath)
	assert.Equal(t, "postgres://user:password@localhost:5432/mydb", config.DBPath)

	// Восстанавливаем старые значения переменных окружения
	if oldServerAddr != "" {
		os.Setenv("SERVER_ADDRESS", oldServerAddr)
	}
	if oldBaseURL != "" {
		os.Setenv("BASE_URL", oldBaseURL)
	}
	if oldLogLevel != "" {
		os.Setenv("FLAG_LOG_LEVEL", oldLogLevel)
	}
	if oldFilePath != "" {
		os.Setenv("FILE_STORAGE_PATH", oldFilePath)
	}
	if oldDBPath != "" {
		os.Setenv("db", oldDBPath)
	}
}
