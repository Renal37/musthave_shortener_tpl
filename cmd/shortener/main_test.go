package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/app"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
)

func TestMain(t *testing.T) {
	// Создаем новый экземпляр хранилища
	storage := storage.NewStorage()

	// Создаем новый экземпляр приложения
	appInstance := app.NewApp(storage, &config.Config{})

	// Создаем новый HTTP сервер для обработки запросов
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			if r.URL.Path == "/start" {
				appInstance.Start() // Предполагается, что Start не блокирует
				http.Error(w, "Application started", http.StatusOK)
			} else if r.URL.Path == "/stop" {
				appInstance.Stop() // Предполагается, что Stop не блокирует
				http.Error(w, "Application stopped", http.StatusOK)
			} else {
				http.Error(w, "Not Found", http.StatusNotFound)
			}
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}))
	defer server.Close()

	// Отправляем запрос для запуска приложения
	resp, err := http.Post(server.URL+"/start", "application/json", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем, что код состояния ответа 200
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got: %d", resp.StatusCode)
	}

	// Отправляем запрос для остановки приложения
	resp, err = http.Post(server.URL+"/stop", "application/json", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем, что код состояния ответа 200
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got: %d", resp.StatusCode)
	}
}
