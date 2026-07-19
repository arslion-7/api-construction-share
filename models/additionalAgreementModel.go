package models

import (
	"time"
)

// AdditionalAgreement represents additional agreements (goşmaça şertnama) for registries
type AdditionalAgreement struct {
	BaseModel
	RegistryID      uint       `gorm:"column:registry_id;type:bigint;not null;index" json:"registry_id"`
	AgreementNumber string     `gorm:"column:agreement_number;type:varchar(255)" json:"agreement_number"`
	AgreementDate   *time.Time `gorm:"column:agreement_date;type:date" json:"agreement_date"`
	Reason          string     `gorm:"column:reason;type:varchar(255)" json:"reason"`
	AdditionalInfo  string     `gorm:"column:additional_info;type:text" json:"additional_info"`
	Registry        Registry   `gorm:"foreignKey:RegistryID;constraint:OnDelete:CASCADE" json:"-"`
	// set when imported from the legacy g2 CSV: source row t_b (unique per import)
	OldG2TB         *int       `gorm:"column:old_g2_tb;index" json:"old_g2_tb"`
	FromOldRegistry bool       `gorm:"column:from_old_registry;default:false;index" json:"from_old_registry"`
}

// TableName specifies the table name for AdditionalAgreement
func (AdditionalAgreement) TableName() string {
	return "additional_agreements"
}
