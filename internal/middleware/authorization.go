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

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

const TOKENEXP = time.Hour * 24
const SECRETKEY = "supersecretkey"

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
			c.Abort()
			return
		}
		c.Set("userID", userInfo.ID)
		c.Set("new", userInfo.New)
	}
}

func getUserIDFromCookie(c *gin.Context) (*user.User, error) {
	token, err := c.Cookie("userID")
	newToken := false
	if err != nil {
		token, err = BuildJWTString()
		newToken = true
		if err != nil {
			return nil, err
		}
		c.SetCookie("userID", token, 3600, "/", "localhost", false, true)
	}
	userID, err := GetUserID(token)
	if err != nil {
		return nil, err
	}
	userInfo := user.NewUser(userID, newToken)

	return userInfo, nil
}

func BuildJWTString() (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKENEXP)),
		},
		// собственное утверждение
		UserID: uuid.New().String(),
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(SECRETKEY))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}

func GetUserID(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SECRETKEY), nil
		})
	if err != nil {
		return "", fmt.Errorf("token is not valid")
	}

	if !token.Valid {
		return "", fmt.Errorf("token is not valid")
	}

	return claims.UserID, nil
}