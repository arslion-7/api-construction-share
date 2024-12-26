package models

import (
	"time"

	"gorm.io/gorm"
)

type Registry struct {
	ID                  uint           `json:"id" gorm:"primarykey"`
	TB                  *int           `gorm:"column:t_b" json:"t_b"`
	CreatedAt           time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt           gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	GeneralContractorID *uint          `gorm:"column:general_contractor_id" json:"general_contractor_id"`
}
