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

func InitDatabase(DatabasePath string) (*StoreDB, error) {
	db, err := sql.Open("pgx", DatabasePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка при открытии базы данных: %w", err)
	}

	storeDB := new(StoreDB)
	storeDB.db = db

	if DatabasePath != "" {
		err = createTable(db) // Создание таблицы, если путь к базе данных не пустой
		if err != nil {
			return nil, fmt.Errorf("ошибка при создании таблицы в базе данных: %w", err)
		}
	}

	return storeDB, nil
}

func (s *StoreDB) Create(originalURL, shortURL string) error {
	if originalURL == "" || shortURL == "" {
		return fmt.Errorf("оригинальный URL и короткий URL не могут быть пустыми")
	}
	query := `
        INSERT INTO urls (short_id, original_url) 
        VALUES ($1, $2)
    `
	_, err := s.db.Exec(query, shortURL, originalURL)
	if err != nil {
		return err
	}
	return nil
}

func createTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS urls (
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
	END $$;`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (s *StoreDB) Get(shortURL string, originalURL string) (string, error) {
	if shortURL == "" && originalURL == "" {
		return "", fmt.Errorf("короткий URL или оригинальный URL должны быть указаны")
	}	
	field1 := "original_url"
	field2 := "short_id"
	field := shortURL
	if shortURL == "" {
		field2 = "original_url"
		field1 = "short_id"
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

func (s *StoreDB) PingStore() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("ошибка при создании записи: %w", err)
	}
	return nil
}