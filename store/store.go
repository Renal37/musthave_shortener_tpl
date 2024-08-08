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

	return storeDB, nil
}

func (s *StoreDB) PingStore() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("pinging db-store: %w", err)
	}
	return nil
}