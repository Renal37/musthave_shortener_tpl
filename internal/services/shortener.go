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

// Интерфейс для работы с хранилищем (БД)
type Store interface {
	PingStore() error
	Create(originalURL, shortURL, UserID string) error
	Get(shortIrl string, originalURL string) (string, error)
	GetFull(userID string, BaseURL string) ([]map[string]string, error)
	DeleteURLs(userID string, shortURL string, updateChan chan<- string) error
}

// Интерфейс для работы с репозиторием (кэш)
type Repository interface {
	Set(shortID string, originalURL string)
	Get(shortID string) (string, bool)
}

// Сервис для сокращения URL
type ShortenerService struct {
	BaseURL   string
	Storage   Repository
	db        Store
	dbDNSTurn bool
}

// Конструктор для создания нового сервиса сокращения URL
func NewShortenerService(BaseURL string, storage Repository, db Store, dbDNSTurn bool) *ShortenerService {
	s := &ShortenerService{
		BaseURL:   BaseURL,
		Storage:   storage,
		db:        db,
		dbDNSTurn: dbDNSTurn,
	}
	return s
}

// Метод для получения уже существующего короткого URL при ошибке уникальности
func (s *ShortenerService) GetExistURL(originalURL string, err error) (string, error) {
	var pgErr *pgconn.PgError
	// Проверяем, если ошибка связана с нарушением уникальности
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		shortID, err := s.GetRep("", originalURL)
		shortURL := fmt.Sprintf("%s/%s", s.BaseURL, shortID)
		return shortURL, err
	}
	return "", err
}

// Метод для сокращения URL
func (s *ShortenerService) Set(userID, originalURL string) (string, error) {
	// Генерируем короткий ID для URL
	shortID := randSeq()
	if s.dbDNSTurn {
		// Если используется база данных, сохраняем URL в БД
		err := s.CreateRep(originalURL, shortID, userID)
		if err != nil {
			return "", err
		}
	} else {
		// Иначе сохраняем URL в репозиторий (например, кэш)
		s.Storage.Set(shortID, originalURL)
	}
	shortURL := fmt.Sprintf("%s/%s", s.BaseURL, shortID)
	return shortURL, nil
}

// Генерация уникального идентификатора
func randSeq() string {
	newUUID := uuid.New()
	return newUUID.String()
}

// Метод для получения оригинального URL по его короткому идентификатору
func (s *ShortenerService) Get(shortID string) (string, error) {
	if s.dbDNSTurn {
		// Если используется база данных, ищем в БД
		originalURL, err := s.GetRep(shortID, "")
		if err != nil {
			return "", err
		}
		return originalURL, nil
	}

	// Ищем в репозитории (кэш)
	originalURL, exists := s.Storage.Get(shortID)
	if !exists {
		err := errors.New("не удалось получить оригинальный URL")
		return "", err
	}
	return originalURL, nil
}

// Проверка доступности базы данных
func (s *ShortenerService) Ping() error {
	return s.db.PingStore()
}

// Метод для сохранения URL в базе данных
func (s *ShortenerService) CreateRep(originalURL, shortURL, UserID string) error {
	return s.db.Create(originalURL, shortURL, UserID)
}

// Метод для получения данных из базы
func (s *ShortenerService) GetRep(shortURL, originalURL string) (string, error) {
	return s.db.Get(shortURL, originalURL)
}

// Метод для получения всех URL пользователя
func (s *ShortenerService) GetFullRep(userID string) ([]map[string]string, error) {
	return s.db.GetFull(userID, s.BaseURL)
}

// Метод для удаления URL пользователя с использованием фанина
func (s *ShortenerService) DeleteURLsRep(userID string, shortURLs []string) error {
	// Канал для результата удаления URL
	updateChan := make(chan string, len(shortURLs))
	// Запускаем обработку удалений в отдельной горутине
	go s.batchURLDeletion(userID, shortURLs, updateChan)
	return nil
}

// Функция для централизованной обработки удалений URL через буферизированный канал
func (s *ShortenerService) batchURLDeletion(userID string, shortURLs []string, updateChan chan string) {
	defer close(updateChan) // Закрываем канал после завершения

	for _, shortURL := range shortURLs {
		// Пытаемся удалить URL
		err := s.db.DeleteURLs(userID, shortURL, updateChan)
		if err != nil {
			logger.Log.Error("Не удалось удалить URL", zap.Error(err))
			continue
		}
		// Отправляем результат успешного удаления в канал
		updateChan <- shortURL
	}
}

// Fan-In воркер для обработки задач удаления в отдельной горутине
func StartURLDeletionWorker(bufferSize int) chan<- string {
	deletionChan := make(chan string, bufferSize)

	go func() {
		batch := []string{}
		for url := range deletionChan {
			batch = append(batch, url)
			// Если накопилось достаточно задач, выполняем их
			if len(batch) >= bufferSize {
				// Здесь выполняется логика обновления/удаления
				fmt.Println("Обрабатываем батч:", batch)
				batch = []string{}
			}
		}
	}()

	return deletionChan
}
