/*
Package middleware предоставляет промежуточное ПО для авторизации пользователей с использованием JWT токенов.
Он включает функции для создания, разбора и проверки JWT токенов, а также обработчик промежуточного ПО для Gin,
который проверяет авторизацию пользователей на основе JWT.

Этот пакет определяет следующие ключевые компоненты:
- Claims: Структура пользовательских утверждений для JWT токенов.
- AuthorizationMiddleware: Функция промежуточного ПО для авторизации пользователей на основе JWT cookie.
- Функции для обработки создания и разбора JWT, включая получение ID пользователя из cookie.
*/

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

// Claims представляет структуру утверждений JWT.
type Claims struct {
	jwt.RegisteredClaims        // Стандартные утверждения для JWT
	UserID               string // ID пользователя, связанный с токеном
}

// TOKENEXP определяет время истечения токена (24 часа).
const TOKENEXP = time.Hour * 24

// SECRETKEY — секретный ключ, используемый для подписи JWT токенов.
const SECRETKEY = "supersecretkey"

// AuthorizationMiddleware возвращает промежуточное ПО Gin, которое проверяет авторизацию пользователя.
// Оно извлекает ID пользователя из cookie и устанавливает его в контексте.
// Если пользователь не авторизован, оно отвечает кодом состояния 401.
func AuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userInfo, err := getUserIDFromCookie(c)
		if err != nil {
			code := http.StatusUnauthorized
			contentType := c.Request.Header.Get("Content-Type")
			if contentType == "application/json" {
				c.Header("Content-Type", "application/json")
				c.JSON(code, gin.H{
					"message": fmt.Sprintf("Unauthorized %s", err),
					"code":    code,
				})
			} else {
				c.String(code, fmt.Sprintf("Unauthorized %s", err))
			}
			c.Abort() // Прерывание обработки запроса
			return
		}
		c.Set("userID", userInfo.ID) // Установка ID пользователя в контексте
		c.Set("new", userInfo.New)   // Установка флага нового пользователя в контексте
	}
}

// getUserIDFromCookie извлекает ID пользователя из cookie.
// Если cookie не существует, он создает новый JWT токен и устанавливает его в cookie.
// Возвращает информацию о пользователе и любую ошибку, возникшую в процессе.
func getUserIDFromCookie(c *gin.Context) (*user.User, error) {
	token, err := c.Cookie("userID")
	newToken := false
	if err != nil {
		token, err = BuildJWTString() // Создание нового токена
		newToken = true
		if err != nil {
			return nil, err
		}
		c.SetCookie("userID", token, 3600, "/", "localhost", false, true) // Установка cookie
	}
	userID, err := GetUserID(token) // Извлечение ID пользователя из токена
	if err != nil {
		return nil, err
	}
	userInfo := user.NewUser(userID, newToken) // Создание нового экземпляра пользователя

	return userInfo, nil
}

// BuildJWTString создает новый JWT токен с заданным временем истечения и уникальным ID пользователя.
// Возвращает подписанную строку токена и любую ошибку, возникшую в процессе.
func BuildJWTString() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKENEXP)), // Установка времени истечения
		},
		UserID: uuid.New().String(), // Присвоение нового ID пользователя
	})

	tokenString, err := token.SignedString([]byte(SECRETKEY)) // Подпись токена
	if err != nil {
		return "", err
	}

	return tokenString, nil // Возвращение строки токена
}

// GetUserID извлекает ID пользователя из переданной строки JWT токена.
// Возвращает ID пользователя и любую ошибку, возникшую при разборе.
func GetUserID(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("неожиданный метод подписи: %v", t.Header["alg"])
			}
			return []byte(SECRETKEY), nil // Возвращение секретного ключа для проверки
		})
	if err != nil {
		return "", fmt.Errorf("токен недействителен: %v", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("токен недействителен")
	}

	return claims.UserID, nil // Возвращение извлеченного ID пользователя
}
