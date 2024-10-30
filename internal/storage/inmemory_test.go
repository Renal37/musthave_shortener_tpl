package storage

import (
    "testing"
)

func TestStorage_SetAndGet(t *testing.T) {
    storage := NewStorage()

    // Тестируем установку и получение значения
    storage.Set("key1", "value1")
    
    // Проверяем, что значение установлено корректно
    value, exists := storage.Get("key1")
    if !exists {
        t.Error("Expected key1 to exist")
    }
    if value != "value1" {
        t.Errorf("Expected value to be 'value1', got '%s'", value)
    }

    // Проверяем, что несуществующий ключ возвращает false
    _, exists = storage.Get("nonexistent")
    if exists {
        t.Error("Expected nonexistent key to not exist")
    }
}

