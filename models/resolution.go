package models

import "time"

type Resolution struct {
	ResolutionCode      *string    `gorm:"column:resolution_code" json:"resolution_code"`
	ResolutionBeginDate *time.Time `gorm:"column:resolution_begin_date" json:"resolution_begin_date"`
	ResolutionEndDate   *time.Time `gorm:"column:resolution_end_date" json:"resolution_end_date"`
}
