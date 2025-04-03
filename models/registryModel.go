package models

import (
	"time"

	"gorm.io/gorm"
)

type RegistryDates struct {
	ReviewedAt   *time.Time `json:"reviewed_at"`
	RegisteredAt *time.Time `json:"registered_at"`
}

type RegistryMail struct {
	MailDate     *time.Time `json:"mail_date"`
	MailNumber   *string    `json:"mail_number"`
	DeliveryDate *time.Time `json:"delivered_date"`
	Count        *int       `json:"count"`
	Queue        *int       `json:"queue"`
	MinToMudDate *time.Time `json:"min_to_mud_date"`
}

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
	BuilderID           *uint              `gorm:"column:builder_id" json:"builder_id"`
	Builder             *Builder           `gorm:"foreignKey:BuilderID" json:"builder"` // Important for preload
	ReceiverID          *uint              `gorm:"column:receiver_id" json:"receiver_id"`
	Receiver            *Receiver          `gorm:"foreignKey:ReceiverID" json:"receiver"` // Important for preload
	ShareholderID       *uint              `gorm:"column:shareholder_id" json:"shareholder_id"`
	Shareholder         *Shareholder       `gorm:"foreignKey:ShareholderID" json:"shareholder"`
	RegistryDates
	RegistryMail
}
