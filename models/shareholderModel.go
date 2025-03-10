package models

type ShareholderAddress struct {
	Areas                 []Area  `gorm:"many2many:shareholder_areas;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"areas"`
	Address               *string `gorm:"type:varchar(510);index" json:"address"`
	AddressAdditionalInfo *string `gorm:"type:varchar(510);index" json:"address_additional_info"`
}

type Phone struct {
	BaseModel
	Kind          *string `gorm:"size:3"`
	Number        *string `gorm:"size:12"`
	Owner         *string `gorm:"size:125"`
	ShareholderID uint
}

type ShareholderDocs struct {
	PassportSeries     *string `json:"passport_series" gorm:"size:6"`
	PassportNumber     *uint   `json:"passport_number" gorm:"type:integer"` // Use uint for positive integers
	PatentSeries       *string `json:"patent_series" gorm:"size:2"`
	PatentNumber       *uint   `json:"patent_number" gorm:"type:integer"`
	CertNumber         *uint   `json:"cert_number" gorm:"type:integer"`
	DocsAdditionalInfo *string `json:"docs_additional_info" gorm:"type:text"`
}

type Shareholder struct {
	BaseModel
	Org
	ShareholderDocs
	ShareholderAddress
	Phones []Phone
}
