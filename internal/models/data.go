// Package models содержит модели данных приложения.
package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Data представляет приватные данные пользователя.
type Data struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	Name      string         `json:"name" gorm:"not null"`
	Metadata  string         `json:"metadata"`          // JSON строка с метаданными
	Login     string         `json:"login"`             // Логин
	Password  string         `json:"-" gorm:"not null"` // Зашифрованный пароль
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName возвращает имя таблицы для модели Data.
func (Data) TableName() string {
	return "data"
}

// BeforeCreate выполняется перед созданием записи данных.
func (d *Data) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

// SetMetadata устанавливает метаданные из структуры.
func (d *Data) SetMetadata(metadata interface{}) error {
	jsonData, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	d.Metadata = string(jsonData)
	return nil
}

// GetMetadata извлекает метаданные в указанную структуру.
func (d *Data) GetMetadata(target interface{}) error {
	return json.Unmarshal([]byte(d.Metadata), target)
}
