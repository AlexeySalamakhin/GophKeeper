// Package main содержит точку входа для клиента GophKeeper.
package main

import (
	"log"

	"github.com/AlexeySalamakhin/GophKeeper/internal/client"
	"github.com/spf13/viper"
)

func main() {
	viper.SetDefault("config.path", ".gophkeeper")

	viper.AutomaticEnv()

	cli := client.New()
	if err := cli.Execute(); err != nil {
		log.Fatalf("Ошибка выполнения клиента: %v", err)
	}
}
