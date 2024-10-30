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
