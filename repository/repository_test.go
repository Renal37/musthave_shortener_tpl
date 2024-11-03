package repository

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreDB_Create_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := &StoreDB{db: db}
	mock.ExpectExec("INSERT INTO urls").
		WithArgs("shortURL", "originalURL", "userID").
		WillReturnError(errors.New("some error"))

	err = store.Create("originalURL", "shortURL", "userID")
	assert.Error(t, err)
	assert.Equal(t, "some error", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreDB_GetFull_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := &StoreDB{db: db}
	mock.ExpectQuery("SELECT short_id, original_url, deletedFlag FROM urls").
		WithArgs("userID").
		WillReturnError(errors.New("query error"))

	result, err := store.GetFull("userID", "http://localhost:8080")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreDB_DeleteURLs_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := &StoreDB{db: db}
	updateChan := make(chan string, 1)

	mock.ExpectExec("UPDATE urls").
		WithArgs("shortURL", "userID").
		WillReturnError(errors.New("delete error"))

	err = store.DeleteURLs("userID", "shortURL", updateChan)
	assert.Error(t, err)
	assert.Equal(t, "delete error", err.Error())
	assert.Empty(t, updateChan)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreDB_Get_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := &StoreDB{db: db}
	mock.ExpectQuery("SELECT original_url, deletedFlag FROM urls WHERE short_id =").
		WithArgs("shortURL").
		WillReturnError(errors.New("get error"))

	originalURL, err := store.Get("shortURL", "")
	assert.Error(t, err)
	assert.Empty(t, originalURL)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreDB_PingStore_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := &StoreDB{db: db}
	mock.ExpectPing().WillReturnError(errors.New("ping error"))

	err = store.PingStore()
	assert.Error(t, err)
	assert.Equal(t, "pinging db-store: ping error", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}
