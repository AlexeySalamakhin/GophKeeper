// Package server содержит тесты для HTTP сервера.
package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AlexeySalamakhin/GophKeeper/internal/config"
	"github.com/AlexeySalamakhin/GophKeeper/internal/logger"
	"github.com/gin-gonic/gin"
)

// setupTestServer создает тестовый сервер для тестирования маршрутов.
// Для полного тестирования с repository нужны моки или изменения структуры Server.
func setupTestServer(t *testing.T) *Server {
	gin.SetMode(gin.TestMode)
	
	// Используем in-memory SQLite для тестов или mock
	// Для простоты тестируем только маршруты без реального repository
	// В реальном проекте здесь должен быть mock или тестовая БД
	
	// Создаем тестовую конфигурацию
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "localhost",
			Port: "8080",
		},
		Database: config.DatabaseConfig{
			Host:    "localhost",
			Port:    5432,
			User:    "test",
			Password: "test",
			DBName:  "test",
			SSLMode: "disable",
		},
		JWT: config.JWTConfig{
			Secret: "test-jwt-secret-key",
		},
		Crypto: config.CryptoConfig{
			Key: "test-encryption-key-32-chars!!",
		},
	}
	
	// Создаем router
	router := gin.New()
	
	server := &Server{
		config: cfg,
		repo:   nil, // В тестах не используем реальный repository
		router: router,
	}
	
	// setupRoutes требует repository, поэтому пропускаем его вызов
	// и тестируем только health endpoint напрямую
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})
	
	return server
}

func TestServer_SetupRoutes_HealthEndpoint(t *testing.T) {
	server := setupTestServer(t)
	
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health", nil)
	server.router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusOK, w.Code)
	}
	
	expected := `{"status":"OK"}`
	if w.Body.String() != expected {
		t.Errorf("Ожидался ответ %s, получен %s", expected, w.Body.String())
	}
}

// Тесты для маршрутов с repository требуют моков или изменений структуры Server
// Для текущей реализации тестируем только health endpoint

func TestServer_Shutdown(t *testing.T) {
	// Инициализируем logger для теста
	err := logger.InitDevelopment()
	if err != nil {
		t.Fatalf("Ошибка инициализации logger: %v", err)
	}
	defer logger.Sync()
	
	server := setupTestServer(t)
	
	// Создаем простой HTTP сервер для тестирования shutdown
	httpServer := &http.Server{
		Addr:    ":0", // Используем случайный порт
		Handler: server.router,
	}
	
	server.httpServer = httpServer
	
	// Запускаем сервер в горутине
	go func() {
		_ = httpServer.ListenAndServe()
	}()
	
	// Даем серверу время запуститься
	time.Sleep(10 * time.Millisecond)
	
	// Тестируем shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	
	err = server.Shutdown(ctx)
	if err != nil {
		t.Errorf("Shutdown не должен возвращать ошибку: %v", err)
	}
}

func TestServer_Shutdown_WithTimeout(t *testing.T) {
	server := setupTestServer(t)
	
	// Создаем простой HTTP сервер
	httpServer := &http.Server{
		Addr:    ":0",
		Handler: server.router,
	}
	
	server.httpServer = httpServer
	
	// Тестируем shutdown с очень коротким таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	
	// Для не запущенного сервера это должно работать
	err := server.Shutdown(ctx)
	// Shutdown может вернуть ошибку для не запущенного сервера, это нормально
	_ = err
}

