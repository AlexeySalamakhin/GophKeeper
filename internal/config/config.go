// Package config содержит конфигурацию приложения.
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config представляет конфигурацию приложения.
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Crypto   CryptoConfig   `mapstructure:"crypto"`
}

// ServerConfig содержит настройки HTTP сервера.
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

// DatabaseConfig содержит настройки базы данных.
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// JWTConfig содержит настройки JWT токенов.
type JWTConfig struct {
	Secret string `mapstructure:"secret"`
}

// CryptoConfig содержит настройки шифрования.
type CryptoConfig struct {
	Key string `mapstructure:"key"`
}

// Load загружает конфигурацию из переменных окружения и файлов.
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используются системные переменные окружения")
	}

	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.sslmode", "disable")

	viper.AutomaticEnv()

	viper.SetEnvPrefix("")
	viper.BindEnv("server.host", "SERVER_HOST")
	viper.BindEnv("server.port", "SERVER_PORT")
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.dbname", "DB_NAME")
	viper.BindEnv("database.sslmode", "DB_SSLMODE")
	viper.BindEnv("jwt.secret", "JWT_SECRET")
	viper.BindEnv("crypto.key", "CRYPTO_KEY")

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	if _, err := os.Stat(configPath); err == nil {
		viper.SetConfigFile(configPath)
		if err := viper.ReadInConfig(); err != nil {
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic("Ошибка парсинга конфигурации: " + err.Error())
	}

	// Валидация критичных полей безопасности
	validateConfig(&config)

	return &config
}

func validateConfig(cfg *Config) {
	if cfg.JWT.Secret == "" {
		log.Fatalf("КРИТИЧЕСКАЯ ОШИБКА: JWT_SECRET не установлен.")
	}

	if cfg.Crypto.Key == "" {
		log.Fatalf("КРИТИЧЕСКАЯ ОШИБКА: CRYPTO_KEY не установлен.")
	}

	if cfg.Database.DBName == "" {
		log.Fatalf("КРИТИЧЕСКАЯ ОШИБКА: DB_NAME не установлен.")
	}

	if cfg.Database.User == "" {
		log.Fatalf("КРИТИЧЕСКАЯ ОШИБКА: DB_USER не установлен.")
	}

	if cfg.Database.Password == "" {
		log.Fatalf("КРИТИЧЕСКАЯ ОШИБКА: DB_PASSWORD не установлен.")
	}
}
