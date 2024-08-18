package store

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"time"
)

type StoreDB struct {
	db *sql.DB
}

// Инициализация базы данных
func InitDatabase(DatabasePath string) (*StoreDB, error) {
	db, err := sql.Open("pgx", DatabasePath) // Открытие соединения с базой данных
	if err != nil {
		return nil, fmt.Errorf("ошибка при открытии базы данных: %w", err)
	}

	storeDB := &StoreDB{db: db}

	if DatabasePath != "" {
		err = createTable(db) // Создание таблицы, если путь к базе данных не пустой
		if err != nil {
			return nil, fmt.Errorf("ошибка при создании таблицы в базе данных: %w", err)
		}
	}

	return storeDB, nil
}

// Создание записи в базе данных
func (s *StoreDB) Create(originalURL, shortURL string) error {
	query := `
        INSERT INTO urls (short_id, original_url) 
        VALUES ($1, $2)
    `
	_, err := s.db.Exec(query, shortURL, originalURL) // Вставка данных в таблицу
	if err != nil {
		return fmt.Errorf("ошибка при создании записи: %w", err)
	}
	return nil
}

// Создание таблицы, если она не существует
func createTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY,
			short_id VARCHAR(256) NOT NULL UNIQUE,
			original_url TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		DO $$ 
		BEGIN 
			IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE tablename = 'urls' AND indexname = 'idx_original_url') THEN
				CREATE UNIQUE INDEX idx_original_url ON urls(original_url);
			END IF;
		END $$;
	`

	_, err := db.Exec(query) // Выполнение SQL-запроса для создания таблицы
	if err != nil {
		return fmt.Errorf("ошибка при создании таблицы: %w", err)
	}

	return nil
}

// Получение оригинального или сокращённого URL из базы данных
func (s *StoreDB) Get(shortURL, originalURL string) (string, error) {
	field1 := "original_url"
	field2 := "short_id"
	field := shortURL
	if shortURL == "" {
		field1, field2 = field2, field1
		field = originalURL
	}

	query := fmt.Sprintf(`
        SELECT %s 
        FROM urls 
        WHERE %s = $1
    `, field1, field2)

	var result string
	err := s.db.QueryRow(query, field).Scan(&result) // Выполнение запроса и сканирование результата
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("запись не найдена: %w", err)
		}
		return "", fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}

	return result, nil
}

// Проверка соединения с базой данных
func (s *StoreDB) PingStore() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("ошибка при проверке соединения с базой данных: %w", err)
	}
	return nil
}
