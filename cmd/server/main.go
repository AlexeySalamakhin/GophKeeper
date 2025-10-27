// Package main содержит точку входа для сервера GophKeeper.
package main

import (
	"log"
	"os"

	"github.com/AlexeySalamakhin/GophKeeper/internal/server"
	"github.com/spf13/viper"
)

func main() {
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("database.path", "gophkeeper.db")
	viper.SetDefault("jwt.secret", "your-secret-key")

	viper.AutomaticEnv()

	srv := server.New()
	if err := srv.Run(); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
		os.Exit(1)
	}
}
