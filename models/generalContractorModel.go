package models

type GeneralContractor struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"` // Correct annotation
	TB *int `gorm:"column:t_b" json:"t_b"`
	Contractor
	Registries []Registry `gorm:"foreignKey:GeneralContractorID;references:ID" json:"registries"`
}
