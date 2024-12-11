package app_test

import (
	"github.com/Renal37/musthave_shortener_tpl.git/internal/app"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/services"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

// MockStorage реализует интерфейс storage.Storage
// для мокирования в тестах.
type MockStorage struct {
	mock.Mock
	storage.Storage // Встраивание интерфейса для удовлетворения типа
}

// MockService реализует интерфейс services.ShortenerService
// для мокирования в тестах.
type MockService struct {
	mock.Mock
	services.ShortenerService // Встраивание интерфейса для удовлетворения типа
}

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

// func TestStop(t *testing.T) {
// 	// Создаем мок-хранилище и конфигурацию
// 	mockStorage := &MockStorage{}
// 	mockService := &MockService{}
// 	mockConfig := &config.Config{
// 		FilePath: "/tmp/test_data.json",
// 	}

// 	// Создаем приложение с мок-объектами
// 	appInstance := app.NewApp(mockStorage, mockService, mockConfig)

// 	// Проверяем вызов метода Stop
// 	appInstance.Stop()
// }
