// Package repository содержит in-memory реализацию репозитория для тестов.
package repository

import (
	"errors"
	"sync"

	"github.com/AlexeySalamakhin/GophKeeper/internal/models"
	"github.com/google/uuid"
)

// MemoryRepository представляет in-memory репозиторий для тестов.
type MemoryRepository struct {
	users map[uuid.UUID]*models.User
	data  map[uuid.UUID]*models.Data
	mutex sync.RWMutex
}

// NewMemoryRepository создает новый in-memory репозиторий.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		users: make(map[uuid.UUID]*models.User),
		data:  make(map[uuid.UUID]*models.Data),
	}
}

// UserRepository содержит методы для работы с пользователями.
type MemoryUserRepository struct {
	repo *MemoryRepository
}

// NewUserRepository создает новый репозиторий пользователей.
func (mr *MemoryRepository) NewUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{repo: mr}
}

// Create создает нового пользователя.
func (mur *MemoryUserRepository) Create(user *models.User) error {
	mur.repo.mutex.Lock()
	defer mur.repo.mutex.Unlock()

	// Проверка на дубликаты
	for _, existingUser := range mur.repo.users {
		if existingUser.Username == user.Username {
			return errors.New("пользователь с таким именем уже существует")
		}
		if existingUser.Email == user.Email {
			return errors.New("пользователь с таким email уже существует")
		}
	}

	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	mur.repo.users[user.ID] = user
	return nil
}

// GetByUsername возвращает пользователя по имени.
func (mur *MemoryUserRepository) GetByUsername(username string) (*models.User, error) {
	mur.repo.mutex.RLock()
	defer mur.repo.mutex.RUnlock()

	for _, user := range mur.repo.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, errors.New("пользователь не найден")
}

// GetByEmail возвращает пользователя по email.
func (mur *MemoryUserRepository) GetByEmail(email string) (*models.User, error) {
	mur.repo.mutex.RLock()
	defer mur.repo.mutex.RUnlock()

	for _, user := range mur.repo.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("пользователь не найден")
}

// GetByID возвращает пользователя по ID.
func (mur *MemoryUserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	mur.repo.mutex.RLock()
	defer mur.repo.mutex.RUnlock()

	user, exists := mur.repo.users[id]
	if !exists {
		return nil, errors.New("пользователь не найден")
	}
	return user, nil
}

// DataRepository содержит методы для работы с данными.
type MemoryDataRepository struct {
	repo *MemoryRepository
}

// NewDataRepository создает новый репозиторий данных.
func (mr *MemoryRepository) NewDataRepository() *MemoryDataRepository {
	return &MemoryDataRepository{repo: mr}
}

// Create создает новую запись данных.
func (mdr *MemoryDataRepository) Create(data *models.Data) error {
	mdr.repo.mutex.Lock()
	defer mdr.repo.mutex.Unlock()

	if data.ID == uuid.Nil {
		data.ID = uuid.New()
	}

	mdr.repo.data[data.ID] = data
	return nil
}

// GetByID возвращает данные по ID.
func (mdr *MemoryDataRepository) GetByID(id uuid.UUID) (*models.Data, error) {
	mdr.repo.mutex.RLock()
	defer mdr.repo.mutex.RUnlock()

	data, exists := mdr.repo.data[id]
	if !exists {
		return nil, errors.New("данные не найдены")
	}
	return data, nil
}

// GetByUserID возвращает все данные пользователя.
func (mdr *MemoryDataRepository) GetByUserID(userID uuid.UUID) ([]models.Data, error) {
	mdr.repo.mutex.RLock()
	defer mdr.repo.mutex.RUnlock()

	var userData []models.Data
	for _, data := range mdr.repo.data {
		if data.UserID == userID {
			userData = append(userData, *data)
		}
	}
	return userData, nil
}

// Update обновляет данные.
func (mdr *MemoryDataRepository) Update(data *models.Data) error {
	mdr.repo.mutex.Lock()
	defer mdr.repo.mutex.Unlock()

	_, exists := mdr.repo.data[data.ID]
	if !exists {
		return errors.New("данные не найдены")
	}

	mdr.repo.data[data.ID] = data
	return nil
}

// Delete удаляет данные.
func (mdr *MemoryDataRepository) Delete(id uuid.UUID) error {
	mdr.repo.mutex.Lock()
	defer mdr.repo.mutex.Unlock()

	_, exists := mdr.repo.data[id]
	if !exists {
		return errors.New("данные не найдены")
	}

	delete(mdr.repo.data, id)
	return nil
}

// CheckUserOwnership проверяет, принадлежат ли данные пользователю.
func (mdr *MemoryDataRepository) CheckUserOwnership(dataID, userID uuid.UUID) error {
	mdr.repo.mutex.RLock()
	defer mdr.repo.mutex.RUnlock()

	data, exists := mdr.repo.data[dataID]
	if !exists {
		return errors.New("данные не найдены")
	}

	if data.UserID != userID {
		return errors.New("данные не принадлежат пользователю")
	}

	return nil
}
