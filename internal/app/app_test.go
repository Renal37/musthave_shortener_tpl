package app

import (
	"testing"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
)

// Тест для метода Start
func TestApp_Start(t *testing.T) {
	// Настройка конфигурации и хранилища для теста
	mockStorage := &storage.Storage{}
	mockConfig := &config.Config{
		DBPath:     "test_db_path",
		FilePath:   "test_file_path",
		ServerAddr: "localhost:8080",
		BaseURL:    "http://localhost",
		LogLevel:   "info",
	}

	// Создаем экземпляр App
	app := NewApp(mockStorage, mockConfig)

	// Вызываем метод Start и проверяем, что он не вызывает паники
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Start вызвал панику: %v", r)
		}
	}()

	app.Start()
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
