package order_history

import (
	"bringeee-capstone/entities"
	"bringeee-capstone/entities/web"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderHistoryRepository struct {
	db *gorm.DB
}

func NewOrderHistoryRepository(db *gorm.DB) *OrderHistoryRepository {
	return &OrderHistoryRepository{
		db: db,
	}
}


/*
 * Delete history
 * -------------------------------
 * Delete history pada sebuah order 
 *
 * @var orderID					ID Order yang akan di cancel
 * @return OrderHistory			response payment 
 * @return error 				error
 */
func (repository OrderHistoryRepository) FindAll(orderID int, sorts []map[string]interface{}) ([]entities.OrderHistory, error) {
	orderHistories := []entities.OrderHistory{}
	builder := repository.db.Where("order_id = ?", orderID)
	// OrderBy Filters
	for _, sort := range sorts {
		builder.Order(clause.OrderByColumn{Column: clause.Column{Name: sort["field"].(string)}, Desc: sort["desc"].(bool)})
	}
	tx := builder.Find(&orderHistories)
	if tx.Error != nil {
		return []entities.OrderHistory{}, web.WebError{Code: 500, ProductionMessage: "Server error", DevelopmentMessage: tx.Error.Error()}
	}
	return orderHistories, nil
}

/*
 * Create history
 * -------------------------------
 * Membuat history baru pada sebuah order 
 *
 * @var orderID					ID Order yang akan ditambahkan history
 * @return OrderHistory			response payment 
 * @return error 				error
 */
func (repository OrderHistoryRepository) Create(orderID int, log string, actor string) (entities.OrderHistory, error) {
	orderHistory := entities.OrderHistory {
		OrderID: uint(orderID),
		Log: log,
		Actor: actor,
	}
	tx := repository.db.Create(&orderHistory)
	if tx.Error != nil {
		return entities.OrderHistory{}, web.WebError{Code: 500, ProductionMessage: "Server error", DevelopmentMessage: tx.Error.Error()}
	}
	return orderHistory, nil
}

/*
 * Delete history
 * -------------------------------
 * Delete history pada sebuah order 
 *
 * @var historyID				history yang akan di hapus
 * @return error 				error
 */
func (repository OrderHistoryRepository) Delete(historyID int) error {
	tx := repository.db.Delete(&entities.OrderHistory{}, historyID)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}