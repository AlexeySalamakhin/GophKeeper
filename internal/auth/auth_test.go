// Package auth содержит тесты для аутентификации.
package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Ошибка хеширования пароля: %v", err)
	}

	if hash == password {
		t.Error("Хеш пароля не должен совпадать с исходным паролем")
	}

	if len(hash) == 0 {
		t.Error("Хеш пароля не должен быть пустым")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testpassword123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Ошибка хеширования пароля: %v", err)
	}

	// Проверка правильного пароля
	if !CheckPasswordHash(password, hash) {
		t.Error("Проверка правильного пароля должна возвращать true")
	}

	// Проверка неправильного пароля
	if CheckPasswordHash("wrongpassword", hash) {
		t.Error("Проверка неправильного пароля должна возвращать false")
	}
}

func TestJWTManager_GenerateToken(t *testing.T) {
	secretKey := "test-secret-key"
	jm := NewJWTManager(secretKey)

	userID := "test-user-id"
	username := "testuser"

	token, err := jm.GenerateToken(userID, username)
	if err != nil {
		t.Fatalf("Ошибка генерации токена: %v", err)
	}

	if token == "" {
		t.Error("Токен не должен быть пустым")
	}
}

func TestJWTManager_ValidateToken(t *testing.T) {
	secretKey := "test-secret-key"
	jm := NewJWTManager(secretKey)

	userID := "test-user-id"
	username := "testuser"

	// Генерация токена
	token, err := jm.GenerateToken(userID, username)
	if err != nil {
		t.Fatalf("Ошибка генерации токена: %v", err)
	}

	// Валидация токена
	claims, err := jm.ValidateToken(token)
	if err != nil {
		t.Fatalf("Ошибка валидации токена: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Ожидался UserID %s, получен %s", userID, claims.UserID)
	}

	if claims.Username != username {
		t.Errorf("Ожидался Username %s, получен %s", username, claims.Username)
	}
}

func TestJWTManager_ValidateInvalidToken(t *testing.T) {
	secretKey := "test-secret-key"
	jm := NewJWTManager(secretKey)

	// Попытка валидации неверного токена
	_, err := jm.ValidateToken("invalid-token")
	if err == nil {
		t.Error("Валидация неверного токена должна возвращать ошибку")
	}
}

func TestJWTManager_ValidateExpiredToken(t *testing.T) {
	secretKey := "test-secret-key"
	jm := NewJWTManager(secretKey)

	// Создание токена с истекшим сроком действия
	claims := &Claims{
		UserID:   "test-user-id",
		Username: "testuser",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Истек час назад
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		t.Fatalf("Ошибка создания токена: %v", err)
	}

	// Попытка валидации истекшего токена
	_, err = jm.ValidateToken(tokenString)
	if err == nil {
		t.Error("Валидация истекшего токена должна возвращать ошибку")
	}
}
