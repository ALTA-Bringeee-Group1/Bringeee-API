package order

import (
	"bringeee-capstone/entities"
	"bringeee-capstone/entities/web"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	order := []entities.Order{}
	builder := repository.db.Limit(limit).Offset(offset).Preload("Destination").Preload("Customer").Preload("Driver").Preload("Driver.TruckType").Preload("Driver.User").Preload("TruckType")
	// Where filters
	for _, filter := range filters {
		builder.Where(filter["field"]+" "+filter["operator"]+" ?", filter["value"])
	}
	// OrderBy Filters
	for _, sort := range sorts {
		builder.Order(clause.OrderByColumn{Column: clause.Column{Name: sort["field"].(string)}, Desc: sort["desc"].(bool)})
	}
	tx := builder.Find(&order)
	if tx.Error != nil {
		return []entities.Order{}, web.WebError{Code: 500, Message: tx.Error.Error()}
	}
	return order, nil
}

/*
 * Find
 * -------------------------------
 * Mencari order tunggal berdasarkan ID
 *
 * @var id 		data id
 */
func (repository OrderRepository) Find(id int) (entities.Order, error) {
	order := entities.Order{}
	tx := repository.db.Preload("Destination").Preload("Customer").Preload("Driver").Preload("Driver.TruckType").Preload("Driver.User").Preload("TruckType").Find(&order, id)
	if tx.Error != nil {
		return entities.Order{}, web.WebError{Code: 500, DevelopmentMessage: "server error", ProductionMessage: "server error"}
	} else if tx.RowsAffected <= 0 {
		return entities.Order{}, web.WebError{Code: 400, DevelopmentMessage: "cannot get truck type data with specified id", ProductionMessage: "data error"}
	}
	return order, nil
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
	order := entities.Order{}
	tx := repository.db.Preload("Destination").Preload("Customer").Preload("Driver").Preload("Driver.TruckType").Preload("Driver.User").Preload("TruckType").Where(field+" = ?", value).First(&order)
	if tx.Error != nil {
		return entities.Order{}, web.WebError{Code: 500, Message: tx.Error.Error()}
	} else if tx.RowsAffected <= 0 {
		return entities.Order{}, web.WebError{Code: 400, ProductionMessage: "The requested ID doesn't match with any record", DevelopmentMessage: "The requested ID doesn't match with any record"}
	}
	return order, nil
}

/*
 * Find First
 * -------------------------------
 * Mengambil data order tunggal berdasarkan filter
 *
 * @var filters	query untuk penyaringan data, { field, operator, value }
 * @return order	order dalam bentuk entity domain
 * @return error	error
 */
func (repository OrderRepository) FindFirst(filters []map[string]string) (entities.Order, error) {
	order := entities.Order{}
	builder := repository.db.Preload("Destination").Preload("Customer").Preload("Driver").Preload("Driver.TruckType").Preload("Driver.User").Preload("TruckType")
	// Where filters
	for _, filter := range filters {
		builder.Where(filter["field"]+" "+filter["operator"]+" ?", filter["value"])
	}
	tx := builder.First(&order)
	if tx.Error != nil {
		return entities.Order{}, web.WebError{Code: 500, Message: tx.Error.Error()}
	}
	return order, nil
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
	var count int64
	builder := repository.db.Model(&entities.Order{})
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
 * Menambahkan data order kedalam database, beserta data destination
 *
 * @var order		single order entity
 * @var destination	single destination entity
 * @return order	single order dalam bentuk entity domain
 */
func (repository OrderRepository) Store(order entities.Order, destination entities.Destination) (entities.Order, error) {
	err := repository.db.Transaction(func(tx *gorm.DB) error {
		destinationTx := repository.db.Create(&destination)
		if destinationTx.Error != nil {
			return destinationTx.Error
		}
		order.DestinationID = destination.ID
		orderTx := repository.db.Create(&order)
		if orderTx.Error != nil {
			return orderTx.Error
		}
		return nil
	})
	if err != nil {
		return entities.Order{}, web.WebError{Code: 500, ProductionMessage: "Server error" ,DevelopmentMessage: err.Error()}
	}
	return order, nil
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
	tx := repository.db.Save(&order)
	if tx.Error != nil {
		return entities.Order{}, web.WebError{Code: 500, Message: tx.Error.Error()}
	}
	return order, nil
}

/*
 * Delete
 * -------------------------------
 * Delete order berdasarkan ID
 *
 * @return error	error	
 */
func (repository OrderRepository) Delete(id int, destinationID int) error {
	err := repository.db.Transaction(func(tx *gorm.DB) error {
		destination := repository.db.Delete(&entities.Destination{}, id)
		if destination.Error != nil {
			return destination.Error
		}
		order := repository.db.Delete(&entities.Order{}, id)
		if order.Error != nil {
			return order.Error
		}
		return nil
	})
	if err != nil {
		return web.WebError{Code: 500, Message: err.Error()}
	}
	return nil
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
	orders := []entities.Order{}
	builder := repository.db
	// Where filters
	for _, filter := range filters {
		builder.Where(filter["field"]+" "+filter["operator"]+" ?", filter["value"])
	}
	builder.Find(&orders)

	for _, order := range orders {
		repository.db.Transaction(func(tx *gorm.DB) error {
			orderTx := tx.Delete(&order)
			if orderTx != nil {
				return orderTx.Error
			}
			destinationTx := tx.Delete(&entities.Destination{}, "id = ?", order.DestinationID)
			if destinationTx != nil {
				return destinationTx.Error
			}
			return nil
		})
	}
	return nil
}