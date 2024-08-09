package api

// Импортируем необходимые пакеты
import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Определяем структуру Request, которая будет хранить URL для сокращения
type Request struct {
	URL string `json:"url"`
}

// Определяем структуру Response, которая будет хранить результат сокращения URL
type Response struct {
	Result string `json:"result"`
}

// Функция ShortenURLHandler обрабатывает запросы на сокращение URL
func (s *RestAPI) ShortenURLHandler(c *gin.Context) {
	// Читаем тело запроса
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		// Если не удалось прочитать тело запроса, отправляем ошибку с кодом 500
		c.String(http.StatusInternalServerError, "Не удалось прочитать тело запроса")
		return
	}
	// Удаляем лишние пробелы в начале и конце строки
	url := strings.TrimSpace(string(body))
	// Создаем сокращенный URL с помощью ShortenerService
	shortURL := s.StructService.Set(url)
	// Устанавливаем заголовок "Content-Type" для текстового ответа
	c.Header("Content-Type", "text/plain")
	// Отправляем сокращенный URL в ответе
	c.String(http.StatusCreated, shortURL)
}

// Функция ShortenURLJSON обрабатывает запросы на сокращение URL в формате JSON
func (s *RestAPI) ShortenURLJSON(c *gin.Context) {
	// Создаем декодер для JSON-тела запроса
	var decoderBody Request
	decoder := json.NewDecoder(c.Request.Body)
	// Декодируем JSON-тело запроса
	err := decoder.Decode(&decoderBody)
	if err != nil {
		// Если не удалось декодировать JSON-тело запроса, отправляем ошибку с кодом 500
		errorMessage := map[string]interface{}{
			"message": "Не удалось прочитать тело запроса",
			"code":    http.StatusInternalServerError,
		}
		answer, _ := json.Marshal(errorMessage)
		// Устанавливаем заголовок "Content-Type" для JSON-ответа
		c.Header("Content-Type", "application/json")
		// Отправляем ошибку в формате JSON
		c.Data(http.StatusInternalServerError, "application/json", answer)
		return
	}
	// Удаляем лишние пробелы в начале и конце строки
	url := strings.TrimSpace(decoderBody.URL)
	// Создаем сокращенный URL с помощью ShortenerService
	shortURL := s.StructService.Set(url)
	// Создаем структуру ответа с результатом сокращения URL
	StructPerformance := Response{Result: shortURL}
	// Конвертируем структуру ответа в JSON
	respJSON, err := json.Marshal(StructPerformance)
	if err != nil {
		// Если не удалось создать JSON-ответ, отправляем ошибку с кодом 500
		errorMessage := map[string]interface{}{
			"message": "Не удалось создать ответ",
			"code":    http.StatusInternalServerError,
		}
		answer, _ := json.Marshal(errorMessage)
		// Устанавливаем заголовок "Content-Type" для JSON-ответа
		c.Header("Content-Type", "application/json")
		// Отправляем ошибку в формате JSON
		c.Data(http.StatusInternalServerError, "application/json", answer)
		return
	}
	// Устанавливаем заголовок "Content-Type" для JSON-ответа
	c.Header("Content-Type", "application/json")
	// Отправляем сокращенный URL в формате JSON
	c.Data(http.StatusCreated, "application/json", respJSON)
}

// Функция RedirectToOriginalURL обрабатывает запросы на переадку на оригинальный URL
func (s *RestAPI) RedirectToOriginalURL(c *gin.Context) {
	// Читаем параметр "id" из URL-адреса запроса
	shortID := c.Param("id")
	// Получаем оригинальный URL из ShortenerService
	originalURL, exists := s.StructService.Get(shortID)
	if !exists {
		// Если оригинальный URL не найден, отправляем ошибку с кодом 404
		c.String(http.StatusTemporaryRedirect, "URL не найден")
		return
	}
	// Устанавливаем заголовок "Location" для переадки на оригинальный URL
	c.Header("Location", originalURL)
	// Отправляем оригинальный URL в ответе
	c.String(http.StatusTemporaryRedirect, originalURL)
}

// Функция Ping обрабатывает запросы на проверку работоспособности сервиса
func (s *RestAPI) Ping(ctx *gin.Context) {
	// Вызываем метод Ping в ShortenerService
	err := s.StructService.Ping()
	if err != nil {
		// Если не удалось выполнить проверку работоспособности сервиса, отправляем ошибку с кодом 500
		ctx.JSON(http.StatusInternalServerError, "")
		return
	}
	// Если проверка выполнена успешно, отправляем пустой JSON-ответ
	ctx.JSON(http.StatusOK, "")
}
