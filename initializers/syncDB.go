package initializers

import (
	"github.com/arslion-7/api-construction-share/models"
)

func SyncDB() {
	if err := DB.AutoMigrate(
		&models.GeneralContractor{},
	); err != nil {
		panic("failed to migrate tables")
	}
}
