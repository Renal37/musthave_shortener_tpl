package middleware

import (
	"bytes"
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"strings"
)

// gzipWriter оборачивает ResponseWriter для поддержки сжатия.
type gzipWriter struct {
	gin.ResponseWriter           // Встраивание gin.ResponseWriter
	Writer             io.Writer // Сжатый вывод
}

// Write записывает данные в сжатый вывод.
func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b) // Запись сжатых данных
}

// CompressMiddleware возвращает middleware для сжатия ответов с помощью Gzip.
func CompressMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверяем, нужно ли сжать ответ
		if c.Request.Header.Get("Content-Type") == "application/json" ||
			c.Request.Header.Get("Content-Type") == "text/html" {

			if strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
				compressWriter := gzip.NewWriter(c.Writer)       // Создаем новый Gzip-Writer
				defer compressWriter.Close()                     // Закрываем writer по завершению
				c.Header("Content-Encoding", "gzip")             // Указываем, что ответ сжат
				c.Writer = &gzipWriter{c.Writer, compressWriter} // Оборачиваем ResponseWriter
			}
		}

		// Обработка входящих запросов с Gzip-сжатием
		if strings.Contains(c.Request.Header.Get("Content-Encoding"), "gzip") {
			compressReader, err := gzip.NewReader(c.Request.Body) // Создаем новый Gzip-Reader
			if err != nil {
				log.Fatalf("error: new reader: %d", err)
				return
			}
			defer compressReader.Close() // Закрываем reader по завершению

			body, err := io.ReadAll(compressReader) // Читаем тело запроса
			if err != nil {
				log.Fatalf("error: read body: %d", err)
				return
			}

			c.Request.Body = io.NopCloser(bytes.NewReader(body)) // Заменяем тело запроса
			c.Request.ContentLength = int64(len(body))           // Устанавливаем новую длину тела запроса
		}
		c.Next() // Переход к следующему обработчику
	}
}
