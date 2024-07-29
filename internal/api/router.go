package api

import (
	"github.com/gin-gonic/gin"
)

// Функция setRoutes - устанавливает маршруты для сервера
func (s *RestAPI) setRoutes(r *gin.Engine) {
	r.POST("/", s.ShortenURLHandler)
	r.POST("/api/shorten", s.ShortenURLHandlerJSON)
	r.GET("/:id", s.RedirectToOriginalURLHandler)
}
