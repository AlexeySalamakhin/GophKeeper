// Package models содержит модели данных приложения.
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User представляет пользователя системы.
type User struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"` // Хеш пароля, не возвращается в JSON
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName возвращает имя таблицы для модели User.
func (User) TableName() string {
	return "users"
}

// BeforeCreate выполняется перед созданием записи пользователя.
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}



