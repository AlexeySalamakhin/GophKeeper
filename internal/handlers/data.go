// Package handlers содержит HTTP обработчики.
package handlers

import (
	"net/http"

	"github.com/AlexeySalamakhin/GophKeeper/internal/crypto"
	"github.com/AlexeySalamakhin/GophKeeper/internal/middleware"
	"github.com/AlexeySalamakhin/GophKeeper/internal/models"
	"github.com/AlexeySalamakhin/GophKeeper/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DataHandler обрабатывает запросы для работы с данными.
type DataHandler struct {
	dataRepo      repository.DataRepositoryInterface
	userRepo      repository.UserRepositoryInterface
	encryptionKey string
}

// NewDataHandler создает новый обработчик данных.
func NewDataHandler(repo *repository.Repository, encryptionKey string) *DataHandler {
	return &DataHandler{
		dataRepo:      repo.NewDataRepository(),
		userRepo:      repo.NewUserRepository(),
		encryptionKey: encryptionKey,
	}
}

// CreateDataRequest представляет запрос создания данных.
type CreateDataRequest struct {
	Name     string      `json:"name"`
	Login    string      `json:"login"`
	Password string      `json:"password"`
	Metadata interface{} `json:"metadata"`
}

// UpdateDataRequest представляет запрос обновления данных.
type UpdateDataRequest struct {
	Name     string      `json:"name"`
	Login    string      `json:"login"`
	Password string      `json:"password"`
	Metadata interface{} `json:"metadata"`
}

// GetData возвращает все данные пользователя.
func (dh *DataHandler) GetData(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не аутентифицирован"})
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID пользователя"})
		return
	}

	data, err := dh.dataRepo.GetByUserID(userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetDataByID возвращает данные по ID.
func (dh *DataHandler) GetDataByID(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не аутентифицирован"})
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID пользователя"})
		return
	}

	idStr := c.Param("id")
	dataID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID данных"})
		return
	}

	data, err := dh.dataRepo.GetByID(dataID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Данные не найдены"})
		return
	}

	if data.UserID != userUUID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Данные не найдены"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// CreateData создает новые данные.
func (dh *DataHandler) CreateData(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не аутентифицирован"})
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID пользователя"})
		return
	}

	if _, err := dh.userRepo.GetByID(userUUID); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
		return
	}

	var req CreateDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON"})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Название обязательно"})
		return
	}

	var encryptedPassword string
	if req.Password != "" {
		var err error
		encryptedPassword, err = crypto.EncryptPassword(req.Password, dh.encryptionKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка шифрования пароля"})
			return
		}
	}

	data := &models.Data{
		UserID:   userUUID,
		Name:     req.Name,
		Login:    req.Login,
		Password: encryptedPassword,
	}

	if req.Metadata != nil {
		if err := data.SetMetadata(req.Metadata); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки метаданных"})
			return
		}
	}

	if err := dh.dataRepo.Create(data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания данных"})
		return
	}

	c.JSON(http.StatusCreated, data)
}

// UpdateData обновляет существующие данные.
func (dh *DataHandler) UpdateData(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не аутентифицирован"})
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID пользователя"})
		return
	}

	idStr := c.Param("id")
	dataID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID данных"})
		return
	}

	if err := dh.dataRepo.CheckUserOwnership(dataID, userUUID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Данные не найдены"})
		return
	}

	var req UpdateDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON"})
		return
	}

	data, err := dh.dataRepo.GetByID(dataID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Данные не найдены"})
		return
	}

	if req.Name != "" {
		data.Name = req.Name
	}
	if req.Login != "" {
		data.Login = req.Login
	}
	if req.Password != "" {
		encryptedPassword, err := crypto.EncryptPassword(req.Password, dh.encryptionKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка шифрования пароля"})
			return
		}
		data.Password = encryptedPassword
	}
	if req.Metadata != nil {
		if err := data.SetMetadata(req.Metadata); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки метаданных"})
			return
		}
	}

	if err := dh.dataRepo.Update(data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления данных"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// DeleteData удаляет данные.
func (dh *DataHandler) DeleteData(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не аутентифицирован"})
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID пользователя"})
		return
	}

	idStr := c.Param("id")
	dataID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID данных"})
		return
	}

	if err := dh.dataRepo.CheckUserOwnership(dataID, userUUID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Данные не найдены"})
		return
	}

	if err := dh.dataRepo.Delete(dataID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления данных"})
		return
	}

	c.Data(http.StatusNoContent, "application/json", nil)
}
