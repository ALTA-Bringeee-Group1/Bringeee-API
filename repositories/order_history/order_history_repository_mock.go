package order_history

import (
	"bringeee-capstone/entities"
	orderRepository "bringeee-capstone/repositories/order"
	"time"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type OrderHistoryRepositoryMock struct {
	Mock *mock.Mock
}

func NewOrderHistoryMock(mock *mock.Mock) *OrderHistoryRepositoryMock {
	return &OrderHistoryRepositoryMock{
		Mock: mock,
	}
}

var OrderHistoryCollection = []entities.OrderHistory {
	{
		Model: gorm.Model{ ID: 1, CreatedAt: time.Now().Add(time.Minute * -10), UpdatedAt: time.Now().Add(time.Minute * -10) },
		Log: "Order dibuat oleh customer",
		Actor: "customer",
		OrderID: 1,
		Order: orderRepository.OrderCollection[0],
	},
	{
		Model: gorm.Model{ ID: 2, CreatedAt: time.Now().Add(time.Minute * -9), UpdatedAt: time.Now().Add(time.Minute * -9) },
		Log: "Order dikonfirmasi oleh admin",
		Actor: "admin",
		OrderID: 1,
		Order: orderRepository.OrderCollection[0],
	},
	{
		Model: gorm.Model{ ID: 3, CreatedAt: time.Now().Add(time.Minute * -7), UpdatedAt: time.Now().Add(time.Minute * -7) },
		Log: "Order sudah dibayarkan dan siap diambil oleh driver",
		Actor: "customer",
		OrderID: 1,
		Order: orderRepository.OrderCollection[0],
	},
	{
		Model: gorm.Model{ ID: 4, CreatedAt: time.Now().Add(time.Minute * -7), UpdatedAt: time.Now().Add(time.Minute * -7) },
		Log: "Order diambil oleh driver untuk diantarkan",
		Actor: "driver",
		OrderID: 1,
		Order: orderRepository.OrderCollection[0],
	},
	{
		Model: gorm.Model{ ID: 5, CreatedAt: time.Now().Add(time.Minute * -5), UpdatedAt: time.Now().Add(time.Minute * -5) },
		Log: "Order tiba pada tujuan",
		Actor: "driver",
		OrderID: 1,
		Order: orderRepository.OrderCollection[0],
	},
	{
		Model: gorm.Model{ ID: 6, CreatedAt: time.Now().Add(time.Hour * -10), UpdatedAt: time.Now().Add(time.Hour * -10) },
		Log: "Order dibuat oleh customer",
		Actor: "customer",
		OrderID: 2,
		Order: orderRepository.OrderCollection[1],
	},
	{
		Model: gorm.Model{ ID: 7, CreatedAt: time.Now().Add(time.Hour * -9), UpdatedAt: time.Now().Add(time.Hour * -9) },
		Log: "Order dikonfirmasi oleh admin",
		Actor: "admin",
		OrderID: 2,
		Order: orderRepository.OrderCollection[1],
	},
	{
		Model: gorm.Model{ ID: 8, CreatedAt: time.Now().Add(time.Hour * -8), UpdatedAt: time.Now().Add(time.Hour * -8) },
		Log: "Order sudah dibayarkan dan siap diambil oleh driver",
		Actor: "customer",
		OrderID: 2,
		Order: orderRepository.OrderCollection[1],
	},
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
func (repo OrderHistoryRepositoryMock) FindAll(orderID int, sort []map[string]interface{}) ([]entities.OrderHistory, error)  {
	param := repo.Mock.Called(orderID, sort)
	return param.Get(0).([]entities.OrderHistory), param.Error(1)
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
func (repo OrderHistoryRepositoryMock) Create(orderID int, log string, actor string) (entities.OrderHistory, error) {
	param := repo.Mock.Called(orderID, log, actor)
	return param.Get(0).(entities.OrderHistory), param.Error(1)
}

/*
 * Delete history
 * -------------------------------
 * Delete history pada sebuah order 
 *
 * @var historyID				history yang akan di hapus
 * @return error 				error
 */
func (repo OrderHistoryRepositoryMock) Delete(historyID int) error {
	param := repo.Mock.Called(historyID)
	return param.Error(0)
}