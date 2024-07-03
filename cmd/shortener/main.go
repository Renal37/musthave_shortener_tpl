package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
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

// mainPage обрабатывает HTTP-запросы для главной страницы и нового эндпоинта.
// Если метод запроса POST, он считывает URL из формы,
// сокращает его с помощью функции ShortenURL и записывает сокращенный URL в ответ.
// Если метод запроса не POST, он записывает HTML-форму в ответ.
// Если путь запроса начинается с "/{id}", он вызывает функцию redirectToOriginalURL для перенаправления пользователя на оригинальный URL.
func getOriginalURL(shortenedURL string) (string, error) {
	originalURLs := map[string]string{
		"EwHXdJfB": "https://example.com/original-url",
		// Add more mappings as needed
	}

	originalURL, ok := originalURLs[shortenedURL]
	if !ok {
		return "", fmt.Errorf("invalid shortened URL: %s", shortenedURL)
	}

	return originalURL, nil
}
func mainPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		url := r.FormValue("url")
		shortenedURL := ShortenURL(url)
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, fmt.Sprintf("<p>Shortened URL: <a href=\"%s\">%s</a></p>", shortenedURL, shortenedURL))
		io.WriteString(w, form)
	} else if len(r.URL.Path) > 1 && r.URL.Path[0] == '/' {
		shortenedURL := r.URL.Path[1:]
		originalURL, err := getOriginalURL(shortenedURL)
		if err != nil {
			http.Error(w, "Invalid shortened URL", http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		io.WriteString(w, form)
	}
}

// main-функция запускает HTTP-сервер на порту 8080 и ожидает входящих запросов.
// Она использует функцию mainPage как обработчик для запросов.
func main() {
	err := http.ListenAndServe(`:8080`, http.HandlerFunc(mainPage))
	if err != nil {
		panic(err)
	}
}
