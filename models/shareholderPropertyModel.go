package models

type ShareholderProperty struct {
	BaseModel
	TB                  *int     `gorm:"column:t_b" json:"t_b"`
	BuildingType        *string  `json:"building_type"`
	Part                *int     `json:"part"`
	Building            *int     `json:"building"`
	Entrance            *int     `json:"entrance"`
	Floor               *int     `json:"floor"`
	Apartment           *int     `json:"apartment"`
	RoomCount           *int     `json:"room_count"`
	Square              *float64 `gorm:"type:decimal(16,6)" json:"square"`
	Price               *float64 `gorm:"type:decimal(16,2)" json:"price"`
	Price1m2            *float64 `gorm:"type:decimal(10,2)" json:"price_1m2"`
	BuildingIdentNumber *int     `json:"building_ident_number"`
	AdditionalInfo      *string  `json:"additional_info"`
	RegistryID          uint     `json:"registry_id"`
}
