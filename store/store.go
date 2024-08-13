package store

// Импортируем необходимые пакеты
import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

// Определяем структуру StoreDB, которая будет хранить подключение к базе данных
type StoreDB struct {
	db *sql.DB
}

// Функция InitDatabase инициализирует подключение к базе данных
func InitDatabase(DatabasePath string) (*StoreDB, error) {
	db, err := sql.Open("pgx", DatabasePath)
	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	storeDB := new(StoreDB)
	storeDB.db = db

	if DatabasePath != "" {
		err = createTable(db)
		if err != nil {
			return nil, fmt.Errorf("error creae table db: %w", err)
		}
	}

	return storeDB, nil
}

// Функция (s *StoreDB).Create создает новый URL-адрес в базе данных
func (s *StoreDB) Create(originalURL, shortURL string) error {
	query := `
        INSERT INTO urls (short_id, original_url) 
        VALUES ($1, $2)
    `
	_, err := s.db.Exec(query, shortURL, originalURL)
	if err != nil {
		return fmt.Errorf("error save URL: %w", err)
	}
	return nil
}

// Функция createTable создает таблицу urls в базе данных
func createTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS urls (
        id SERIAL PRIMARY KEY,
        short_id VARCHAR(256) NOT NULL UNIQUE,
        original_url TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// Функция (s *StoreDB).Get получает URL-адрес по его короткому

func (s *StoreDB) Get(shortURL string, originalURL string) (string, error) {
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

	var answer string
	err := s.db.QueryRow(query, field).Scan(&answer)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", err
		}
		return "", err
	}

	return answer, err
}

// Функция (s *StoreDB).PingStore проверяет подключение к базе данных
func (s *StoreDB) PingStore() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("pinging db-store: %w", err)
	}
	return nil
}
