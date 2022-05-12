package orderHistory

import "bringeee-capstone/entities"

type OrderHistoryRepositoryInterface interface {

	/*
	 * Delete history
	 * -------------------------------
	 * Delete history pada sebuah order 
	 *
	 * @var orderID					ID Order yang akan di cancel
	 * @return OrderHistory			response payment 
	 * @return error 				error
	 */
	 FindAll(orderID int, sort []map[string]interface{}) ([]entities.OrderHistory, error) 

	/*
	 * Create history
	 * -------------------------------
	 * Membuat history baru pada sebuah order 
	 *
	 * @var orderID					ID Order yang akan ditambahkan history
	 * @return OrderHistory			response payment 
	 * @return error 				error
	 */
	 Create(orderID int, log string, actor string) (entities.OrderHistory, error)

	/*
	 * Delete history
	 * -------------------------------
	 * Delete history pada sebuah order 
	 *
	 * @var historyID				history yang akan di hapus
	 * @return error 				error
	 */
	 Delete(historyID int) error
}