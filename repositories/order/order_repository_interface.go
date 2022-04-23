package order

import "bringeee-capstone/entities"

type OrderRepositoryInterface interface {
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
	FindAll(limit int, offset int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.Order, error)

	/*
	 * Find
	 * -------------------------------
	 * Mencari order tunggal berdasarkan ID
	 *
	 * @var id 		data id
	 */
	Find(id int) (entities.Order, error)

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
	FindBy(field string, value string) (entities.Order, error)

	/*
	 * CountAll
	 * -------------------------------
	 * Menghitung semua orders (ini digunakan untuk pagination di service)
	 *
	 * @return order	single order dalam bentuk entity domain
	 * @return error	error	
	 */
	CountAll(filters []map[string]string) (int64, error)

	/*
	 * Store
	 * -------------------------------
	 * Menambahkan data order kedalam database
	 *
	 * @var order		single order entity
	 * @return order	single order dalam bentuk entity domain
	 */
	Store(order entities.Order) (entities.Order, error)

	/*
	 * Update
	 * -------------------------------
	 * Mengupdate data order berdasarkan ID
	 *
	 * @var order		single order entity
	 * @return order	single order dalam bentuk entity domain
	 * @return error	error
	 */
	Update(order entities.Order, id int) (entities.Order, error)

	/*
	 * Delete
	 * -------------------------------
	 * Delete order berdasarkan ID
	 *
	 * @return error	error	
	 */
	Delete(id int) error

	/*
	 * Delete Batch
	 * -------------------------------
	 * Delete multiple order berdasarkan filter tertentu
	 *
	 * @var filters	query untuk penyaringan data, { field, operator, value }
	 * @return error	error	
	 */
	DeleteBatch(filters []map[string]string) error
}