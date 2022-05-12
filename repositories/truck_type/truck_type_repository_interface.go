package truckType

import "bringeee-capstone/entities"

type TruckTypeRepositoryInterface interface {

	/*
	 * Find All
	 * -------------------------------
	 * Mengambil data truckType berdasarkan filters dan sorts
	 *
	 * @var limit 			batas limit hasil query
	 * @var offset 			offset hasil query
	 * @var filters			query untuk penyaringan data, { field, operator, value }
	 * @var sorts			pengurutan data, { field, value[bool] }
	 * @return truckType	list truckType dalam bentuk entity domain
	 * @return error		error
	 */
	FindAll(limit int, offset int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.TruckType, error)

	/*
	 * Find
	 * -------------------------------
	 * Mencari truckType tunggal berdasarkan ID
	 *
	 * @var id 		data id
	 */
	Find(id int) (entities.TruckType, error)

	/*
	 * Find User
	 * -------------------------------
	 * Mencari truckType berdasarkan field tertentu
	 *
	 * @var field 			kolom data
	 * @var value			nilai data
	 * @return truckType	single truckType dalam bentuk entity domain
	 * @return error		error	
	 */
	FindBy(field string, value string) (entities.TruckType, error)

	/*
	 * CountAll
	 * -------------------------------
	 * Menghitung semua truckTypes (ini digunakan untuk pagination di service)
	 *
	 * @return truckType	single truckType dalam bentuk entity domain
	 * @return error		error	
	 */
	CountAll(filters []map[string]string) (int64, error)

	/*
	 * Store
	 * -------------------------------
	 * Menambahkan data truckType kedalam database
	 *
	 * @var truckType		single truckType entity
	 * @return truckType	single truckType dalam bentuk entity domain
	 */
	Store(truckType entities.TruckType) (entities.TruckType, error)

	/*
	 * Update
	 * -------------------------------
	 * Mengupdate data truckType berdasarkan ID
	 *
	 * @var truckType		single truckType entity
	 * @return truckType	single truckType dalam bentuk entity domain
	 * @return error		error
	 */
	Update(truckType entities.TruckType, id int) (entities.TruckType, error)

	/*
	 * Delete
	 * -------------------------------
	 * Delete truckType berdasarkan ID
	 *
	 * @return error		error	
	 */
	Delete(id int) error

	/*
	 * Delete Batch
	 * -------------------------------
	 * Delete multiple truckType berdasarkan filter tertentu
	 *
	 * @var filters	query untuk penyaringan data, { field, operator, value }
	 *
	 * @return error		error	
	 */
	DeleteBatch(filters []map[string]string) error

}