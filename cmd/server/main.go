// Package main содержит точку входа для сервера GophKeeper.
package main

import (
	"os"

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

	srv := server.New()
	if err := srv.Run(); err != nil {
		logger.Logger.Fatal("Ошибка запуска сервера",
			zap.Error(err),
		)
		os.Exit(1)
	}
}
