package truck_type

import (
	"bringeee-capstone/entities"
	"bringeee-capstone/entities/web"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
func (repository TruckTypeRepository) FindAll(limit int, offset int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.TruckType, error) {
	truckType := []entities.TruckType{}
	builder := repository.db.Limit(limit).Offset(offset)
	// Where filters
	for _, filter := range filters {
		builder.Where(filter["field"]+" "+filter["operator"]+" ?", filter["value"])
	}
	// OrderBy Filters
	for _, sort := range sorts {
		builder.Order(clause.OrderByColumn{Column: clause.Column{Name: sort["field"].(string)}, Desc: sort["desc"].(bool)})
	}
	tx := builder.Find(&truckType)
	if tx.Error != nil {
		return []entities.TruckType{}, web.WebError{Code: 500, Message: tx.Error.Error()}
	}
	return truckType, nil
}
/*
 * Find
 * -------------------------------
 * Mencari truckType tunggal berdasarkan ID
 *
 * @var id 		data id
 */
func (repository TruckTypeRepository) Find(id int) (entities.TruckType, error) {
	truckType := entities.TruckType{}
	tx := repository.db.Find(&truckType, id)
	if tx.Error != nil {
		return entities.TruckType{}, web.WebError{Code: 500, DevelopmentMessage: "server error", ProductionMessage: "server error"}
	} else if tx.RowsAffected <= 0 {
		return entities.TruckType{}, web.WebError{Code: 400, DevelopmentMessage: "cannot get truck type data with specified id", ProductionMessage: "data error"}
	}
	return truckType, nil
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
	truckType := entities.TruckType{}
	tx := repository.db.Where(field+" = ?", value).First(&truckType)
	if tx.Error != nil {
		return entities.TruckType{}, web.WebError{Code: 500, Message: tx.Error.Error()}
	} else if tx.RowsAffected <= 0 {
		return entities.TruckType{}, web.WebError{Code: 400, ProductionMessage: "The requested ID doesn't match with any record", DevelopmentMessage: "The requested ID doesn't match with any record"}
	}
	return truckType, nil
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
	var count int64
	builder := repository.db.Model(&entities.TruckType{})
	// Where filters
	for _, filter := range filters {
		builder.Where(filter["field"]+" "+filter["operator"]+" ?", filter["value"])
	}
	tx := builder.Count(&count)
	if tx.Error != nil {
		return -1, web.WebError{Code: 400, ProductionMessage: tx.Error.Error(), DevelopmentMessage: tx.Error.Error()}
	}
	return count, nil
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
	tx := repository.db.Preload("User").Preload("Category").Create(&truckType)
	if tx.Error != nil {
		return entities.TruckType{}, web.WebError{Code: 500, Message: tx.Error.Error()}
	}
	return truckType, nil
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
	tx := repository.db.Save(&truckType)
	if tx.Error != nil {
		return entities.TruckType{}, web.WebError{Code: 500, Message: tx.Error.Error()}
	}
	return truckType, nil
}
/*
 * Delete
 * -------------------------------
 * Delete truckType berdasarkan ID
 *
 * @return error		error	
 */
func (repository TruckTypeRepository) Delete(id int) error {
	tx := repository.db.Delete(&entities.TruckType{}, id)
	if tx.Error != nil {
		return web.WebError{Code: 500, Message: tx.Error.Error()}
	}
	return nil
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
	builder := repository.db
	tx := builder.Delete(&entities.TruckType{}, filters[0]["field"]+" "+filters[0]["operator"]+" ?", filters[0]["value"])
	if tx.Error != nil {
		return web.WebError{Code: 400, Message: tx.Error.Error()}
	}
	return nil
}