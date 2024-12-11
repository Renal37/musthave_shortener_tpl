package app

import (
	"context"
	"testing"
	"time"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/services"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStorage is a mock implementation of the storage.Storage interface
type MockStorage struct {
	mock.Mock
	storage.Storage // Embed the actual interface to satisfy the type requirement
}

// MockService is a mock implementation of the services.ShortenerService interface
type MockService struct {
	mock.Mock
	services.ShortenerService // Embed the actual interface to satisfy the type requirement
}

// MockDump is a mock implementation of the dump package functions
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

func TestAppStartAndStop(t *testing.T) {
	// Create mock instances
	mockStorage := &MockStorage{}
	mockService := &MockService{}
	mockDump := &MockDump{}
	mockConfig := &config.Config{
		DBPath:      "", // Ensure this is empty to trigger UseDatabase() == true
		FilePath:    "/tmp/test_data.json",
		ServerAddr:  ":8080",
		BaseURL:     "http://localhost:8080",
		LogLevel:    "info",
		EnableHTTPS: false,
		CertFile:    "",
		KeyFile:     "",
	}

	// Set expectations for mockDump
	mockDump.On("FillFromStorage", mock.Anything, mockConfig.FilePath).Return(nil)
	mockDump.On("Set", mock.Anything, mockConfig.FilePath).Return(nil)

	// Create an instance of the App with mock functions
	app := &App{
		storageInstance:  &mockStorage.Storage,
		servicesInstance: &mockService.ShortenerService,
		config:           mockConfig,
		fillFromStorage:  mockDump.FillFromStorage,
		set:              mockDump.Set,
	}

	// Create a context with cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the application in a separate goroutine
	go func() {
		err := app.Start(ctx)
		assert.NoError(t, err, "App should start without error")
	}()

	// Give the server time to start
	time.Sleep(1 * time.Second)

	// Check if the application is using the database
	assert.True(t, app.UseDatabase(), "App should use the database")

	// Cancel the context to stop the application
	cancel()

	// Give the server time to stop
	time.Sleep(1 * time.Second)

	// Verify that the Stop method was called
	mockDump.AssertCalled(t, "Set", mock.Anything, mockConfig.FilePath)
}
