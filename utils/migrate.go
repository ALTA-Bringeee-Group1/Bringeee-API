package utils

import (
	"bringeee-capstone/entities"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&entities.User{},
		&entities.TruckType{},
		&entities.Driver{},
		&entities.Destination{},
		&entities.OrderHistory{},
		&entities.Order{},
		&entities.District{},
		&entities.City{},
		&entities.Province{},
	)
}
