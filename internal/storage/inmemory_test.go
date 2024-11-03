package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStorage(t *testing.T) {
	// Создаем новое хранилище
	storage := NewStorage()

	// Проверяем, что хранилище инициализировано и карта URL-адресов не nil
	assert.NotNil(t, storage)
	assert.NotNil(t, storage.URLs)
	assert.Empty(t, storage.URLs) // Проверяем, что карта пустая
}

func TestSetAndGet(t *testing.T) {
	// Создаем новое хранилище
	storage := NewStorage()

	// Устанавливаем значение
	key := "example"
	value := "http://example.com"
	storage.Set(key, value)

	// Проверяем, что значение корректно сохраняется
	retrievedValue, exists := storage.Get(key)
	assert.True(t, exists)                     // Проверяем, что ключ существует
	assert.Equal(t, value, retrievedValue)     // Проверяем, что возвращаемое значение совпадает с сохраненным
}

func TestGet_NonExistentKey(t *testing.T) {
	// Создаем новое хранилище
	storage := NewStorage()

	// Проверяем, что получение несуществующего ключа возвращает false
	retrievedValue, exists := storage.Get("nonexistent")

	assert.False(t, exists)                  // Проверяем, что ключ не существует
	assert.Empty(t, retrievedValue)           // Проверяем, что возвращаемое значение пустое
}

// package storage

// import (
// 	"fmt"
// 	"testing"
// )

// func BenchmarkStorageSet(b *testing.B) {
// 	storage := NewStorage()
// 	for i := 0; i < b.N; i++ {
// 		storage.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
// 	}
// }

// func BenchmarkStorageGet(b *testing.B) {
// 	storage := NewStorage()
// 	// Заполним хранилище 1000 значениями
// 	for i := 0; i < 1000; i++ {
// 		storage.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
// 	}

// 	b.ResetTimer() // Сбрасываем таймер

// 	for i := 0; i < b.N; i++ {
// 		storage.Get(fmt.Sprintf("key%d", i%1000))
// 	}
// }

// func BenchmarkStorageGetNonExistent(b *testing.B) {
// 	storage := NewStorage()
// 	for i := 0; i < 1000; i++ {
// 		storage.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
// 	}

// 	b.ResetTimer()

// 	for i := 0; i < b.N; i++ {
// 		storage.Get("nonexistent")
// 	}
// }

