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

// Интерфейс для работы с хранилищем (БД).
type Store interface {
	PingStore() error                                                          // Проверка соединения с хранилищем
	Create(originalURL, shortURL, UserID string) error                         // Создание новой записи в хранилище
	Get(shortIrl string, originalURL string) (string, error)                   // Получение оригинальной ссылки по сокращенной
	GetFull(userID string, BaseURL string) ([]map[string]string, error)        // Получение всех ссылок для пользователя
	DeleteURLs(userID string, shortURL string, updateChan chan<- string) error // Удаление ссылки
}

// Интерфейс для работы с кэш-памятью.
type Repository interface {
	Set(shortID string, originalURL string) // Сохранение данных в кэш
	Get(shortID string) (string, bool)      // Получение данных из кэша
}

// Структура сервиса сокращения ссылок.
type ShortenerService struct {
	BaseURL   string     // Базовый URL для коротких ссылок
	Storage   Repository // Хранилище для кэша
	db        Store      // База данных для работы с хранилищем
	dbDNSTurn bool       // Флаг, указывающий использовать ли БД
}

// Конструктор для создания сервиса сокращения ссылок.
func NewShortenerService(BaseURL string, storage Repository, db Store, dbDNSTurn bool) *ShortenerService {
	s := &ShortenerService{
		BaseURL:   BaseURL,
		Storage:   storage,
		db:        db,
		dbDNSTurn: dbDNSTurn,
	}
	return s
}

// Обработка ошибки уникального ограничения при создании ссылки.
func (s *ShortenerService) GetExistURL(originalURL string, err error) (string, error) {
	var pgErr *pgconn.PgError
	// Если возникла ошибка уникальности (duplicate key error)
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		// Получаем уже существующую короткую ссылку для данного originalURL
		shortID, err := s.GetRep("", originalURL)
		shortURL := fmt.Sprintf("%s/%s", s.BaseURL, shortID)
		return shortURL, err
	}
	return "", err
}

// Создание новой сокращенной ссылки.
func (s *ShortenerService) Set(userID, originalURL string) (string, error) {
	// Генерация уникального идентификатора для короткой ссылки
	shortID := randSeq()
	if s.dbDNSTurn {
		// Если используется база данных, сохраняем данные в хранилище
		err := s.CreateRep(originalURL, shortID, userID)
		if err != nil {
			return "", err
		}
	} else {
		// Иначе сохраняем в кэш
		s.Storage.Set(shortID, originalURL)
	}
	// Формируем полную короткую ссылку
	shortURL := fmt.Sprintf("%s/%s", s.BaseURL, shortID)
	return shortURL, nil
}

// Генерация случайной последовательности (UUID).
func randSeq() string {
	newUUID := uuid.New()
	return newUUID.String()
}

// Получение оригинальной ссылки по короткому идентификатору.
func (s *ShortenerService) Get(shortID string) (string, error) {
	// Если включен флаг использования базы данных
	if s.dbDNSTurn {
		originalURL, err := s.GetRep(shortID, "")
		if err != nil {
			return "", err
		}
		return originalURL, nil
	}

	// Получаем из кэша
	originalURL, exists := s.Storage.Get(shortID)
	if !exists {
		// Если не удалось найти ссылку
		err := errors.New("не удалось получить оригинальную ссылку")
		return "", err
	}
	return originalURL, nil
}

// Проверка доступности базы данных.
func (s *ShortenerService) Ping() error {
	return s.db.PingStore()
}

// Создание записи в хранилище (БД).
func (s *ShortenerService) CreateRep(originalURL, shortURL, UserID string) error {
	return s.db.Create(originalURL, shortURL, UserID)
}

// Получение записи из хранилища (БД) по короткой или оригинальной ссылке.
func (s *ShortenerService) GetRep(shortURL, originalURL string) (string, error) {
	return s.db.Get(shortURL, originalURL)
}

// Получение всех ссылок для конкретного пользователя.
func (s *ShortenerService) GetFullRep(userID string) ([]map[string]string, error) {
	return s.db.GetFull(userID, s.BaseURL)
}

// Удаление нескольких ссылок через централизованный воркер.
func (s *ShortenerService) DeleteURLsRep(userID string, shortURLs []string) error {
	updateChan := make(chan string, len(shortURLs)) // Канал для обновления URL
	workerChan := make(chan string, len(shortURLs)) // Канал задач для воркера

	// Запуск воркера в отдельной горутине для обработки удалений
	go func() {
		for shortURL := range workerChan {
			err := s.db.DeleteURLs(userID, shortURL, updateChan)
			if err != nil {
				// Логируем ошибку, если не удалось удалить ссылку
				logger.Log.Error("Не удалось удалить ссылку", zap.Error(err))
			}
		}
		close(updateChan) // Закрываем канал обновлений по завершении
	}()

	// Добавляем задачи на удаление в канал воркера
	go func() {
		for _, shortURL := range shortURLs {
			workerChan <- shortURL
		}
		close(workerChan) // Закрываем канал задач, когда все URL добавлены
	}()

	return nil
}
