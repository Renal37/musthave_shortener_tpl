package middleware

import (
	"fmt"
	"net/http"
	"time"

	user "github.com/Renal37/musthave_shortener_tpl.git/internal/users"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// Claims описывает данные, которые будут храниться в JWT.
type Claims struct {
	jwt.RegisteredClaims
	UserID string // ID пользователя
}

const (
	TOKENEXP  = time.Hour * 24        // Время жизни токена
	SECRETKEY = "supersecretkey"      // Секретный ключ для подписи токена
)

// AuthorizationMiddleware возвращает middleware для авторизации пользователей.
func AuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем информацию о пользователе из cookie
		userInfo, err := getUserIDFromCookie(c)
		if err != nil {
			code := http.StatusUnauthorized
			contentType := c.Request.Header.Get("Content-Type")
			// Обрабатываем ответ в зависимости от типа контента
			if contentType == "application/json" {
				c.Header("Content-Type", "application/json")
				c.JSON(code, gin.H{
					"message": fmt.Sprintf("Unauthorized %s", err),
					"code":    code,
				})
			} else {
				c.String(code, fmt.Sprintf("Unauthorized %s", err))
			}
			c.Abort() // Прерываем выполнение следующего обработчика
			return
		}
		c.Set("userID", userInfo.ID) // Устанавливаем userID в контексте
		c.Set("new", userInfo.New)    // Устанавливаем признак нового пользователя
	}
}

// getUserIDFromCookie получает ID пользователя из cookie и, при необходимости, создает новый токен.
func getUserIDFromCookie(c *gin.Context) (*user.User, error) {
	token, err := c.Cookie("userID") // Получаем токен из cookie
	newToken := false
	if err != nil {
		token, err = BuildJWTString() // Создаем новый токен, если он отсутствует
		newToken = true
		if err != nil {
			return nil, err
		}
		c.SetCookie("userID", token, 3600, "/", "localhost", false, true) // Устанавливаем новый токен в cookie
	}
	userID, err := GetUserID(token) // Извлекаем ID пользователя из токена
	if err != nil {
		return nil, err
	}
	userInfo := user.NewUser(userID, newToken) // Создаем объект пользователя

	return userInfo, nil
}

// BuildJWTString создает и возвращает новый JWT-токен.
func BuildJWTString() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKENEXP)), // Устанавливаем время истечения токена
		},
		UserID: uuid.New().String(), // Генерируем новый UUID для UserID
	})

	// Создаем строку токена
	tokenString, err := token.SignedString([]byte(SECRETKEY))
	if err != nil {
		return "", err
	}

	return tokenString, nil // Возвращаем строку токена
}

// GetUserID извлекает ID пользователя из JWT-токена.
func GetUserID(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SECRETKEY), nil // Возвращаем секретный ключ
		})
	if err != nil {
		return "", fmt.Errorf("token is not valid") // Проверка на валидность токена
	}

	if !token.Valid {
		return "", fmt.Errorf("token is not valid") // Токен недействителен
	}

	return claims.UserID, nil // Возвращаем ID пользователя
}
