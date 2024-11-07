package middleware_test

import (
	"bytes"
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/middleware" 
)

// Helper-функция для сжатия данных с Gzip
func gzipCompress(data []byte) []byte {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	_, _ = writer.Write(data)
	writer.Close()
	return buf.Bytes()
}

// Тест на обработку ответа с сжатием Gzip
func TestCompressMiddleware_ResponseCompression(t *testing.T) {
	// Создаем новый gin роутер и добавляем middleware
	router := gin.New()
	router.Use(middleware.CompressMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"}) // Ответ JSON, который будет сжат
	})

	// Создаем HTTP-запрос с заголовком Accept-Encoding: gzip
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем, что ответ сжат
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "gzip", w.Header().Get("Content-Encoding"))

	// Декодируем сжатый ответ
	reader, err := gzip.NewReader(w.Body)
	assert.NoError(t, err)

	// Читаем распакованный ответ
	body, err := io.ReadAll(reader)
	assert.NoError(t, err)
	reader.Close()

	// Проверяем правильность содержимого ответа
	assert.Contains(t, string(body), `"message":"test"`)
}

// Тест на обработку сжатого запроса с Gzip
func TestCompressMiddleware_RequestDecompression(t *testing.T) {
	// Создаем новый gin роутер и добавляем middleware
	router := gin.New()
	router.Use(middleware.CompressMiddleware())
	router.POST("/test", func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)
		c.String(http.StatusOK, string(body)) // Возвращаем то, что получили в теле запроса
	})

	// Создаем сжатое тело запроса
	originalBody := `{"message": "test"}`
	compressedBody := gzipCompress([]byte(originalBody))

	// Создаем HTTP-запрос с заголовком Content-Encoding: gzip
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/test", bytes.NewReader(compressedBody))
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем, что ответ равен исходному телу запроса
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, originalBody, w.Body.String())
}

// Тест на случай, когда Content-Type не является application/json или text/html
func TestCompressMiddleware_NoCompressionForOtherContentTypes(t *testing.T) {
	// Создаем новый gin роутер и добавляем middleware
	router := gin.New()
	router.Use(middleware.CompressMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "no compression for this content type")
	})

	// Создаем HTTP-запрос с заголовком Accept-Encoding: gzip
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Type", "text/plain") // Некомпрессируемый тип

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем, что ответ не сжат
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "", w.Header().Get("Content-Encoding"))
	assert.Equal(t, "no compression for this content type", w.Body.String())
}
