package entities

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email       string `gorm:"unique"`
	Name        string
	Password    string
	DOB         time.Time
	Gender      string
	Address     string
	PhoneNumber string
	Avatar      string
	Role        string
}

type Driver struct {
	gorm.Model
	UserID            uint
	TruckTypeID       int
	KtpFile           string
	StnkFile          string
	DriverLicenseFile string
	Age               int
	VehicleIdentifier string
	NIK               string
	AccountStatus     string
	Status            string
	VehiclePicture    string
	User              User      `gorm:"foreignKey:UserID;references:ID"`
	TruckType         TruckType `gorm:"foreignKey:TruckTypeID;references:ID"`
	Orders 		 	  []Order 	`gorm:"foreignKey:DriverID;references:ID;constraint:OnUpdate:SET NULL,OnDelete:SET NULL"`

}

// File request: avatar
type CreateCustomerRequest struct {
	Email       string `form:"email" validate:"required,email"`
	Password    string `form:"password" validate:"required"`
	Name        string `form:"name" validate:"required"`
	DOB         string `form:"dob" validate:"required"`
	Gender      string `form:"gender" validate:"required"`
	Address     string `form:"address" validate:"required"`
	PhoneNumber string `form:"phone_number" validate:"required"`
}

type UpdateCustomerRequest struct {
	Email       string `form:"email"`
	Password    string `form:"password"`
	Name        string `form:"name"`
	DOB         string `form:"dob"`
	Gender      string `form:"gender"`
	Address     string `form:"address"`
	PhoneNumber string `form:"phone_number"`
}
type CustomerResponse struct {
	ID          uint      `json:"id"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	DOB         time.Time `json:"dob"`
	Gender      string    `json:"gender"`
	Address     string    `json:"address"`
	PhoneNumber string    `json:"phone_number"`
	Avatar      string    `json:"avatar"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AdminResponse struct {
	ID          uint      `json:"id"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	DOB         time.Time `json:"dob"`
	Gender      string    `json:"gender"`
	Address     string    `json:"address"`
	PhoneNumber string    `json:"phone_number"`
	Avatar      string    `json:"avatar"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// File request: avatar, ktp_file, stnk_file, driver_license_file, vehicle_picture
type CreateDriverRequest struct {
	Email             string `form:"email" validate:"required,email"`
	Password          string `form:"password" validate:"required"`
	Name              string `form:"name" validate:"required"`
	DOB               string `form:"dob" validate:"required"`
	Gender            string `form:"gender" validate:"required"`
	Address           string `form:"address" validate:"required"`
	PhoneNumber       string `form:"phone_number" validate:"required"`
	TruckTypeID       uint   `form:"truck_type_id" validate:"required"`
	Age               int    `form:"age" validate:"required"`
	VehicleIdentifier string `form:"vehicle_identifier" validate:"required"`
	NIK               string `form:"nik" validate:"required"`
}

type UpdateDriverRequest struct {
	Email             string `form:"email"`
	Password          string `form:"password"`
	Name              string `form:"name"`
	DOB               string `form:"dob"`
	Gender            string `form:"gender"`
	Address           string `form:"address"`
	PhoneNumber       string `form:"phone_number"`
	TruckTypeID       uint   `form:"truck_type_id"`
	Age               int    `form:"age"`
	VehicleIdentifier string `form:"vehicle_identifier"`
	NIK               string `form:"nik"`
}

type DriverResponse struct {
	ID                uint              `json:"id"`
	UserID            string            `json:"user_id"`
	User              CustomerResponse  `json:"user"`
	TruckTypeID       uint              `json:"truck_type_id"`
	TruckType         TruckTypeResponse `json:"truck_type"`
	KtpFile           string            `json:"ktp_file"`
	StnkFile          string            `json:"stnk_file"`
	DriverLicenseFile string            `json:"driver_license_file"`
	Age               int            	`json:"age"`
	VehicleIdentifier string            `json:"vehicle_identifier"`
	NIK               string            `json:"nik"`
	VehiclePicture    string            `json:"vehicle_picture"`
	Status            string            `json:"status"`
	AccountStatus     string            `json:"account_status"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}
