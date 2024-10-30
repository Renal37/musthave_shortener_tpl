package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthorizationMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Создаем тестовый обработчик с middleware
	r := gin.New()
	r.Use(AuthorizationMiddleware())
	r.GET("/test", func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if exists {
			c.String(http.StatusOK, userID.(string))
		} else {
			c.String(http.StatusUnauthorized, "Unauthorized")
		}
	})

	// Запрос без куки
	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Запрос с установленной кукой
	req, _ = http.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{Name: "userID", Value: "validToken"})
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
