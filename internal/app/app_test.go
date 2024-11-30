package app

import (
	"os"
	"fmt"
	"syscall"
	"testing"
	"time"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
)

// Моковая структура для базы данных
type MockDatabase struct{}

func (m *MockDatabase) InitDatabase(dbPath string) error {
	if dbPath == "invalid_db_path" {
		return fmt.Errorf("Ошибка инициализации базы данных")
	}
	return nil
}

// Тест для создания нового приложения
func TestNewApp(t *testing.T) {
	storageInstance := &storage.Storage{}
	configInstance := &config.Config{
		DBPath:      "test_db_path",
		FilePath:    "test_file_path",
		ServerAddr:  "localhost:8080",
		BaseURL:     "http://localhost",
		LogLevel:    "debug",
		EnableHTTPS: false,
		CertFile:    "",
		KeyFile:     "",
	}

	app := NewApp(storageInstance, configInstance)

	if app.storageInstance != storageInstance {
		t.Errorf("Expected storageInstance to be %v, got %v", storageInstance, app.storageInstance)
	}

	if app.config != configInstance {
		t.Errorf("Expected config to be %v, got %v", configInstance, app.config)
	}
}

// Тест для метода Start
func TestApp_Start(t *testing.T) {
	mockStorage := &storage.Storage{}
	mockConfig := &config.Config{
		DBPath:     "test_db_path",
		FilePath:   "test_file_path",
		ServerAddr: "localhost:8080",
		BaseURL:    "http://localhost",
		LogLevel:   "info",
	}

	app := NewApp(mockStorage, mockConfig)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Start вызвал панику: %v", r)
		}
	}()

	app.Start()
}

// Тест для метода Stop
func TestApp_Stop(t *testing.T) {
	mockStorage := &storage.Storage{}
	mockConfig := &config.Config{
		DBPath:   "test_db_path",
		FilePath: "test_file_path",
	}

	app := NewApp(mockStorage, mockConfig)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Stop вызвал панику: %v", r)
		}
	}()

	app.Stop()
}

// Тест для метода Start с ошибкой инициализации базы данных
func TestApp_Start_ErrorInitDatabase(t *testing.T) {
	mockConfig := &config.Config{
		DBPath:     "invalid_db_path",
		FilePath:   "test_file_path",
		ServerAddr: "localhost:8080",
		BaseURL:    "http://localhost",
		LogLevel:   "info",
	}

	mockStorage := storage.NewStorage()

	app := NewApp(mockStorage, mockConfig)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Start вызвал панику: %v", r)
		}
	}()

	err := app.Start()
	if err == nil {
		t.Errorf("Ожидалась ошибка при инициализации базы данных, но её не было")
	}
}

// Тест для метода UseDatabase
func TestApp_UseDatabase(t *testing.T) {
	tests := []struct {
		name       string
		dbPath     string
		wantResult bool
	}{
		{
			name:       "Database Path Provided",
			dbPath:     "some/db/path",
			wantResult: false,
		},
		{
			name:       "No Database Path",
			dbPath:     "",
			wantResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockConfig := &config.Config{DBPath: tt.dbPath}
			app := NewApp(nil, mockConfig)

			got := app.UseDatabase()
			if got != tt.wantResult {
				t.Errorf("UseDatabase() = %v, want %v", got, tt.wantResult)
			}
		})
	}
}

// Тест обработки сигналов завершения
func TestApp_SignalHandling(t *testing.T) {
	mockStorage := &storage.Storage{}
	mockConfig := &config.Config{
		DBPath:     "test_db_path",
		FilePath:   "test_file_path",
		ServerAddr: "localhost:8080",
		BaseURL:    "http://localhost",
		LogLevel:   "info",
	}

	app := NewApp(mockStorage, mockConfig)

	signalChan := make(chan os.Signal, 1)

	go func() {
		err := app.Start()
		if err != nil {
			t.Errorf("Start вернул ошибку: %v", err)
		}
	}()

	signalChan <- syscall.SIGINT

	select {
	case <-signalChan:
	case <-time.After(2 * time.Second):
		t.Error("Сигнал завершения не был обработан")
	}
}
