package order_history

import (
	"bringeee-capstone/entities"

	"gorm.io/gorm"
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
func (repository OrderHistoryRepository) FindAll(orderID int, sort []map[string]interface{}) ([]entities.OrderHistory, error) {
	panic("implement me")
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
func (repository OrderHistoryRepository) Create(orderID int) (entities.OrderHistory, error) {
	panic("implement me")
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
	panic("implement me")
}