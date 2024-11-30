package main

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/app"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
)

func TestInitializeAndStartApp(t *testing.T) {
	// Создаем фиктивную конфигурацию и хранилище
	mockConfig := &config.Config{
		FilePath:   "test_file_path",
		ServerAddr: "localhost:8080",
		BaseURL:    "http://localhost",
		LogLevel:   "info",
	}
	mockStorage := storage.NewStorage()

	// Создаем экземпляр приложения
	appInstance := app.NewApp(mockStorage, mockConfig)

	// Запускаем приложение в отдельной горутине
	go func() {
		if err := appInstance.Start(); err != nil {
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
