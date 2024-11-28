package services

import (
	"errors"
	"fmt"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"go.uber.org/zap"
	"strings"
)

type Store interface {
	PingStore() error
	Create(originalURL, shortURL, UserID string) error
	Get(shortID string, originalURL string) (string, error)
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
	if !strings.HasPrefix(BaseURL, "http://") && !strings.HasPrefix(BaseURL, "https://") {
		logger.Log.Error("Invalid BaseURL, must start with http:// or https://")
		return nil
	}
	return &ShortenerService{
		BaseURL:   BaseURL,
		Storage:   storage,
		db:        db,
		dbDNSTurn: dbDNSTurn,
	}
}

func (s *ShortenerService) GetExistURL(originalURL string, err error) (string, error) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		shortID, err := s.GetRep("", originalURL)
		if err != nil {
			return "", err
		}
		shortURL := fmt.Sprintf("%s/%s", s.BaseURL, shortID)
		return shortURL, nil
	}
	return "", err
}

func (s *ShortenerService) Set(userID, originalURL string) (string, error) {
	if s.db == nil && s.Storage == nil {
		return "", errors.New("database and storage are not initialized")
	}

	shortID := randSeq()
	if s.dbDNSTurn {
		if err := s.CreateRep(originalURL, shortID, userID); err != nil {
			return "", err
		}
	} else {
		s.Storage.Set(shortID, originalURL)
	}
	shortURL := fmt.Sprintf("%s/%s", s.BaseURL, shortID)
	return shortURL, nil
}

func randSeq() string {
	return uuid.New().String()
}

func (s *ShortenerService) Get(shortID string) (string, error) {
	if s.dbDNSTurn {
		return s.GetRep(shortID, "")
	}
	originalURL, exists := s.Storage.Get(shortID)
	if !exists {
		return "", errors.New("не удалось получить оригинальную ссылку")
	}
	return originalURL, nil
}

func (s *ShortenerService) Ping() error {
	if s.db == nil {
		return errors.New("database not initialized")
	}
	return s.db.PingStore()
}

func (s *ShortenerService) CreateRep(originalURL, shortURL, UserID string) error {
	if s.db == nil {
		return errors.New("database not initialized")
	}
	return s.db.Create(originalURL, shortURL, UserID)
}

func (s *ShortenerService) GetRep(shortURL, originalURL string) (string, error) {
	if s.db == nil {
		return "", errors.New("database not initialized")
	}
	return s.db.Get(shortURL, originalURL)
}

func (s *ShortenerService) GetFullRep(userID string) ([]map[string]string, error) {
	if s.db == nil {
		return nil, errors.New("database not initialized")
	}
	return s.db.GetFull(userID, s.BaseURL)
}

func (s *ShortenerService) DeleteURLsRep(userID string, shortURLs []string) error {
	if s.db == nil {
		return errors.New("database not initialized")
	}

	updateChan := make(chan string, len(shortURLs))
	workerChan := make(chan string, len(shortURLs))

	go func() {
		for shortURL := range workerChan {
			if err := s.db.DeleteURLs(userID, shortURL, updateChan); err != nil {
				logger.Log.Error("Не удалось удалить ссылку", zap.Error(err))
			}
		}
		close(updateChan)
	}()

	go func() {
		for _, shortURL := range shortURLs {
			workerChan <- shortURL
		}
		close(workerChan)
	}()

	return nil
}
