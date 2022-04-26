package user

import "bringeee-capstone/entities"

type UserRepositoryInterface interface {
	// FindAllCustomer(limit int, offset int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.User, error)
	// FindCustomer(id int) (entities.User, error)
	// FindByCustomer(field string, value string) (entities.User, error)
	// CountAllCustomer(filters []map[string]string) (int64, error)
	StoreCustomer(user entities.User) (entities.User, error)
	// UpdateCustomer(user entities.User, id int) (entities.User, error)
	// DeleteCustomer(id int) error
	// FindAllDriver(limit int, offset int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.Driver, error)
	// FindDriver(id int) (entities.Driver, error)
	// FindByDriver(field string, value string) (entities.Driver, error)
	// CountAllDriver(filters []map[string]string) (int64, error)
	StoreDriver(driver entities.Driver) (entities.Driver, error)
	// UpdateDriver(driver entities.Driver, id int) (entities.Driver, error)
	// DeleteDriver(id int) error
}
