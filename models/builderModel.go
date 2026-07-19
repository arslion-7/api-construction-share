package models

type BuilderAddress struct {
	Areas                 []Area  `gorm:"many2many:builder_areas;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"areas"`
	Address               *string `gorm:"type:varchar(510);index" json:"address"`
	AddressAdditionalInfo *string `gorm:"type:varchar(510);index" json:"address_additional_info"`
}

type Builder struct {
	BaseModel
	BuilderAddress
	Org
	// true when this record was created by the old-registries migration
	FromOldRegistry bool `gorm:"column:from_old_registry;default:false;index" json:"from_old_registry"`
}
