// Package handlers содержит HTTP обработчики.
package handlers

import (
	"net/http"

	"github.com/AlexeySalamakhin/GophKeeper/internal/auth"
	"github.com/AlexeySalamakhin/GophKeeper/internal/models"
	"github.com/AlexeySalamakhin/GophKeeper/internal/repository"
	"github.com/gin-gonic/gin"
)

// AuthHandler обрабатывает запросы аутентификации.
type AuthHandler struct {
	userRepo  repository.UserRepositoryInterface
	jwtSecret string
}

// NewAuthHandler создает новый обработчик аутентификации.
func NewAuthHandler(repo *repository.Repository, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		userRepo:  repo.NewUserRepository(),
		jwtSecret: jwtSecret,
	}
}

// RegisterRequest представляет запрос регистрации.
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest представляет запрос входа.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse представляет ответ аутентификации.
type AuthResponse struct {
	Token string `json:"token"`
	User  struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	} `json:"user"`
}

// Register обрабатывает регистрацию нового пользователя.
func (ah *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON"})
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Все поля обязательны"})
		return
	}

	if _, err := ah.userRepo.GetByUsername(req.Username); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Пользователь с таким именем уже существует"})
		return
	}

	if _, err := ah.userRepo.GetByEmail(req.Email); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Пользователь с таким email уже существует"})
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки пароля"})
		return
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := ah.userRepo.Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания пользователя"})
		return
	}

	jwtManager := auth.NewJWTManager(ah.jwtSecret)
	token, err := jwtManager.GenerateToken(user.ID.String(), user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
		return
	}

	response := AuthResponse{
		Token: token,
		User: struct {
			ID       string `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
		}{
			ID:       user.ID.String(),
			Username: user.Username,
			Email:    user.Email,
		},
	}

	c.JSON(http.StatusCreated, response)
}

// Login обрабатывает вход пользователя.
func (ah *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON"})
		return
	}

	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Имя пользователя и пароль обязательны"})
		return
	}

	user, err := ah.userRepo.GetByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные учетные данные"})
		return
	}

	if !auth.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные учетные данные"})
		return
	}

	jwtManager := auth.NewJWTManager(ah.jwtSecret)
	token, err := jwtManager.GenerateToken(user.ID.String(), user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
		return
	}

	response := AuthResponse{
		Token: token,
		User: struct {
			ID       string `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
		}{
			ID:       user.ID.String(),
			Username: user.Username,
			Email:    user.Email,
		},
	}

	c.JSON(http.StatusOK, response)
}
