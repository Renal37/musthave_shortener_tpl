package middleware

import (
	"bytes"
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"strings"
)

// gzipWriter оборачивает gin.ResponseWriter и добавляет Writer для сжатия
type gzipWriter struct {
	gin.ResponseWriter
	Writer io.Writer
}

// Write выполняет запись через сжимающий Writer
func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// CompressRequest возвращает middleware для обработки сжатия запросов и ответов
func CompressMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверяем, что Content-Type запроса либо JSON, либо HTML
		if c.Request.Header.Get("Content-Type") == "application/json" ||
			c.Request.Header.Get("Content-Type") == "text/html" {
			// Если клиент поддерживает gzip сжатие, сжимаем ответ
			if strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
				compressWriter := gzip.NewWriter(c.Writer)
				defer compressWriter.Close()
				c.Header("Content-Encoding", "gzip")
				c.Writer = &gzipWriter{c.Writer, compressWriter}
			}
		}
		// Если запрос сжат с использованием gzip, разжимаем тело запроса
		if strings.Contains(c.Request.Header.Get("Content-Encoding"), "gzip") {
			compressReader, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				log.Printf("ошибка: создание нового gzip reader: %v", err)
				c.AbortWithStatusJSON(400, gin.H{"ошибка": "неверное gzip кодирование"})
				return
			}
			defer compressReader.Close()

			body, err := io.ReadAll(compressReader)
			if err != nil {
				log.Printf("ошибка: чтение тела запроса: %v", err)
				c.AbortWithStatusJSON(400, gin.H{"ошибка": "невозможно прочитать тело gzip"})
				return
			}

			c.Request.Body = io.NopCloser(bytes.NewReader(body))
			c.Request.ContentLength = int64(len(body))
		}

		// Переходим к следующему обработчику
		c.Next()
	}
}

// containsGzip проверяет, присутствует ли "gzip" в списке строк
func containsGzip(content []string) bool {
	for _, v := range content {
		if v == "gzip" {
			return true
		}
	}
	return false
}
