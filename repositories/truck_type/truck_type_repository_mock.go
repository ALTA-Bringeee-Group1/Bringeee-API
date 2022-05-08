package truck_type

import (
	"bringeee-capstone/entities"
	"bringeee-capstone/entities/web"
	"reflect"
	"sort"
	"strings"
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
var truckTypeCollection = []entities.TruckType {
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

	// Mock data source
	args := repo.Mock.Called()
	switch args.Get(0).(type) {
	case []entities.TruckType:
		truckTypeCollection = args.Get(0).([]entities.TruckType)
	}

	// Mock errors
	if args.String(1) == "SERVER_ERROR" {
		return []entities.TruckType{}, web.WebError{Code: 500, DevelopmentMessage: "server error", ProductionMessage: "server error"}
	}

	// Return all data if filter is empty
	var filteredCollection []entities.TruckType
	if len(filters) <= 0 {
		return truckTypeCollection, nil
	}
	for _, item := range truckTypeCollection {
		// Filtering
		for _, filter := range filters {
			if filter["operator"] == "LIKE" {
				value := reflect.Indirect(reflect.ValueOf(item)).FieldByName(colMapper[filter["field"]]).String()
				if strings.Contains(value, filter["value"]) {
					filteredCollection = append(filteredCollection, item)
				}
			}
		}
		// Sort
		for _, sortItem := range sorts {
			sort.Slice(filteredCollection, func(i, j int) bool {
				value := reflect.Indirect(reflect.ValueOf(filteredCollection[0])).FieldByName(colMapper[sortItem["field"].(string)]).String()
				value2 := reflect.Indirect(reflect.ValueOf(filteredCollection[1])).FieldByName(colMapper[sortItem["field"].(string)]).String()
				if sortItem["desc"].(bool) {
					return value > value2 
				} else {
					return value < value2 
				}
			})
		}
	}
	return filteredCollection, nil
}

/*
 * Find
 * -------------------------------
 * Mencari truckType tunggal berdasarkan ID
 *
 * @var id 		data id
 */
func (repo TruckTypeRepositoryMock) Find(id int) (entities.TruckType, error) {
	// Mock data source
	args := repo.Mock.Called(id)
	switch args.Get(0).(type) {
	case []entities.TruckType:
		truckTypeCollection = args.Get(0).([]entities.TruckType)
	}

	// Mock errors
	if args.String(1) == "SERVER_ERROR" {
		return entities.TruckType{}, web.WebError{Code: 500, DevelopmentMessage: "server error", ProductionMessage: "server error"}
	}

	for _, truckType := range truckTypeCollection {
		if truckType.ID == uint(id) {
			return truckType, nil
		}
	}
	return entities.TruckType{}, web.WebError{Code: 400, DevelopmentMessage: "cannot get truck type data with specified id", ProductionMessage: "data error"}
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
	panic("Implement Me")

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
	// Mock data source
	args := repo.Mock.Called()
	switch args.Get(0).(type) {
	case []entities.TruckType:
		truckTypeCollection = args.Get(0).([]entities.TruckType)
	}

	// Mock errors
	if args.String(1) == "SERVER_ERROR" {
		return 0, web.WebError{Code: 500, DevelopmentMessage: "server error", ProductionMessage: "server error"}
	}

	// Return all data if filter is empty
	var filteredCollection []entities.TruckType
	if len(filters) <= 0 {
		return int64(len(truckTypeCollection)), nil
	}
	for _, item := range truckTypeCollection {
		// Filtering
		for _, filter := range filters {
			if filter["operator"] == "LIKE" {
				value := reflect.Indirect(reflect.ValueOf(item)).FieldByName(colMapper[filter["field"]]).String()
				if strings.Contains(value, filter["value"]) {
					filteredCollection = append(filteredCollection, item)
				}
			}
		}
	}
	return int64(len(filteredCollection)), nil

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
	panic("Implement Me")

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
	panic("Implement Me")

}
/*
 * Delete
 * -------------------------------
 * Delete truckType berdasarkan ID
 *
 * @return error		error	
 */
func (repo TruckTypeRepositoryMock) Delete(id int) error {
	panic("Gud")
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
	panic("Gud")
}