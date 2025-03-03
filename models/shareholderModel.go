package models

type ShareholderAddress struct {
	Areas                 []Area  `gorm:"many2many:shareholder_areas;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"areas"`
	Address               *string `gorm:"type:varchar(510);index" json:"address"`
	AddressAdditionalInfo *string `gorm:"type:varchar(510);index" json:"address_additional_info"`
}

type Phone struct {
	BaseModel
	Kind          string `gorm:"size:3"`
	Number        string `gorm:"size:12"`
	Owner         string `gorm:"size:125"`
	ShareholderID uint
}

type Shareholder struct {
	BaseModel
	Org
	ShareholderAddress
	PassportSeries     string `gorm:"size:6"`
	PassportNumber     uint   `gorm:"type:integer"` // Use uint for positive integers
	PatentSeries       string `gorm:"size:2"`
	PatentNumber       uint   `gorm:"type:integer"`
	CertNumber         uint   `gorm:"type:integer"`
	DocsAdditionalInfo string `gorm:"type:text"`
	Phones             []Phone
}
