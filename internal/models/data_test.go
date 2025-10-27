// Package models содержит тесты для моделей данных.
package models

import (
	"encoding/json"
	"testing"
)

func TestData_SetMetadata(t *testing.T) {
	data := &Data{}

	metadata := map[string]interface{}{
		"url":   "https://example.com",
		"notes": "Test notes",
	}

	err := data.SetMetadata(metadata)
	if err != nil {
		t.Fatalf("Ошибка установки метаданных: %v", err)
	}

	if data.Metadata == "" {
		t.Error("Метаданные не должны быть пустыми")
	}

	var retrievedMetadata map[string]interface{}
	err = data.GetMetadata(&retrievedMetadata)
	if err != nil {
		t.Fatalf("Ошибка получения метаданных: %v", err)
	}

	if retrievedMetadata["url"] != metadata["url"] {
		t.Errorf("Ожидался URL %s, получен %s", metadata["url"], retrievedMetadata["url"])
	}

	if retrievedMetadata["notes"] != metadata["notes"] {
		t.Errorf("Ожидались Notes %s, получены %s", metadata["notes"], retrievedMetadata["notes"])
	}
}

func TestData_GetMetadata(t *testing.T) {
	data := &Data{}

	originalMetadata := map[string]interface{}{
		"number":      "1234567890123456",
		"expiry_date": "12/25",
		"cvv":         "123",
		"cardholder":  "John Doe",
		"bank":        "Test Bank",
		"notes":       "Test card",
	}

	jsonData, err := json.Marshal(originalMetadata)
	if err != nil {
		t.Fatalf("Ошибка маршалинга метаданных: %v", err)
	}

	data.Metadata = string(jsonData)

	// Получение метаданных
	var retrievedMetadata map[string]interface{}
	err = data.GetMetadata(&retrievedMetadata)
	if err != nil {
		t.Fatalf("Ошибка получения метаданных: %v", err)
	}

	if retrievedMetadata["number"] != originalMetadata["number"] {
		t.Errorf("Ожидался Number %s, получен %s", originalMetadata["number"], retrievedMetadata["number"])
	}

	if retrievedMetadata["cardholder"] != originalMetadata["cardholder"] {
		t.Errorf("Ожидался Cardholder %s, получен %s", originalMetadata["cardholder"], retrievedMetadata["cardholder"])
	}
}

func TestData_BasicFields(t *testing.T) {
	data := &Data{
		Name:     "Test Data",
		Login:    "testuser",
		Password: "encrypted_password",
	}

	if data.Name != "Test Data" {
		t.Errorf("Ожидалось имя %s, получено %s", "Test Data", data.Name)
	}

	if data.Login != "testuser" {
		t.Errorf("Ожидался логин %s, получен %s", "testuser", data.Login)
	}

	if data.Password != "encrypted_password" {
		t.Errorf("Ожидался пароль %s, получен %s", "encrypted_password", data.Password)
	}
}
