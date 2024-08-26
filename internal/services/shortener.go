package services

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/logger"
	"go.uber.org/zap"
)

type Store interface {
	PingStore() error
	Create(originalURL, shortURL, UserID string) error
	Get(shortIrl string, originalURL string) (string, error)
	GetFull(userID string, BaseURL string) ([]map[string]string, error)
	DeleteURLs(userID string, shortURL string, updateChan chan<- string) error
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

func (s *ShortenerService) GetExistURL(originalURL string, err error) (string, error) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		shortID, err := s.GetRep("", originalURL)
		shortURL := fmt.Sprintf("%s/%s", s.BaseURL, shortID)
		return shortURL, err
	}
	return "", err
}

func (s *ShortenerService) Set(userID, originalURL string) (string, error) {
	shortID := randSeq()
	if s.dbDNSTurn {
		err := s.CreateRep(originalURL, shortID, userID)
		if err != nil {
			return "", err
		}
	} else {
		s.Storage.Set(shortID, originalURL)
	}
	shortURL := fmt.Sprintf("%s/%s", s.BaseURL, shortID)
	return shortURL, nil
}

func randSeq() string {
	newUUID := uuid.New()
	return newUUID.String()
}

func (s *ShortenerService) Get(shortID string) (string, error) {
	if s.dbDNSTurn {
		originalURL, err := s.GetRep(shortID, "")
		if err != nil {
			return "", err
		}
		return originalURL, nil
	}

	originalURL, exists := s.Storage.Get(shortID)
	if !exists {
		err := errors.New("failed get original url")
		return "", err
	}
	return originalURL, nil
}

func (s *ShortenerService) Ping() error {
	return s.db.PingStore()
}

func (s *ShortenerService) CreateRep(originalURL, shortURL, UserID string) error {
	return s.db.Create(originalURL, shortURL, UserID)
}

func (s *ShortenerService) GetRep(shortURL, originalURL string) (string, error) {
	return s.db.Get(shortURL, originalURL)
}

func (s *ShortenerService) GetFullRep(userID string) ([]map[string]string, error) {
	return s.db.GetFull(userID, s.BaseURL)
}

func (s *ShortenerService) DeleteURLsRep(userID string, shorURLs []string) error {
	resultChan := make(chan string)
	updateChan := make(chan string, len(shorURLs))

	go func() {
		for _, shortURL := range shorURLs {
			err := s.db.DeleteURLs(userID, shortURL, updateChan)
			if err != nil {
				logger.Log.Error("Failed to delete URLs", zap.Error(err))
			}
		}
	}()

	go func() {
		for updateShortID := range updateChan {
			resultChan <- updateShortID
		}
		close(resultChan)
	}()
	return nil
}

func (s *ShortenerService) GetDeletedFlagType() {

}