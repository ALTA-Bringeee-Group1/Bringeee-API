package order_test

import (
	"bringeee-capstone/entities"
	"bringeee-capstone/entities/web"
	_distanceMatrixRepository "bringeee-capstone/repositories/distance_matrix"
	_orderRepository "bringeee-capstone/repositories/order"
	_orderHistoryRepository "bringeee-capstone/repositories/order_history"
	_paymentRepository "bringeee-capstone/repositories/payment"
	_truckTypeRepository "bringeee-capstone/repositories/truck_type"
	_userRepository "bringeee-capstone/repositories/user"
	_orderService "bringeee-capstone/services/order"
	_storageProvider "bringeee-capstone/services/storage"
	"mime/multipart"
	"testing"

	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/guregu/null.v4"
)

/*
 * Collection
 * ---------------------------
 * kumpulan mock data untuk dapat dilakukan
 * query dan command sama halnya dengan repository sebenarnya
 */

func TestFindAll(t *testing.T) {
	
	t.Run("success-all", func(t *testing.T) {
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("FindAll", 0, 0, []map[string]interface{}{}, []map[string]interface{}{}).Return(_orderRepository.OrderCollection, nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			_orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{}),
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		actual, err := orderService.FindAll(0, 0, []map[string]interface{}{}, []map[string]interface{}{})

		// Konversi ke order response
		expected := []entities.OrderResponse{}
		copier.Copy(&expected, _orderRepository.OrderCollection)
		for i, order := range _orderRepository.OrderCollection {
			copier.Copy(&expected[i], &order.Destination)
			copier.Copy(&expected[i].Driver, &order.Driver.User)
			expected[i].ID = order.ID // fix: overlap destinationID vs orderID
		}

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("with-filter", func(t *testing.T) {
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("FindAll", 0, 0, []map[string]interface{}{
			{ "field": "status", "operator": "=", "value": "DELIVERED" },
		}, []map[string]interface{}{}).Return([]entities.Order{_orderRepository.OrderCollection[0]}, nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			_orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{}),
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		actual, err := orderService.FindAll(0, 0, []map[string]interface{}{
			{ "field": "status", "operator": "=", "value": "DELIVERED" },
		}, []map[string]interface{}{})

		// Konversi ke order response
		expected := []entities.OrderResponse{}
		copier.Copy(&expected, []entities.Order{_orderRepository.OrderCollection[0]})
		for i, order := range []entities.Order{_orderRepository.OrderCollection[0]} {
			copier.Copy(&expected[i], &order.Destination)
			copier.Copy(&expected[i].Driver, &order.Driver.User)
			expected[i].ID = order.ID // fix: overlap destinationID vs orderID
		}

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("error", func(t *testing.T) {
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("FindAll", 0, 0, []map[string]interface{}{}, []map[string]interface{}{}).Return(
			[]entities.Order{}, web.WebError{ Code: 500, Message: "Internal server error"},
		)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			_orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{}),
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		actual, err := orderService.FindAll(0, 0, []map[string]interface{}{}, []map[string]interface{}{})

		// Konversi ke order response
		expected := []entities.OrderResponse{}

		assert.Error(t, err)
		assert.Equal(t, expected, actual)
	})
}


func TestGetPagination(t *testing.T) {
	t.Run("all", func(t *testing.T) {
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("CountAll", []map[string]interface{}{}).Return(2, nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			_orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{}),
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		actual, err := orderService.GetPagination(2, 1, []map[string]interface{}{})
		expected := web.Pagination {
			Page: 1,
			Limit: 2,
			TotalPages: 1,
			TotalRecords: 2,
		}
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("single-row", func(t *testing.T) {
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("CountAll", []map[string]interface{}{}).Return(2, nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			_orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{}),
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		actual, err := orderService.GetPagination(1, 1, []map[string]interface{}{})
		expected := web.Pagination {
			Page: 1,
			Limit: 1,
			TotalPages: 2,
			TotalRecords: 2,
		}
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("error", func(t *testing.T) {
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("CountAll", []map[string]interface{}{}).Return(0, web.WebError{ Code: 500, Message: "Error" })

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			_orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{}),
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		actual, err := orderService.GetPagination(1, 1, []map[string]interface{}{})
		expected := web.Pagination{}
		assert.Error(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestFindFirst(t *testing.T) {
	t.Run("all", func(t *testing.T) {
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("FindFirst",[]map[string]interface{}{
			{ "field": "driver_id",  "operator": "=", "value": "3"},
			{ "field": "status",  "operator": "=", "value": "DELIVERED"},
		}).Return(_orderRepository.OrderCollection[0], nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			_orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{}),
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		actual, err := orderService.FindFirst([]map[string]interface{}{
			{ "field": "driver_id",  "operator": "=", "value": "3"},
			{ "field": "status",  "operator": "=", "value": "DELIVERED"},
		})
		
		expected := entities.OrderResponse{}
		copier.Copy(&expected, &_orderRepository.OrderCollection[0])
		copier.Copy(&expected, &_orderRepository.OrderCollection[0].Destination)
		copier.Copy(&expected.Driver, &_orderRepository.OrderCollection[0].Driver.User)
		expected.ID = _orderRepository.OrderCollection[0].ID // fix: overlap destinationID vs orderID

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("error", func(t *testing.T) {
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("FindFirst",[]map[string]interface{}{
			{ "field": "driver_id",  "operator": "=", "value": "3"},
			{ "field": "status",  "operator": "=", "value": "DELIVERED"},
		}).Return(entities.Order{}, web.WebError{ Code: 400, Message: "Not found" })

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			_orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{}),
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		actual, err := orderService.FindFirst([]map[string]interface{}{
			{ "field": "driver_id",  "operator": "=", "value": "3"},
			{ "field": "status",  "operator": "=", "value": "DELIVERED"},
		})
		expected := entities.OrderResponse{}

		assert.Error(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestFind(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		orderOutput := _orderRepository.OrderCollection[0]
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", 1).Return(orderOutput, nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			_orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{}),
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		actual, err := orderService.Find(int(orderOutput.ID))

		// convert to response
		expected := entities.OrderResponse{}
		copier.Copy(&expected, &orderOutput)
		copier.Copy(&expected, &orderOutput.Destination)
		copier.Copy(&expected.Driver, &orderOutput.Driver.User)
		expected.ID = orderOutput.ID // fix: overlap destinationID vs orderID

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("error", func(t *testing.T) {
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", 1).Return(entities.Order{}, web.WebError{ Code: 500, Message: "Error"})

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			_orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{}),
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		actual, err := orderService.Find(1)
		assert.Error(t, err)
		assert.Equal(t, entities.OrderResponse{}, actual)
	})
}


func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Request
		user := _userRepository.UserCollection[0]
		orderRequest := entities.CustomerCreateOrderRequest {
			DestinationStartProvince: "DI YOGYAKARTA",
			DestinationStartCity: "YOGYAKARTA",
			DestinationStartDistrict: "DANUREJAN",
			DestinationStartAddress: "Suryatmajan, Kec. Danurejan, Kota Yogyakarta, Daerah Istimewa Yogyakarta",
			DestinationStartPostal: "55213",
			DestinationStartLat: "-7.793050394271023",
			DestinationStartLong: "110.36756312713727",
			DestinationEndProvince: "JAWA TENGAH",
			DestinationEndCity: "SURAKARTA (SOLO)",
			DestinationEndDistrict: "JEBRES",
			DestinationEndAddress: "Tegalharjo, Kec. Jebres, Kota Surakarta, Jawa Tengah",
			DestinationEndPostal: "57129",
			DestinationEndLat: "-7.561160260537138",
			DestinationEndLong: "110.83655443176414",
			TruckTypeID: 1, 
			TotalVolume: 8000, 
			TotalWeight: 200, 
			Description: "", 
		}
		files := map[string]*multipart.FileHeader{
			"order_picture": {
				Filename: "order.png",
				Size: 1024 * 56,
			},
		}

		// Truck Type Exists
		truckTypeOutput := _truckTypeRepository.TruckTypeCollection[0]
		truckTypeRepository := _truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{})
		truckTypeRepository.Mock.On("Find", 1).Return(truckTypeOutput, nil)

		// Storage Upload Mock
		storage := _storageProvider.NewStorageMock(&mock.Mock{})
		storageOutput := "example.com/orders/order_picture/" + files["order_picture"].Filename
		storage.Mock.On("UploadFromRequest").Return(storageOutput, nil)

		// Distance Matrix Mock
		distanceMatrixOutput := _distanceMatrixRepository.DistanceMatrixCollection[0]
		distanceMatrixRepository := _distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{})
		distanceMatrixRepository.Mock.On(
			"EstimateShortest", 
			orderRequest.DestinationStartLat,
			orderRequest.DestinationStartLong,
			orderRequest.DestinationEndLat,
			orderRequest.DestinationEndLong,
		).Return(
			distanceMatrixOutput,
			nil,
		)

		// order repository mock
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderOutput := _orderRepository.OrderCollection[0]
		orderRepositoryMock.Mock.On("Store").Return(orderOutput, nil)
		orderRepositoryMock.Mock.On("Find", int(orderOutput.ID)).Return(orderOutput, nil)
		
		// order history mock
		orderHistoryOutput := _orderHistoryRepository.OrderHistoryCollection[0]
		orderHistoryRepositoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})
		orderHistoryRepositoryMock.Mock.On(
			"Create", 
			int(orderOutput.ID), 
			"Order dibuat dan diajukan oleh customer", 
			"customer",
		).Return(orderHistoryOutput, nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryRepositoryMock,
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			distanceMatrixRepository,
			truckTypeRepository,
		)

		actual, err := orderService.Create(orderRequest, files, int(user.ID), storage)
		// convert to response
		expected := entities.OrderResponse{}
		copier.Copy(&expected, &orderOutput)
		copier.Copy(&expected, &orderOutput.Destination)
		copier.Copy(&expected.Driver, &orderOutput.Driver.User)
		expected.ID = orderOutput.ID // fix: overlap destinationID vs orderID

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
}


func TestSetFixOrder(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		fixPriceSample := int64(800000)
		setPriceRequest := entities.AdminSetPriceOrderRequest{
			FixedPrice: uint64(fixPriceSample),
		}

		orderSample := _orderRepository.OrderCollection[2]
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)
		orderUpdateSample := orderSample
		orderUpdateSample.FixPrice = fixPriceSample
		orderUpdateSample.Status = "NEED_CUSTOMER_CONFIRM"
		orderRepositoryMock.Mock.On("Update", orderUpdateSample).Return(orderUpdateSample, nil)

		orderHistoryRepositoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})
		historyOutput := _orderHistoryRepository.OrderHistoryCollection[1]
		orderHistoryRepositoryMock.Mock.On("Create", int(orderSample.ID), historyOutput.Log, historyOutput.Actor).Return(historyOutput, nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryRepositoryMock,
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		err := orderService.SetFixOrder(int(orderSample.ID), setPriceRequest)
		assert.Nil(t, err)
	})

	t.Run("confirmed-status", func(t *testing.T) {
		fixPriceSample := int64(800000)
		setPriceRequest := entities.AdminSetPriceOrderRequest{
			FixedPrice: uint64(fixPriceSample),
		}

		orderSample := _orderRepository.OrderCollection[0]
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)
		orderUpdateSample := orderSample
		orderUpdateSample.FixPrice = fixPriceSample
		orderUpdateSample.Status = "NEED_CUSTOMER_CONFIRM"
		orderRepositoryMock.Mock.On("Update", orderUpdateSample).Return(orderUpdateSample, nil)

		orderHistoryRepositoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})
		historyOutput := _orderHistoryRepository.OrderHistoryCollection[1]
		orderHistoryRepositoryMock.Mock.On("Create", int(orderSample.ID), historyOutput.Log, historyOutput.Actor).Return(historyOutput, nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryRepositoryMock,
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		err := orderService.SetFixOrder(int(orderSample.ID), setPriceRequest)
		assert.Error(t, err)
	})
	t.Run("failed-repository", func(t *testing.T) {
		fixPriceSample := int64(800000)
		setPriceRequest := entities.AdminSetPriceOrderRequest{
			FixedPrice: uint64(fixPriceSample),
		}

		orderSample := _orderRepository.OrderCollection[0]
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)
		orderUpdateSample := orderSample
		orderUpdateSample.FixPrice = fixPriceSample
		orderUpdateSample.Status = "NEED_CUSTOMER_CONFIRM"
		orderRepositoryMock.Mock.On("Update", orderUpdateSample).Return(entities.Order{}, web.WebError{ Code: 500, Message: "Error" })

		orderHistoryRepositoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})
		historyOutput := _orderHistoryRepository.OrderHistoryCollection[1]
		orderHistoryRepositoryMock.Mock.On("Create", int(orderSample.ID), historyOutput.Log, historyOutput.Actor).Return(historyOutput, nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryRepositoryMock,
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		err := orderService.SetFixOrder(int(orderSample.ID), setPriceRequest)
		assert.NotNil(t, err)
	})
}

func TestConfirmOrder(t *testing.T) {
	t.Run("customer-order-requested", func(t *testing.T) {		
		orderSample := _orderRepository.OrderCollection[2]
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)

		orderHistoryRepositoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})
		historyOutput := _orderHistoryRepository.OrderHistoryCollection[2]
		orderHistoryRepositoryMock.Mock.On("Create", int(orderSample.ID), historyOutput.Log, historyOutput.Actor).Return(historyOutput, nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryRepositoryMock,
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		err := orderService.ConfirmOrder(int(orderSample.ID), 1, false)
		assert.Error(t, err)
	})
	t.Run("customer-order-need-confirm-success", func(t *testing.T) {		
		orderSample := _orderRepository.OrderCollection[2]
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderSample.Status = "NEED_CUSTOMER_CONFIRM"
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)
		orderOutput := orderSample
		orderOutput.Status = "CONFIRMED"
		orderOutput.FixPrice = 3000000
		orderRepositoryMock.Mock.On("Update", orderOutput).Return(orderOutput, nil)

		orderHistoryRepositoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})
		historyOutput := _orderHistoryRepository.OrderHistoryCollection[2]
		orderHistoryRepositoryMock.Mock.On("Create", int(orderSample.ID), historyOutput.Log, historyOutput.Actor).Return(historyOutput, nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryRepositoryMock,
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		err := orderService.ConfirmOrder(int(orderSample.ID), 1, false)
		assert.Nil(t, err)
	})
	t.Run("order-notfound", func(t *testing.T) {		
		orderSample := _orderRepository.OrderCollection[2]
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(entities.Order{}, web.WebError{ Code: 500 })

		orderHistoryRepositoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})
		historyOutput := _orderHistoryRepository.OrderHistoryCollection[2]
		orderHistoryRepositoryMock.Mock.On("Create", int(orderSample.ID), historyOutput.Log, historyOutput.Actor).Return(historyOutput, nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryRepositoryMock,
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		err := orderService.ConfirmOrder(int(orderSample.ID), 1, false)
		assert.NotNil(t, err)
	})
	t.Run("order-onprocess", func(t *testing.T) {		
		orderSample := _orderRepository.OrderCollection[2]
		orderSample.Status = "ON_PROCESS"
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)

		orderHistoryRepositoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})
		historyOutput := _orderHistoryRepository.OrderHistoryCollection[2]
		orderHistoryRepositoryMock.Mock.On("Create", int(orderSample.ID), historyOutput.Log, historyOutput.Actor).Return(historyOutput, nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryRepositoryMock,
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		err := orderService.ConfirmOrder(int(orderSample.ID), 1, false)
		assert.NotNil(t, err)
	})
	t.Run("customer-doesnt-belong-order", func(t *testing.T) {		
		orderSample := _orderRepository.OrderCollection[2]
		orderSample.Status = "NEED_CUSTOMER_CONFIRM"
		orderSample.CustomerID = 2
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)

		orderHistoryRepositoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})
		historyOutput := _orderHistoryRepository.OrderHistoryCollection[2]
		orderHistoryRepositoryMock.Mock.On("Create", int(orderSample.ID), historyOutput.Log, historyOutput.Actor).Return(historyOutput, nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryRepositoryMock,
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		err := orderService.ConfirmOrder(int(orderSample.ID), 1, false)
		assert.NotNil(t, err)
	})
	t.Run("admin-need-customer-confirm", func(t *testing.T) {		
		orderSample := _orderRepository.OrderCollection[2]
		orderSample.Status = "NEED_CUSTOMER_CONFIRM"
		orderSample.CustomerID = 2
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)

		orderHistoryRepositoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})
		historyOutput := _orderHistoryRepository.OrderHistoryCollection[2]
		orderHistoryRepositoryMock.Mock.On("Create", int(orderSample.ID), historyOutput.Log, historyOutput.Actor).Return(historyOutput, nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryRepositoryMock,
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		err := orderService.ConfirmOrder(int(orderSample.ID), 1, true)
		assert.NotNil(t, err)
	})
	t.Run("failed-repository", func(t *testing.T) {		
		orderSample := _orderRepository.OrderCollection[2]
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderSample.Status = "NEED_CUSTOMER_CONFIRM"
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)
		orderOutput := orderSample
		orderOutput.Status = "CONFIRMED"
		orderOutput.FixPrice = 3000000
		orderRepositoryMock.Mock.On("Update", orderOutput).Return(entities.Order{}, web.WebError{})

		orderHistoryRepositoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})
		historyOutput := _orderHistoryRepository.OrderHistoryCollection[2]
		orderHistoryRepositoryMock.Mock.On("Create", int(orderSample.ID), historyOutput.Log, historyOutput.Actor).Return(historyOutput, nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryRepositoryMock,
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		err := orderService.ConfirmOrder(int(orderSample.ID), 1, false)
		assert.NotNil(t, err)
	})
}

func TestCancelOrder(t *testing.T) {
	t.Run("customer-cancel-success", func(t *testing.T) {		
		orderSample := _orderRepository.OrderCollection[2]
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderSample.Status = "CONFIRMED"
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)
		orderOutput := orderSample
		orderOutput.Status = "CANCELLED"
		orderRepositoryMock.Mock.On("Update", orderOutput).Return(orderOutput, nil)

		orderHistoryRepositoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})
		historyOutput := _orderHistoryRepository.OrderHistoryCollection[8]
		orderHistoryRepositoryMock.Mock.On("Create", int(orderSample.ID), historyOutput.Log, historyOutput.Actor).Return(historyOutput, nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryRepositoryMock,
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		err := orderService.CancelOrder(int(orderSample.ID), 1, false)
		assert.Nil(t, err)
	})
	t.Run("customer-cancel-doesnt-belong", func(t *testing.T) {		
		orderSample := _orderRepository.OrderCollection[2]
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderSample.Status = "CONFIRMED"
		orderSample.CustomerID = 2
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)

		orderHistoryRepositoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})
		historyOutput := _orderHistoryRepository.OrderHistoryCollection[8]
		orderHistoryRepositoryMock.Mock.On("Create", int(orderSample.ID), historyOutput.Log, historyOutput.Actor).Return(historyOutput, nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryRepositoryMock,
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		err := orderService.CancelOrder(int(orderSample.ID), 1, false)
		assert.Error(t, err)
	})
	t.Run("invalid-order", func(t *testing.T) {		
		orderSample := _orderRepository.OrderCollection[2]
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(entities.Order{}, web.WebError{})
		orderHistoryRepositoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})
		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryRepositoryMock,
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		err := orderService.CancelOrder(int(orderSample.ID), 1, false)
		assert.Error(t, err)
	})
	t.Run("order-on-process", func(t *testing.T) {		
		orderSample := _orderRepository.OrderCollection[2]
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderSample.Status = "ON_PROCESS"
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)
		orderHistoryRepositoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})
		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryRepositoryMock,
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		err := orderService.CancelOrder(int(orderSample.ID), 1, false)
		assert.Error(t, err)
	})
	t.Run("failed-update-repo", func(t *testing.T) {		
		orderSample := _orderRepository.OrderCollection[2]
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderSample.Status = "CONFIRMED"
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)
		orderOutput := orderSample
		orderOutput.Status = "CANCELLED"
		orderRepositoryMock.Mock.On("Update", orderOutput).Return(entities.Order{}, web.WebError{})

		orderHistoryRepositoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})
		historyOutput := _orderHistoryRepository.OrderHistoryCollection[8]
		orderHistoryRepositoryMock.Mock.On("Create", int(orderSample.ID), historyOutput.Log, historyOutput.Actor).Return(historyOutput, nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryRepositoryMock,
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		err := orderService.CancelOrder(int(orderSample.ID), 1, false)
		assert.Error(t, err)
	})
}

func TestCreatePayment(t *testing.T) {
	requestSample := entities.CreatePaymentRequest {}
	paymentRepositoryMock := _paymentRepository.NewPaymentRepositoryMock(&mock.Mock{})
	paymentRepositoryMock.Mock.On("CreateBankTransferBNI").Return(_paymentRepository.PaymentResponseCollection[0], nil)
	paymentRepositoryMock.Mock.On("CreateBankTransferBCA").Return(_paymentRepository.PaymentResponseCollection[1], nil)
	paymentRepositoryMock.Mock.On("CreateBankTransferBRI").Return(_paymentRepository.PaymentResponseCollection[2], nil)
	paymentRepositoryMock.Mock.On("CreateBankTransferPermata").Return(_paymentRepository.PaymentResponseCollection[3], nil)
	paymentRepositoryMock.Mock.On("CreateBankTransferMandiri").Return(_paymentRepository.PaymentResponseCollection[4], nil)
	t.Run("success-bni", func(t *testing.T) {
		requestSample.PaymentMethod = "BANK_TRANSFER_BNI"
		orderSample := _orderRepository.OrderCollection[2]
		orderSample.Status = "CONFIRMED"
		orderSample.TransactionID = ""
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)
		orderOutput := orderSample
		paymentOutput := _paymentRepository.PaymentResponseCollection[0]
		orderOutput.PaymentMethod = paymentOutput.PaymentMethod
		orderOutput.TransactionID = paymentOutput.TransactionID
		orderRepositoryMock.Mock.On("Update", orderOutput).Return(orderOutput, nil)

		orderHistoryRepositoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})
		historyOutput := _orderHistoryRepository.OrderHistoryCollection[7]
		orderHistoryRepositoryMock.Mock.On("Create", int(orderSample.ID), historyOutput.Log, historyOutput.Actor).Return(historyOutput, nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryRepositoryMock,
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			paymentRepositoryMock,
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		_, err := orderService.CreatePayment(int(orderSample.ID), requestSample)
		assert.Nil(t, err)
	})
	t.Run("success-bri", func(t *testing.T) {
		requestSample.PaymentMethod = "BANK_TRANSFER_BRI"
		orderSample := _orderRepository.OrderCollection[2]
		orderSample.Status = "CONFIRMED"
		orderSample.TransactionID = ""
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)
		orderOutput := orderSample
		paymentOutput := _paymentRepository.PaymentResponseCollection[2]
		orderOutput.PaymentMethod = paymentOutput.PaymentMethod
		orderOutput.TransactionID = paymentOutput.TransactionID
		orderRepositoryMock.Mock.On("Update", orderOutput).Return(orderOutput, nil)

		orderHistoryRepositoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})
		historyOutput := _orderHistoryRepository.OrderHistoryCollection[7]
		orderHistoryRepositoryMock.Mock.On("Create", int(orderSample.ID), historyOutput.Log, historyOutput.Actor).Return(historyOutput, nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryRepositoryMock,
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			paymentRepositoryMock,
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		_, err := orderService.CreatePayment(int(orderSample.ID), requestSample)
		assert.Nil(t, err)
	})
	t.Run("success-bca", func(t *testing.T) {
		requestSample.PaymentMethod = "BANK_TRANSFER_BCA"
		orderSample := _orderRepository.OrderCollection[2]
		orderSample.Status = "CONFIRMED"
		orderSample.TransactionID = ""
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)
		orderOutput := orderSample
		paymentOutput := _paymentRepository.PaymentResponseCollection[1]
		orderOutput.PaymentMethod = paymentOutput.PaymentMethod
		orderOutput.TransactionID = paymentOutput.TransactionID
		orderRepositoryMock.Mock.On("Update", orderOutput).Return(orderOutput, nil)

		orderHistoryRepositoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})
		historyOutput := _orderHistoryRepository.OrderHistoryCollection[7]
		orderHistoryRepositoryMock.Mock.On("Create", int(orderSample.ID), historyOutput.Log, historyOutput.Actor).Return(historyOutput, nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryRepositoryMock,
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			paymentRepositoryMock,
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		_, err := orderService.CreatePayment(int(orderSample.ID), requestSample)
		assert.Nil(t, err)
	})
	t.Run("success-permata", func(t *testing.T) {
		requestSample.PaymentMethod = "BANK_TRANSFER_PERMATA"
		orderSample := _orderRepository.OrderCollection[2]
		orderSample.Status = "CONFIRMED"
		orderSample.TransactionID = ""
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)
		orderOutput := orderSample
		paymentOutput := _paymentRepository.PaymentResponseCollection[3]
		orderOutput.PaymentMethod = paymentOutput.PaymentMethod
		orderOutput.TransactionID = paymentOutput.TransactionID
		orderRepositoryMock.Mock.On("Update", orderOutput).Return(orderOutput, nil)

		orderHistoryRepositoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})
		historyOutput := _orderHistoryRepository.OrderHistoryCollection[7]
		orderHistoryRepositoryMock.Mock.On("Create", int(orderSample.ID), historyOutput.Log, historyOutput.Actor).Return(historyOutput, nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryRepositoryMock,
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			paymentRepositoryMock,
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		_, err := orderService.CreatePayment(int(orderSample.ID), requestSample)
		assert.Nil(t, err)
	})
	t.Run("success-mandiri", func(t *testing.T) {
		requestSample.PaymentMethod = "BANK_TRANSFER_MANDIRI"
		orderSample := _orderRepository.OrderCollection[2]
		orderSample.Status = "CONFIRMED"
		orderSample.TransactionID = ""
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)
		orderOutput := orderSample
		paymentOutput := _paymentRepository.PaymentResponseCollection[4]
		orderOutput.PaymentMethod = paymentOutput.PaymentMethod
		orderOutput.TransactionID = paymentOutput.TransactionID
		orderRepositoryMock.Mock.On("Update", orderOutput).Return(orderOutput, nil)

		orderHistoryRepositoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})
		historyOutput := _orderHistoryRepository.OrderHistoryCollection[7]
		orderHistoryRepositoryMock.Mock.On("Create", int(orderSample.ID), historyOutput.Log, historyOutput.Actor).Return(historyOutput, nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryRepositoryMock,
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			paymentRepositoryMock,
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		_, err := orderService.CreatePayment(int(orderSample.ID), requestSample)
		assert.Nil(t, err)
	})
	t.Run("invalid-payment-method", func(t *testing.T) {
		requestSample.PaymentMethod = "sss"
		orderSample := _orderRepository.OrderCollection[2]
		orderSample.Status = "CONFIRMED"
		orderSample.TransactionID = ""
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)
		orderHistoryRepositoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryRepositoryMock,
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			paymentRepositoryMock,
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		_, err := orderService.CreatePayment(int(orderSample.ID), requestSample)
		assert.Error(t, err)
	})
	t.Run("error-payment-repo", func(t *testing.T) {
		requestSample.PaymentMethod = "BANK_TRANSFER_BNI"
		orderSample := _orderRepository.OrderCollection[2]
		orderSample.Status = "CONFIRMED"
		orderSample.TransactionID = ""
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)
		orderOutput := orderSample
		paymentOutput := _paymentRepository.PaymentResponseCollection[0]
		orderOutput.PaymentMethod = paymentOutput.PaymentMethod
		orderOutput.TransactionID = paymentOutput.TransactionID
		orderRepositoryMock.Mock.On("Update", orderOutput).Return(orderOutput, nil)

		paymentRepositoryMock2 := _paymentRepository.NewPaymentRepositoryMock(&mock.Mock{})
		paymentRepositoryMock2.Mock.On("CreateBankTransferBNI").Return(entities.PaymentResponse{}, web.WebError{})

		orderHistoryRepositoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})
		historyOutput := _orderHistoryRepository.OrderHistoryCollection[7]
		orderHistoryRepositoryMock.Mock.On("Create", int(orderSample.ID), historyOutput.Log, historyOutput.Actor).Return(historyOutput, nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryRepositoryMock,
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			paymentRepositoryMock2,
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		_, err := orderService.CreatePayment(int(orderSample.ID), requestSample)
		assert.Error(t, err)
	})
	t.Run("validation-error", func(t *testing.T) {
		requestSample.PaymentMethod = ""
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		paymentRepositoryMock2 := _paymentRepository.NewPaymentRepositoryMock(&mock.Mock{})
		orderHistoryRepositoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryRepositoryMock,
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			paymentRepositoryMock2,
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		_, err := orderService.CreatePayment(1, requestSample)
		assert.Error(t, err)
	})
	t.Run("invalid-order", func(t *testing.T) {
		requestSample.PaymentMethod = "TRANSFER_BANK_BNI"
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", 1).Return(entities.Order{}, web.WebError{})
		paymentRepositoryMock2 := _paymentRepository.NewPaymentRepositoryMock(&mock.Mock{})
		orderHistoryRepositoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryRepositoryMock,
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			paymentRepositoryMock2,
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		_, err := orderService.CreatePayment(1, requestSample)
		assert.Error(t, err)
	})
	t.Run("invalid-order-status", func(t *testing.T) {
		requestSample.PaymentMethod = "TRANSFER_BANK_BNI"
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})

		orderSample := _orderRepository.OrderCollection[2]
		orderSample.Status = "ON_PROGRESS"
		orderRepositoryMock.Mock.On("Find", 1).Return(orderSample, nil)
		paymentRepositoryMock2 := _paymentRepository.NewPaymentRepositoryMock(&mock.Mock{})
		orderHistoryRepositoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryRepositoryMock,
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			paymentRepositoryMock2,
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		_, err := orderService.CreatePayment(1, requestSample)
		assert.Error(t, err)
	})
	t.Run("transaction-not-empty", func(t *testing.T) {
		requestSample.PaymentMethod = "TRANSFER_BANK_BNI"
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})

		orderSample := _orderRepository.OrderCollection[2]
		orderSample.Status = "CONFIRMED"
		orderSample.TransactionID = "asdasdasd"
		orderRepositoryMock.Mock.On("Find", 1).Return(orderSample, nil)
		paymentRepositoryMock2 := _paymentRepository.NewPaymentRepositoryMock(&mock.Mock{})
		orderHistoryRepositoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryRepositoryMock,
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			paymentRepositoryMock2,
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		_, err := orderService.CreatePayment(1, requestSample)
		assert.Error(t, err)
	})
	t.Run("repo-update-failed", func(t *testing.T) {
		requestSample.PaymentMethod = "BANK_TRANSFER_BRI"
		orderSample := _orderRepository.OrderCollection[2]
		orderSample.Status = "CONFIRMED"
		orderSample.TransactionID = ""
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)
		orderOutput := orderSample
		paymentOutput := _paymentRepository.PaymentResponseCollection[2]
		orderOutput.PaymentMethod = paymentOutput.PaymentMethod
		orderOutput.TransactionID = paymentOutput.TransactionID
		orderRepositoryMock.Mock.On("Update", orderOutput).Return(entities.Order{}, web.WebError{})

		orderHistoryRepositoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})
		historyOutput := _orderHistoryRepository.OrderHistoryCollection[7]
		orderHistoryRepositoryMock.Mock.On("Create", int(orderSample.ID), historyOutput.Log, historyOutput.Actor).Return(historyOutput, nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryRepositoryMock,
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			paymentRepositoryMock,
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		_, err := orderService.CreatePayment(int(orderSample.ID), requestSample)
		assert.Error(t, err)
	})
}

func TestGetPayment(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		paymentSample := _paymentRepository.PaymentResponseCollection[0]
		paymentRepositoryMock := _paymentRepository.NewPaymentRepositoryMock(&mock.Mock{})
		paymentRepositoryMock.Mock.On("GetPaymentStatus").Return(paymentSample, nil)

		orderSample := _orderRepository.OrderCollection[2]
		orderSample.PaymentMethod = paymentSample.PaymentMethod
		orderSample.TransactionID = paymentSample.TransactionID
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			_orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{}),
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			paymentRepositoryMock,
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		_, err := orderService.GetPayment(int(orderSample.ID))
		assert.Nil(t, err)
	})
	t.Run("invalid-order", func(t *testing.T) {
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", 1).Return(entities.Order{}, web.WebError{})
		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			_orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{}),
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		_, err := orderService.GetPayment(1)
		assert.Error(t, err)
	})
	t.Run("repo-failed", func(t *testing.T) {
		paymentSample := _paymentRepository.PaymentResponseCollection[0]
		paymentRepositoryMock := _paymentRepository.NewPaymentRepositoryMock(&mock.Mock{})
		paymentRepositoryMock.Mock.On("GetPaymentStatus").Return(entities.PaymentResponse{}, web.WebError{})

		orderSample := _orderRepository.OrderCollection[2]
		orderSample.PaymentMethod = paymentSample.PaymentMethod
		orderSample.TransactionID = paymentSample.TransactionID
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			_orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{}),
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			paymentRepositoryMock,
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		_, err := orderService.GetPayment(int(orderSample.ID))
		assert.Error(t, err)
	})
}

func TestCancelPayment(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		orderSample := _orderRepository.OrderCollection[1]
		orderSample.Status = "CONFIRMED"
		
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)

		paymentRepositoryMock := _paymentRepository.NewPaymentRepositoryMock(&mock.Mock{})
		paymentRepositoryMock.Mock.On("CancelPayment").Return(nil)

		orderUpdateOutput := orderSample
		orderUpdateOutput.TransactionID = ""
		orderUpdateOutput.PaymentMethod = ""
		orderRepositoryMock.Mock.On("Update", orderUpdateOutput).Return(orderUpdateOutput, nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			_orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{}),
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			paymentRepositoryMock,
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		err := orderService.CancelPayment(int(orderSample.ID))
		assert.Nil(t, err)
	})
	t.Run("invalid-order", func(t *testing.T) {
		orderSample := _orderRepository.OrderCollection[1]
		orderSample.Status = "CONFIRMED"
		
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(entities.Order{}, web.WebError{})
		paymentRepositoryMock := _paymentRepository.NewPaymentRepositoryMock(&mock.Mock{})

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			_orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{}),
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			paymentRepositoryMock,
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		err := orderService.CancelPayment(int(orderSample.ID))
		assert.Error(t, err)
	})
	t.Run("invalid-order-status", func(t *testing.T) {
		orderSample := _orderRepository.OrderCollection[1]
		orderSample.Status = "ON_PROCESS"
		
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)
		paymentRepositoryMock := _paymentRepository.NewPaymentRepositoryMock(&mock.Mock{})

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			_orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{}),
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			paymentRepositoryMock,
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		err := orderService.CancelPayment(int(orderSample.ID))
		assert.Error(t, err)
	})
	t.Run("transaction-empty", func(t *testing.T) {
		orderSample := _orderRepository.OrderCollection[1]
		orderSample.Status = "CONFIRMED"
		orderSample.TransactionID = ""
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)

		paymentRepositoryMock := _paymentRepository.NewPaymentRepositoryMock(&mock.Mock{})

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			_orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{}),
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			paymentRepositoryMock,
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		err := orderService.CancelPayment(int(orderSample.ID))
		assert.Error(t, err)
	})
	t.Run("cancel-failed", func(t *testing.T) {
		orderSample := _orderRepository.OrderCollection[1]
		orderSample.Status = "CONFIRMED"
		
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)

		paymentRepositoryMock := _paymentRepository.NewPaymentRepositoryMock(&mock.Mock{})
		paymentRepositoryMock.Mock.On("CancelPayment").Return(web.WebError{})

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			_orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{}),
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			paymentRepositoryMock,
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		err := orderService.CancelPayment(int(orderSample.ID))
		assert.Error(t, err)
	})
	t.Run("success", func(t *testing.T) {
		orderSample := _orderRepository.OrderCollection[1]
		orderSample.Status = "CONFIRMED"
		
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)

		paymentRepositoryMock := _paymentRepository.NewPaymentRepositoryMock(&mock.Mock{})
		paymentRepositoryMock.Mock.On("CancelPayment").Return(nil)

		orderUpdateOutput := orderSample
		orderUpdateOutput.TransactionID = ""
		orderUpdateOutput.PaymentMethod = ""
		orderRepositoryMock.Mock.On("Update", orderUpdateOutput).Return(entities.Order{}, web.WebError{})

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			_orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{}),
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			paymentRepositoryMock,
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		err := orderService.CancelPayment(int(orderSample.ID))
		assert.Error(t, err)
	})
}

func TestFindAllHistory(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		orderHistoryRepositoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})
		orderHistoryRepositoryMock.Mock.On("FindAll", 1, []map[string]interface{}{}).Return(_orderHistoryRepository.OrderHistoryCollection, nil)

		orderService := _orderService.NewOrderService(
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
			orderHistoryRepositoryMock,
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		actual, err := orderService.FindAllHistory(1, []map[string]interface{}{})
		expected := []entities.OrderHistoryResponse{}
        copier.Copy(&expected, &_orderHistoryRepository.OrderHistoryCollection)

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("error", func(t *testing.T) {
		orderHistoryRepositoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})
		orderHistoryRepositoryMock.Mock.On("FindAll", 1, []map[string]interface{}{}).Return([]entities.OrderHistory{}, web.WebError{})

		orderService := _orderService.NewOrderService(
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
			orderHistoryRepositoryMock,
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		_, err := orderService.FindAllHistory(1, []map[string]interface{}{})
		assert.Error(t, err)
	})
}

func TestPaymentWebhook(t *testing.T) {
	t.Run("success-settlement", func(t *testing.T) {
		orderSample := _orderRepository.OrderCollection[1]
		orderSample.Status = "CONFIRMED"
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)
		
		orderUpdateOutput := orderSample
		orderUpdateOutput.Status = "MANIFESTED"
		orderRepositoryMock.Mock.On("Update", orderUpdateOutput).Return(orderUpdateOutput, nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			_orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{}),
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		err := orderService.PaymentWebhook(int(orderSample.ID), "settlement")
		assert.Nil(t, err)
	})
	t.Run("invalid-order", func(t *testing.T) {
		orderSample := _orderRepository.OrderCollection[1]
		orderSample.Status = "CONFIRMED"
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(entities.Order{}, web.WebError{})
	
		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			_orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{}),
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		err := orderService.PaymentWebhook(int(orderSample.ID), "settlement")
		assert.Error(t, err)
	})
	t.Run("deny", func(t *testing.T) {
		orderSample := _orderRepository.OrderCollection[1]
		orderSample.Status = "CONFIRMED"
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)
		
		orderUpdateOutput := orderSample
		orderUpdateOutput.PaymentMethod = ""
		orderUpdateOutput.TransactionID = ""
		orderRepositoryMock.Mock.On("Update", orderUpdateOutput).Return(orderUpdateOutput, nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			_orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{}),
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		err := orderService.PaymentWebhook(int(orderSample.ID), "deny")
		assert.Nil(t, err)
	})
	t.Run("failed-repo", func(t *testing.T) {
		orderSample := _orderRepository.OrderCollection[1]
		orderSample.Status = "CONFIRMED"
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)
		
		orderUpdateOutput := orderSample
		orderUpdateOutput.Status = "MANIFESTED"
		orderRepositoryMock.Mock.On("Update", orderUpdateOutput).Return(entities.Order{}, web.WebError{})

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			_orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{}),
			_userRepository.NewUserRepositoryMock(&mock.Mock{}),
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		err := orderService.PaymentWebhook(int(orderSample.ID), "settlement")
		assert.Error(t, err)
	})
}

func TestTakeOrder(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		orderSample := _orderRepository.OrderCollection[1]
		driverSample := _userRepository.DriverCollection[1]

		orderSample.Status = "MANIFESTED"
		orderSample.DriverID = null.IntFromPtr(nil)
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)

		driverSample.Status = "IDLE"
		driverSample.TruckTypeID = orderSample.TruckTypeID
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindByDriver").Return(driverSample, nil)

		orderUpdateOutput := orderSample
		driverID := int64(driverSample.ID)
        orderUpdateOutput.DriverID = null.IntFromPtr(&driverID)
        orderUpdateOutput.Status = "ON_PROCESS"
		orderRepositoryMock.Mock.On("Update", orderUpdateOutput).Return(orderUpdateOutput, nil)

		driverUpdateOutput := driverSample
		driverUpdateOutput.Status = "IDLE"
		userRepositoryMock.Mock.On("UpdateDriver").Return(driverUpdateOutput, nil)

		orderHistorySample := _orderHistoryRepository.OrderHistoryCollection[5]
		orderHistorySample.Log = "Order diambil oleh " + driverSample.User.Name
		orderHistorySample.Actor = "driver"
		orderHistoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})
		orderHistoryMock.Mock.On("Create", int(orderSample.ID) ,orderHistorySample.Log, orderHistorySample.Actor).Return(orderHistorySample, nil)

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryMock,
			userRepositoryMock,
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		err := orderService.TakeOrder(int(orderSample.ID), int(driverSample.UserID))
		assert.Nil(t, err)
	})
	t.Run("invalid-order", func(t *testing.T) {
		orderSample := _orderRepository.OrderCollection[1]
		driverSample := _userRepository.DriverCollection[1]

		orderSample.Status = "MANIFESTED"
		orderSample.DriverID = null.IntFromPtr(nil)
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(entities.Order{}, web.WebError{})

		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		orderHistoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryMock,
			userRepositoryMock,
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		err := orderService.TakeOrder(int(orderSample.ID), int(driverSample.UserID))
		assert.Error(t, err)
	})
	t.Run("driver-not-found", func(t *testing.T) {
		orderSample := _orderRepository.OrderCollection[1]
		driverSample := _userRepository.DriverCollection[1]

		orderSample.Status = "MANIFESTED"
		orderSample.DriverID = null.IntFromPtr(nil)
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)

		driverSample.Status = "IDLE"
		driverSample.TruckTypeID = orderSample.TruckTypeID
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindByDriver").Return(entities.Driver{}, web.WebError{})
		orderHistoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})

		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryMock,
			userRepositoryMock,
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		err := orderService.TakeOrder(int(orderSample.ID), int(driverSample.UserID))
		assert.Error(t, err)
	})
	t.Run("taken-by-driver", func(t *testing.T) {
		orderSample := _orderRepository.OrderCollection[1]
		driverSample := _userRepository.DriverCollection[1]

		orderSample.Status = "MANIFESTED"
		orderSample.DriverID = null.IntFrom(999)
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)

		driverSample.Status = "IDLE"
		driverSample.TruckTypeID = 999
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindByDriver").Return(driverSample, nil)
		orderHistoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})
		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryMock,
			userRepositoryMock,
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		err := orderService.TakeOrder(int(orderSample.ID), int(driverSample.UserID))
		assert.Error(t, err)
	})
	t.Run("invalid-truck", func(t *testing.T) {
		orderSample := _orderRepository.OrderCollection[1]
		driverSample := _userRepository.DriverCollection[1]

		orderSample.Status = "MANIFESTED"
		orderSample.DriverID = null.IntFromPtr(nil)
		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("Find", int(orderSample.ID)).Return(orderSample, nil)

		driverSample.Status = "IDLE"
		driverSample.TruckTypeID = 999
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindByDriver").Return(driverSample, nil)
		orderHistoryMock := _orderHistoryRepository.NewOrderHistoryMock(&mock.Mock{})
		orderService := _orderService.NewOrderService(
			orderRepositoryMock,
			orderHistoryMock,
			userRepositoryMock,
			_paymentRepository.NewPaymentRepositoryMock(&mock.Mock{}),
			_distanceMatrixRepository.NewDistanceMatrixRepositoryMock(&mock.Mock{}),
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
		)
		err := orderService.TakeOrder(int(orderSample.ID), int(driverSample.UserID))
		assert.Error(t, err)
	})
}