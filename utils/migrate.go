package utils

import (
	"bringeee-capstone/entities"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&entities.User{},
	)
}
