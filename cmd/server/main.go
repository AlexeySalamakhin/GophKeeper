// Package main содержит точку входа для сервера GophKeeper.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AlexeySalamakhin/GophKeeper/internal/logger"
	"github.com/AlexeySalamakhin/GophKeeper/internal/server"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	if err := logger.Init(logLevel); err != nil {
		panic("Ошибка инициализации логгера: " + err.Error())
	}
	defer logger.Sync()

	viper.AutomaticEnv()

	// Создаем контекст для graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := server.New()

	// Запускаем сервер
	errChan := srv.Run(ctx)

	// Ждем сигнала о завершении или ошибки от сервера
	select {
	case err := <-errChan:
		if err != nil {
			logger.Logger.Fatal("Ошибка запуска сервера",
				zap.Error(err),
			)
			os.Exit(1)
		}
	case <-ctx.Done():
		logger.Logger.Info("Получен сигнал о завершении работы")
	}

	// Graceful shutdown с таймаутом
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Logger.Error("Ошибка при завершении работы сервера",
			zap.Error(err),
		)
		os.Exit(1)
	}

	logger.Logger.Info("Сервер успешно завершил работу")
}
