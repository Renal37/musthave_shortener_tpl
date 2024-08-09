package services

import (
	"fmt"
	"math/rand"
)

type Store interface {
	PingStore() error
	Create(originalURL, shortURL string) error
	Get(shortIrl string) (string, error)
}
type Repository interface {
	Set(shortID string, originalURL string)
	Get(shortID string) (string, bool)
}
type ShortenerService struct {
	BaseURL   string
	Storage   Repository
	db        Store
	dbDNSTurn bool
}

func NewShortenerService(BaseURL string, storage Repository, db Store, dbDNSTurn bool) *ShortenerService {
	s := &ShortenerService{
		BaseURL:   BaseURL,
		Storage:   storage,
		db:        db,
		dbDNSTurn: dbDNSTurn,
	}
	return s
}

// Функция GetShortURL - генерирует и возвращает короткую ссылку для переданной оригинальной ссылки
func (s *ShortenerService) GetShortURL(originalURL string) string {
	shortID := randSeq(8)
	if s.dbDNSTurn {
		err := s.CreateRep(originalURL, shortID)
		if err != nil {
			return ""
		}
	} else {
		s.Storage.Set(shortID, originalURL)
	}
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

// Функция Get - возвращает оригинальную ссылку по короткой ссылке
func (s *ShortenerService) Get(shortID string) (string, bool) {
	if s.dbDNSTurn {
		originalURL, err := s.GetRep(shortID)
		if err != nil {
			return "", false
		}
		return originalURL, true
	}

	return s.Storage.Get(shortID)
}
func (s *ShortenerService) Ping() error {
	return s.db.PingStore()
}

func (s *ShortenerService) CreateRep(originalURL, shortURL string) error {
	return s.db.Create(originalURL, shortURL)
}

func (s *ShortenerService) GetRep(shortURL string) (string, error) {
	return s.db.Get(shortURL)
}
