package rest

import (
	"github.com/gin-gonic/gin"
)

// Публичный метод SetRoutes
func (s *RestAPI) SetRoutes(r *gin.Engine) {
	r.POST("/", s.ShortenURLHandler)
	r.POST("/api/shorten", s.ShortenURLJSON)
	r.GET("/:id", s.RedirectToOriginalURL)
	r.GET("/ping", s.Ping)
	r.POST("/api/shorten/batch", s.ShortenURLsJSON)
	r.GET("/api/user/urls", s.UserURLsHandler)
	r.DELETE("/api/user/urls", s.DeleteUserUrls)
	r.GET("/api/internal/stats", s.StatsHandler)
}
