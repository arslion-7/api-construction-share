package models

type Receiver struct {
	BaseModel
	TB             *int   `gorm:"column:t_b" json:"t_b"`
	CitizenStatus  string `json:"citizen_status"`
	OrgName        string `json:"org_name"`
	Department     string `json:"department"`
	Position       string `json:"position"`
	Firstname      string `json:"firstname"`
	Lastname       string `json:"lastname"`
	Patronymic     string `json:"patronymic"`
	AdditionalInfo string `json:"additional_info"`
}
