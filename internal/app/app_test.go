package app

import (
	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"testing"
)

// Мок для конфигурации
type MockConfig struct {
	DBPath     string
	FilePath   string
	ServerAddr string
	BaseURL    string
	LogLevel   string
}

func TestNewApp(t *testing.T) {
	// Создаем экземпляры хранилища и конфигурации для теста
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

	// Создаем новое приложение
	app := NewApp(storageInstance, configInstance)

	// Проверяем, что приложение было инициализировано корректно
	if app.storageInstance != storageInstance {
		t.Errorf("Expected storageInstance to be %v, got %v", storageInstance, app.storageInstance)
	}

	if app.config != configInstance {
		t.Errorf("Expected config to be %v, got %v", configInstance, app.config)
	}
}

// Тест для метода Start с реальным хранилищем
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

// Тест для метода Start с ошибкой инициализации базы данных
func TestApp_Start_ErrorInitDatabase(t *testing.T) {
	// Настроим конфигурацию для теста
	mockConfig := &config.Config{
		DBPath:     "invalid_db_path", // Ошибочный путь для базы данных
		FilePath:   "test_file_path",
		ServerAddr: "localhost:8080",
		BaseURL:    "http://localhost",
		LogLevel:   "info",
	}

	// Создаем экземпляр хранилища
	mockStorage := storage.NewStorage()

	// Создаем экземпляр App
	app := NewApp(mockStorage, mockConfig)

	// Проверка того, что метод Start не вызывает панику
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Start вызвал панику: %v", r)
		}
	}()

	// Запуск приложения и проверка ошибки
	err := app.Start()
	if err == nil {
		t.Errorf("Ожидалась ошибка при инициализации базы данных, но её не было")
	}
}

// Тест для метода Stop с реальным хранилищем
func TestApp_Stop(t *testing.T) {
	// Настройка конфигурации для теста
	mockConfig := &config.Config{
		DBPath:     "test_db_path",
		FilePath:   "test_file_path",
		ServerAddr: "localhost:8080",
		BaseURL:    "http://localhost",
		LogLevel:   "info",
	}

	// Создаем экземпляр хранилища через storage.NewStorage()
	mockStorage := storage.NewStorage()

	// Создаем экземпляр App
	app := NewApp(mockStorage, mockConfig)

	// Проверка того, что метод Stop не вызывает панику
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Stop вызвал панику: %v", r)
		}
	}()

	// Запуск и остановка приложения
	app.Stop()
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
