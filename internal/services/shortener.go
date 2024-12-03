package services

import (
	"errors"
	"fmt"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"go.uber.org/zap"
)

// Store определяет интерфейс взаимодействия с хранилищем URL.
type Store interface {
	PingStore() error                                                          // Проверяет соединение с хранилищем
	Create(originalURL, shortURL, UserID string) error                         // Создаёт новую запись URL
	Get(shortID string, originalURL string) (string, error)                    // Извлекает оригинальный URL по сокращенному
	GetFull(userID string, BaseURL string) ([]map[string]string, error)        // Извлекает все URL пользователя
	DeleteURLs(userID string, shortURL string, updateChan chan<- string) error // Удаляет URL
	GetURLCount() (int, error)                                                 // Возвращает общее количество URL
	GetUserCount() (int, error)                                                // Возвращает количество уникальных пользователей
}

// Repository определяет интерфейс для работы с кэшем.
type Repository interface {
	Set(shortID string, originalURL string) // Сохраняет URL в кэш
	Get(shortID string) (string, bool)      // Извлекает URL из кэша
}

// ShortenerService предоставляет функционал для создания и управления короткими ссылками.
type ShortenerService struct {
	BaseURL   string     // Базовый URL для генерации коротких ссылок
	Storage   Repository // Кэш-хранилище ссылок
	db        Store      // Хранилище данных (БД)
	dbDNSTurn bool       // Флаг использования БД для хранения ссылок
}

// NewShortenerService создаёт и возвращает новый экземпляр сервиса сокращения ссылок.
func NewShortenerService(BaseURL string, storage Repository, db Store, dbDNSTurn bool) *ShortenerService {
	return &ShortenerService{
		BaseURL:   BaseURL,
		Storage:   storage,
		db:        db,
		dbDNSTurn: dbDNSTurn,
	}
}

// GetExistURL проверяет наличие ошибки уникальности и возвращает существующую короткую ссылку, если таковая уже имеется.
func (s *ShortenerService) GetExistURL(originalURL string, err error) (string, error) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		shortID, err := s.GetRep("", originalURL)
		shortURL := fmt.Sprintf("%s/%s", s.BaseURL, shortID)
		return shortURL, err
	}
	return "", err
}

// Set генерирует короткую ссылку для заданного originalURL и сохраняет её в хранилище.
func (s *ShortenerService) Set(userID, originalURL string) (string, error) {
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

// randSeq генерирует уникальный идентификатор (UUID) для короткой ссылки.
func randSeq() string {
	return uuid.New().String()
}

// Get возвращает оригинальный URL, используя короткий идентификатор, проверяя сначала БД, затем кэш.
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

// Ping проверяет доступность соединения с базой данных.
func (s *ShortenerService) Ping() error {
	return s.db.PingStore()
}

// CreateRep сохраняет запись URL в базе данных.
func (s *ShortenerService) CreateRep(originalURL, shortURL, UserID string) error {
	return s.db.Create(originalURL, shortURL, UserID)
}

// GetRep извлекает запись из базы данных по короткому или оригинальному URL.
func (s *ShortenerService) GetRep(shortURL, originalURL string) (string, error) {
	return s.db.Get(shortURL, originalURL)
}

// GetFullRep извлекает все URL пользователя по userID.
func (s *ShortenerService) GetFullRep(userID string) ([]map[string]string, error) {
	return s.db.GetFull(userID, s.BaseURL)
}

// DeleteURLsRep удаляет несколько URL для пользователя, используя централизованный воркер.
func (s *ShortenerService) DeleteURLsRep(userID string, shortURLs []string) error {
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

// GetURLCount возвращает общее количество сокращённых URL.
func (s *ShortenerService) GetURLCount() (int, error) {
	if s.dbDNSTurn {
		return s.db.GetURLCount()
	}
	return 0, errors.New("метод не поддерживается для текущего хранилища")
}

// GetUserCount возвращает количество уникальных пользователей.
func (s *ShortenerService) GetUserCount() (int, error) {
	if s.dbDNSTurn {
		return s.db.GetUserCount()
	}
	return 0, errors.New("метод не поддерживается для текущего хранилища")
}
