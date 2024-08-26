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
func InitDatabase(DatabasePath string) (*StoreDB, error) {
	if DatabasePath == "" {
		return nil, fmt.Errorf("путь к базе данных не может быть пустым")
	}

	// Подключаемся к базе данных
	db, err := sql.Open("pgx", DatabasePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка при открытии базы данных: %w", err)
	}

	// Создаем экземпляр StoreDB и инициализируем его
	storeDB := &StoreDB{db: db}

	// Создаем таблицу, если путь к базе данных не пустой
	err = createTable(db)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании таблицы в базе данных: %w", err)
	}

	return storeDB, nil
}

// Create добавляет новую запись в таблицу с коротким и оригинальным URL
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
		return fmt.Errorf("ошибка при добавлении записи в базу данных: %w", err)
	}
	return nil
}

// createTable создает таблицу urls, если она еще не существует
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
		return fmt.Errorf("ошибка при создании таблицы: %w", err)
	}

	return nil
}

// Get возвращает оригинальный или короткий URL на основе переданного значения
func (s *StoreDB) Get(shortURL string, originalURL string) (string, error) {
	if shortURL == "" && originalURL == "" {
		return "", fmt.Errorf("короткий URL или оригинальный URL должны быть указаны")
	}

	field1 := "original_url"
	field2 := "short_id"
	field := shortURL

	// Определяем, по какому полю будет происходить поиск
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

	var answer string
	err := s.db.QueryRow(query, field).Scan(&answer)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("запись не найдена")
		}
		return "", fmt.Errorf("ошибка при получении данных из базы: %w", err)
	}

	return answer, nil
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
