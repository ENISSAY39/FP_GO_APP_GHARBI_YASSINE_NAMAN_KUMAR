package initializers

import (
	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR.git/models"
)

func SyncDataBase() {
	// Migrate the schema
	DB.AutoMigrate(&models.User{})
}