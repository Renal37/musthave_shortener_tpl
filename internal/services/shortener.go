package services

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

// Store определяет интерфейс для работы с хранилищем, использующимся для взаимодействия с базой данных.
type Store interface {
	PingStore() error
	Create(originalURL, shortURL string) error
	Get(shortID string, originalURL string) (string, error)
	GetUserURLsByAuth(authHeader string) ([]ResponseBodyURLs, error)
}

// Repository определяет интерфейс для локального хранилища коротких URL.
type Repository interface {
	Set(shortID string, originalURL string)
	Get(shortID string) (string, bool)
}

// ResponseBodyURLs определяет структуру данных ответа для URL.
type ResponseBodyURLs struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
	OriginalURL   string `json:"original_url"`
}

// ShortenerService представляет собой сервис для сокращения URL.
type ShortenerService struct {
	BaseURL   string
	Storage   Repository
	db        Store
	dbDNSTurn bool
}

// NewShortenerService создает новый экземпляр ShortenerService.
func NewShortenerService(BaseURL string, storage Repository, db Store, dbDNSTurn bool) *ShortenerService {
	return &ShortenerService{
		BaseURL:   BaseURL,
		Storage:   storage,
		db:        db,
		dbDNSTurn: dbDNSTurn,
	}
}

func (s *ShortenerService) GetUserURLs(authHeader string) ([]ResponseBodyURLs, error) {
    if authHeader == "" {
        return nil, fmt.Errorf("missing auth header")
    }

    if s.dbDNSTurn {
        urls, err := s.db.GetUserURLsByAuth(authHeader)
        if err != nil {
            return nil, err
        }
        return urls, nil
    }

    return []ResponseBodyURLs{}, nil
}



// GetExistURL возвращает существующий короткий URL, если оригинальный URL уже есть в базе данных.
func (s *ShortenerService) GetExistURL(originalURL string, err error) (string, error) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		shortID, err := s.GetRep("", originalURL)
		shortURL := fmt.Sprintf("%s/%s", s.BaseURL, shortID)
		return shortURL, err
	}
	return "", err
}

// Set создает новый короткий URL для заданного оригинального URL.
func (s *ShortenerService) Set(originalURL string) (string, error) {
	shortID := randSeq()
	if s.dbDNSTurn {
		err := s.CreateRep(originalURL, shortID)
		if err != nil {
			return "", err
		}
	} else {
		s.Storage.Set(shortID, originalURL)
	}
	shortURL := fmt.Sprintf("%s/%s", s.BaseURL, shortID)
	return shortURL, nil
}

// randSeq генерирует случайную строку, используемую в качестве короткого ID для URL.
func randSeq() string {
	newUUID := uuid.New()
	return newUUID.String()
}

// Get возвращает оригинальный URL по заданному короткому ID.
func (s *ShortenerService) Get(shortID string) (string, bool) {
	if s.dbDNSTurn {
		originalURL, err := s.GetRep(shortID, "")
		if err != nil {
			return "", false
		}
		return originalURL, true
	}
	return s.Storage.Get(shortID)
}

// Ping проверяет доступность базы данных, возвращая ошибку при проблемах.
func (s *ShortenerService) Ping() error {
	return s.db.PingStore()
}

// CreateRep сохраняет оригинальный URL и его короткий аналог в базу данных.
func (s *ShortenerService) CreateRep(originalURL, shortURL string) error {
	return s.db.Create(originalURL, shortURL)
}

// GetRep получает оригинальный URL по его короткому аналогу из базы данных.
func (s *ShortenerService) GetRep(shortURL, originalURL string) (string, error) {
	return s.db.Get(shortURL, originalURL)
}
