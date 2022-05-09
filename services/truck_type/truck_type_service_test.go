package truck_type_test

import (
	"bringeee-capstone/entities"
	truckTypeRepository "bringeee-capstone/repositories/truck_type"
	truckTypeService "bringeee-capstone/services/truck_type"
	"testing"
	"time"

	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Truck Repository Mock
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
		TruckType: "Truck - C", 
		MaxVolume: 8000000,
		MaxWeight: 2000,
		PricePerDistance: 3000,
	},
}

func TestFindAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		truckTypeRepositoryMock := truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{})
		truckTypeRepositoryMock.Mock.On("FindAll").Return(truckTypeCollection, "")

		service := truckTypeService.NewTruckTypeService(truckTypeRepositoryMock)
		data, err := service.FindAll(0, 0, []map[string]string{}, []map[string]interface{}{})

		// Konversi expected data ke response
		truckTypesRes := []entities.TruckTypeResponse{}
		copier.Copy(&truckTypesRes, &truckTypeCollection)

		assert.Nil(t, err, "Service FindAll error must be nil")
		assert.Equal(t, truckTypesRes, data, "Truck Type data collection is not returned")
	})

	t.Run("with-filter", func(t *testing.T) {
		truckTypeRepositoryMock := truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{})
		truckTypeRepositoryMock.Mock.On("FindAll").Return(truckTypeCollection, "")

		service := truckTypeService.NewTruckTypeService(truckTypeRepositoryMock)
		data, err := service.FindAll(0, 0, []map[string]string{
			{ "field" : "truck_type", "operator" : "LIKE", "value" : "Pickup" },
		}, []map[string]interface{}{})

		// Konversi expected data ke response
		truckTypesRes := []entities.TruckTypeResponse{}
		copier.Copy(&truckTypesRes, &[]entities.TruckType{truckTypeCollection[0], truckTypeCollection[1]})

		assert.Nil(t, err, "Service FindAll error must be nil")
		assert.Equal(t, truckTypesRes, data, "Truck Type data collection is not returned")
	})

	t.Run("with-error", func(t *testing.T) {
		truckTypeRepositoryMock := truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{})
		truckTypeRepositoryMock.Mock.On("FindAll").Return(truckTypeCollection, "SERVER_ERROR")

		service := truckTypeService.NewTruckTypeService(truckTypeRepositoryMock)
		data, err := service.FindAll(0, 0, []map[string]string{}, []map[string]interface{}{})
		
		assert.Error(t, err, "error data is not returned")
		assert.Equal(t, []entities.TruckTypeResponse{}, data)
	})
}


func TestCountTruck(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		truckTypeRepositoryMock := truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{})
		truckTypeRepositoryMock.Mock.On("CountAll").Return(truckTypeCollection, "")

		service := truckTypeService.NewTruckTypeService(truckTypeRepositoryMock)
		count, err := service.CountTruck([]map[string]string{})

		assert.Nil(t, err, "error is not nil")
		assert.Equal(t, len(truckTypeCollection), count, "Truck Type data collection is not returned")
	})

	t.Run("with-filter", func(t *testing.T) {
		truckTypeRepositoryMock := truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{})
		truckTypeRepositoryMock.Mock.On("CountAll").Return(truckTypeCollection, "")

		service := truckTypeService.NewTruckTypeService(truckTypeRepositoryMock)
		count, err := service.CountTruck([]map[string]string{
			{ "field" : "truck_type", "operator" : "LIKE", "value" : "Pickup" },
		})

		assert.Nil(t, err, "error is not nil")
		assert.Equal(t, 2, count, "Truck Type data collection is not returned")
	})

	t.Run("with-error", func(t *testing.T) {
		truckTypeRepositoryMock := truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{})
		truckTypeRepositoryMock.Mock.On("CountAll").Return(truckTypeCollection, "SERVER_ERROR")

		service := truckTypeService.NewTruckTypeService(truckTypeRepositoryMock)
		count, err := service.CountTruck([]map[string]string{})
		
		assert.Error(t, err, "error data is not returned")
		assert.Equal(t, 0, count)
	})
}