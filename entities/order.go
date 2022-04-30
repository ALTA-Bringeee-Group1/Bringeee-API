package entities

import (
	"time"

	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	DriverID null.Int `gorm:"default:null"`
	CustomerID uint
	DestinationID uint
	TruckTypeID int
	OrderPicture string
	TotalVolume int
	TotalWeight int
	Distance int `gorm:"default:null"`
	EstimatedPrice int64 `gorm:"default:null"`
	FixPrice int64 `gorm:"default:null"`
	TransactionID string `gorm:"default:null"`
	PaymentMethod string `gorm:"default:null"`
	Status string
	Description string `gorm:"default:null"`
	ArrivedPicture string `gorm:"default:null"`
	Destination Destination `gorm:"foreignkey:DestinationID;references:ID"`
	TruckType TruckType `gorm:"foreignkey:TruckTypeID;references:ID"`
	Customer User `gorm:"foreignKey:CustomerID;references:ID"`
	Driver Driver `gorm:"foreignKey:DriverID;references:ID;constraint:OnUpdate:SET NULL,OnDelete:SET NULL"`
}

type Destination struct {
	gorm.Model
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
	gorm.Model
	Log string
	Actor string
	OrderID uint
	Order Order `gorm:"foreignKey:OrderID;references:ID"`
}

type OrderHistoryResponse struct {
	ID uint 		 	`json:"id"`
	Log string		 	`json:"log"`
	Actor string	 	`json:"actor"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
	TotalVolume int `form:"total_volume" validate:"required"`
	TotalWeight int `form:"total_weight" validate:"required"`
	Description string `form:"description"`
}

type CreatePaymentRequest struct {
	PaymentMethod string `form:"payment_method"`
}

type PaymentResponse struct {
	OrderID string 				`json:"order_id"`
	TransactionID string 		`json:"transaction_id"`
	PaymentMethod string 		`json:"payment_method"`
	BillNumber string 			`json:"bill_number"`
	Bank string 				`json:"bank"`
	GrossAmount int64 			`json:"gross_amount"`
	TransactionTime time.Time	`json:"transaction_time"`
	TransactionExpire time.Time `json:"transaction_expire"`
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
	TotalVolume int	`json:"total_volume"`
	TotalWeight int	`json:"total_weight"`
	Distance int	`json:"distance"`
	EstimatedPrice int64	`json:"estimated_price"`
	FixPrice int64	`json:"fix_price"`
	TransactionID interface{}	`json:"transaction_id"`
	Status string	`json:"status"`
	Description string	`json:"description"`
	ArrivedPicture string	`json:"arrived_picture"`
	CreatedAt time.Time	`json:"created_at"`
	UpdatedAt time.Time	`json:"updated_at"`
}


































