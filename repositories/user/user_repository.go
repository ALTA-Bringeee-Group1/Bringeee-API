package user

import (
	"bringeee-capstone/entities"
	"bringeee-capstone/entities/web"

	"gorm.io/gorm"
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
		return entities.User{}, web.WebError{Code: 500, Message: tx.Error.Error()}
	}
	return user, nil
}

func (repo UserRepository) FindByCustomer(field string, value string) (entities.User, error) {
	user := entities.User{}
	tx := repo.db.Where(field+" = ?", value).Find(&user)
	if tx.Error != nil {

		// return kode 500 jika terjadi error
		return entities.User{}, web.WebError{Code: 500, Message: "server error"}
	} else if tx.RowsAffected <= 0 {

		// return kode 400 jika tidak ditemukan
		return entities.User{}, web.WebError{Code: 400, Message: "bad Request"}
	}
	return user, nil
}

func (repo UserRepository) FindCustomer(id int) (entities.User, error) {
	// Get user dari database
	user := entities.User{}
	tx := repo.db.Find(&user, id)
	if tx.Error != nil {

		// Return error dengan code 500
		return entities.User{}, web.WebError{Code: 500, Message: "server error"}
	} else if tx.RowsAffected <= 0 {

		// Return error dengan code 400 jika tidak ditemukan
		return entities.User{}, web.WebError{Code: 400, Message: "bad request"}
	}
	return user, nil
}

func (repo UserRepository) FindByDriver(field string, value string) (entities.Driver, error) {
	driver := entities.Driver{}
	tx := repo.db.Preload("User").Preload("TruckType").Where(field+" = ?", value).Find(&driver)
	if tx.Error != nil {

		// return kode 500 jika terjadi error
		return entities.Driver{}, web.WebError{Code: 500, Message: "server error"}
	} else if tx.RowsAffected <= 0 {

		// return kode 400 jika tidak ditemukan
		return entities.Driver{}, web.WebError{Code: 400, Message: "bad Request"}
	}
	return driver, nil
}

func (repo UserRepository) FindDriver(id int) (entities.Driver, error) {
	// Get driver dari database
	driver := entities.Driver{}
	tx := repo.db.Preload("User").Preload("TruckType").Find(&driver, id)
	if tx.Error != nil {

		// Return error dengan code 500
		return entities.Driver{}, web.WebError{Code: 500, Message: "server error"}
	} else if tx.RowsAffected <= 0 {

		// Return error dengan code 400 jika tidak ditemukan
		return entities.Driver{}, web.WebError{Code: 400, Message: "bad request"}
	}
	return driver, nil
}
