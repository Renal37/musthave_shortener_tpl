package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"
   "strings"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type StoreDB struct {
	db *sql.DB
}

// Функция инициализации базы данных
func InitDatabase(DatabasePath string) (*StoreDB, error) {
	db, err := sql.Open("pgx", DatabasePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка при открытии базы данных: %w", err)
	}

	storeDB := &StoreDB{db: db}

	// Если указан путь к базе данных, то проверяем наличие таблицы
	if DatabasePath != "" {
		err = migrateSchema(db)
		if err != nil {
			return nil, fmt.Errorf("ошибка при создании таблицы в базе данных: %w", err)
		}
	}

	return storeDB, nil
}

// Функция создания URL
func (s *StoreDB) Create(originalURL, shortURL, userID string) error {
	query := `
        INSERT INTO urls (short_id, original_url, userID) 
        VALUES ($1, $2, $3)
    `
	_, err := s.db.Exec(query, shortURL, originalURL, userID)
	if err != nil {
		return fmt.Errorf("ошибка при создании URL: %w", err)
	}
	return nil
}

// Миграция схемы базы данных
func migrateSchema(db *sql.DB) error {
	// Создание таблицы, если она еще не существует
	query := `CREATE TABLE IF NOT EXISTS urls (
		id SERIAL PRIMARY KEY,
		short_id VARCHAR(256) NOT NULL UNIQUE,
		original_url TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    	userID VARCHAR(360),
    	deletedFlag BOOLEAN DEFAULT FALSE
	);`

	// Выполнение запроса создания таблицы
	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("ошибка при создании таблицы: %w", err)
	}

	// Создание индекса, если его еще нет
	indexQuery := `DO $$ 
	BEGIN 
   		IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE tablename = 'urls' AND indexname = 'idx_original_url') THEN
        	CREATE UNIQUE INDEX idx_original_url ON urls(original_url);
    	END IF;
	END $$;`

	if _, err := db.Exec(indexQuery); err != nil {
		return fmt.Errorf("ошибка при создании индекса: %w", err)
	}

	return nil
}

// Получение всех URL для пользователя
func (s *StoreDB) GetFull(userID string, BaseURL string) ([]map[string]string, error) {
	query := `SELECT short_id, original_url, deletedFlag FROM urls WHERE userID = $1`
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить ссылки: %w", err)
	}
	defer rows.Close()

	var urls []map[string]string
	for rows.Next() {
		var (
			shortID     string
			originalURL string
			deletedFlag bool
		)
		if err = rows.Scan(&shortID, &originalURL, &deletedFlag); err != nil {
			return nil, fmt.Errorf("ошибка при чтении строки: %w", err)
		}
		if deletedFlag {
			return nil, errors.New(http.StatusText(http.StatusGone))
		}
		shortURL := fmt.Sprintf("%s/%s", BaseURL, shortID)
		urlMap := map[string]string{"short_url": shortURL, "original_url": originalURL}
		urls = append(urls, urlMap)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации строк: %w", err)
	}

	return urls, nil
}
func (s *StoreDB) DeleteURLs(userID string, shortURLs string, updateChan chan<- string) error {
    // Преобразование строки URL в срез строк
    urls := strings.Split(shortURLs, ",")

    if len(urls) == 0 {
        return nil // Нечего удалять
    }

    // Создание плейсхолдеров для запроса
    placeholders := make([]string, len(urls))
    for i := range urls {
        placeholders[i] = fmt.Sprintf("$%d", i+1)
    }

    query := fmt.Sprintf(`
        UPDATE urls
        SET deletedFlag = true
        WHERE short_id IN (%s) AND userID = $%d
    `, strings.Join(placeholders, ", "), len(urls)+1)

    // Создание аргументов для запроса
    args := make([]interface{}, len(urls)+1)
    for i, url := range urls {
        args[i] = url
    }
    args[len(urls)] = userID

    // Выполнение запроса
    _, err := s.db.Exec(query, args...)
    if err != nil {
        return fmt.Errorf("ошибка при удалении URL-ов: %w", err)
    }

    // Отправка уведомлений
    for _, url := range urls {
        updateChan <- url
    }

    return nil
}


// Получение URL или короткой ссылки
func (s *StoreDB) Get(shortURL, originalURL string) (string, error) {
	var field1, field2, field string
	if shortURL != "" {
		field1, field2, field = "original_url", "short_id", shortURL
	} else {
		field1, field2, field = "short_id", "original_url", originalURL
	}

	query := fmt.Sprintf(`
        SELECT %s, deletedFlag 
        FROM urls 
        WHERE %s = $1
    `, field1, field2)

	var (
		answer      string
		deletedFlag bool
	)
	err := s.db.QueryRow(query, field).Scan(&answer, &deletedFlag)
	if err != nil {
		return "", fmt.Errorf("ошибка при получении данных: %w", err)
	}

	if deletedFlag {
		return "", errors.New(http.StatusText(http.StatusGone))
	}

	return answer, nil
}

// Проверка подключения к базе данных
func (s *StoreDB) PingStore() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("ошибка при пинге базы данных: %w", err)
	}
	return nil
}