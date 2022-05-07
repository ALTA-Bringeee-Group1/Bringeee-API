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
	FindAll(limit int, offset int, filters []map[string]interface{}, sorts []map[string]interface{}) ([]entities.Order, error)

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
	 * Find First
	 * -------------------------------
	 * Mencari order tunggal berdasarkan filter
	 *
	 * @var field 	kolom data
	 * @var value	nilai data
	 * @return order	single order dalam bentuk entity domain
	 * @return error	error
	 */
	FindFirst(filters []map[string]interface{}) (entities.Order, error)

	/*
	 * CountAll
	 * -------------------------------
	 * Menghitung semua orders (ini digunakan untuk pagination di service)
	 *
	 * @return order	single order dalam bentuk entity domain
	 * @return error	error
	 */
	CountAll(filters []map[string]interface{}) (int64, error)

	/*
	 * Store
	 * -------------------------------
	 * Menambahkan data order kedalam database
	 *
	 * @var order		single order entity
	 * @var destination	single destination entity
	 * @return order	single order dalam bentuk entity domain
	 */
	Store(order entities.Order, destination entities.Destination) (entities.Order, error)

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
	Delete(id int, destinationID int) error

	/*
	 * Delete Batch
	 * -------------------------------
	 * Delete multiple order berdasarkan filter tertentu
	 *
	 * @var filters	query untuk penyaringan data, { field, operator, value }
	 * @return error	error
	 */
	DeleteBatch(filters []map[string]interface{}) error

	/*
	 * Find By Date
	 * -------------------------------
	 * Mencari order tunggal berdasarkan tanggal pembuatan
	 *
	 * @var day 		rentang waktu pencarian
	 */
	FindByDate(day int) ([]map[string]interface{}, error)

	/*
	 * Find By Date
	 * -------------------------------
	 * Mencari order tunggal berdasarkan tanggal pembuatan
	 *
	 * @var month 		bulan pencarian
	 * @var year 		tahun pencarian
	 */
	FindByMonth(mont int, year int) ([]entities.Order, error)
}
