package initializers

import (
	"github.com/arslion-7/api-construction-share/models"
)

func SyncDB() {
	if err := DB.AutoMigrate(
		&models.Area{},
		&models.User{},
		&models.Builder{},
		&models.Building{},
		&models.GeneralContractor{},
		&models.Registry{},
		&models.Receiver{},
		&models.Shareholder{},
		&models.Phone{},
		&models.ShareholderProperty{},
		&models.OldRegistry{},
		&models.AdditionalAgreement{},
	); err != nil {
		panic("failed to migrate tables: " + err.Error())
	}
}
