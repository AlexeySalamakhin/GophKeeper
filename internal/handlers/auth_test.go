// Package handlers содержит тесты для HTTP обработчиков.
package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlexeySalamakhin/GophKeeper/internal/auth"
	"github.com/AlexeySalamakhin/GophKeeper/internal/models"
	"github.com/AlexeySalamakhin/GophKeeper/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// setupTestAuthHandler создает AuthHandler с memory repository для тестов.
func setupTestAuthHandler(t *testing.T) *AuthHandler {
	gin.SetMode(gin.TestMode)
	memRepo := repository.NewMemoryRepository()

	handler := &AuthHandler{
		userRepo:  memRepo.NewUserRepository(),
		jwtSecret: "test-secret-key",
	}

	return handler
}

func TestAuthHandler_Register_Success(t *testing.T) {
	handler := setupTestAuthHandler(t)

	req := RegisterRequest{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "password123",
	}

	jsonData, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Register(c)

	if w.Code != http.StatusCreated {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusCreated, w.Code)
	}

	var response AuthResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Ошибка парсинга ответа: %v", err)
	}

	if response.Token == "" {
		t.Error("Токен не должен быть пустым")
	}

	if response.User.Username != "newuser" {
		t.Errorf("Ожидался username %s, получен %s", "newuser", response.User.Username)
	}
}

func TestAuthHandler_Register_DuplicateUsername(t *testing.T) {
	handler := setupTestAuthHandler(t)

	// Создаем первого пользователя
	userID := uuid.New()
	hashedPassword, _ := auth.HashPassword("password123")
	user := &models.User{
		ID:       userID,
		Username: "existinguser",
		Email:    "existing@example.com",
		Password: hashedPassword,
	}

	userRepo := handler.userRepo.(*repository.MemoryUserRepository)
	userRepo.Create(user)

	// Пытаемся зарегистрировать пользователя с таким же username
	req := RegisterRequest{
		Username: "existinguser",
		Email:    "newemail@example.com",
		Password: "password123",
	}

	jsonData, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Register(c)

	if w.Code != http.StatusConflict {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusConflict, w.Code)
	}
}

func TestAuthHandler_Register_DuplicateEmail(t *testing.T) {
	handler := setupTestAuthHandler(t)

	// Создаем первого пользователя
	userID := uuid.New()
	hashedPassword, _ := auth.HashPassword("password123")
	user := &models.User{
		ID:       userID,
		Username: "user1",
		Email:    "existing@example.com",
		Password: hashedPassword,
	}

	userRepo := handler.userRepo.(*repository.MemoryUserRepository)
	userRepo.Create(user)

	// Пытаемся зарегистрировать пользователя с таким же email
	req := RegisterRequest{
		Username: "user2",
		Email:    "existing@example.com",
		Password: "password123",
	}

	jsonData, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Register(c)

	if w.Code != http.StatusConflict {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusConflict, w.Code)
	}
}

func TestAuthHandler_Login_Success(t *testing.T) {
	handler := setupTestAuthHandler(t)

	// Создаем пользователя
	userID := uuid.New()
	hashedPassword, _ := auth.HashPassword("password123")
	user := &models.User{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
		Password: hashedPassword,
	}

	userRepo := handler.userRepo.(*repository.MemoryUserRepository)
	userRepo.Create(user)

	req := LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	jsonData, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	if w.Code != http.StatusOK {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusOK, w.Code)
	}

	var response AuthResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Ошибка парсинга ответа: %v", err)
	}

	if response.Token == "" {
		t.Error("Токен не должен быть пустым")
	}

	if response.User.Username != "testuser" {
		t.Errorf("Ожидался username %s, получен %s", "testuser", response.User.Username)
	}
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	handler := setupTestAuthHandler(t)

	// Создаем пользователя
	userID := uuid.New()
	hashedPassword, _ := auth.HashPassword("password123")
	user := &models.User{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
		Password: hashedPassword,
	}

	userRepo := handler.userRepo.(*repository.MemoryUserRepository)
	userRepo.Create(user)

	req := LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	jsonData, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthHandler_Login_UserNotFound(t *testing.T) {
	handler := setupTestAuthHandler(t)

	req := LoginRequest{
		Username: "nonexistentuser",
		Password: "password123",
	}

	jsonData, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusUnauthorized, w.Code)
	}
}
