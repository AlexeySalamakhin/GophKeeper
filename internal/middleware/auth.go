// Package middleware содержит HTTP middleware.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/AlexeySalamakhin/GophKeeper/internal/auth"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware создает middleware для проверки JWT токенов.
func AuthMiddleware(secretKey string) gin.HandlerFunc {
	jwtManager := auth.NewJWTManager(secretKey)

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Отсутствует заголовок Authorization"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный формат заголовка Authorization"})
			c.Abort()
			return
		}

		token := parts[1]

		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен: " + err.Error()})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("claims", claims)

		c.Next()
	}
}

// GetUserID извлекает ID пользователя из контекста Gin.
func GetUserID(c *gin.Context) (string, bool) {
	userID, ok := c.Get("user_id")
	if !ok {
		return "", false
	}
	userIDStr, ok := userID.(string)
	return userIDStr, ok
}

// GetUsername извлекает имя пользователя из контекста Gin.
func GetUsername(c *gin.Context) (string, bool) {
	username, ok := c.Get("username")
	if !ok {
		return "", false
	}
	usernameStr, ok := username.(string)
	return usernameStr, ok
}

// GetClaims извлекает claims из контекста Gin.
func GetClaims(c *gin.Context) (*auth.Claims, bool) {
	claims, ok := c.Get("claims")
	if !ok {
		return nil, false
	}
	claimsObj, ok := claims.(*auth.Claims)
	return claimsObj, ok
}

// GetUserIDFromContext извлекает ID пользователя из контекста (для обратной совместимости).
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value("user_id").(string)
	return userID, ok
}

// GetUsernameFromContext извлекает имя пользователя из контекста (для обратной совместимости).
func GetUsernameFromContext(ctx context.Context) (string, bool) {
	username, ok := ctx.Value("username").(string)
	return username, ok
}

// GetClaimsFromContext извлекает claims из контекста (для обратной совместимости).
func GetClaimsFromContext(ctx context.Context) (*auth.Claims, bool) {
	claims, ok := ctx.Value("claims").(*auth.Claims)
	return claims, ok
}
