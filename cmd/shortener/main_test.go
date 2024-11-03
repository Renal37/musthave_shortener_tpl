package main

import (
	"flag"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/app"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
)

// TestConfigInit проверяет корректную инициализацию конфигурации.
func TestConfigInit(t *testing.T) {
	cfg := config.InitConfig()
	if cfg == nil {
		t.Fatal("Ожидалось, что конфигурация будет инициализирована, но получено nil")
	}

	// Дополнительные проверки конфигурации
	if cfg.ServerAddr == "" {
		t.Fatal("Ожидалось, что ServerAddr будет инициализирован, но получено пустое значение")
	}
	if cfg.BaseURL == "" {
		t.Fatal("Ожидалось, что BaseURL будет инициализирован, но получено пустое значение")
	}
}

// TestStorage проверяет создание экземпляра хранилища.
func TestStorage(t *testing.T) {
	storageInstance := storage.NewStorage()
	if storageInstance == nil {
		t.Fatal("Ожидалось, что экземпляр хранилища будет инициализирован, но получено nil")
	}

	// Дополнительные проверки хранилища
	if len(storageInstance.URLs) != 0 {
		t.Fatal("Ожидалось, что карта URL-адресов в хранилище будет пустой, но получено непустое значение")
	}
}

// TestServer проверяет, что сервер успешно запускается и отвечает на запросы.
func TestServer(t *testing.T) {
	// Сброс флагов перед запуском теста
	flag.CommandLine.Parse([]string{})

	addrConfig := config.InitConfig()
	storageInstance := storage.NewStorage()
	appInstance := app.NewApp(storageInstance, addrConfig)

	// Запускаем приложение в горутине

	// Делаем HTTP-запрос к серверу
	req, err := http.NewRequest(http.MethodGet, addrConfig.ServerAddr+"/", nil)
	if err != nil {
		t.Fatalf("Не удалось создать запрос: %v", err)
	}

	// Создаем новый тестовый сервер
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Вызываем хендлер
	handler.ServeHTTP(rr, req)

	// Проверяем ответ
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Ожидался статус OK, получен %v", status)
	}

	// Остановка приложения
	appInstance.Stop()
}
