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
}
