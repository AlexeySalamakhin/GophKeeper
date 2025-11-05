// Package handlers содержит тесты для HTTP обработчиков данных.
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

// setupTestDataHandler создает DataHandler с memory repository для тестов.
func setupTestDataHandler(t *testing.T) (*DataHandler, *repository.MemoryRepository, uuid.UUID) {
	gin.SetMode(gin.TestMode)
	memRepo := repository.NewMemoryRepository()
	
	// Создаем тестового пользователя
	userID := uuid.New()
	hashedPassword, _ := auth.HashPassword("testpassword")
	user := &models.User{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
		Password: hashedPassword,
	}
	
	userRepo := memRepo.NewUserRepository()
	userRepo.Create(user)
	
	handler := &DataHandler{
		dataRepo:      memRepo.NewDataRepository(),
		userRepo:      userRepo,
		encryptionKey: "test-encryption-key-32-chars!!",
	}
	
	return handler, memRepo, userID
}

// createAuthenticatedContext создает контекст с аутентифицированным пользователем.
func createAuthenticatedContext(userID uuid.UUID) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", userID.String())
	return c, w
}

func TestDataHandler_GetData_Success(t *testing.T) {
	handler, _, userID := setupTestDataHandler(t)
	
	// Создаем тестовые данные
	dataRepo := handler.dataRepo.(*repository.MemoryDataRepository)
	testData := &models.Data{
		UserID:   userID,
		Name:     "Test Data",
		Login:    "testlogin",
		Password: "encrypted",
	}
	dataRepo.Create(testData)
	
	c, w := createAuthenticatedContext(userID)
	c.Request = httptest.NewRequest("GET", "/api/v1/data", nil)
	
	handler.GetData(c)
	
	if w.Code != http.StatusOK {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusOK, w.Code)
	}
	
	var response []models.Data
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Ошибка парсинга ответа: %v", err)
	}
	
	if len(response) != 1 {
		t.Errorf("Ожидалось 1 запись данных, получено %d", len(response))
	}
}

func TestDataHandler_GetData_Unauthorized(t *testing.T) {
	handler, _, _ := setupTestDataHandler(t)
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/data", nil)
	// Не устанавливаем user_id
	
	handler.GetData(c)
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusUnauthorized, w.Code)
	}
}

func TestDataHandler_GetData_InvalidUserID(t *testing.T) {
	handler, _, _ := setupTestDataHandler(t)
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", "invalid-uuid")
	c.Request = httptest.NewRequest("GET", "/api/v1/data", nil)
	
	handler.GetData(c)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusBadRequest, w.Code)
	}
}

func TestDataHandler_GetDataByID_Success(t *testing.T) {
	handler, _, userID := setupTestDataHandler(t)
	
	// Создаем тестовые данные
	dataRepo := handler.dataRepo.(*repository.MemoryDataRepository)
	testData := &models.Data{
		UserID:   userID,
		Name:     "Test Data",
		Login:    "testlogin",
		Password: "encrypted",
	}
	dataRepo.Create(testData)
	
	c, w := createAuthenticatedContext(userID)
	c.Request = httptest.NewRequest("GET", "/api/v1/data/"+testData.ID.String(), nil)
	c.Params = gin.Params{gin.Param{Key: "id", Value: testData.ID.String()}}
	
	handler.GetDataByID(c)
	
	if w.Code != http.StatusOK {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusOK, w.Code)
	}
	
	var response models.Data
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Ошибка парсинга ответа: %v", err)
	}
	
	if response.Name != "Test Data" {
		t.Errorf("Ожидалось название %s, получено %s", "Test Data", response.Name)
	}
}

func TestDataHandler_GetDataByID_NotFound(t *testing.T) {
	handler, _, userID := setupTestDataHandler(t)
	
	nonExistentID := uuid.New()
	c, w := createAuthenticatedContext(userID)
	c.Request = httptest.NewRequest("GET", "/api/v1/data/"+nonExistentID.String(), nil)
	c.Params = gin.Params{gin.Param{Key: "id", Value: nonExistentID.String()}}
	
	handler.GetDataByID(c)
	
	if w.Code != http.StatusNotFound {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusNotFound, w.Code)
	}
}

func TestDataHandler_GetDataByID_WrongUser(t *testing.T) {
	handler, _, userID := setupTestDataHandler(t)
	
	// Создаем данные для другого пользователя
	dataRepo := handler.dataRepo.(*repository.MemoryDataRepository)
	otherUserID := uuid.New()
	testData := &models.Data{
		UserID:   otherUserID,
		Name:     "Other User Data",
		Login:    "otherlogin",
		Password: "encrypted",
	}
	dataRepo.Create(testData)
	
	c, w := createAuthenticatedContext(userID)
	c.Request = httptest.NewRequest("GET", "/api/v1/data/"+testData.ID.String(), nil)
	c.Params = gin.Params{gin.Param{Key: "id", Value: testData.ID.String()}}
	
	handler.GetDataByID(c)
	
	if w.Code != http.StatusNotFound {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusNotFound, w.Code)
	}
}

func TestDataHandler_CreateData_Success(t *testing.T) {
	handler, _, userID := setupTestDataHandler(t)
	
	req := CreateDataRequest{
		Name:     "New Data",
		Login:    "newlogin",
		Password: "newpassword",
		Metadata: map[string]string{"key": "value"},
	}
	
	jsonData, _ := json.Marshal(req)
	c, w := createAuthenticatedContext(userID)
	c.Request = httptest.NewRequest("POST", "/api/v1/data", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")
	
	handler.CreateData(c)
	
	if w.Code != http.StatusCreated {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusCreated, w.Code)
	}
	
	var response models.Data
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Ошибка парсинга ответа: %v", err)
	}
	
	if response.Name != "New Data" {
		t.Errorf("Ожидалось название %s, получено %s", "New Data", response.Name)
	}
}

func TestDataHandler_CreateData_InvalidJSON(t *testing.T) {
	handler, _, userID := setupTestDataHandler(t)
	
	c, w := createAuthenticatedContext(userID)
	c.Request = httptest.NewRequest("POST", "/api/v1/data", bytes.NewBufferString("invalid json"))
	c.Request.Header.Set("Content-Type", "application/json")
	
	handler.CreateData(c)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusBadRequest, w.Code)
	}
}

func TestDataHandler_CreateData_EmptyName(t *testing.T) {
	handler, _, userID := setupTestDataHandler(t)
	
	req := CreateDataRequest{
		Name:     "",
		Login:    "testlogin",
		Password: "testpassword",
	}
	
	jsonData, _ := json.Marshal(req)
	c, w := createAuthenticatedContext(userID)
	c.Request = httptest.NewRequest("POST", "/api/v1/data", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")
	
	handler.CreateData(c)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusBadRequest, w.Code)
	}
}

func TestDataHandler_UpdateData_Success(t *testing.T) {
	handler, _, userID := setupTestDataHandler(t)
	
	// Создаем тестовые данные
	dataRepo := handler.dataRepo.(*repository.MemoryDataRepository)
	testData := &models.Data{
		UserID:   userID,
		Name:     "Original Name",
		Login:    "originallogin",
		Password: "encrypted",
	}
	dataRepo.Create(testData)
	
	req := UpdateDataRequest{
		Name:     "Updated Name",
		Login:    "updatedlogin",
		Password: "newpassword",
	}
	
	jsonData, _ := json.Marshal(req)
	c, w := createAuthenticatedContext(userID)
	c.Request = httptest.NewRequest("PUT", "/api/v1/data/"+testData.ID.String(), bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{gin.Param{Key: "id", Value: testData.ID.String()}}
	
	handler.UpdateData(c)
	
	if w.Code != http.StatusOK {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusOK, w.Code)
	}
	
	var response models.Data
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Ошибка парсинга ответа: %v", err)
	}
	
	if response.Name != "Updated Name" {
		t.Errorf("Ожидалось название %s, получено %s", "Updated Name", response.Name)
	}
}

func TestDataHandler_UpdateData_Forbidden(t *testing.T) {
	handler, _, userID := setupTestDataHandler(t)
	
	// Создаем данные для другого пользователя
	dataRepo := handler.dataRepo.(*repository.MemoryDataRepository)
	otherUserID := uuid.New()
	testData := &models.Data{
		UserID:   otherUserID,
		Name:     "Other User Data",
		Login:    "otherlogin",
		Password: "encrypted",
	}
	dataRepo.Create(testData)
	
	req := UpdateDataRequest{
		Name: "Updated Name",
	}
	
	jsonData, _ := json.Marshal(req)
	c, w := createAuthenticatedContext(userID)
	c.Request = httptest.NewRequest("PUT", "/api/v1/data/"+testData.ID.String(), bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{gin.Param{Key: "id", Value: testData.ID.String()}}
	
	handler.UpdateData(c)
	
	if w.Code != http.StatusForbidden {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusForbidden, w.Code)
	}
}

func TestDataHandler_DeleteData_Success(t *testing.T) {
	handler, _, userID := setupTestDataHandler(t)
	
	// Создаем тестовые данные
	dataRepo := handler.dataRepo.(*repository.MemoryDataRepository)
	testData := &models.Data{
		UserID:   userID,
		Name:     "Test Data",
		Login:    "testlogin",
		Password: "encrypted",
	}
	dataRepo.Create(testData)
	
	c, w := createAuthenticatedContext(userID)
	c.Request = httptest.NewRequest("DELETE", "/api/v1/data/"+testData.ID.String(), nil)
	c.Params = gin.Params{gin.Param{Key: "id", Value: testData.ID.String()}}
	
	handler.DeleteData(c)
	
	if w.Code != http.StatusNoContent {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusNoContent, w.Code)
	}
}

func TestDataHandler_DeleteData_Forbidden(t *testing.T) {
	handler, _, userID := setupTestDataHandler(t)
	
	// Создаем данные для другого пользователя
	dataRepo := handler.dataRepo.(*repository.MemoryDataRepository)
	otherUserID := uuid.New()
	testData := &models.Data{
		UserID:   otherUserID,
		Name:     "Other User Data",
		Login:    "otherlogin",
		Password: "encrypted",
	}
	dataRepo.Create(testData)
	
	c, w := createAuthenticatedContext(userID)
	c.Request = httptest.NewRequest("DELETE", "/api/v1/data/"+testData.ID.String(), nil)
	c.Params = gin.Params{gin.Param{Key: "id", Value: testData.ID.String()}}
	
	handler.DeleteData(c)
	
	if w.Code != http.StatusForbidden {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusForbidden, w.Code)
	}
}

func TestDataHandler_DeleteData_InvalidID(t *testing.T) {
	handler, _, userID := setupTestDataHandler(t)
	
	c, w := createAuthenticatedContext(userID)
	c.Request = httptest.NewRequest("DELETE", "/api/v1/data/invalid-id", nil)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "invalid-id"}}
	
	handler.DeleteData(c)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusBadRequest, w.Code)
	}
}
