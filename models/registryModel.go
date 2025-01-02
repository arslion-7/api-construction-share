package models

import (
	"time"

	"gorm.io/gorm"
)

type Registry struct {
	ID                  uint               `json:"id" gorm:"primarykey"`
	TB                  *int               `gorm:"column:t_b" json:"t_b"`
	CreatedAt           time.Time          `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time          `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt           gorm.DeletedAt     `json:"deleted_at" gorm:"index"`
	UserID              *uint              `gorm:"column:user_id" json:"user_id"`
	User                *User              `gorm:"foreignKey:UserID" json:"user"`
	GeneralContractorID *uint              `gorm:"column:general_contractor_id" json:"general_contractor_id"`
	GeneralContractor   *GeneralContractor `gorm:"foreignKey:GeneralContractorID" json:"general_contractor"` // Important for preload
	BuildingID          *uint              `gorm:"column:building_id" json:"building_id"`
	Building            *Building          `gorm:"foreignKey:BuildingID" json:"building"` // Important for preload

}
