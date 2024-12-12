package app_test

import (
	"github.com/Renal37/musthave_shortener_tpl.git/internal/app"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/stretchr/testify/assert"
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
