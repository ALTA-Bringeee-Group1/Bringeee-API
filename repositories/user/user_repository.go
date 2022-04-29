package user

import (
	"bringeee-capstone/entities"
	"bringeee-capstone/entities/web"
	"bringeee-capstone/utils"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	db *gorm.DB
}

// Constructor
func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{
		db: db,
	}
}

/*
 * Store
 * -------------------------------
 * Menambahkan user tunggal kedalam database
 */
func (repo UserRepository) StoreCustomer(user entities.User) (entities.User, error) {

	// insert user ke database
	tx := repo.db.Create(&user)
	if tx.Error != nil {
		// return kode 500 jika error
		return entities.User{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "servernya error"}
	}
	return user, nil
}

/*
 * Store
 * -------------------------------
 * Menambahkan driver tunggal kedalam database
 */
func (repo UserRepository) StoreDriver(driver entities.Driver) (entities.Driver, error) {

	// insert driver ke database
	tx := repo.db.Create(&driver)
	if tx.Error != nil {

		// return kode 500 jika error
		return entities.Driver{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "servernya error"}
	}
	return driver, nil
}

func (repo UserRepository) FindByCustomer(field string, value string) (entities.User, error) {
	user := entities.User{}
	tx := repo.db.Where(field+" = ?", value).Find(&user)
	if tx.Error != nil {

		// return kode 500 jika terjadi error
		return entities.User{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "servernya error"}
	} else if tx.RowsAffected <= 0 {

		// return kode 400 jika tidak ditemukan
		return entities.User{}, web.WebError{Code: 400, ProductionMessage: "bad request", DevelopmentMessage: "permintaan buruk"}
	}
	return user, nil
}

func (repo UserRepository) FindCustomer(id int) (entities.User, error) {
	// Get user dari database
	user := entities.User{}
	tx := repo.db.Find(&user, id)
	if tx.Error != nil {

		// Return error dengan code 500
		return entities.User{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "servernya error"}
	} else if tx.RowsAffected <= 0 {

		// Return error dengan code 400 jika tidak ditemukan
		return entities.User{}, web.WebError{Code: 400, ProductionMessage: "bad request", DevelopmentMessage: "permintaan buruk"}
	}
	return user, nil
}

func (repo UserRepository) FindByDriver(field string, value string) (entities.Driver, error) {
	driver := entities.Driver{}
	tx := repo.db.Preload("User").Preload("TruckType").Where(field+" = ?", value).Find(&driver)
	if tx.Error != nil {

		// return kode 500 jika terjadi error
		return entities.Driver{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "servernya error"}
	} else if tx.RowsAffected <= 0 {

		// return kode 400 jika tidak ditemukan
		return entities.Driver{}, web.WebError{Code: 400, ProductionMessage: "bad request", DevelopmentMessage: "permintaan buruk"}
	}
	return driver, nil
}

func (repo UserRepository) FindDriver(id int) (entities.Driver, error) {
	// Get driver dari database
	driver := entities.Driver{}
	tx := repo.db.Preload("User").Preload("TruckType").Find(&driver, id)
	if tx.Error != nil {

		// Return error dengan code 500
		return entities.Driver{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "servernya error"}
	} else if tx.RowsAffected <= 0 {

		// Return error dengan code 400 jika tidak ditemukan
		return entities.Driver{}, web.WebError{Code: 400, ProductionMessage: "bad request", DevelopmentMessage: "permintaan buruk"}

	}
	return driver, nil
}

func (repo UserRepository) UpdateCustomer(user entities.User, id int) (entities.User, error) {
	tx := repo.db.Save(&user)
	if tx.Error != nil {
		// return Kode 500 jika error
		return entities.User{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "servernya error"}
	}
	return user, nil
}

func (repo UserRepository) UpdateDriver(driver entities.Driver, id int) (entities.Driver, error) {
	fmt.Println(utils.JsonEncode(driver))
	tx := repo.db.Save(&driver)
	if tx.Error != nil {
		// return Kode 500 jika error
		return entities.Driver{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "servernya error"}
	}
	return driver, nil
}

func (repo UserRepository) DeleteCustomer(id int) error {

	// Delete from database
	tx := repo.db.Delete(&entities.User{}, id)
	if tx.Error != nil {

		// return kode 500 jika error
		return web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "servernya error"}
	}
	return nil
}

func (repo UserRepository) FindAllCustomer(limit int, offset int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.User, error) {
	users := []entities.User{}
	builder := repo.db.Limit(limit).Offset(offset)
	// Where filters
	for _, filter := range filters {
		builder.Where(filter["field"]+" "+filter["operator"]+" ?", filter["value"])
	}
	// OrderBy Filters
	for _, sort := range sorts {
		builder.Order(clause.OrderByColumn{Column: clause.Column{Name: sort["field"].(string)}, Desc: sort["desc"].(bool)})
	}
	tx := builder.Find(&users)
	if tx.Error != nil {
		return []entities.User{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "servernya error"}
	}
	return users, nil
}

func (repo UserRepository) FindAllDriver(limit int, offset int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.Driver, error) {
	drivers := []entities.Driver{}
	builder := repo.db.Joins("User").Preload("User").Preload("TruckType").Limit(limit).Offset(offset)
	// Where filters
	for _, filter := range filters {
		builder.Where(filter["field"]+" "+filter["operator"]+" ?", filter["value"])
	}
	// OrderBy Filters
	for _, sort := range sorts {
		builder.Order(clause.OrderByColumn{Column: clause.Column{Name: sort["field"].(string)}, Desc: sort["desc"].(bool)})
	}
	tx := builder.Find(&drivers)
	if tx.Error != nil {
		return []entities.Driver{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "servernya error"}
	}
	return drivers, nil
}

func (repo UserRepository) CountAllCustomer(filters []map[string]string) (int64, error) {
	var count int64
	builder := repo.db.Model(&entities.User{})
	// Where filters
	for _, filter := range filters {
		builder.Where(filter["field"]+" "+filter["operator"]+" ?", filter["value"])
	}
	tx := builder.Count(&count)
	if tx.Error != nil {
		return -1, web.WebError{Code: 400, ProductionMessage: "bad request", DevelopmentMessage: "permintaan buruk"}
	}
	return count, nil
}

func (repo UserRepository) CountAllDriver(filters []map[string]string) (int64, error) {
	var count int64
	builder := repo.db.Model(&entities.Driver{}).Joins("User")
	// Where filters
	for _, filter := range filters {
		builder.Where(filter["field"]+" "+filter["operator"]+" ?", filter["value"])
	}
	tx := builder.Count(&count)
	if tx.Error != nil {
		return -1, web.WebError{Code: 400, ProductionMessage: "bad request", DevelopmentMessage: "permintaan buruk"}
	}
	return count, nil
}

func (repo UserRepository) DeleteDriver(id int) error {

	// Delete from database
	tx := repo.db.Delete(&entities.Driver{}, id)
	if tx.Error != nil {

		// return kode 500 jika error
		return web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "servernya error"}
	}
	return nil
}
