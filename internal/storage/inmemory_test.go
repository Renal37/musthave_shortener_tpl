package storage

import (
	"fmt"
	"testing"
)

func BenchmarkStorageSet(b *testing.B) {
	storage := NewStorage()
	for i := 0; i < b.N; i++ {
		storage.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}
}

func BenchmarkStorageGet(b *testing.B) {
	storage := NewStorage()
	// Заполним хранилище 1000 значениями
	for i := 0; i < 1000; i++ {
		storage.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}

	b.ResetTimer() // Сбрасываем таймер

	for i := 0; i < b.N; i++ {
		storage.Get(fmt.Sprintf("key%d", i%1000))
	}
}

func BenchmarkStorageGetNonExistent(b *testing.B) {
	storage := NewStorage()
	for i := 0; i < 1000; i++ {
		storage.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		storage.Get("nonexistent")
	}
}

// func TestStorage_SetAndGet(t *testing.T) {
//     storage := NewStorage()

//     // Тестируем установку и получение значения
//     storage.Set("key1", "value1")

//     // Проверяем, что значение установлено корректно
//     value, exists := storage.Get("key1")
//     if !exists {
//         t.Error("Expected key1 to exist")
//     }
//     if value != "value1" {
//         t.Errorf("Expected value to be 'value1', got '%s'", value)
//     }

//     // Проверяем, что несуществующий ключ возвращает false
//     _, exists = storage.Get("nonexistent")
//     if exists {
//         t.Error("Expected nonexistent key to not exist")
//     }
// }
