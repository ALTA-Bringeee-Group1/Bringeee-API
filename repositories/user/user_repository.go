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
		return entities.Driver{}, web.WebError{Code: 500, Message: tx.Error.Error()}
	}
	return driver, nil
}
