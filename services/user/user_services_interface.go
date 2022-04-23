package user

import (
	"bringeee-capstone/entities"
	web "bringeee-capstone/entities/web"
	"mime/multipart"
)

type UserServiceInterface interface {
	FindAllCustomer(limit, page int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.CustomerResponse, error)
	CreateCustomer(customerRequest entities.CreateCustomerRequest, avatar *multipart.FileHeader) (entities.CustomerAuthResponse, error)
	UpdateCustomer(customerRequest entities.UpdateCustomerRequest, avatar *multipart.FileHeader) (entities.CustomerResponse, error)
	GetPaginationCustomer(limit, page int, filters []map[string]string) (web.Pagination, error)
	FindCustomer(id int) (entities.CustomerResponse, error)
	DeleteCustomer(id int) error
	FindAllDriver(limit, page int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.DriverAuthResponse, error)
	CreateDriver(driverRequest entities.CreateDriverRequest, avatar, ktp, lisence, stnk, vehicle_picture *multipart.FileHeader) (entities.DriverResponse, error)
	UpdateDriver(driverRequest entities.UpdateDriverRequest, avatar, ktp, lisence, stnk, vehicle_picture *multipart.FileHeader) (entities.DriverResponse, error)
	GetPaginationDriver(limit, page int, filters []map[string]string) (web.Pagination, error)
	FindDriver(id int) (entities.CustomerResponse, error)
	DeleteDriver(id int) error
}
