package truckType_test

import (
	"bringeee-capstone/entities"
	"bringeee-capstone/entities/web"
	truckTypeRepository "bringeee-capstone/repositories/truck_type"
	truckTypeService "bringeee-capstone/services/truck_type"
	"testing"

	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)


func TestFindAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		truckSamples := truckTypeRepository.TruckTypeCollection
		truckTypeRepositoryMock := truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{})
		truckTypeRepositoryMock.Mock.On(
			"FindAll",
			0, 0,
			[]map[string]string{}, 
			[]map[string]interface{}{},
		).Return(truckSamples, nil)

		service := truckTypeService.NewTruckTypeService(truckTypeRepositoryMock)
		data, err := service.FindAll(0, 0, []map[string]string{}, []map[string]interface{}{})

		// Konversi expected data ke response
		truckTypesRes := []entities.TruckTypeResponse{}
		copier.Copy(&truckTypesRes, &truckSamples)

		assert.Nil(t, err)
		assert.Equal(t, truckTypesRes, data)
	})
	t.Run("success", func(t *testing.T) {
		truckTypeRepositoryMock := truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{})
		truckTypeRepositoryMock.Mock.On(
			"FindAll",
			0, 0,
			[]map[string]string{}, 
			[]map[string]interface{}{},
		).Return([]entities.TruckType{}, web.WebError{})

		service := truckTypeService.NewTruckTypeService(truckTypeRepositoryMock)
		data, err := service.FindAll(0, 0, []map[string]string{}, []map[string]interface{}{})

		// Konversi expected data ke response
		truckTypesRes := []entities.TruckTypeResponse{}

		assert.Error(t, err)
		assert.Equal(t, truckTypesRes, data)
	})
}


func TestCountTruck(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		truckSamples := truckTypeRepository.TruckTypeCollection
		truckTypeRepositoryMock := truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{})
		truckTypeRepositoryMock.Mock.On("CountAll").Return(len(truckSamples), nil)

		service := truckTypeService.NewTruckTypeService(truckTypeRepositoryMock)
		count, err := service.CountTruck([]map[string]string{})

		assert.Nil(t, err, "error is not nil")
		assert.Equal(t, len(truckSamples), count, "Truck Type data collection is not returned")
	})

	t.Run("with-filter", func(t *testing.T) {
		truckSamples := truckTypeRepository.TruckTypeCollection
		truckTypeRepositoryMock := truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{})
		truckTypeRepositoryMock.Mock.On("CountAll").Return(len(truckSamples), nil)

		service := truckTypeService.NewTruckTypeService(truckTypeRepositoryMock)
		count, err := service.CountTruck([]map[string]string{
			{ "field" : "truck_type", "operator" : "LIKE", "value" : "Pickup" },
		})

		assert.Nil(t, err, "error is not nil")
		assert.Equal(t, 3, count, "Truck Type data collection is not returned")
	})

	t.Run("with-error", func(t *testing.T) {
		truckSamples := truckTypeRepository.TruckTypeCollection
		truckTypeRepositoryMock := truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{})
		truckTypeRepositoryMock.Mock.On("CountAll").Return(len(truckSamples), web.WebError{})

		service := truckTypeService.NewTruckTypeService(truckTypeRepositoryMock)
		count, err := service.CountTruck([]map[string]string{})
		
		assert.Error(t, err, "error data is not returned")
		assert.Equal(t, 0, count)
	})
}