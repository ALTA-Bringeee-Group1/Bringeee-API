package entities

import (
	"time"

	"gorm.io/gorm"
)

type TruckType struct {
	gorm.DB
	TruckType string
	MaxVolume int
	MaxWeight int
	PricePerDistance int64
}

type CreateTruckTypeRequest struct {
	TruckType string 		`form:"truck_type" validate:"required"`
	MaxVolume int 			`form:"max_volume" validate:"required"`
	MaxWeight int 			`form:"max_weight" validate:"required"`
	PricePerDistance int64 	`form:"price_per_distance" validate:"required"`
}

type UpdateTruckTypeRequest struct {
	TruckType string 		`form:"truck_type"`
	MaxVolume int 			`form:"max_volume"`
	MaxWeight int 			`form:"max_weight"`
	PricePerDistance int64 	`form:"price_per_distance"`
}

type TruckTypeResponse struct {
	ID uint					`json:"id"`
	TruckType string		`json:"truck_type"`
	MaxVolume int			`json:"max_volume"`
	MaxWeight int			`json:"max_weight"`
	PricePerDistance int64	`json:"price_per_distance"`
	CreatedAt time.Time		`json:"created_at"`
	UpdatedAt time.Time		`json:"updated_at"`
}