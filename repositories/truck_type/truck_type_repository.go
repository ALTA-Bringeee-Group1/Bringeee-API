package truck_type

import (
	"bringeee-capstone/entities"

	"gorm.io/gorm"
)

type TruckTypeRepository struct {
	db *gorm.DB
}

func NewTruckTypeRepository(db *gorm.DB) *TruckTypeRepository {
	return &TruckTypeRepository{
		db: db,
	}
}

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
func (repository TruckTypeRepository) FindAll(limit int, offset int, filters []map[string]string, sorts []map[string]interface{}) {
	panic("implement me")
}
/*
 * Find
 * -------------------------------
 * Mencari truckType tunggal berdasarkan ID
 *
 * @var id 		data id
 */
func (repository TruckTypeRepository) Find(id int) (entities.TruckType, error) {
	panic("implement me")
}
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
func (repository TruckTypeRepository) FindBy(field string, value string) (entities.TruckType, error) {
	panic("implement me")
}
/*
 * CountAll
 * -------------------------------
 * Menghitung semua truckTypes (ini digunakan untuk pagination di service)
 *
 * @return truckType	single truckType dalam bentuk entity domain
 * @return error		error	
 */
func (repository TruckTypeRepository) CountAll(filters []map[string]string) (int64, error) {
	panic("implement me")
}
/*
 * Store
 * -------------------------------
 * Menambahkan data truckType kedalam database
 *
 * @var truckType		single truckType entity
 * @return truckType	single truckType dalam bentuk entity domain
 */
func (repository TruckTypeRepository) Store(truckType entities.TruckType) (entities.TruckType, error) {
	panic("implement me")
}
/*
 * Update
 * -------------------------------
 * Mengupdate data truckType berdasarkan ID
 *
 * @var truckType		single truckType entity
 * @return truckType	single truckType dalam bentuk entity domain
 * @return error		error
 */
func (repository TruckTypeRepository) Update(truckType entities.TruckType, id int) (entities.TruckType, error) {
	panic("implement me")
}
/*
 * Delete
 * -------------------------------
 * Delete truckType berdasarkan ID
 *
 * @return error		error	
 */
func (repository TruckTypeRepository) Delete(id int) error {
	panic("implement me")
}
/*
 * Delete Batch
 * -------------------------------
 * Delete multiple truckType berdasarkan filter tertentu
 *
 * @var filters	query untuk penyaringan data, { field, operator, value }
 *
 * @return error		error	
 */
func (repository TruckTypeRepository) DeleteBatch(filters []map[string]string) error {
	panic("implement me")
}