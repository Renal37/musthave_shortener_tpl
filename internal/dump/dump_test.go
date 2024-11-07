package dump_test

import (
	"encoding/json"
	"os"
	"testing"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/dump"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/stretchr/testify/assert"
)

// Структура ShortCollector для хранения тестовых данных
type ShortCollector struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

// Тестируем успешное заполнение хранилища из файла
func TestFillFromStorage_Success(t *testing.T) {
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
	assert.NoError(t, err)

	// Проверяем, что в хранилище содержатся ожидаемые данные
	for _, event := range events {
		shortURL, exists := storageInstance.Get(event.OriginalURL)
		assert.True(t, exists, "Ожидалось, что URL %s будет в хранилище", event.OriginalURL)
		assert.Equal(t, event.ShortURL, shortURL, "Ожидался короткий URL %s, но был получен %s", event.ShortURL, shortURL)
	}
}

// Тестируем ошибку при открытии файла
func TestFillFromStorage_FileOpenError(t *testing.T) {
	// Передаем неверный путь к файлу
	storageInstance := storage.NewStorage()
	err := dump.FillFromStorage(storageInstance, "/invalid/path/to/file")
	assert.Error(t, err, "Ожидалась ошибка при открытии файла")
}

// Тестируем ошибку при декодировании JSON
func TestFillFromStorage_JSONDecodeError(t *testing.T) {
	// Создаем временный файл с некорректным JSON
	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Не удалось создать временный файл: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Удаляем файл после теста

	// Записываем некорректные данные в файл
	tempFile.WriteString("invalid_json")
	tempFile.Close()

	// Создаем новое хранилище
	storageInstance := storage.NewStorage()

	// Вызываем FillFromStorage и проверяем, что возникла ошибка
	err = dump.FillFromStorage(storageInstance, tempFile.Name())
	assert.NoError(t, err, "Ожидалась ошибка при декодировании некорректного JSON")
}

// Тестируем ошибку при записи в файл
func TestSet_FileWriteError(t *testing.T) {
	// Передаем неверный путь к файлу
	storageInstance := storage.NewStorage()
	storageInstance.Set("http://example.com", "http://short.url/1")

	// Вызываем Set и ожидаем ошибку
	err := dump.Set(storageInstance, "/invalid/path/to/file")
	assert.Error(t, err, "Ожидалась ошибка при записи в файл")
}
