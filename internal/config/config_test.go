// Package config содержит тесты для конфигурации.
package config

import (
	"os"
	"testing"
)

func TestConfig_Load_DefaultValues(t *testing.T) {
	// Сохранение оригинальных переменных окружения
	originalHost := os.Getenv("SERVER_HOST")
	originalPort := os.Getenv("SERVER_PORT")
	originalDBHost := os.Getenv("DB_HOST")
	originalDBPort := os.Getenv("DB_PORT")
	originalDBUser := os.Getenv("DB_USER")
	originalDBPassword := os.Getenv("DB_PASSWORD")
	originalDBName := os.Getenv("DB_NAME")
	originalDBSSLMode := os.Getenv("DB_SSLMODE")
	originalJWTSecret := os.Getenv("JWT_SECRET")

	// Очистка переменных окружения для теста
	os.Unsetenv("SERVER_HOST")
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_SSLMODE")
	os.Unsetenv("JWT_SECRET")

	// Восстановление переменных окружения после теста
	defer func() {
		if originalHost != "" {
			os.Setenv("SERVER_HOST", originalHost)
		}
		if originalPort != "" {
			os.Setenv("SERVER_PORT", originalPort)
		}
		if originalDBHost != "" {
			os.Setenv("DB_HOST", originalDBHost)
		}
		if originalDBPort != "" {
			os.Setenv("DB_PORT", originalDBPort)
		}
		if originalDBUser != "" {
			os.Setenv("DB_USER", originalDBUser)
		}
		if originalDBPassword != "" {
			os.Setenv("DB_PASSWORD", originalDBPassword)
		}
		if originalDBName != "" {
			os.Setenv("DB_NAME", originalDBName)
		}
		if originalDBSSLMode != "" {
			os.Setenv("DB_SSLMODE", originalDBSSLMode)
		}
		if originalJWTSecret != "" {
			os.Setenv("JWT_SECRET", originalJWTSecret)
		}
	}()

	config := Load()

	// Проверка значений по умолчанию
	if config.Server.Host != "localhost" {
		t.Errorf("Ожидался host 'localhost', получен '%s'", config.Server.Host)
	}

	if config.Server.Port != "8080" {
		t.Errorf("Ожидался port '8080', получен '%s'", config.Server.Port)
	}

	if config.Database.Host != "localhost" {
		t.Errorf("Ожидался database host 'localhost', получен '%s'", config.Database.Host)
	}

	if config.Database.Port != 5432 {
		t.Errorf("Ожидался database port 5432, получен %d", config.Database.Port)
	}

	if config.Database.User != "gophkeeper" {
		t.Errorf("Ожидался database user 'gophkeeper', получен '%s'", config.Database.User)
	}

	if config.JWT.Secret == "" {
		t.Error("JWT Secret не должен быть пустым")
	}
}

func TestConfig_Load_EnvironmentVariables(t *testing.T) {
	// Установка переменных окружения
	os.Setenv("SERVER_HOST", "test-host")
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("DB_HOST", "test-db-host")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "test-user")
	os.Setenv("DB_PASSWORD", "test-password")
	os.Setenv("DB_NAME", "test-db")
	os.Setenv("DB_SSLMODE", "require")
	os.Setenv("JWT_SECRET", "test-secret")

	// Восстановление после теста
	defer func() {
		os.Unsetenv("SERVER_HOST")
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_SSLMODE")
		os.Unsetenv("JWT_SECRET")
	}()

	config := Load()

	// Проверка значений из переменных окружения
	if config.Server.Host != "test-host" {
		t.Errorf("Ожидался host 'test-host', получен '%s'", config.Server.Host)
	}

	if config.Server.Port != "9090" {
		t.Errorf("Ожидался port '9090', получен '%s'", config.Server.Port)
	}

	if config.Database.Host != "test-db-host" {
		t.Errorf("Ожидался database host 'test-db-host', получен '%s'", config.Database.Host)
	}

	if config.Database.Port != 5433 {
		t.Errorf("Ожидался database port 5433, получен %d", config.Database.Port)
	}

	if config.JWT.Secret != "test-secret" {
		t.Errorf("Ожидался JWT secret 'test-secret', получен '%s'", config.JWT.Secret)
	}
}

func TestConfig_Load_ConfigFile(t *testing.T) {
	// Создание временного файла конфигурации
	configContent := `
server:
  host: "file-host"
  port: "7777"
database:
  host: "file-db-host"
  port: 5434
  user: "file-user"
  password: "file-password"
  dbname: "file-db"
  sslmode: "require"
jwt:
  secret: "file-secret"
`

	configFile := "test_config.yaml"
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Ошибка создания файла конфигурации: %v", err)
	}

	// Восстановление после теста
	defer os.Remove(configFile)

	// Установка переменной окружения для пути к конфигурации
	os.Setenv("CONFIG_PATH", configFile)
	defer os.Unsetenv("CONFIG_PATH")

	config := Load()

	// Проверка значений из файла конфигурации
	if config.Server.Host != "file-host" {
		t.Errorf("Ожидался host 'file-host', получен '%s'", config.Server.Host)
	}

	if config.Server.Port != "7777" {
		t.Errorf("Ожидался port '7777', получен '%s'", config.Server.Port)
	}

	if config.Database.Host != "file-db-host" {
		t.Errorf("Ожидался database host 'file-db-host', получен '%s'", config.Database.Host)
	}

	if config.Database.Port != 5434 {
		t.Errorf("Ожидался database port 5434, получен %d", config.Database.Port)
	}

	if config.JWT.Secret != "file-secret" {
		t.Errorf("Ожидался JWT secret 'file-secret', получен '%s'", config.JWT.Secret)
	}
}
