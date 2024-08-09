package services

import (
	"fmt"
	"github.com/google/uuid"
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

func (s *ShortenerService) Set(originalURL string) string {
	shortID := randSeq()
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

func randSeq() string {
	newUUID := uuid.New()
	return newUUID.String()
}

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