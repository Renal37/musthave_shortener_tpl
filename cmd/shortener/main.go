package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"github.com/Renal37/musthave_shortener_tpl.git/config"
	"github.com/gorilla/mux"
)

// ShortenURL принимает URL в качестве входных данных и возвращает сокращенную версию.
func ShortenURL(url string, storage *URLStorage) (string, int, error) {
	hasher := sha256.New()
	hasher.Write([]byte(url))
	hash := hex.EncodeToString(hasher.Sum(nil))[:8]

	// Проверяем, есть ли уже такой хеш в хранилище
	if originalURL, exists := storage.originalURLs[hash]; exists {
		// Если хеш уже существует, возвращаем его и статус конфликта
		return hash, http.StatusConflict, fmt.Errorf("collision detected: %s", originalURL)
	}

	// Если хеш не найден, сохраняем URL под этим хешем
	storage.originalURLs[hash] = url
	return hash, http.StatusCreated, nil
}

type URLStorage struct {
	originalURLs map[string]string
}

func NewURLStorage() *URLStorage {
	return &URLStorage{
		originalURLs: make(map[string]string),
	}
}

func (s *URLStorage) SaveURL(shortenedURL, originalURL string) error {
	s.originalURLs[shortenedURL] = originalURL
	return nil
}

func (s *URLStorage) GetURL(shortenedURL string) (string, bool) {
	url, ok := s.originalURLs[shortenedURL]
	return url, ok
}

func mainPage(baseURL string, storage *URLStorage) http.HandlerFunc {
	const form = `<html>
        <head>
        <title></title>
        </head>
        <body>
            <form action="/" method="post">
                <label>Введите сюда URL Который хотите сократить <input type="text" name="url"></label>
                <input type="submit" value="Сократить">
            </form>
        </body>
    </html>`

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Error parsing form", http.StatusBadRequest)
				return
			}
			url := r.FormValue("url")
			if url == "" {
				http.Error(w, "URL не может быть пустым", http.StatusBadRequest)
				return
			}

			shortenedURL, status, err := ShortenURL(url, storage)
			if err != nil {
				http.Error(w, "Error creating shortened URL: "+err.Error(), http.StatusInternalServerError)
				return
			}

			if status == http.StatusConflict {
				http.Error(w, fmt.Sprintf("Хеш уже существует для: %s/%s", baseURL, shortenedURL), status)
				return
			}

			if err := storage.SaveURL(shortenedURL, url); err != nil {
				http.Error(w, "Error saving URL: "+err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(status)
			io.WriteString(w, fmt.Sprintf("%s/%s", baseURL, shortenedURL))
		} else {
			io.WriteString(w, form)
		}
	}
}


func redirectHandler(storage *URLStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		shortenedURL := vars["id"]
		originalURL, ok := storage.GetURL(shortenedURL)
		if !ok {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
	}
}

func main() {
	cfg := config.InitConfig()

	storage := NewURLStorage()

	r := mux.NewRouter()
	r.HandleFunc("/", mainPage(cfg.BaseURL, storage)).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/{id}", redirectHandler(storage)).Methods(http.MethodGet)

	err := http.ListenAndServe(cfg.ServerAddress, r)
	if err != nil {
		log.Fatal(err)
	}
}
