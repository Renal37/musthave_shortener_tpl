package api

import (
	"github.com/gin-gonic/gin"
)

// Функция setRoutes - устанавливает маршруты для сервера
func (s *RestAPI) setRoutes(r *gin.Engine) {
	r.POST("/", s.ShortenURLHandler)
	r.POST("/api/shorten", s.ShortenURLJSON)
	r.GET("/:id", s.RedirectToOriginalURL)
	r.GET("/ping", s.Ping)
	r.POST("/api/shorten/batch", s.ShortenURLsJSON)

}
