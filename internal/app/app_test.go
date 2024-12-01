package app

import (
	"context"
	"testing"
	"time"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/stretchr/testify/assert"
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
	mockStorage := storage.NewStorage()
	mockConfig := &config.Config{
		FilePath:   "/tmp/test_data.json",
		DBPath:     "",
		ServerAddr: ":8080",
		BaseURL:    "http://localhost:8080",
		LogLevel:   "debug",
	}

	app := NewApp(mockStorage, mockConfig)

	// Используем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- app.Start(ctx)
	}()

	select {
	case err := <-done:
		// Тест завершился до таймаута
		assert.Nil(t, err, "App.Start вернул ошибку")
	case <-time.After(3 * time.Second): // даём дополнительное время
		t.Fatal("Тест завершён по таймауту")
	}
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
