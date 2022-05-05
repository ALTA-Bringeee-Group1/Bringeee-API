package order

import (
	"bringeee-capstone/entities"
	"bringeee-capstone/entities/web"
	"fmt"
	"time"

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
func (repository OrderRepository) FindAll(limit int, offset int, filters []map[string]interface{}, sorts []map[string]interface{}) ([]entities.Order, error) {
	order := []entities.Order{}
	builder := repository.db.Limit(limit).Offset(offset).Preload("Destination").Preload("Customer").Preload("Driver").Preload("Driver.TruckType").Preload("Driver.User").Preload("TruckType")
	// Where filters
	for _, filter := range filters {
		orGroup, orGroupExist := filter["or"]
		if orGroupExist {
			orGroupMap := orGroup.([]map[string]string)
			orBuilder := repository.db
			for _, orQuery := range orGroupMap {
				orBuilder = orBuilder.Or(fmt.Sprintf("%s %s ?", orQuery["field"], orQuery["operator"]), orQuery["value"])
			}
			builder.Where(orBuilder)
		} else {
			builder.Where(filter["field"].(string)+" "+filter["operator"].(string)+" ?", filter["value"].(string))
		}
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
		return entities.Order{}, web.WebError{Code: 400, DevelopmentMessage: "cannot get order data with specified id", ProductionMessage: "data error"}
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
func (repository OrderRepository) FindFirst(filters []map[string]interface{}) (entities.Order, error) {
	order := entities.Order{}
	builder := repository.db.Preload("Destination").Preload("Customer").Preload("Driver").Preload("Driver.TruckType").Preload("Driver.User").Preload("TruckType")
	// Where filters
	for _, filter := range filters {
		orGroup, orGroupExist := filter["or"]
		if orGroupExist {
			orGroupMap := orGroup.([]map[string]string)
			orBuilder := repository.db
			for _, orQuery := range orGroupMap {
				orBuilder = orBuilder.Or(fmt.Sprintf("%s %s ?", orQuery["field"], orQuery["operator"]), orQuery["value"])
			}
			builder.Where(orBuilder)
		} else {
			builder.Where(filter["field"].(string)+" "+filter["operator"].(string)+" ?", filter["value"].(string))
		}
	}
	tx := builder.First(&order)
	if tx.RowsAffected <= 0 {
		return entities.Order{}, web.WebError{Code: 200, Message: "no result"}
	} else if tx.Error != nil {
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
func (repository OrderRepository) CountAll(filters []map[string]interface{}) (int64, error) {
	var count int64
	builder := repository.db.Model(&entities.Order{})
	// Where filters
	for _, filter := range filters {
		orGroup, orGroupExist := filter["or"]
		if orGroupExist {
			orGroupMap := orGroup.([]map[string]string)
			orBuilder := repository.db
			for _, orQuery := range orGroupMap {
				orBuilder = orBuilder.Or(fmt.Sprintf("%s %s ?", orQuery["field"], orQuery["operator"]), orQuery["value"])
			}
			builder.Where(orBuilder)
		} else {
			builder.Where(filter["field"].(string)+" "+filter["operator"].(string)+" ?", filter["value"].(string))
		}
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
		return entities.Order{}, web.WebError{Code: 500, ProductionMessage: "Server error", DevelopmentMessage: err.Error()}
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
func (repository OrderRepository) DeleteBatch(filters []map[string]interface{}) error {
	orders := []entities.Order{}
	builder := repository.db
	// Where filters
	for _, filter := range filters {
		orGroup, orGroupExist := filter["or"]
		if orGroupExist {
			orGroupMap := orGroup.([]map[string]string)
			orBuilder := repository.db
			for _, orQuery := range orGroupMap {
				orBuilder = orBuilder.Or(fmt.Sprintf("%s %s ?", orQuery["field"], orQuery["operator"]), orQuery["value"])
			}
			builder = builder.Where(orBuilder)
		} else {
			builder = builder.Where(filter["field"].(string)+" "+filter["operator"].(string)+" ?", filter["value"].(string))
		}
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

func (repository OrderRepository) FindByDate(day int) ([]map[string]interface{}, error) {
	result := []map[string]interface{}{}
	var count int64
	start := time.Now().AddDate(0, 0, -(day))
	end := time.Now()
	for d := start; d.After(end) == false; d = d.AddDate(0, 0, 1) {
		orders := []entities.Order{}
		repository.db.Where("created_at LIKE ?", "%"+d.Format("2006-01-02")+"%").Find(&orders).Count(&count)
		result = append(result, map[string]interface{}{
			"label": d.Format("2006-01-02"),
			"value": int(count),
		})
	}
	return result, nil
}
