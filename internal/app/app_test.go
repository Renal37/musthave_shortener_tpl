package app

import (
	"context"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"os"
	"syscall"
	"testing"
	"time"
)

// TestNewApp проверяет создание нового экземпляра приложения
func TestNewApp(t *testing.T) {
	mockStorage := storage.NewStorage() // Предположим, есть функция создания пустого хранилища
	mockConfig := &config.Config{
		FilePath:    "/tmp/test_data.json",
		DBPath:      "",
		ServerAddr:  ":8080",
		BaseURL:     "http://localhost:8080",
		LogLevel:    "debug",
		EnableHTTPS: false,
		CertFile:    "",
		KeyFile:     "",
	}

	app := NewApp(mockStorage, mockConfig)
	assert.NotNil(t, app)
	assert.Equal(t, mockStorage, app.storageInstance)
	assert.Equal(t, mockConfig, app.config)
}

// TestUseDatabase проверяет функцию UseDatabase
func TestUseDatabase(t *testing.T) {
	mockStorage := storage.NewStorage()
	mockConfig := &config.Config{
		DBPath: "",
	}

	app := NewApp(mockStorage, mockConfig)
	assert.True(t, app.UseDatabase())

	mockConfig.DBPath = "/path/to/db"
	app = NewApp(mockStorage, mockConfig)
	assert.False(t, app.UseDatabase())
}
func TestStart(t *testing.T) {
	// Создаем фиктивную конфигурацию и хранилище
	mockConfig := &config.Config{
		DBPath:     "",
		FilePath:   "test_file_path",
		ServerAddr: "localhost:8080",
		BaseURL:    "http://localhost",
		LogLevel:   "info",
	}
	mockStorage := storage.NewStorage()

	// Создаем экземпляр приложения
	app := NewApp(mockStorage, mockConfig)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Запускаем приложение в отдельной горутине
	go func() {
		if err := app.Start(ctx); err != nil {
			t.Errorf("Ошибка при запуске приложения: %v", err)
		}
	}()

	// Даем серверу время для запуска
	time.Sleep(2 * time.Second)

	// Отправляем сигнал завершения
	process, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("Не удалось найти процесс: %v", err)
	}

	if err := process.Signal(syscall.SIGINT); err != nil {
		t.Fatalf("Не удалось отправить сигнал завершения: %v", err)
	}

	// Даем серверу время для завершения
	time.Sleep(2 * time.Second)
}

// TestStop проверяет остановку приложения
func TestStop(t *testing.T) {
	mockStorage := storage.NewStorage()
	mockConfig := &config.Config{
		FilePath: "/tmp/test_data.json",
	}

	app := NewApp(mockStorage, mockConfig)
	app.Stop()
}
