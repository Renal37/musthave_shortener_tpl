package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"github.com/gorilla/mux"
)

// form - HTML-форма для ввода пользователем URL.
const form = `<html>
    <head>
    <title></title>
    </head>
    <body>
        <form action="/" method="post">
            <label>URl <input type="text" name="url"></label>
            <input type="submit" value="Login">
        </form>
    </body>
</html>`

// ShortenURL принимает URL в качестве входных данных и возвращает сокращенную версию.
// Она использует алгоритм хеширования SHA256 для создания уникального хеша для URL.
// Возвращаются первые 8 символов хеша как сокращенный URL.
func ShortenURL(url string) string {
	hasher := sha256.New()
	hasher.Write([]byte(url))
	return hex.EncodeToString(hasher.Sum(nil))[:8]
}

var originalURLs = map[string]string{
	"1395ec37":"https://vk.com",
	"3c0a9a5c":"https://practicum.yandex.ru/profile/go-advanced/",
}

// mainPage обрабатывает HTTP-запросы для главной страницы и нового эндпоинта.
// Если метод запроса POST, он считывает URL из формы,
// сокращает его с помощью функции ShortenURL и записывает сокращенный URL в ответ.
// Если метод запроса не POST, он записывает HTML-форму в ответ.
func mainPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		url := r.FormValue("url")
		shortenedURL := ShortenURL(url)
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, fmt.Sprintf("<p>Shortened URL: %s</p>", shortenedURL))
		io.WriteString(w, form)
	} else {
		io.WriteString(w, form)
	}
}

// redirectHandler обрабатывает перенаправления с сокращенного URL на оригинальный URL.
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortenedURL := vars["id"]
	originalURL, ok := originalURLs[shortenedURL]
	if !ok {
		http.Error(w, "Invalid shortened URL", http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// main-функция запускает HTTP-сервер на порту 8080 и ожидает входящих запросов.
// Она использует пакет gorilla/mux для маршрутизации.
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", mainPage).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/{id}", redirectHandler).Methods(http.MethodGet)

	http.Handle("/", r)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
