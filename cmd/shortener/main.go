package main

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "io"
    "net/http"
    "github.com/gorilla/mux"
	"github.com/Renal37/musthave_shortener_tpl.git/config"
)

// form - HTML-форма для ввода пользователем URL.
const form = `<html>
    <head>
    <title></title>
    </head>
    <body>
        <form action="/" method="post">
            <label>URL <input type="text" name="url"></label>
            <input type="submit" value="Login">
        </form>
    </body>
</html>`

// ShortenURL принимает URL в качестве входных данных и возвращает сокращенную версию.
func ShortenURL(url string) string {
    hasher := sha256.New()
    hasher.Write([]byte(url))
    return hex.EncodeToString(hasher.Sum(nil))[:8]
}

var originalURLs = map[string]string{
    "EwHXdJfB": "https://example.com/original-url",
}

// mainPage обрабатывает HTTP-запросы для главной страницы и нового эндпоинта.
func mainPage(baseURL string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            url := r.FormValue("url")
            shortenedURL := ShortenURL(url)
            originalURLs[shortenedURL] = url // Сохранение оригинального URL в карту
            w.WriteHeader(http.StatusCreated)
            io.WriteString(w, fmt.Sprintf(`<p>Original URL: <a href="%s">%s</a></p>`, url, url))
            io.WriteString(w, fmt.Sprintf(`<p>Shortened URL: <a href="%s/%s">%s/%s</a></p>`, baseURL, shortenedURL, baseURL, shortenedURL))
            io.WriteString(w, form)
        } else {
            io.WriteString(w, form)
        }
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

// main-функция запускает HTTP-сервер с использованием конфигурации.
func main() {
    cfg := config.InitConfig()

    r := mux.NewRouter()
    r.HandleFunc("/", mainPage(cfg.BaseURL)).Methods(http.MethodGet, http.MethodPost)
    r.HandleFunc("/{id}", redirectHandler).Methods(http.MethodGet)

    http.Handle("/", r)
    err := http.ListenAndServe(cfg.ServerAddress, nil)
    if err != nil {
        panic(err)
    }
}
