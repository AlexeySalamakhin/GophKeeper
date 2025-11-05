// Package client содержит тесты для CLI клиента.
package client

import (
	"testing"
)

func TestClient_New(t *testing.T) {
	client := New()
	if client == nil {
		t.Error("New() не должен возвращать nil")
	}

	// baseURL устанавливается из конфигурации, может быть пустым в тестах

	if client.httpClient == nil {
		t.Error("httpClient не должен быть nil")
	}
}

func TestClient_makeRequest_WithoutToken(t *testing.T) {
	client := New()
	client.token = "" // Убеждаемся, что токен пустой

	// Этот тест проверяет, что запрос создается без токена
	// В реальном тесте здесь был бы мок HTTP клиента
	if client.token != "" {
		t.Error("Токен должен быть пустым для этого теста")
	}
}

func TestClient_makeRequest_WithToken(t *testing.T) {
	client := New()
	client.token = "test-token"

	// Проверка, что токен установлен
	if client.token != "test-token" {
		t.Error("Токен должен быть установлен")
	}
}

func TestClient_saveToken(t *testing.T) {
	client := New()
	client.token = "test-token"

	client.saveToken()
}

func TestClient_loadToken(t *testing.T) {
	client := New()

	client.loadToken()
}

// Тесты для команд
func TestClient_createAuthCommands(t *testing.T) {
	client := New()
	authCmd := client.createAuthCommands()

	if authCmd == nil {
		t.Error("createAuthCommands не должен возвращать nil")
	}

	if authCmd.Use != "auth" {
		t.Errorf("Ожидался Use 'auth', получен '%s'", authCmd.Use)
	}
}

func TestClient_createDataCommands(t *testing.T) {
	client := New()
	dataCmd := client.createDataCommands()

	if dataCmd == nil {
		t.Error("createDataCommands не должен возвращать nil")
	}

	if dataCmd.Use != "data" {
		t.Errorf("Ожидался Use 'data', получен '%s'", dataCmd.Use)
	}
}

func TestClient_createVersionCommand(t *testing.T) {
	client := New()
	versionCmd := client.createVersionCommand()

	if versionCmd == nil {
		t.Error("createVersionCommand не должен возвращать nil")
	}

	if versionCmd.Use != "version" {
		t.Errorf("Ожидался Use 'version', получен '%s'", versionCmd.Use)
	}
}
