// Package repository содержит тесты для репозитория.
package repository

import (
	"testing"

	"github.com/AlexeySalamakhin/GophKeeper/internal/models"
	"github.com/google/uuid"
)

func setupTestDB(t *testing.T) *MemoryRepository {

	repo := NewMemoryRepository()
	return repo
}

func cleanupTestDB(t *testing.T, repo *MemoryRepository) {
	// In-memory репозиторий не требует очистки
}

func TestUserRepository_Create(t *testing.T) {
	repo := setupTestDB(t)
	defer cleanupTestDB(t, repo)

	userRepo := repo.NewUserRepository()

	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	err := userRepo.Create(user)
	if err != nil {
		t.Fatalf("Ошибка создания пользователя: %v", err)
	}

	if user.ID == uuid.Nil {
		t.Error("ID пользователя не должен быть пустым")
	}
}

func TestUserRepository_GetByUsername(t *testing.T) {
	repo := setupTestDB(t)
	defer cleanupTestDB(t, repo)

	userRepo := repo.NewUserRepository()

	// Создание тестового пользователя
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	err := userRepo.Create(user)
	if err != nil {
		t.Fatalf("Ошибка создания пользователя: %v", err)
	}

	// Поиск пользователя по имени
	foundUser, err := userRepo.GetByUsername("testuser")
	if err != nil {
		t.Fatalf("Ошибка поиска пользователя: %v", err)
	}

	if foundUser.Username != "testuser" {
		t.Errorf("Ожидался username %s, получен %s", "testuser", foundUser.Username)
	}

	if foundUser.Email != "test@example.com" {
		t.Errorf("Ожидался email %s, получен %s", "test@example.com", foundUser.Email)
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	repo := setupTestDB(t)
	defer cleanupTestDB(t, repo)

	userRepo := repo.NewUserRepository()

	// Создание тестового пользователя
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	err := userRepo.Create(user)
	if err != nil {
		t.Fatalf("Ошибка создания пользователя: %v", err)
	}

	// Поиск пользователя по email
	foundUser, err := userRepo.GetByEmail("test@example.com")
	if err != nil {
		t.Fatalf("Ошибка поиска пользователя: %v", err)
	}

	if foundUser.Email != "test@example.com" {
		t.Errorf("Ожидался email %s, получен %s", "test@example.com", foundUser.Email)
	}
}

func TestDataRepository_Create(t *testing.T) {
	repo := setupTestDB(t)
	defer cleanupTestDB(t, repo)

	dataRepo := repo.NewDataRepository()

	userID := uuid.New()
	data := &models.Data{
		UserID:   userID,
		Name:     "Test Data",
		Login:    "testuser",
		Password: "encrypted_password",
	}

	err := dataRepo.Create(data)
	if err != nil {
		t.Fatalf("Ошибка создания данных: %v", err)
	}

	if data.ID == uuid.Nil {
		t.Error("ID данных не должен быть пустым")
	}
}

func TestDataRepository_GetByUserID(t *testing.T) {
	repo := setupTestDB(t)
	defer cleanupTestDB(t, repo)

	dataRepo := repo.NewDataRepository()
	userID := uuid.New()

	// Создание тестовых данных
	data1 := &models.Data{
		UserID:   userID,
		Name:     "Test Data 1",
		Login:    "user1",
		Password: "encrypted_password_1",
	}

	data2 := &models.Data{
		UserID:   userID,
		Name:     "Test Data 2",
		Login:    "user2",
		Password: "encrypted_password_2",
	}

	err := dataRepo.Create(data1)
	if err != nil {
		t.Fatalf("Ошибка создания данных 1: %v", err)
	}

	err = dataRepo.Create(data2)
	if err != nil {
		t.Fatalf("Ошибка создания данных 2: %v", err)
	}

	userData, err := dataRepo.GetByUserID(userID)
	if err != nil {
		t.Fatalf("Ошибка получения данных пользователя: %v", err)
	}

	if len(userData) != 2 {
		t.Errorf("Ожидалось 2 записи данных, получено %d", len(userData))
	}
}

func TestDataRepository_CheckUserOwnership(t *testing.T) {
	repo := setupTestDB(t)
	defer cleanupTestDB(t, repo)

	dataRepo := repo.NewDataRepository()
	userID := uuid.New()
	otherUserID := uuid.New()

	data := &models.Data{
		UserID:   userID,
		Name:     "Test Data",
		Login:    "testuser",
		Password: "encrypted_password",
	}

	err := dataRepo.Create(data)
	if err != nil {
		t.Fatalf("Ошибка создания данных: %v", err)
	}

	err = dataRepo.CheckUserOwnership(data.ID, userID)
	if err != nil {
		t.Errorf("Ошибка проверки принадлежности данных владельцу: %v", err)
	}

	err = dataRepo.CheckUserOwnership(data.ID, otherUserID)
	if err == nil {
		t.Error("Проверка принадлежности данных другому пользователю должна возвращать ошибку")
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	repo := setupTestDB(t)
	defer cleanupTestDB(t, repo)

	userRepo := repo.NewUserRepository()

	// Создание тестового пользователя
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	err := userRepo.Create(user)
	if err != nil {
		t.Fatalf("Ошибка создания пользователя: %v", err)
	}

	// Поиск пользователя по ID
	foundUser, err := userRepo.GetByID(user.ID)
	if err != nil {
		t.Fatalf("Ошибка поиска пользователя: %v", err)
	}

	if foundUser.ID != user.ID {
		t.Errorf("Ожидался ID %s, получен %s", user.ID, foundUser.ID)
	}

	if foundUser.Username != "testuser" {
		t.Errorf("Ожидался username %s, получен %s", "testuser", foundUser.Username)
	}
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	repo := setupTestDB(t)
	defer cleanupTestDB(t, repo)

	userRepo := repo.NewUserRepository()

	nonExistentID := uuid.New()
	_, err := userRepo.GetByID(nonExistentID)
	if err == nil {
		t.Error("Ожидалась ошибка для несуществующего пользователя")
	}
}

func TestDataRepository_GetByID(t *testing.T) {
	repo := setupTestDB(t)
	defer cleanupTestDB(t, repo)

	dataRepo := repo.NewDataRepository()
	userID := uuid.New()

	data := &models.Data{
		UserID:   userID,
		Name:     "Test Data",
		Login:    "testuser",
		Password: "encrypted_password",
	}

	err := dataRepo.Create(data)
	if err != nil {
		t.Fatalf("Ошибка создания данных: %v", err)
	}

	// Поиск данных по ID
	foundData, err := dataRepo.GetByID(data.ID)
	if err != nil {
		t.Fatalf("Ошибка поиска данных: %v", err)
	}

	if foundData.ID != data.ID {
		t.Errorf("Ожидался ID %s, получен %s", data.ID, foundData.ID)
	}

	if foundData.Name != "Test Data" {
		t.Errorf("Ожидалось название %s, получено %s", "Test Data", foundData.Name)
	}
}

func TestDataRepository_GetByID_NotFound(t *testing.T) {
	repo := setupTestDB(t)
	defer cleanupTestDB(t, repo)

	dataRepo := repo.NewDataRepository()

	nonExistentID := uuid.New()
	_, err := dataRepo.GetByID(nonExistentID)
	if err == nil {
		t.Error("Ожидалась ошибка для несуществующих данных")
	}
}

func TestDataRepository_Update(t *testing.T) {
	repo := setupTestDB(t)
	defer cleanupTestDB(t, repo)

	dataRepo := repo.NewDataRepository()
	userID := uuid.New()

	data := &models.Data{
		UserID:   userID,
		Name:     "Original Name",
		Login:    "originallogin",
		Password: "encrypted_password",
	}

	err := dataRepo.Create(data)
	if err != nil {
		t.Fatalf("Ошибка создания данных: %v", err)
	}

	// Обновление данных
	data.Name = "Updated Name"
	data.Login = "updatedlogin"

	err = dataRepo.Update(data)
	if err != nil {
		t.Fatalf("Ошибка обновления данных: %v", err)
	}

	// Проверка обновленных данных
	updatedData, err := dataRepo.GetByID(data.ID)
	if err != nil {
		t.Fatalf("Ошибка получения обновленных данных: %v", err)
	}

	if updatedData.Name != "Updated Name" {
		t.Errorf("Ожидалось название %s, получено %s", "Updated Name", updatedData.Name)
	}

	if updatedData.Login != "updatedlogin" {
		t.Errorf("Ожидался логин %s, получен %s", "updatedlogin", updatedData.Login)
	}
}

func TestDataRepository_Delete(t *testing.T) {
	repo := setupTestDB(t)
	defer cleanupTestDB(t, repo)

	dataRepo := repo.NewDataRepository()
	userID := uuid.New()

	data := &models.Data{
		UserID:   userID,
		Name:     "Test Data",
		Login:    "testuser",
		Password: "encrypted_password",
	}

	err := dataRepo.Create(data)
	if err != nil {
		t.Fatalf("Ошибка создания данных: %v", err)
	}

	// Удаление данных
	err = dataRepo.Delete(data.ID)
	if err != nil {
		t.Fatalf("Ошибка удаления данных: %v", err)
	}

	// Проверка, что данные удалены
	_, err = dataRepo.GetByID(data.ID)
	if err == nil {
		t.Error("Ожидалась ошибка для удаленных данных")
	}
}

func TestDataRepository_Delete_NotFound(t *testing.T) {
	repo := setupTestDB(t)
	defer cleanupTestDB(t, repo)

	dataRepo := repo.NewDataRepository()

	nonExistentID := uuid.New()
	err := dataRepo.Delete(nonExistentID)
	// Удаление несуществующей записи может не возвращать ошибку в зависимости от реализации
	// Проверяем, что метод выполняется без паники
	_ = err
}