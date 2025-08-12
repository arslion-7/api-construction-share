package models

import (
	"time"
)

// OldRegistry represents the legacy mainpayly table from MySQL
type OldRegistry struct {
	BaseModel
	TB                   uint       `gorm:"column:t_b;type:bigint;index" json:"t_b"`
	MinHat               *string    `gorm:"column:min_hat;type:varchar(150)" json:"min_hat"`
	SeneHatMinToMud      *time.Time `gorm:"column:sene_hat_min_to_mud;type:date" json:"sene_hat_min_to_mud"`
	Gurujy               *string    `gorm:"column:gurujy;type:varchar(255)" json:"gurujy"`
	Paychy               *string    `gorm:"column:paychy;type:varchar(100)" json:"paychy"`
	SertnamaGurujyPaychy *string    `gorm:"column:sertnama_gurujy_paychy;type:varchar(100)" json:"sertnama_gurujy_paychy"`
	Desga                *string    `gorm:"column:desga;type:varchar(275)" json:"desga"`
	BahaUmumy            *string    `gorm:"column:baha_umumy;type:varchar(20);default:'0.00'" json:"baha_umumy"`
	MeydanUmumy          *string    `gorm:"column:meydan_umumy;type:varchar(125)" json:"meydan_umumy"`
	KepResminama         *string    `gorm:"column:kep_resminama;type:varchar(225)" json:"kep_resminama"`
	EmlakPaychy          *string    `gorm:"column:emlak_paychy;type:varchar(200)" json:"emlak_paychy"`
	BahaPaychy           *string    `gorm:"column:baha_paychy;type:varchar(35)" json:"baha_paychy"`
	Baha1m2Paychy        *string    `gorm:"column:baha_1m2_paychy;type:varchar(25)" json:"baha_1m2_paychy"`
	SalgyDesga           *string    `gorm:"column:salgy_desga;type:varchar(255)" json:"salgy_desga"`
	SalgyGurujy          *string    `gorm:"column:salgy_gurujy;type:varchar(150)" json:"salgy_gurujy"`
	SalgyPaychy          *string    `gorm:"column:salgy_paychy;type:varchar(150)" json:"salgy_paychy"`
	BashPotr             *string    `gorm:"column:bash_potr;type:varchar(100)" json:"bash_potr"`
	SertnamaGurPotr      *string    `gorm:"column:sertnama_gur_potr;type:varchar(125)" json:"sertnama_gur_potr"`
	PotratchyKomek       *string    `gorm:"column:potratchy_komek;type:varchar(125)" json:"potratchy_komek"`
	Shahadatnama         *string    `gorm:"column:shahadatnama;type:varchar(175)" json:"shahadatnama"`
	Ygtyyarnama          *string    `gorm:"column:ygtyyarnama;type:varchar(255)" json:"ygtyyarnama"`
	PatentPasport        *string    `gorm:"column:patent_pasport;type:varchar(255)" json:"patent_pasport"`
	SeneBashySongy       *string    `gorm:"column:sene_bashy_songy;type:varchar(50)" json:"sene_bashy_songy"`
	SeneSeredilen        *time.Time `gorm:"column:sene_seredilen;type:date" json:"sene_seredilen"`
	SeneHasabaAlnan      *string    `gorm:"column:sene_hasaba_alnan;type:varchar(125)" json:"sene_hasaba_alnan"`
	WezipeAlanAdam       *string    `gorm:"column:wezipe_alan_adam;type:varchar(125)" json:"wezipe_alan_adam"`
	AdyAlanAdam          *string    `gorm:"column:ady_alan_adam;type:varchar(125)" json:"ady_alan_adam"`
	SeneSanSertnama      *string    `gorm:"column:sene_san_sertnama;type:varchar(55)" json:"sene_san_sertnama"`
	AdyPaychyAlan        *string    `gorm:"column:ady_paychy_alan;type:varchar(125)" json:"ady_paychy_alan"`
	SenePaychyAlan       *string    `gorm:"column:sene_paychy_alan;type:varchar(255)" json:"sene_paychy_alan"`
	Login                *string    `gorm:"column:login;type:varchar(30)" json:"login"`
}

// TableName specifies the table name for OldRegistry
func (OldRegistry) TableName() string {
	return "old_registries"
}
