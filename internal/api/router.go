package api

import (
	"github.com/gin-gonic/gin"
)

func (s *RestAPI) setRoutes(r *gin.Engine) {
	r.POST("/", s.ShortenURLHandler)
	r.POST("/api/shorten", s.ShortenURLJSON)
	r.GET("/:id", s.RedirectToOriginalURL)
	r.GET("/ping", s.Ping)
	r.POST("/api/shorten/batch", s.ShortenURLsJSON)
	r.GET("/api/user/urls", s.UserURLsHandler)
	r.DELETE("/api/user/urls", s.DeleteUserUrls)
}
