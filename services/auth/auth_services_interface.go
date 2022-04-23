package auth

import (
	"bringeee-capstone/entities"
)

type AuthServiceInterface interface {
	LoginCustomer(AuthReq entities.AuthRequest) (entities.CustomerAuthResponse, error)
	LoginDriver(AuthReq entities.AuthRequest) (entities.DriverAuthResponse, error)
	LoginAdmin(AuthReq entities.AuthRequest) (entities.AdminAuthResponse, error)
	CustomerMe(Id int) (entities.CustomerAuthResponse, error)
	DriverMe(Id int) (entities.DriverAuthResponse, error)
	AdminMe(Id int) (entities.AdminAuthResponse, error)
}
