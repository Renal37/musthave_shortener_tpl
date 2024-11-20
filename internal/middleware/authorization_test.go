package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestAuthorizationMiddleware тестирует AuthorizationMiddleware
func TestAuthorizationMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Тестируемый роут
	r := gin.Default()
	r.Use(middleware.AuthorizationMiddleware())
	r.GET("/test", func(c *gin.Context) {
		userID, _ := c.Get("userID")
		c.String(http.StatusOK, fmt.Sprintf("Success, userID: %v", userID))
	})

	// Генерируем валидный токен
	validToken, err := middleware.BuildJWTString()
	assert.NoError(t, err)

	tests := []struct {
		name           string
		cookie         string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid token",
			cookie:         validToken,
			expectedStatus: http.StatusOK,
			expectedBody:   "Success, userID:", // Так как ID пользователя уникален
		},
		{
			name:           "Missing token",
			cookie:         "",
			expectedStatus: http.StatusOK,
			expectedBody:   "Success, userID:", // Новый токен создается
		},
		{
			name:           "Invalid token",
			cookie:         "invalid.token.string",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Unauthorized токен недействителен",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)

			if tt.cookie != "" {
				req.AddCookie(&http.Cookie{Name: "userID", Value: tt.cookie})
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}

// TestGetUserIDFromCookie тестирует функцию getUserIDFromCookie
func TestGetUserIDFromCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	// Генерируем валидный токен
	validToken, err := middleware.BuildJWTString()
	assert.NoError(t, err)

	tests := []struct {
		name        string
		cookieValue string
		expectError bool
	}{
		{
			name:        "Valid cookie",
			cookieValue: validToken,
			expectError: false,
		},
		{
			name:        "Invalid cookie",
			cookieValue: "invalid.token",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Устанавливаем cookie
			ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
			ctx.Request.AddCookie(&http.Cookie{Name: "userID", Value: tt.cookieValue})

			userInfo, err := middleware.GetUserIDFromCookie(ctx)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, userInfo)
				assert.NotEmpty(t, userInfo.ID)
			}
		})
	}
}

// TestBuildJWTString тестирует функцию BuildJWTString
func TestBuildJWTString(t *testing.T) {
	tokenString, err := middleware.BuildJWTString()
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Проверим, что токен можно распарсить
	userID, err := middleware.GetUserID(tokenString)
	assert.NoError(t, err)
	assert.NotEmpty(t, userID)
}

// TestGetUserID тестирует функцию GetUserID
func TestGetUserID(t *testing.T) {
	// Создаем валидный токен
	tokenString, err := middleware.BuildJWTString()
	assert.NoError(t, err)

	tests := []struct {
		name         string
		token        string
		expectError  bool
		expectedUser string
	}{
		{
			name:         "Valid token",
			token:        tokenString,
			expectError:  false,
			expectedUser: "", // Уникальный UUID, проверим на не пустое значение
		},
		{
			name:        "Invalid token",
			token:       "invalid.token",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := middleware.GetUserID(tt.token)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, userID)
			}
		})
	}
}
