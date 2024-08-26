package store

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"time"
)

// StoreDB представляет структуру базы данных
type StoreDB struct {
	db *sql.DB
}

// InitDatabase инициализирует базу данных по заданному пути
func InitDatabase(databasePath string) (*StoreDB, error) {
	if databasePath == "" {
		return nil, fmt.Errorf("путь к базе данных не может быть пустым")
	}

	// Подключаемся к базе данных
	db, err := sql.Open("pgx", databasePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка при открытии базы данных: %w", err)
	}

	// Создаем таблицу, если требуется
	if err := createTable(db); err != nil {
		return nil, fmt.Errorf("ошибка при создании таблицы в базе данных: %w", err)
	}

	return &StoreDB{db: db}, nil
}

// Create добавляет новую запись с оригинальным и коротким URL
func (s *StoreDB) Create(originalURL, shortURL string) error {
	if originalURL == "" || shortURL == "" {
		return fmt.Errorf("оригинальный URL и короткий URL не могут быть пустыми")
	}

	query := `INSERT INTO urls (short_id, original_url) VALUES ($1, $2)`
	if _, err := s.db.Exec(query, shortURL, originalURL); err != nil {
		return fmt.Errorf("ошибка при добавлении записи в базу данных: %w", err)
	}
	return nil
}

// createTable создает таблицу urls, если она еще не существует
func createTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY,
			short_id VARCHAR(256) NOT NULL UNIQUE,
			original_url TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_original_url ON urls(original_url);
	`
	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("ошибка при создании таблицы: %w", err)
	}
	return nil
}

// Get возвращает оригинальный или короткий URL на основе переданного значения
func (s *StoreDB) Get(shortURL, originalURL string) (string, error) {
	if shortURL == "" && originalURL == "" {
		return "", fmt.Errorf("короткий URL или оригинальный URL должны быть указаны")
	}

	// Определяем, по какому полю будет происходить поиск
	field1, field2, fieldValue := "original_url", "short_id", shortURL
	if shortURL == "" {
		field1, field2, fieldValue = "short_id", "original_url", originalURL
	}

	query := fmt.Sprintf(`SELECT %s FROM urls WHERE %s = $1`, field1, field2)

	var result string
	if err := s.db.QueryRow(query, fieldValue).Scan(&result); err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("запись не найдена")
		}
		return "", fmt.Errorf("ошибка при получении данных из базы: %w", err)
	}

	return result, nil
}

// PingStore проверяет доступность базы данных
func (s *StoreDB) PingStore() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("ошибка при проверке доступности базы данных: %w", err)
	}

	return nil
}
