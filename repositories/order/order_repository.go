package order

import (
	"bringeee-capstone/entities"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

/*
 * Find All
 * -------------------------------
 * Mengambil data order berdasarkan filters dan sorts
 *
 * @var limit 	batas limit hasil query
 * @var offset 	offset hasil query
 * @var filters	query untuk penyaringan data, { field, operator, value }
 * @var sorts	pengurutan data, { field, value[bool] }
 * @return order	list order dalam bentuk entity domain
 * @return error	error
 */
func (repository OrderRepository) FindAll(limit int, offset int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.Order, error) {
	panic("implement me")
}

/*
 * Find
 * -------------------------------
 * Mencari order tunggal berdasarkan ID
 *
 * @var id 		data id
 */
func (repository OrderRepository) Find(id int) (entities.Order, error) {
	panic("implement me")
}

/*
 * Find User
 * -------------------------------
 * Mencari order berdasarkan field tertentu
 *
 * @var field 	kolom data
 * @var value	nilai data
 * @return order	single order dalam bentuk entity domain
 * @return error	error	
 */
func (repository OrderRepository) FindBy(field string, value string) (entities.Order, error) {
	panic("implement me")
}

/*
 * CountAll
 * -------------------------------
 * Menghitung semua orders (ini digunakan untuk pagination di service)
 *
 * @return order	single order dalam bentuk entity domain
 * @return error	error	
 */
func (repository OrderRepository) CountAll(filters []map[string]string) (int64, error) {
	panic("implement me")
}

/*
 * Store
 * -------------------------------
 * Menambahkan data order kedalam database
 *
 * @var order		single order entity
 * @return order	single order dalam bentuk entity domain
 */
func (repository OrderRepository) Store(order entities.Order) (entities.Order, error) {
	panic("implement me")
}

/*
 * Update
 * -------------------------------
 * Mengupdate data order berdasarkan ID
 *
 * @var order		single order entity
 * @return order	single order dalam bentuk entity domain
 * @return error	error
 */
func (repository OrderRepository) Update(order entities.Order, id int) (entities.Order, error) {
	panic("implement me")
}

/*
 * Delete
 * -------------------------------
 * Delete order berdasarkan ID
 *
 * @return error	error	
 */
func (repository OrderRepository) Delete(id int) error {
	panic("implement me")
}

/*
 * Delete Batch
 * -------------------------------
 * Delete multiple order berdasarkan filter tertentu
 *
 * @var filters	query untuk penyaringan data, { field, operator, value }
 * @return error	error	
 */
func (repository OrderRepository) DeleteBatch(filters []map[string]string) error {
	panic("implement me")
}