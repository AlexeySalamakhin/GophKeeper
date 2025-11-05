// Package repository содержит интерфейсы репозиториев для тестирования.
package repository

import (
	"github.com/AlexeySalamakhin/GophKeeper/internal/models"
	"github.com/google/uuid"
)

// UserRepositoryInterface определяет интерфейс для работы с пользователями.
type UserRepositoryInterface interface {
	Create(user *models.User) error
	GetByUsername(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByID(id uuid.UUID) (*models.User, error)
}

// DataRepositoryInterface определяет интерфейс для работы с данными.
type DataRepositoryInterface interface {
	Create(data *models.Data) error
	GetByID(id uuid.UUID) (*models.Data, error)
	GetByUserID(userID uuid.UUID) ([]models.Data, error)
	Update(data *models.Data) error
	Delete(id uuid.UUID) error
	CheckUserOwnership(dataID, userID uuid.UUID) error
}

