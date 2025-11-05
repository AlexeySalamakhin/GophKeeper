// Package logger содержит middleware для Gin с zap логированием.
package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinLoggerMiddleware создает middleware для логирования HTTP запросов через zap.
func GinLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		userAgent := c.Request.UserAgent()

		fields := []zap.Field{
			zap.Int("status", statusCode),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", clientIP),
			zap.String("user-agent", userAgent),
			zap.Duration("latency", latency),
			zap.Int("size", c.Writer.Size()),
		}

		if len(c.Errors) > 0 {
			fields = append(fields, zap.Strings("errors", c.Errors.Errors()))
		}

		switch {
		case statusCode >= 500:
			Logger.Error("HTTP Request", fields...)
		case statusCode >= 400:
			Logger.Warn("HTTP Request", fields...)
		default:
			Logger.Info("HTTP Request", fields...)
		}
	}
}

// GinRecoveryMiddleware создает middleware для восстановления после паники с логированием через zap.
func GinRecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				Logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.String("ip", c.ClientIP()),
				)

				c.JSON(500, gin.H{
					"error": "Internal server error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
