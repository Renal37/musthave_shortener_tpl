package dump_test

import (
	"encoding/json"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/dump"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strconv"
	"testing"
)

// Определение структуры ShortCollector для тестовых данных
type ShortCollector struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

// Тест для несуществующего файла
func TestFillFromStorage_NonExistentFile(t *testing.T) {
	storageInstance := storage.NewStorage()
	err := dump.FillFromStorage(storageInstance, "/non/existent/file")
	assert.Error(t, err, "Ожидалась ошибка при попытке открыть несуществующий файл")
}

// Тест для большого набора данных
func TestFillFromStorage_LargeDataSet(t *testing.T) {
	tempFile, err := os.CreateTemp("", "testfile")
	require.NoError(t, err, "Не удалось создать временный файл")
	defer os.Remove(tempFile.Name())

	encoder := json.NewEncoder(tempFile)
	for i := 0; i < 10000; i++ {
		event := ShortCollector{
			OriginalURL: "http://example.com/" + strconv.Itoa(i),
			ShortURL:    "http://short.url/" + strconv.Itoa(i),
		}
		require.NoError(t, encoder.Encode(event), "Не удалось закодировать событие")
	}
	tempFile.Close()

	storageInstance := storage.NewStorage()
	err = dump.FillFromStorage(storageInstance, tempFile.Name())
	assert.NoError(t, err, "Не ожидалось ошибки при обработке большого объема данных")
}

// Тест для крайних случаев в Set
func TestSet_EdgeCases(t *testing.T) {
	tempFile, err := os.CreateTemp("", "testfile")
	require.NoError(t, err, "Не удалось создать временный файл")
	defer os.Remove(tempFile.Name())

	storageInstance := storage.NewStorage()

	// Тест пустого URL
	storageInstance.Set("", "")
	err = dump.Set(storageInstance, tempFile.Name())
	assert.NoError(t, err, "Не ожидалось ошибки при записи пустого URL")

	// Тест очень длинного URL
	longURL := "http://example.com/" + string(make([]byte, 10000))
	storageInstance.Set(longURL, "http://short.url/long")
	err = dump.Set(storageInstance, tempFile.Name())
	assert.NoError(t, err, "Не ожидалось ошибки при записи длинного URL")
}
