package utils

import "gorm.io/gorm"

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&entities.User{},
	)
}
