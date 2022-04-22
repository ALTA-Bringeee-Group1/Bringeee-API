package entities

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.DB
	DriverID uint
	CustomerID uint
	DestinationID uint
	TruckTypeID int
	OrderPicture string
	TotalVolume int
	TotalWeight int
	Distance int
	EstimatedPrice int64
	FixPrice int64
	TransactionID int64
	Status string
	Description string
	ArrivedPicture string
	Destination Destination `gorm:"foreignkey:DestinationID;references:ID"`
	TruckType TruckType `gorm:"foreignkey:TruckTypeID;references:ID"`
	Customer User `gorm:"foreignKey:CustomerID;references:ID"`
}

type Destination struct {
	gorm.DB
	DestinationStartProvince string
	DestinationStartCity string
	DestinationStartDistrict string
	DestinationStartAddress string
	DestinationStartPostal string
	DestinationStartLat string
	DestinationStartLong string
	DestinationEndProvince string
	DestinationEndCity string
	DestinationEndDistrict string
	DestinationEndAddress string
	DestinationEndPostal string
	DestinationEndLat string
	DestinationEndLong string
}

type OrderHistory struct {
	gorm.DB
	Log string
	Actor string
	OrderID uint
	Order Order `gorm:"foreignKey:OrderID;references:ID"`
}


// file request: order_picture
type CustomerCreateOrderRequest struct {
	DestinationStartProvince string `form:"destination_start_province" validate:"required"`
	DestinationStartCity string `form:"destination_start_city" validate:"required"`
	DestinationStartDistrict string `form:"destination_start_district" validate:"required"`
	DestinationStartAddress string `form:"destination_start_address" validate:"required"`
	DestinationStartPostal string `form:"destination_start_postal" validate:"required"`
	DestinationStartLat string `form:"destination_start_lat" validate:"required"`
	DestinationStartLong string `form:"destination_start_long" validate:"required"`
	DestinationEndProvince string `form:"destination_end_province" validate:"required"`
	DestinationEndCity string `form:"destination_end_city" validate:"required"`
	DestinationEndDistrict string `form:"destination_end_district" validate:"required"`
	DestinationEndAddress string `form:"destination_end_address" validate:"required"`
	DestinationEndPostal string `form:"destination_end_postal" validate:"required"`
	DestinationEndLat string `form:"destination_end_lat" validate:"required"`
	DestinationEndLong string `form:"destination_end_long" validate:"required"`
	TruckTypeID int `form:"truck_type_id" validate:"required"`
	TotalVolume string `form:"total_volume" validate:"required"`
	TotalWeight string `form:"total_weight" validate:"required"`
	Description string `form:"description"`
}


// file request: arrived_picture
type DriverFinishOrderRequest struct {

}

type AdminSetPriceOrderRequest struct {
	FixedPrice int64 `form:"fixed_price" validate:"required"`
}

type OrderResponse struct {
	ID uint	`json:"id"`
	DriverID uint	`json:"driver_id"`
	Driver DriverResponse	`json:"driver"`
	CustomerID uint	`json:"customer_id"`
	Customer CustomerResponse	`json:"customer"`
	DestinationStartProvince string	`json:"destination_start_province"`
	DestinationStartCity string	`json:"destination_start_city"`
	DestinationStartDistrict string	`json:"destination_start_district"`
	DestinationStartAddress string	`json:"destination_start_address"`
	DestinationStartPostal string	`json:"destination_start_postal"`
	DestinationStartLat string	`json:"destination_start_lat"`
	DestinationStartLong string	`json:"destination_start_long"`
	DestinationEndProvince string	`json:"destination_end_province"`
	DestinationEndCity string	`json:"destination_end_city"`
	DestinationEndDistrict string	`json:"destination_end_district"`
	DestinationEndAddress string	`json:"destination_end_address"`
	DestinationEndPostal string	`json:"destination_end_postal"`
	DestinationEndLat string	`json:"destination_end_lat"`
	DestinationEndLong string	`json:"destination_end_long"`
	TruckTypeID uint	`json:"truck_type_id"`
	TruckType TruckTypeResponse	`json:"truck_type"`
	OrderPicture string	`json:"order_picture"`
	TotalVolume string	`json:"total_volume"`
	TotalWeight string	`json:"total_weight"`
	Distance string	`json:"distance"`
	EstimatedPrice int64	`json:"estimated_price"`
	FixedPrice int64	`json:"fixed_price"`
	TransactionID interface{}	`json:"transaction_id"`
	Status string	`json:"status"`
	Description string	`json:"description"`
	ArrivedPicture string	`json:"arrived_picture"`
	CreatedAt time.Time	`json:"created_at"`
	UpdatedAt time.Time	`json:"updated_at"`
}


































