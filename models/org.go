package models

import "time"

type Org struct {
	TB                *int      `gorm:"column:t_b;unique" json:"tb"`
	OrgType           *string   `gorm:"column:org_type;size:255" json:"org_type"`
	OrgName           *string   `gorm:"column:org_name;size:255" json:"org_name"`
	HeadPosition      *string   `gorm:"column:head_position;size:50" json:"head_position"`
	HeadFullName      *string   `gorm:"column:head_full_name;size:255" json:"head_full_name"`
	OrgAdditionalInfo *string   `gorm:"column:org_additional_info;type:text" json:"org_additional_info"`
	CreatedAt         time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}
