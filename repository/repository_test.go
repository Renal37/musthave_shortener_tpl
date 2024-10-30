package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreDB_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := &StoreDB{db: db}
	mock.ExpectExec("INSERT INTO urls").
		WithArgs("shortURL", "originalURL", "userID").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = store.Create("originalURL", "shortURL", "userID")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreDB_GetFull(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := &StoreDB{db: db}
	rows := sqlmock.NewRows([]string{"short_id", "original_url", "deletedFlag"}).
		AddRow("shortURL1", "originalURL1", false).
		AddRow("shortURL2", "originalURL2", false)

	mock.ExpectQuery("SELECT short_id, original_url, deletedFlag FROM urls").
		WithArgs("userID").
		WillReturnRows(rows)

	result, err := store.GetFull("userID", "http://localhost:8080")
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "http://localhost:8080/shortURL1", result[0]["short_url"])
}

func TestStoreDB_DeleteURLs(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := &StoreDB{db: db}
	updateChan := make(chan string, 1)

	mock.ExpectExec("UPDATE urls").
		WithArgs("shortURL", "userID").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = store.DeleteURLs("userID", "shortURL", updateChan)
	assert.NoError(t, err)
	assert.Equal(t, "shortURL", <-updateChan)
}

func TestStoreDB_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := &StoreDB{db: db}
	rows := sqlmock.NewRows([]string{"original_url", "deletedFlag"}).AddRow("originalURL", false)

	mock.ExpectQuery("SELECT original_url, deletedFlag FROM urls WHERE short_id =").
		WithArgs("shortURL").
		WillReturnRows(rows)

	originalURL, err := store.Get("shortURL", "")
	assert.NoError(t, err)
	assert.Equal(t, "originalURL", originalURL)
}

func TestStoreDB_PingStore(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := &StoreDB{db: db}
	mock.ExpectPing()

	err = store.PingStore()
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
