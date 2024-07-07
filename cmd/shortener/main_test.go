package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/Renal37/musthave_shortener_tpl.git/config"
	"github.com/gorilla/mux"
)

// setupRouter создает роутер для тестирования
func setupRouter(cfg *config.Config) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", mainPage(cfg.BaseURL)).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/{id}", redirectHandler).Methods(http.MethodGet)
	return r
}

// TestMainPageGet тестирует GET-запрос к главной странице.
func TestMainPageGet(t *testing.T) {
	cfg := &config.Config{BaseURL: "http://localhost:8080"}
	r := setupRouter(cfg)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус OK; получен %v", resp.StatusCode)
	}
	if !strings.Contains(string(body), "<form") {
		t.Errorf("Ожидалась форма в теле ответа; получено %v", string(body))
	}
}

// TestMainPagePost тестирует POST-запрос к главной странице.
func TestMainPagePost(t *testing.T) {
	cfg := &config.Config{BaseURL: "http://localhost:8080"}
	r := setupRouter(cfg)

	formData := "url=https://example.com"
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(formData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Ожидался статус Created; получен %v", resp.StatusCode)
	}
	if !strings.Contains(string(body), "Shortened URL") {
		t.Errorf("Ожидалось наличие Shortened URL в теле ответа; получено %v", string(body))
	}
	if !strings.Contains(string(body), "Original URL") {
		t.Errorf("Ожидалось наличие Original URL в теле ответа; получено %v", string(body))
	}
}

// TestMainPageRedirect тестирует перенаправление с сокращенного URL на оригинальный URL.
func TestMainPageRedirect(t *testing.T) {
	cfg := &config.Config{BaseURL: "http://localhost:8080"}
	r := setupRouter(cfg)

	req := httptest.NewRequest(http.MethodGet, "/EwHXdJfB", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusTemporaryRedirect {
		t.Errorf("Ожидался статус TemporaryRedirect; получен %v", resp.StatusCode)
	}
	if location := resp.Header.Get("Location"); location != "https://example.com/original-url" {
		t.Errorf("Ожидалось перенаправление на https://example.com/original-url; получено %v", location)
	}
}

// TestMainPageInvalidShortURL тестирует ответ на некорректный сокращенный URL.
func TestMainPageInvalidShortURL(t *testing.T) {
	cfg := &config.Config{BaseURL: "http://localhost:8080"}
	r := setupRouter(cfg)

	req := httptest.NewRequest(http.MethodGet, "/invalidURL", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Ожидался статус BadRequest; получен %v", resp.StatusCode)
	}
	if body, _ := ioutil.ReadAll(resp.Body); !strings.Contains(string(body), "Invalid shortened URL") {
		t.Errorf("Ожидалось сообщение об ошибке в теле ответа; получено %v", string(body))
	}
}
