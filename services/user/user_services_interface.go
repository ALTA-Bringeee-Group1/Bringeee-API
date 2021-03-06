package user

import (
	"bringeee-capstone/entities"
	web "bringeee-capstone/entities/web"
	storageProvider "bringeee-capstone/services/storage"
	"mime/multipart"
)

type UserServiceInterface interface {
	FindAllCustomer(limit, page int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.CustomerResponse, error)
	CreateCustomer(customerRequest entities.CreateCustomerRequest, files map[string]*multipart.FileHeader, storageProvider storageProvider.StorageInterface) (entities.CustomerAuthResponse, error)
	UpdateCustomer(customerRequest entities.UpdateCustomerRequest, id int, files map[string]*multipart.FileHeader, storageProvider storageProvider.StorageInterface) (entities.CustomerResponse, error)
	GetPaginationCustomer(limit, page int, filters []map[string]string) (web.Pagination, error)
	FindByCustomer(field string, value string) (entities.CustomerResponse, error)
	FindCustomer(id int) (entities.CustomerResponse, error)
	DeleteCustomer(id int, storageProvider storageProvider.StorageInterface) error
	FindAllDriver(limit, page int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.DriverResponse, error)
	CreateDriver(driverRequest entities.CreateDriverRequest, files map[string]*multipart.FileHeader, storageProvider storageProvider.StorageInterface) (entities.DriverAuthResponse, error)
	UpdateDriver(driverRequest entities.UpdateDriverRequest, id int, files map[string]*multipart.FileHeader, storageProvider storageProvider.StorageInterface) (entities.DriverResponse, error)
	UpdateDriverByAdmin(driverRequest entities.UpdateDriverByAdminRequest, id int, files map[string]*multipart.FileHeader, storageProvider storageProvider.StorageInterface) (entities.DriverResponse, error)
	GetPaginationDriver(limit, page int, filters []map[string]string) (web.Pagination, error)
	FindByDriver(field string, value string) (entities.DriverResponse, error)
	FindDriver(id int) (entities.DriverResponse, error)
	DeleteDriver(id int, storageProvider storageProvider.StorageInterface) error
	VerifiedDriverAccount(id int) error
	CountCustomer(filters []map[string]string) (int, error)
	CountDriver(filters []map[string]string) (int, error)
}
