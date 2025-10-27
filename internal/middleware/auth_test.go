// Package middleware содержит тесты для HTTP middleware.
package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlexeySalamakhin/GophKeeper/internal/auth"
	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtManager := auth.NewJWTManager("test-secret")
	token, err := jwtManager.GenerateToken("test-user-id", "testuser")
	if err != nil {
		t.Fatalf("Ошибка генерации токена: %v", err)
	}

	middleware := AuthMiddleware("test-secret")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	handler := func(c *gin.Context) {
		userID, ok := GetUserID(c)
		if !ok {
			t.Error("UserID не найден в контексте")
		}
		if userID != "test-user-id" {
			t.Errorf("Ожидался UserID %s, получен %s", "test-user-id", userID)
		}

		username, ok := GetUsername(c)
		if !ok {
			t.Error("Username не найден в контексте")
		}
		if username != "testuser" {
			t.Errorf("Ожидался Username %s, получен %s", "testuser", username)
		}

		claims, ok := GetClaims(c)
		if !ok {
			t.Error("Claims не найдены в контексте")
		}
		if claims.UserID != "test-user-id" {
			t.Errorf("Ожидался UserID в claims %s, получен %s", "test-user-id", claims.UserID)
		}

		c.Status(http.StatusOK)
	}

	middleware(c)
	if !c.IsAborted() {
		handler(c)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusOK, w.Code)
	}
}

func TestAuthMiddleware_NoToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	middleware := AuthMiddleware("test-secret")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/test", nil)

	middleware(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	middleware := AuthMiddleware("test-secret")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer invalid-token")

	middleware(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_WrongFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	middleware := AuthMiddleware("test-secret")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "InvalidFormat token")

	middleware(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusUnauthorized, w.Code)
	}

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "just-token")

	middleware(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusUnauthorized, w.Code)
	}
}

func TestGetUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	userID, ok := GetUserID(c)
	if ok {
		t.Error("GetUserID должен возвращать false для пустого контекста")
	}

	c.Set("user_id", "test-user-id")
	userID, ok = GetUserID(c)
	if !ok {
		t.Error("GetUserID должен возвращать true для контекста с UserID")
	}
	if userID != "test-user-id" {
		t.Errorf("Ожидался UserID %s, получен %s", "test-user-id", userID)
	}
}

func TestGetUsername(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	username, ok := GetUsername(c)
	if ok {
		t.Error("GetUsername должен возвращать false для пустого контекста")
	}

	c.Set("username", "testuser")
	username, ok = GetUsername(c)
	if !ok {
		t.Error("GetUsername должен возвращать true для контекста с Username")
	}
	if username != "testuser" {
		t.Errorf("Ожидался Username %s, получен %s", "testuser", username)
	}
}

func TestGetClaims(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	claims, ok := GetClaims(c)
	if ok {
		t.Error("GetClaims должен возвращать false для пустого контекста")
	}

	testClaims := &auth.Claims{
		UserID:   "test-user-id",
		Username: "testuser",
	}
	c.Set("claims", testClaims)
	claims, ok = GetClaims(c)
	if !ok {
		t.Error("GetClaims должен возвращать true для контекста с Claims")
	}
	if claims.UserID != "test-user-id" {
		t.Errorf("Ожидался UserID в claims %s, получен %s", "test-user-id", claims.UserID)
	}
}
