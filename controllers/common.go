package controllers

type AddressInput struct {
	// TB    *int   `json:"t_b"`
	Areas                 []uint  `json:"areas"` // IDs of the associated Areas
	Address               *string `gorm:"type:varchar(510);index" json:"address"`
	AddressAdditionalInfo *string `gorm:"type:varchar(510);index" json:"address_additional_info"`
}
