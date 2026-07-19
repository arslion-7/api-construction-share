package models

type GeneralContractor struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"` // Correct annotation
	TB *int `gorm:"column:t_b" json:"t_b"`
	Contractor
	Registries []Registry `gorm:"foreignKey:GeneralContractorID;references:ID" json:"registries"`
	// true when this record was created by the old-registries migration
	FromOldRegistry bool `gorm:"column:from_old_registry;default:false;index" json:"from_old_registry"`
}
