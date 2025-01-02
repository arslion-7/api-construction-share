package models

// Area represents the hierarchical structure similar to AreaMptt in Django.
type Area struct {
	Code     uint   `gorm:"primaryKey;autoIncrement"`
	Name     string `gorm:"size:255;not null"`
	ParentID *uint  `gorm:"index"` // Foreign key to self
	Children []Area `gorm:"foreignKey:ParentID"`
}
