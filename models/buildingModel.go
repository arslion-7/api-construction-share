package models

import (
	"time"
)

type BuildingAddress struct {
	Areas  []Area  `gorm:"many2many:building_areas;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"areas"`
	Street *string `gorm:"type:varchar(510);index" json:"street"`
}

type BuildingMain struct {
	Kind       *string    `gorm:"type:varchar(510);index" json:"kind"`
	StartDate  *time.Time `json:"start_date"`
	EndDate    *time.Time `json:"end_date"`
	Price      *float64   `gorm:"type:decimal(14,2)" json:"price"`
	Percentage *float64   `json:"percentage"`
}

type BuildingOrder struct {
	OrderWhoseWhat      *string    `gorm:"type:varchar(255)" json:"order_whose_what"`
	OrderDate           *time.Time `json:"order_date"`
	OrderCode           *string    `gorm:"type:varchar(50);index" json:"order_code"`
	OrderAdditionalInfo *string    `gorm:"type:text" json:"order_additional_info"`
}

type BuildingCert struct {
	CertName  *string    `gorm:"type:varchar(50)" json:"cert_name"`
	Cert1Date *time.Time `json:"cert_1_date"`
	Cert1Code *string    `gorm:"type:varchar(50);index" json:"cert_1_code"`
	Cert2Date *time.Time `json:"cert_2_date"`
	Cert2Code *string    `gorm:"type:varchar(50);index" json:"cert_2_code"`
}

type BuildingSquare struct {
	Square1              *float64 `json:"square_1"`
	Square1Name          *string  `gorm:"type:varchar(5)" json:"square_1_name"`
	Square2              *float64 `json:"square_2"`
	Square2Name          *string  `gorm:"type:varchar(5)" json:"square_2_name"`
	Square3              *float64 `json:"square_3"`
	Square3Name          *string  `gorm:"type:varchar(5)" json:"square_3_name"`
	SquareAdditionalInfo *string  `gorm:"type:text" json:"square_additional_info"`
}

type Building struct {
	BaseModel
	BuildingAddress
	BuildingMain
	BuildingOrder
	BuildingCert
	BuildingSquare
	TB          *int `gorm:"uniqueIndex" json:"t_b"`
	IdentNumber *int `json:"ident_number"`
}
