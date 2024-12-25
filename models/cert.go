package models

import "time"

type Cert struct {
	CertNumber *int       `gorm:"column:cert_number" json:"cert_number"`
	CertDate   *time.Time `gorm:"column:cert_date" json:"cert_date"`
}
