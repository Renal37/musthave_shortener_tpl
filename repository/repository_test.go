package repository

import (
	"errors"
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
	assert.Equal(t, "originalURL1", result[0]["original_url"])
	assert.Equal(t, "http://localhost:8080/shortURL2", result[1]["short_url"])
	assert.Equal(t, "originalURL2", result[1]["original_url"])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreDB_GetFull_DeletedURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := &StoreDB{db: db}
	rows := sqlmock.NewRows([]string{"short_id", "original_url", "deletedFlag"}).
		AddRow("shortURL1", "originalURL1", true)

	mock.ExpectQuery("SELECT short_id, original_url, deletedFlag FROM urls").
		WithArgs("userID").
		WillReturnRows(rows)

	result, err := store.GetFull("userID", "http://localhost:8080")
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Equal(t, "Gone", err.Error()) // Проверяем правильность статуса ошибки
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
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreDB_Get_Deleted(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := &StoreDB{db: db}
	rows := sqlmock.NewRows([]string{"original_url", "deletedFlag"}).AddRow("originalURL", true)

	mock.ExpectQuery("SELECT original_url, deletedFlag FROM urls WHERE short_id =").
		WithArgs("shortURL").
		WillReturnRows(rows)

	originalURL, err := store.Get("shortURL", "")
	assert.Error(t, err)
	assert.Empty(t, originalURL)
	assert.Equal(t, "Gone", err.Error()) // Проверяем правильность статуса ошибки
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

func TestStoreDB_PingStore_Error(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)
	defer db.Close()

	store := &StoreDB{db: db}

	// Установите ожидание на пинг, который вернет ошибку
	mock.ExpectPing().WillReturnError(errors.New("ping error"))

	err = store.PingStore()
	assert.Error(t, err)
	assert.Equal(t, "pinging db-store: ping error", err.Error()) // Проверяем текст ошибки
	assert.NoError(t, mock.ExpectationsWereMet())                // Проверяем, что все ожидания выполнены
}
