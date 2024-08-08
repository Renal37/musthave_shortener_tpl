package services

import (
	"fmt"
	"math/rand"
)

type Store interface {
	PingStore() error
}
type Repository interface {
	Set(shortID string, originalURL string)
	Get(shortID string) (string, bool)
}
type ShortenerService struct {
	BaseURL string
	Storage Repository
	bd      Store
}

func NewShortenerService(BaseURL string, storage Repository, bd Store) *ShortenerService {
	s := &ShortenerService{
		BaseURL: BaseURL,
		Storage: storage,
		bd:      bd,
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
func (s *ShortenerService) Ping() error {
	return s.bd.PingStore()
}
