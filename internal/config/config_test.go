package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitConfig_Defaults(t *testing.T) {
	// Создаем новый экземпляр Config с начальными значениями
	config := InitConfig()

	assert.Equal(t, "localhost:8080", config.ServerAddr)
	assert.Equal(t, "http://localhost:8080", config.BaseURL)
	assert.Equal(t, "info", config.LogLevel)
	assert.Equal(t, "short-url-db.json", config.FilePath)
	assert.Equal(t, "", config.DBPath)
	assert.Equal(t, "false", config.EnablePprof)
	assert.Equal(t, false, config.EnableHTTPS)
	assert.Equal(t, "cert.pem", config.CertFile)
	assert.Equal(t, "key.pem", config.KeyFile)
}

func TestInitConfig_WithEnvVars(t *testing.T) {
	// Сохраняем старые значения переменных окружения
	oldServerAddr := os.Getenv("SERVER_ADDRESS")
	oldBaseURL := os.Getenv("BASE_URL")
	oldLogLevel := os.Getenv("FLAG_LOG_LEVEL")
	oldFilePath := os.Getenv("FILE_STORAGE_PATH")
	oldDBPath := os.Getenv("DB_PATH")

	// Устанавливаем переменные окружения для проверки их приоритета
	os.Setenv("SERVER_ADDRESS", "127.0.0.1:9090")
	os.Setenv("BASE_URL", "http://127.0.0.1:9090")
	os.Setenv("FLAG_LOG_LEVEL", "debug")
	os.Setenv("FILE_STORAGE_PATH", "custom-file.json")
	os.Setenv("DB_PATH", "postgres://user:password@localhost:5432/mydb")

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
		os.Setenv("DB_PATH", oldDBPath)
	} else {
		os.Unsetenv("DB_PATH")
	}
}

func TestInitConfig_WithConfigFile(t *testing.T) {
	// Создаем временный файл с JSON-конфигурацией
	tempFile, err := os.CreateTemp("", "config-*.json")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// Записываем тестовые данные в файл
	configData := `{
		"server_address": "192.168.1.1:8080",
		"base_url": "http://192.168.1.1:8080",
		"enable_https": true,
		"cert_file": "custom-cert.pem",
		"key_file": "custom-key.pem"
	}`
	_, err = tempFile.WriteString(configData)
	assert.NoError(t, err)
	tempFile.Close()

	// Устанавливаем переменную окружения для пути к файлу конфигурации
	os.Setenv("CONFIG", tempFile.Name())

	// Создаем новый экземпляр Config и парсим файл конфигурации
	config := InitConfig()

	assert.Equal(t, "192.168.1.1:8080", config.ServerAddr)
	assert.Equal(t, "http://192.168.1.1:8080", config.BaseURL)
	assert.Equal(t, true, config.EnableHTTPS)
	assert.Equal(t, "custom-cert.pem", config.CertFile)
	assert.Equal(t, "custom-key.pem", config.KeyFile)

	// Удаляем переменную окружения
	os.Unsetenv("CONFIG")
}
