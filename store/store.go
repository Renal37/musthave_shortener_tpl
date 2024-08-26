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

func (s *StoreDB) PingStore() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("pinging db-store: %w", err)
	}
	return nil
}