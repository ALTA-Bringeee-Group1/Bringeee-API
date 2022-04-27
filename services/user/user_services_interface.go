package user

import (
	"bringeee-capstone/entities"
	web "bringeee-capstone/entities/web"
	"mime/multipart"
)

type UserServiceInterface interface {
	FindAllCustomer(limit, page int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.CustomerResponse, error)
	CreateCustomer(customerRequest entities.CreateCustomerRequest, files map[string]*multipart.FileHeader) (entities.CustomerAuthResponse, error)
	UpdateCustomer(customerRequest entities.UpdateCustomerRequest, id int, files map[string]*multipart.FileHeader) (entities.CustomerResponse, error)
	GetPaginationCustomer(limit, page int, filters []map[string]string) (web.Pagination, error)
	FindCustomer(id int) (entities.CustomerResponse, error)
	DeleteCustomer(id int) error
	FindAllDriver(limit, page int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.DriverAuthResponse, error)
	CreateDriver(driverRequest entities.CreateDriverRequest, files map[string]*multipart.FileHeader) (entities.DriverAuthResponse, error)
	UpdateDriver(driverRequest entities.UpdateDriverRequest, id int, files []*multipart.FileHeader) (entities.DriverResponse, error)
	UpdateDriverByAdmin(driverRequest entities.UpdateDriverByAdminRequest, id int, files []*multipart.FileHeader) (entities.DriverResponse, error)
	GetPaginationDriver(limit, page int, filters []map[string]string) (web.Pagination, error)
	FindDriver(id int) (entities.DriverResponse, error)
	DeleteDriver(id int) error
}
