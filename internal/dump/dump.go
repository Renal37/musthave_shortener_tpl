package dump

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"io"
	"os"
	"strconv"
)

type ShortCollector struct {
	NumberUUID  string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// FillFromStorage загружает данные из файла в хранилище
func FillFromStorage(storageInstance *storage.Storage, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	newDecoder := json.NewDecoder(file)
	maxUUID := 0

	for {
		var event ShortCollector
		if err := newDecoder.Decode(&event); err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println("Ошибка при декодировании JSON:", err)
				break
			}
		}
		maxUUID++
		storageInstance.Set(event.OriginalURL, event.ShortURL)
	}
	return nil
}

// Set сохраняет данные из хранилища в файл
func Set(storageInstance *storage.Storage, filePath string, BaseURL string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	maxUUID := 0
	for shortURL, originalURL := range storageInstance.URLs {
		maxUUID++
		collector := ShortCollector{
			NumberUUID:  strconv.Itoa(maxUUID),
			ShortURL:    shortURL,
			OriginalURL: originalURL,
		}
		writer := bufio.NewWriter(file)
		err = writeEvent(&collector, writer)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeEvent(collector *ShortCollector, writer *bufio.Writer) error {
	data, err := json.Marshal(collector)
	if err != nil {
		return err
	}

	// Записываем событие в буфер
	if _, err := writer.Write(data); err != nil {
		return err
	}

	// Добавляем перенос строки
	if err := writer.WriteByte('\n'); err != nil {
		return err
	}

	// Сбрасываем буфер в файл
	return writer.Flush()
}
