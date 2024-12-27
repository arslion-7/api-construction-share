package models

type User struct {
	BaseModel
	Email       string  `json:"email" gorm:"unique;not null" binding:"required" `
	Password    string  `json:"-" gorm:"not null"`
	FullName    *string `json:"full_name" gorm:"not null"`
	PhoneNumber *string `json:"phone_number"`
	Role        *string `json:"role"`
}
