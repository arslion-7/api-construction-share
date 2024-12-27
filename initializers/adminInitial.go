package initializers

import (
	"errors"

	"github.com/arslion-7/api-construction-share/models"
	"github.com/arslion-7/api-construction-share/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func AdminInitial() {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("010203"), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	if DB.Migrator().HasTable(&models.User{}) {
		if err := DB.First(&models.User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {

			admin := &models.User{
				Email:    "admin@mail.ru",
				Password: string(hashedPassword),
				FullName: &utils.Role.Admin,
				Role:     &utils.Role.Admin,
			}

			DB.Create(admin)
		}
	}
}
