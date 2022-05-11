package user_test

import (
	"bringeee-capstone/entities"
	orderRepository "bringeee-capstone/repositories/order"
	truckRepository "bringeee-capstone/repositories/truck_type"
	userRepository "bringeee-capstone/repositories/user"
	userService "bringeee-capstone/services/user"
	"testing"

	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFindAllCustomer(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		truckRepositoryMock := truckRepository.NewTruckTypeRepositoryMock(&mock.Mock{})
		orderRepositoryMock := orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindAllCustomer", 0, 0, []map[string]string{}, []map[string]interface{}{}).Return(userRepository.UserCollection[0], nil)

		service := userService.NewUserService(userRepositoryMock, truckRepositoryMock, orderRepositoryMock)
		data, err := service.FindAllCustomer(0, 0, []map[string]string{}, []map[string]interface{}{})

		// Konversi expected data ke response
		userRes := []entities.CustomerResponse{}
		copier.Copy(&userRes, &userRepository.UserCollection[0])

		assert.Nil(t, err, "Service FindAllCustomer error must be nil")
		assert.Equal(t, userRes, data, "User data collection is not returned")
	})

	t.Run("with-filter", func(t *testing.T) {
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		truckRepositoryMock := truckRepository.NewTruckTypeRepositoryMock(&mock.Mock{})
		orderRepositoryMock := orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindAllCustomer", 0, 0, []map[string]string{
			{"field": "name", "operator": "LIKE", "value": "test1"},
		}, []map[string]interface{}{}).Return(userRepository.UserCollection[1], nil)

		service := userService.NewUserService(userRepositoryMock, truckRepositoryMock, orderRepositoryMock)
		data, err := service.FindAllCustomer(0, 0, []map[string]string{
			{"field": "name", "operator": "LIKE", "value": "test1"},
		}, []map[string]interface{}{})

		// Konversi expected data ke response
		userRes := []entities.CustomerResponse{}
		copier.Copy(&userRes, &userRepository.UserCollection[0])

		assert.Nil(t, err, "Service FindAllCustomer error must be nil")
		assert.Equal(t, userRes, data, "User data collection is not returned")
	})

	t.Run("with-error", func(t *testing.T) {
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		truckRepositoryMock := truckRepository.NewTruckTypeRepositoryMock(&mock.Mock{})
		orderRepositoryMock := orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindAllCustomer", 0, 0, []map[string]string{}, []map[string]interface{}{}).Return(userRepository.UserCollection[0], "SERVER_ERROR")

		service := userService.NewUserService(userRepositoryMock, truckRepositoryMock, orderRepositoryMock)
		data, err := service.FindAllCustomer(0, 0, []map[string]string{}, []map[string]interface{}{})

		assert.Error(t, err, "error data is not returned")
		assert.Equal(t, []entities.CustomerResponse{}, data)
	})
}

func TestCountCustomer(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		truckRepositoryMock := truckRepository.NewTruckTypeRepositoryMock(&mock.Mock{})
		orderRepositoryMock := orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("CountAll").Return(userRepository.UserCollection, "")

		service := userService.NewUserService(userRepositoryMock, truckRepositoryMock, orderRepositoryMock)
		count, err := service.CountCustomer([]map[string]string{})

		assert.Nil(t, err, "error is not nil")
		assert.Equal(t, len(userRepository.UserCollection), count, "User data collection is not returned")
	})

	t.Run("with-filter", func(t *testing.T) {
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		truckRepositoryMock := truckRepository.NewTruckTypeRepositoryMock(&mock.Mock{})
		orderRepositoryMock := orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("CountAll").Return(userRepository.UserCollection, "")

		service := userService.NewUserService(userRepositoryMock, truckRepositoryMock, orderRepositoryMock)
		count, err := service.CountCustomer([]map[string]string{
			{"field": "truck_type", "operator": "LIKE", "value": "Pickup"},
		})

		assert.Nil(t, err, "error is not nil")
		assert.Equal(t, 2, count, "User data collection is not returned")
	})

	t.Run("with-error", func(t *testing.T) {
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		truckRepositoryMock := truckRepository.NewTruckTypeRepositoryMock(&mock.Mock{})
		orderRepositoryMock := orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("CountAll").Return(userRepository.UserCollection, "SERVER_ERROR")

		service := userService.NewUserService(userRepositoryMock, truckRepositoryMock, orderRepositoryMock)
		count, err := service.CountCustomer([]map[string]string{})

		assert.Error(t, err, "error data is not returned")
		assert.Equal(t, 0, count)
	})
}
