package models

// Area represents the hierarchical structure similar to AreaMptt in Django.
type Area struct {
	Code     uint   `gorm:"primaryKey;autoIncrement" json:"code"`
	Name     string `gorm:"size:255;not null" json:"name"`
	ParentID *uint  `gorm:"index" json:"parent_id"` // Foreign key to self
	Children []Area `gorm:"foreignKey:ParentID" json:"children"`
}
