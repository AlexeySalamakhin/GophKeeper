// Package repository содержит слой доступа к данным.
package repository

import (
	"errors"

	"github.com/AlexeySalamakhin/GophKeeper/internal/models"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Repository представляет слой доступа к данным.
type Repository struct {
	db *gorm.DB
}

// New создает новый экземпляр репозитория.
func New(databaseURL string) *Repository {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		panic("Ошибка подключения к базе данных: " + err.Error())
	}

	// Автомиграция схемы
	if err := db.AutoMigrate(&models.User{}, &models.Data{}); err != nil {
		panic("Ошибка миграции базы данных: " + err.Error())
	}

	return &Repository{db: db}
}

// UserRepository содержит методы для работы с пользователями.
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository создает новый репозиторий пользователей.
func (r *Repository) NewUserRepository() *UserRepository {
	return &UserRepository{db: r.db}
}

// Create создает нового пользователя.
func (ur *UserRepository) Create(user *models.User) error {
	return ur.db.Create(user).Error
}

// GetByUsername возвращает пользователя по имени.
func (ur *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := ur.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail возвращает пользователя по email.
func (ur *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := ur.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByID возвращает пользователя по ID.
func (ur *UserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := ur.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// DataRepository содержит методы для работы с данными пользователей.
type DataRepository struct {
	db *gorm.DB
}

// NewDataRepository создает новый репозиторий данных.
func (r *Repository) NewDataRepository() *DataRepository {
	return &DataRepository{db: r.db}
}

// Create создает новую запись данных.
func (dr *DataRepository) Create(data *models.Data) error {
	return dr.db.Create(data).Error
}

// GetByID возвращает данные по ID.
func (dr *DataRepository) GetByID(id uuid.UUID) (*models.Data, error) {
	var data models.Data
	err := dr.db.Where("id = ?", id).First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// GetByUserID возвращает все данные пользователя.
func (dr *DataRepository) GetByUserID(userID uuid.UUID) ([]models.Data, error) {
	var data []models.Data
	err := dr.db.Where("user_id = ?", userID).Find(&data).Error
	return data, err
}

// Update обновляет данные.
func (dr *DataRepository) Update(data *models.Data) error {
	return dr.db.Save(data).Error
}

// Delete удаляет данные.
func (dr *DataRepository) Delete(id uuid.UUID) error {
	return dr.db.Delete(&models.Data{}, id).Error
}

// GetLatestVersion возвращает последнюю версию данных.
func (dr *DataRepository) GetLatestVersion(userID uuid.UUID) (int, error) {
	return 0, nil
}

// GetUpdatedSince возвращает данные, обновленные после указанной версии.
func (dr *DataRepository) GetUpdatedSince(userID uuid.UUID, version int) ([]models.Data, error) {
	return []models.Data{}, nil
}

// CheckUserOwnership проверяет, принадлежат ли данные пользователю.
func (dr *DataRepository) CheckUserOwnership(dataID, userID uuid.UUID) error {
	var count int64
	err := dr.db.Model(&models.Data{}).
		Where("id = ? AND user_id = ?", dataID, userID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("данные не принадлежат пользователю")
	}
	return nil
}
