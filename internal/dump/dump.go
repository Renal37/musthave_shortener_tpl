package dump

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/services"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
)

// Memory представляет собой структуру для хранения сервисов и файла.
type Memory struct {
	Storage *services.ShortenerService // Указатель на сервис сокращения URL
	File    *os.File                   // Указатель на файл для хранения данных
}

// ShortCollector представляет собой структуру для хранения данных о сокращенных URL.
type ShortCollector struct {
	NumberUUID  string `json:"uuid"`         // UUID
	ShortURL    string `json:"short_url"`    // Сокращенный URL
	OriginalURL string `json:"original_url"` // Оригинальный URL
}

// FillFromStorage заполняет хранилище данными из указанного файла.
func FillFromStorage(storageInstance *storage.Storage, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666) // Открываем файл
	if err != nil {
		return err // Возвращаем ошибку, если не удалось открыть файл
	}
	defer file.Close() // Закрываем файл по завершении
	newDecoder := json.NewDecoder(file) // Создаем новый декодер JSON
	maxUUID := 0 // Переменная для отслеживания максимального UUID

	// Читаем данные из файла
	for {
		var event ShortCollector
		if err := newDecoder.Decode(&event); err != nil {
			if err == io.EOF {
				break // Прерываем цикл, если достигнут конец файла
			} else {
				fmt.Println("error decode JSON:", err)
				break // Прерываем цикл, если произошла ошибка
			}
		}
		maxUUID += 1 // Увеличиваем счетчик UUID
		storageInstance.Set(event.OriginalURL, event.ShortURL) // Сохраняем данные в хранилище
	}
	return nil
}

// Set сохраняет данные из хранилища в указанный файл.
func Set(storageInstance *storage.Storage, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666) // Открываем файл
	if err != nil {
		return err // Возвращаем ошибку, если не удалось открыть файл
	}
	defer file.Close() // Закрываем файл по завершении
	maxUUID := 0 // Переменная для отслеживания максимального UUID

	// Сохраняем данные из хранилища в файл
	for shortURL, originalURL := range storageInstance.URLs {
		maxUUID += 1 // Увеличиваем счетчик UUID
		ShortCollector := ShortCollector{
			strconv.Itoa(maxUUID), // Преобразуем UUID в строку
			shortURL,
			originalURL,
		}
		writer := bufio.NewWriter(file) // Создаем буферизованный писатель
		err = writeEvent(&ShortCollector, writer) // Записываем событие в файл
	}
	return err
}

// writeEvent записывает событие в буферизованный писатель.
func writeEvent(ShortCollector *ShortCollector, writer *bufio.Writer) error {
	data, err := json.Marshal(&ShortCollector) // Кодируем структуру в JSON
	if err != nil {
		return err // Возвращаем ошибку, если произошла ошибка кодирования
	}

	// Записываем событие в буфер
	if _, err := writer.Write(data); err != nil {
		return err // Возвращаем ошибку, если произошла ошибка записи
	}

	// Добавляем перенос строки
	if err := writer.WriteByte('\n'); err != nil {
		return err // Возвращаем ошибку, если произошла ошибка записи
	}

	// Записываем буфер в файл
	return writer.Flush()
}
