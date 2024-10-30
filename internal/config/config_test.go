package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitConfig_Defaults(t *testing.T) {
	// Создаем новый экземпляр Config с начальными значениями
	config := &Config{
		ServerAddr: "localhost:8080",
		BaseURL:    "http://localhost:8080",
		LogLevel:   "info",
		FilePath:   "short-url-db.json",
		DBPath:     "",
	}

	assert.Equal(t, "localhost:8080", config.ServerAddr)
	assert.Equal(t, "http://localhost:8080", config.BaseURL)
	assert.Equal(t, "info", config.LogLevel)
	assert.Equal(t, "short-url-db.json", config.FilePath)
	assert.Equal(t, "", config.DBPath)
}

func TestInitConfig_WithEnvVars(t *testing.T) {
	// Сохраняем старые значения переменных окружения
	oldServerAddr := os.Getenv("SERVER_ADDRESS")
	oldBaseURL := os.Getenv("BASE_URL")
	oldLogLevel := os.Getenv("FLAG_LOG_LEVEL")
	oldFilePath := os.Getenv("FILE_STORAGE_PATH")
	oldDBPath := os.Getenv("db")

	// Устанавливаем переменные окружения для проверки их приоритета
	os.Setenv("SERVER_ADDRESS", "127.0.0.1:9090")
	os.Setenv("BASE_URL", "http://127.0.0.1:9090")
	os.Setenv("FLAG_LOG_LEVEL", "debug")
	os.Setenv("FILE_STORAGE_PATH", "custom-file.json")
	os.Setenv("db", "postgres://user:password@localhost:5432/mydb")

	// Создаем новый экземпляр Config и парсим окружение
	config := InitConfig()

	assert.Equal(t, "127.0.0.1:9090", config.ServerAddr)
	assert.Equal(t, "http://127.0.0.1:9090", config.BaseURL)
	assert.Equal(t, "debug", config.LogLevel)
	assert.Equal(t, "custom-file.json", config.FilePath)
	assert.Equal(t, "postgres://user:password@localhost:5432/mydb", config.DBPath)

	// Восстанавливаем старые значения переменных окружения
	if oldServerAddr != "" {
		os.Setenv("SERVER_ADDRESS", oldServerAddr)
	} else {
		os.Unsetenv("SERVER_ADDRESS")
	}
	if oldBaseURL != "" {
		os.Setenv("BASE_URL", oldBaseURL)
	} else {
		os.Unsetenv("BASE_URL")
	}
	if oldLogLevel != "" {
		os.Setenv("FLAG_LOG_LEVEL", oldLogLevel)
	} else {
		os.Unsetenv("FLAG_LOG_LEVEL")
	}
	if oldFilePath != "" {
		os.Setenv("FILE_STORAGE_PATH", oldFilePath)
	} else {
		os.Unsetenv("FILE_STORAGE_PATH")
	}
	if oldDBPath != "" {
		os.Setenv("db", oldDBPath)
	} else {
		os.Unsetenv("db")
	}
}
