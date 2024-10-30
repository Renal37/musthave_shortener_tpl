package storage

import (
    "testing"
	"fmt"
)

// Benchmark for the Set method
func BenchmarkStorage_Set(b *testing.B) {
    storage := NewStorage()
    for i := 0; i < b.N; i++ {
        // Используем Sprintf для корректного преобразования int в string
        storage.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
    }
}

// Benchmark for the Get method
func BenchmarkStorage_Get(b *testing.B) {
    storage := NewStorage()
    // Заполним хранилище 1000 значениями
    for i := 0; i < 1000; i++ {
        storage.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i)) // Используем Sprintf
    }

    b.ResetTimer() // Сбрасываем таймер, чтобы не включать время на заполнение

    for i := 0; i < b.N; i++ {
        storage.Get(fmt.Sprintf("key%d", i%1000)) // Получаем ключи по кругу
    }
}

// Benchmark for the Get method for nonexistent keys
func BenchmarkStorage_GetNonExistent(b *testing.B) {
    storage := NewStorage()
    // Заполним хранилище 1000 значениями
    for i := 0; i < 1000; i++ {
        storage.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i)) // Используем Sprintf
    }

    b.ResetTimer() // Сбрасываем таймер

    for i := 0; i < b.N; i++ {
        storage.Get("nonexistent") // Запрашиваем несуществующий ключ
    }
}
