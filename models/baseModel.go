package models

import (
	"time"

	"gorm.io/gorm"
)

// gorm.Model definition
type BaseModel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
