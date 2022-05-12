package truckType

import (
	"bringeee-capstone/entities"
	"time"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type TruckTypeRepositoryMock struct {
	Mock *mock.Mock
}

/* 
 * Collection
 * ---------------------------
 * kumpulan mock data jenis truk untuk dapat dilakukan
 * query dan command sama halnya dengan repository sebenarnya 
 */
var TruckTypeCollection = []entities.TruckType {
	{
		Model: gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		TruckType: "Pickup Truck - A", 
		MaxVolume: 16000000,
		MaxWeight: 1000,
		PricePerDistance: 2000,
	},
	{
		Model: gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		TruckType: "Pickup Truck - B", 
		MaxVolume: 8000000,
		MaxWeight: 2000,
		PricePerDistance: 3000,
	},
	{
		Model: gorm.Model{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		TruckType: "Pickup Truck - B", 
		MaxVolume: 8000000,
		MaxWeight: 2000,
		PricePerDistance: 3000,
	},
}

/*
 * Column to Struct Mapper
 * -----------------------------
 * memetakan kolom database ke struct
 */
var colMapper = map[string]string {
	"id": "ID",
	"created_at": "CreatedAt",
	"updated_at": "UpdatedAt",
	"deleted_at": "DeletedAt",
	"truck_type": "TruckType",
	"max_volume": "MaxVolume",
	"max_weight": "MaxWeight",
	"price_per_distance": "PricePerDistance",
}



func NewTruckTypeRepositoryMock(mock *mock.Mock) *TruckTypeRepositoryMock {
	return &TruckTypeRepositoryMock{
		Mock: mock,
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
func (repo TruckTypeRepositoryMock) FindAll(limit int, offset int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.TruckType, error) {
	param := repo.Mock.Called(limit, offset, filters, sorts)
	return param.Get(0).([]entities.TruckType), param.Error(1) 
}

/*
 * Find
 * -------------------------------
 * Mencari truckType tunggal berdasarkan ID
 *
 * @var id 		data id
 */
func (repo TruckTypeRepositoryMock) Find(id int) (entities.TruckType, error) {
	param := repo.Mock.Called(id)
	return param.Get(0).(entities.TruckType), param.Error(1) 
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
func (repo TruckTypeRepositoryMock) FindBy(field string, value string) (entities.TruckType, error) {
	param := repo.Mock.Called(field, value)
	return param.Get(0).(entities.TruckType), param.Error(1) 
}
/*
 * CountAll
 * -------------------------------
 * Menghitung semua truckTypes (ini digunakan untuk pagination di service)
 *
 * @return truckType	single truckType dalam bentuk entity domain
 * @return error		error
 */
func (repo TruckTypeRepositoryMock) CountAll(filters []map[string]string) (int64, error) {
	param := repo.Mock.Called()
	return int64(param.Int(0)), param.Error(1)
}
/*
 * Store
 * -------------------------------
 * Menambahkan data truckType kedalam database
 *
 * @var truckType		single truckType entity
 * @return truckType	single truckType dalam bentuk entity domain
 */
func (repo TruckTypeRepositoryMock) Store(truckType entities.TruckType) (entities.TruckType, error) {
	param := repo.Mock.Called(truckType)
	return param.Get(0).(entities.TruckType), param.Error(1)
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
func (repo TruckTypeRepositoryMock) Update(truckType entities.TruckType, id int) (entities.TruckType, error) {
	param := repo.Mock.Called(truckType)
	return param.Get(0).(entities.TruckType), param.Error(1)
}
/*
 * Delete
 * -------------------------------
 * Delete truckType berdasarkan ID
 *
 * @return error		error	
 */
func (repo TruckTypeRepositoryMock) Delete(id int) error {
	param := repo.Mock.Called(id)
	return param.Error(1)
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
func (repo TruckTypeRepositoryMock) DeleteBatch(filters []map[string]string) error {
	param := repo.Mock.Called(filters)
	return param.Error(1)
}