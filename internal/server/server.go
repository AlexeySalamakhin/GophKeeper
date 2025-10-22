// Package server содержит логику HTTP сервера GophKeeper.
package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/AlexeySalamakhin/GophKeeper/internal/config"
	"github.com/AlexeySalamakhin/GophKeeper/internal/handlers"
	"github.com/AlexeySalamakhin/GophKeeper/internal/middleware"
	"github.com/AlexeySalamakhin/GophKeeper/internal/repository"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
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

	// 1. Сначала создать базу данных, если не существует:
	ensureDatabase(&cfg.Database)
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
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	return &Server{
		config: cfg,
		repo:   repo,
		router: router,
	}
}

// Run запускает HTTP сервер.
func (s *Server) Run() error {
	s.setupRoutes()

	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port),
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Сервер запущен на %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Завершение работы сервера...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

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

func ensureDatabase(cfg *config.DatabaseConfig) {
	dbName := cfg.DBName
	adminDSN := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=postgres sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.SSLMode,
	)
	targetDSN := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, dbName, cfg.SSLMode,
	)

	// Пробуем открыть целевую БД
	db, err := sql.Open("postgres", targetDSN)
	if err == nil {
		err = db.Ping()
		if err == nil {
			db.Close()
			return // База существует
		}
	}
	// if err != nil && !strings.Contains(err.Error(), "does not exist") && !strings.Contains(err.Error(), "3D000") {
	// 	// Не ошибка существования БД — паникуем
	// 	panic("Ошибка подключения к БД: " + err.Error())
	// }

	// Открываем соединение к postgres и создаём БД
	adminDB, err := sql.Open("postgres", adminDSN)
	if err != nil {
		panic("Ошибка подключения к postgres для создания БД: " + err.Error())
	}
	defer adminDB.Close()

	_, err = adminDB.Exec("CREATE DATABASE " + dbName)
	// Игнорируем ошибку, если БД уже есть
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		panic("Ошибка создания базы данных: " + err.Error())
	}
	// Проверяем, что теперь можем подключиться
	db2, err := sql.Open("postgres", targetDSN)
	if err != nil {
		panic("После создания БД: ошибка подключения к целевой базе: " + err.Error())
	}
	db2.Close()
}
