package services

import (
    "fmt"
    "github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
    "math/rand"
)

type ShortenerService struct {
    BaseURL string
    Storage *storage.Storage
}

func NewShortenerService(BaseURL string, storage *storage.Storage) *ShortenerService {
    s := &ShortenerService{
        BaseURL: BaseURL,
        Storage: storage,
    }
    return s
}

// Функция GetShortURL - генерирует и возвращает короткую ссылку для переданной оригинальной ссылки
func (s *ShortenerService) GetShortURL(originalURL string) string {
    shortID := randSeq(8)
    s.Storage.Set(shortID, originalURL)
    shortURL := fmt.Sprintf("%s/%s", s.BaseURL, shortID)
    return shortURL
}

// Функция randSeq - генерирует случайную последовательность из заданного количества символов
func randSeq(n int) string {
    letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

// Функция GetOriginalURL - возвращает оригинальную ссылку по короткой ссылке
func (s *ShortenerService) GetOriginalURL(shortID string) (string, bool) {
    return s.Storage.Get(shortID)
}