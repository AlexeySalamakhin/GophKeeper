// Package handlers содержит тесты для HTTP обработчиков.
package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuthHandler_Register_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := &AuthHandler{
		userRepo:  nil,
		jwtSecret: "test-secret",
	}

	// Создание Gin контекста для теста
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Установка невалидного JSON
	c.Request = httptest.NewRequest("POST", "/api/v1/register", bytes.NewBufferString("invalid json"))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Register(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusBadRequest, w.Code)
	}
}

func TestAuthHandler_Login_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := &AuthHandler{
		userRepo:  nil,
		jwtSecret: "test-secret",
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("POST", "/api/v1/login", bytes.NewBufferString("invalid json"))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusBadRequest, w.Code)
	}
}

func TestAuthHandler_Register_EmptyFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := &AuthHandler{
		userRepo:  nil,
		jwtSecret: "test-secret",
	}

	req := RegisterRequest{
		Username: "",
		Email:    "",
		Password: "",
	}

	jsonData, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Register(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusBadRequest, w.Code)
	}
}

func TestAuthHandler_Login_EmptyFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := &AuthHandler{
		userRepo:  nil,
		jwtSecret: "test-secret",
	}

	req := LoginRequest{
		Username: "",
		Password: "",
	}

	jsonData, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusBadRequest, w.Code)
	}
}
