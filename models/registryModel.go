package models

import (
	"time"

	"gorm.io/gorm"
)

type RegistryDates struct {
	ReviewedAt   *time.Time `json:"reviewed_at"`
	RegisteredAt *time.Time `json:"registered_at"`
}

type RegistryMail struct {
	MailDate     *time.Time `json:"mail_date"`
	MailNumber   *string    `json:"mail_number"`
	DeliveryDate *time.Time `json:"delivered_date"`
	Count        *int       `json:"count"`
	Queue        *int       `json:"queue"`
	MinToMudDate *time.Time `json:"min_to_mud_date"`
}

// type ContractBuilderShareholderAddress struct {
// 	ContractBuilderShareholderAreas  []Area  `gorm:"many2many:registry_builder_shareholder_areas;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"contract_builder_shareholder_areas"`
// 	ContractBuilderShareholderStreet *string `gorm:"type:varchar(510);index" json:"contract_builder_shareholder_street"`
// }

// type ContractBuilderContractorAddress struct {
// 	ContractBuilderContractorAreas  []Area  `gorm:"many2many:registry_builder_contractor_areas;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"contract_builder_contractor_areas"`
// 	ContractBuilderContractorStreet *string `gorm:"type:varchar(510);index" json:"contract_builder_contractor_street"`
// }

type RegistryContract struct {
	BuilderShareholderNumber         *string    `json:"builder_shareholder_number"`
	BuilderShareholderDate           *time.Time `json:"builder_shareholder_date"`
	BuilderShareholderAdditionalInfo *string    `json:"builder_shareholder_additional_info"`
	// ContractBuilderShareholderAddress
	BuilderContractorNumber         *string    `json:"builder_contractor_number"`
	BuilderContractorDate           *time.Time `json:"builder_contractor_date"`
	BuilderContractorAdditionalInfo *string    `json:"builder_contractor_additional_info"`
	// ContractBuilderContractorAddress
}

type RegistryDenial struct {
	DenialReason         *string    `json:"denial_reason"`
	DenialDate           *time.Time `json:"denial_date"`
	DenialAdditionalInfo *string    `json:"denial_additional_info"`
}

// RegistryOldData keeps the verbatim legacy values copied from the
// old_registries table (pre-2019 MySQL `mainpayly`). A non-nil OldRegistryID
// marks the registry as migrated from old data and makes migration idempotent.
type RegistryOldData struct {
	OldRegistryID           *uint      `gorm:"column:old_registry_id;index" json:"old_registry_id"`
	OldTB                   *uint      `gorm:"column:old_t_b" json:"old_t_b"`
	OldMinHat               *string    `gorm:"column:old_min_hat;type:varchar(150)" json:"old_min_hat"`
	OldSeneHatMinToMud      *time.Time `gorm:"column:old_sene_hat_min_to_mud;type:date" json:"old_sene_hat_min_to_mud"`
	OldGurujy               *string    `gorm:"column:old_gurujy;type:varchar(255)" json:"old_gurujy"`
	OldPaychy               *string    `gorm:"column:old_paychy;type:varchar(100)" json:"old_paychy"`
	OldSertnamaGurujyPaychy *string    `gorm:"column:old_sertnama_gurujy_paychy;type:varchar(100)" json:"old_sertnama_gurujy_paychy"`
	OldDesga                *string    `gorm:"column:old_desga;type:varchar(275)" json:"old_desga"`
	OldBahaUmumy            *string    `gorm:"column:old_baha_umumy;type:varchar(20)" json:"old_baha_umumy"`
	OldMeydanUmumy          *string    `gorm:"column:old_meydan_umumy;type:varchar(125)" json:"old_meydan_umumy"`
	OldKepResminama         *string    `gorm:"column:old_kep_resminama;type:varchar(225)" json:"old_kep_resminama"`
	OldEmlakPaychy          *string    `gorm:"column:old_emlak_paychy;type:varchar(200)" json:"old_emlak_paychy"`
	OldBahaPaychy           *string    `gorm:"column:old_baha_paychy;type:varchar(35)" json:"old_baha_paychy"`
	OldBaha1m2Paychy        *string    `gorm:"column:old_baha_1m2_paychy;type:varchar(25)" json:"old_baha_1m2_paychy"`
	OldSalgyDesga           *string    `gorm:"column:old_salgy_desga;type:varchar(255)" json:"old_salgy_desga"`
	OldSalgyGurujy          *string    `gorm:"column:old_salgy_gurujy;type:varchar(150)" json:"old_salgy_gurujy"`
	OldSalgyPaychy          *string    `gorm:"column:old_salgy_paychy;type:varchar(150)" json:"old_salgy_paychy"`
	OldBashPotr             *string    `gorm:"column:old_bash_potr;type:varchar(100)" json:"old_bash_potr"`
	OldSertnamaGurPotr      *string    `gorm:"column:old_sertnama_gur_potr;type:varchar(125)" json:"old_sertnama_gur_potr"`
	OldPotratchyKomek       *string    `gorm:"column:old_potratchy_komek;type:varchar(125)" json:"old_potratchy_komek"`
	OldShahadatnama         *string    `gorm:"column:old_shahadatnama;type:varchar(175)" json:"old_shahadatnama"`
	OldYgtyyarnama          *string    `gorm:"column:old_ygtyyarnama;type:varchar(255)" json:"old_ygtyyarnama"`
	OldPatentPasport        *string    `gorm:"column:old_patent_pasport;type:varchar(255)" json:"old_patent_pasport"`
	OldSeneBashySongy       *string    `gorm:"column:old_sene_bashy_songy;type:varchar(50)" json:"old_sene_bashy_songy"`
	OldSeneSeredilen        *time.Time `gorm:"column:old_sene_seredilen;type:date" json:"old_sene_seredilen"`
	OldSeneHasabaAlnan      *string    `gorm:"column:old_sene_hasaba_alnan;type:varchar(125)" json:"old_sene_hasaba_alnan"`
	OldWezipeAlanAdam       *string    `gorm:"column:old_wezipe_alan_adam;type:varchar(125)" json:"old_wezipe_alan_adam"`
	OldAdyAlanAdam          *string    `gorm:"column:old_ady_alan_adam;type:varchar(125)" json:"old_ady_alan_adam"`
	OldSeneSanSertnama      *string    `gorm:"column:old_sene_san_sertnama;type:varchar(55)" json:"old_sene_san_sertnama"`
	OldAdyPaychyAlan        *string    `gorm:"column:old_ady_paychy_alan;type:varchar(125)" json:"old_ady_paychy_alan"`
	OldSenePaychyAlan       *string    `gorm:"column:old_sene_paychy_alan;type:varchar(255)" json:"old_sene_paychy_alan"`
	OldLogin                *string    `gorm:"column:old_login;type:varchar(30)" json:"old_login"`
}

type Registry struct {
	ID                  uint               `json:"id" gorm:"primarykey"`
	TB                  *int               `gorm:"column:t_b" json:"t_b"`
	CreatedAt           time.Time          `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time          `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt           gorm.DeletedAt     `json:"deleted_at" gorm:"index"`
	UserID              *uint              `gorm:"column:user_id" json:"user_id"`
	User                *User              `gorm:"foreignKey:UserID" json:"user"`
	GeneralContractorID *uint              `gorm:"column:general_contractor_id" json:"general_contractor_id"`
	GeneralContractor   *GeneralContractor `gorm:"foreignKey:GeneralContractorID" json:"general_contractor"` // Important for preload
	BuildingID          *uint              `gorm:"column:building_id" json:"building_id"`
	Building            *Building          `gorm:"foreignKey:BuildingID" json:"building"` // Important for preload
	BuilderID           *uint              `gorm:"column:builder_id" json:"builder_id"`
	Builder             *Builder           `gorm:"foreignKey:BuilderID" json:"builder"` // Important for preload
	ReceiverID          *uint              `gorm:"column:receiver_id" json:"receiver_id"`
	Receiver            *Receiver          `gorm:"foreignKey:ReceiverID" json:"receiver"` // Important for preload
	ShareholderID       *uint              `gorm:"column:shareholder_id" json:"shareholder_id"`
	Shareholder         *Shareholder       `gorm:"foreignKey:ShareholderID" json:"shareholder"`
	RegistryDates
	RegistryMail
	RegistryContract
	RegistryDenial
	RegistryOldData

	// Computed field for UI: Shareholder description (Org + DocsAdditionalInfo)
	ShareholderDescription string `gorm:"-" json:"shareholder_description"`

	// Computed fields for UI: Other entity descriptions
	GeneralContractorDescription string `gorm:"-" json:"general_contractor_description"`
	BuildingDescription          string `gorm:"-" json:"building_description"`
	BuilderDescription           string `gorm:"-" json:"builder_description"`
}
