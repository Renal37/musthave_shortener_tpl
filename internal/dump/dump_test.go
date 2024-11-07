package dump_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/dump"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
)

// Структура ShortCollector для хранения тестовых данных
type ShortCollector struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

func TestFillFromStorage(t *testing.T) {
	// Создаем временный файл для симуляции входного файла
	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Не удалось создать временный файл: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Удаляем файл после теста

	// Подготавливаем тестовые данные
	events := []ShortCollector{
		{OriginalURL: "http://example.com", ShortURL: "http://short.url/1"},
		{OriginalURL: "http://example.org", ShortURL: "http://short.url/2"},
	}

	// Записываем тестовые данные в временный файл
	encoder := json.NewEncoder(tempFile)
	for _, event := range events {
		if err := encoder.Encode(event); err != nil {
			t.Fatalf("Не удалось закодировать событие: %v", err)
		}
	}
	tempFile.Close() // Закрываем временный файл для завершения записи

	// Создаем новый экземпляр хранилища
	storageInstance := storage.NewStorage()

	// Вызываем тестируемую функцию
	err = dump.FillFromStorage(storageInstance, tempFile.Name())
	if err != nil {
		t.Fatalf("FillFromStorage вернул ошибку: %v", err)
	}

	// Проверяем, что в хранилище содержатся ожидаемые данные
	for _, event := range events {
		shortURL, exists := storageInstance.Get(event.OriginalURL)
		if !exists {
			t.Errorf("Ожидалось, что URL %s будет в хранилище", event.OriginalURL)
		}
		if shortURL != event.ShortURL {
			t.Errorf("Ожидался короткий URL %s, но был получен %s", event.ShortURL, shortURL)
		}
	}
}
