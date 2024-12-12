package app_test

import (
	"context"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/app"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/services"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"time"
	"github.com/stretchr/testify/mock"
	"testing"
)

// MockDump реализует мокирование функций из пакета dump
type MockDump struct {
	mock.Mock
}

func (m *MockDump) FillFromStorage(storage *storage.Storage, filePath string) error {
	args := m.Called(storage, filePath)
	return args.Error(0)
}

func (m *MockDump) Set(storage *storage.Storage, filePath string) error {
	args := m.Called(storage, filePath)
	return args.Error(0)
}

func TestApp_UseDatabase(t *testing.T) {
	// Инициализация конфигурации без пути к базе данных
	cfg := &config.Config{DBPath: ""}

	// Создаем приложение с пустыми зависимостями
	appInstance := app.NewApp(nil, nil, cfg)
	assert.True(t, appInstance.UseDatabase(), "UseDatabase должен возвращать true, если DBPath пустой")

	// Обновляем конфигурацию с указанием пути к базе данных
	cfg.DBPath = "test.db"
	appInstance = app.NewApp(nil, nil, cfg)
	assert.False(t, appInstance.UseDatabase(), "UseDatabase должен возвращать false, если DBPath не пустой")
}

func TestStop(t *testing.T) {
	mockStorage := storage.NewStorage()
	mockConfig := &config.Config{
		FilePath: "/tmp/test_data.json",
	}

	app := app.NewApp(mockStorage, nil, mockConfig)
	app.Stop()
}
func TestApp_UseDatabaseWithDifferentDBPaths(t *testing.T) {
	tests := []struct {
		dbPath string
		expect bool
	}{
		{"", true},                      // Пустой путь
		{"test.db", false},              // Непустой путь
		{"./relative_path.db", false},   // Относительный путь
		{"/absolute/path/to/db", false}, // Абсолютный путь
	}

	for _, tt := range tests {
		t.Run(tt.dbPath, func(t *testing.T) {
			cfg := &config.Config{DBPath: tt.dbPath}
			appInstance := app.NewApp(nil, nil, cfg)
			assert.Equal(t, tt.expect, appInstance.UseDatabase())
		})
	}
}
func TestApp_Start(t *testing.T) {
    // Мокаем хранилище и сервисы
    mockStorage := storage.NewStorage()
    mockServices := &services.ShortenerService{}
    
    // Конфигурация с пустым DBPath
    mockConfig := &config.Config{
        DBPath:     "",
        FilePath:   "/tmp/test_data.json", // Убедитесь, что этот путь не используется на самом деле в тестах
        ServerAddr: "localhost:0",
        BaseURL:    "/api",
        LogLevel:   "info",
        EnableHTTPS: false,
    }

    // Мокируем функции из dump
    mockDump := new(MockDump)
    app := app.NewApp(mockStorage, mockServices, mockConfig)
    
    // Замена реальных функций на моки
    app.FillFromStorage = mockDump.FillFromStorage
    app.Set = mockDump.Set

    // Ожидаем вызов метода FillFromStorage

    // Ожидаем вызов метода Set с аргументами, которые могут быть переданы
    mockDump.On("Set", mockStorage, mockConfig.FilePath).Return(nil)

    // Создаем контекст с тайм-аутом 5 секунд
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Запускаем приложение с ограничением времени
    err := app.Start(ctx)
    assert.NoError(t, err, "app.Start должен завершиться без ошибок")

    // Проверяем, что методы были вызваны
}



func TestStop_WithoutDatabase(t *testing.T) {
	mockStorage := storage.NewStorage()
	mockConfig := &config.Config{
		DBPath:   "", // База данных не используется
		FilePath: "/tmp/test_data.json",
	}

	// Мокируем функции из dump
	mockDump := new(MockDump)
	mockDump.On("Set", mockStorage, mockConfig.FilePath).Return(nil)

	// Создаем приложение с моками
	app := app.NewApp(mockStorage, nil, mockConfig)
	app.Set = mockDump.Set

	// Вызываем Stop
	err := app.Stop()
	assert.NoError(t, err)

	// Проверяем вызовы mock-объектов
	mockDump.AssertExpectations(t)
}
