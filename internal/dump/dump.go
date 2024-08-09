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

type Memory struct {
	Storage *services.ShortenerService
	File    *os.File
}
type ShortCollector struct {
	NumberUUID  string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

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
				fmt.Println("error decode JSON:", err)
				break
			}
		}
		maxUUID += 1
		storageInstance.Set(event.OriginalURL, event.ShortURL)
	}
	return nil
}

func Set(storageInstance *storage.Storage, filePath string, BaseURL string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	maxUUID := 0
	for shortURL, originalURL := range storageInstance.URLs {
		// shortURL := fmt.Sprintf("%s/%s", BaseURL, key)
		maxUUID += 1
		ShortCollector := ShortCollector{
			strconv.Itoa(maxUUID),
			shortURL,
			originalURL,
		}
		writer := bufio.NewWriter(file)
		err = writeEvent(&ShortCollector, writer)
	}
	return err
}

func writeEvent(ShortCollector *ShortCollector, writer *bufio.Writer) error {
	data, err := json.Marshal(&ShortCollector)
	if err != nil {
		return err
	}

	// записываем событие в буфер
	if _, err := writer.Write(data); err != nil {
		return err
	}

	// добавляем перенос строки
	if err := writer.WriteByte('\n'); err != nil {
		return err
	}

	// записываем буфер в файл
	return writer.Flush()
}
