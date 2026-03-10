package http

import (
	logger "github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Error("Authorization header is empty", "authHeader", authHeader)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Пустой заголовок"})
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Error("Authorization header format is invalid", "authHeader", authHeader)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Формат заголовка неверный"})
			c.Abort()
			return
		}
		tokenStr := parts[1]

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			logger.Error("Authorization token is invalid", "token", token)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен"})
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Next()

	}
}
