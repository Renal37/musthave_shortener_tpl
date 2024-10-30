package storage

import (
	"fmt"
	"testing"
)

func TestStorage(t *testing.T) {
	// Создаем новый экземпляр хранилища
	storage := NewStorage()

	// Тестируем метод Set
	storage.Set("key1", "value1") // Удалили переменную ошибки

	// Проверяем, что значение для key1 действительно равно value1
	value, exists := storage.Get("key1")
	if !exists {
		t.Errorf("Ожидалось, что key1 существует")
	}
	if value != "value1" {
		t.Errorf("Ожидалось, что значение будет: %s, а получилось: %s", "value1", value)
	}

	// Тестируем метод Get для несуществующего ключа
	value, exists = storage.Get("key2")
	if exists {
		t.Errorf("Ожидалось, что key2 не существует")
	}

	// Тестируем несколько операций Set и Get
	for i := 0; i < 10; i++ {
		storage.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i)) // Удалили переменную ошибки
	}
	for i := 0; i < 10; i++ {
		value, exists = storage.Get(fmt.Sprintf("key%d", i))
		if !exists {
			t.Errorf("Ожидалось, что key%d существует", i)
		}
		if value != fmt.Sprintf("value%d", i) {
			t.Errorf("Ожидалось, что значение будет: %s, а получилось: %s", fmt.Sprintf("value%d", i), value)
		}
	}
}
