// Package server содержит логику HTTP сервера GophKeeper.
package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/AlexeySalamakhin/GophKeeper/internal/config"
	"github.com/AlexeySalamakhin/GophKeeper/internal/handlers"
	"github.com/AlexeySalamakhin/GophKeeper/internal/logger"
	"github.com/AlexeySalamakhin/GophKeeper/internal/middleware"
	"github.com/AlexeySalamakhin/GophKeeper/internal/repository"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx"
	"go.uber.org/zap"
)

// Server представляет HTTP сервер приложения.
type Server struct {
	httpServer *http.Server
	config     *config.Config
	repo       *repository.Repository
	router     *gin.Engine
}

// New создает новый экземпляр сервера.
func New() *Server {
	cfg := config.Load()

	// Формирование URL для подключения к PostgreSQL
	databaseURL := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	repo := repository.New(databaseURL)

	// Настройка Gin роутера
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(logger.GinLoggerMiddleware())
	router.Use(logger.GinRecoveryMiddleware())

	return &Server{
		config: cfg,
		repo:   repo,
		router: router,
	}
}

// Run запускает HTTP сервер в горутине и возвращает канал ошибок.
func (s *Server) Run(ctx context.Context) <-chan error {
	s.setupRoutes()

	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port),
		Handler:      s.router,
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	errChan := make(chan error, 1)

	go func() {
		logger.Logger.Info("Сервер запущен",
			zap.String("address", s.httpServer.Addr),
		)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
		close(errChan)
	}()

	return errChan
}

// Shutdown выполняет graceful shutdown сервера.
func (s *Server) Shutdown(ctx context.Context) error {
	logger.Logger.Info("Завершение работы сервера...")
	return s.httpServer.Shutdown(ctx)
}

// setupRoutes настраивает маршруты HTTP сервера.
func (s *Server) setupRoutes() {
	authHandler := handlers.NewAuthHandler(s.repo, s.config.JWT.Secret)
	dataHandler := handlers.NewDataHandler(s.repo, s.config.Crypto.Key)

	api := s.router.Group("/api/v1")
	{
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(s.config.JWT.Secret))
		{
			protected.GET("/data", dataHandler.GetData)
			protected.GET("/data/:id", dataHandler.GetDataByID)
			protected.POST("/data", dataHandler.CreateData)
			protected.PUT("/data/:id", dataHandler.UpdateData)
			protected.DELETE("/data/:id", dataHandler.DeleteData)
		}
	}

	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})
}
