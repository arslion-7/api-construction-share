package models

import (
	"time"
)

type GeneralContractor struct {
	Org
	CertNumber          *int       `gorm:"column:cert_number" json:"cert_number"`
	CertDate            *time.Time `gorm:"column:cert_date" json:"cert_date"`
	ResolutionCode      *string    `gorm:"column:resolution_code" json:"resolution_code"`
	ResolutionBeginDate *time.Time `gorm:"column:resolution_begin_date" json:"resolution_begin_date"`
	ResolutionEndDate   *time.Time `gorm:"column:resolution_end_date" json:"resolution_end_date"`
}
