package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// TestMainPageHandler тестирует обработку главной страницы
func TestMainPageHandler(t *testing.T) {
	storage := NewURLStorage()
	handler := mainPage("http://localhost:8080", storage)

	t.Run("GET запрос", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "Ожидался статус 200 OK")
		assert.Contains(t, rr.Body.String(), `<form action="/" method="post">`, "Ожидалось, что HTML-форма будет присутствовать в ответе")
	})

	t.Run("POST запрос с корректным URL", func(t *testing.T) {
		form := strings.NewReader("url=https://example.com")
		req, err := http.NewRequest(http.MethodPost, "/", form)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code, "Ожидался статус 201 Created")
		assert.Contains(t, rr.Body.String(), "http://localhost:8080/", "Ожидалось, что сокращенный URL будет присутствовать в ответе")
	})

	t.Run("POST запрос с пустым URL", func(t *testing.T) {
		form := strings.NewReader("url=")
		req, err := http.NewRequest(http.MethodPost, "/", form)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code, "Ожидался статус 400 Bad Request")
		assert.Equal(t, "URL не может быть пустым\n", rr.Body.String(), "Ожидалось сообщение 'URL не может быть пустым'")
	})
}

// TestRedirectHandler тестирует обработку перенаправлений
func TestRedirectHandler(t *testing.T) {
	storage := NewURLStorage()
	shortURL, _, _ := ShortenURL("https://example.com", storage)
	handler := redirectHandler(storage)

	t.Run("Перенаправление на оригинальный URL", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/"+shortURL, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/{id}", handler).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusTemporaryRedirect, rr.Code, "Ожидался статус 307 Temporary Redirect")
		assert.Equal(t, "https://example.com", rr.Header().Get("Location"), "Ожидалось перенаправление на https://example.com")
	})

	t.Run("Короткий URL не найден", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/nonexistent", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/{id}", handler).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code, "Ожидался статус 404 Not Found")
	})
}
