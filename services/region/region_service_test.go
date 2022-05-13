package region_test

import (
	"bringeee-capstone/entities"
	"bringeee-capstone/entities/web"
	_regionRepository "bringeee-capstone/repositories/region"
	_regionService "bringeee-capstone/services/region"
	"testing"

	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFindAllProvince(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		provinceSamples := _regionRepository.ProvinceCollection
		regionRepositoryMock := _regionRepository.NewRegionRepositoryMock(&mock.Mock{})
		regionRepositoryMock.Mock.On("FindAllProvince").Return(provinceSamples, nil)

		regionService := _regionService.NewRegionService(regionRepositoryMock)
		actual, err := regionService.FindAllProvince([]map[string]interface{}{})

		expected := []entities.ProvinceResponse{}
		copier.Copy(&expected, &provinceSamples)

		assert.Nil(t, err)
		assert.Equal(t, actual, expected)
	})
	t.Run("error", func(t *testing.T) {
		regionRepositoryMock := _regionRepository.NewRegionRepositoryMock(&mock.Mock{})
		regionRepositoryMock.Mock.On("FindAllProvince").Return([]entities.Province{}, web.WebError{})

		regionService := _regionService.NewRegionService(regionRepositoryMock)
		actual, err := regionService.FindAllProvince([]map[string]interface{}{})

		assert.Error(t, err)
		assert.Equal(t, []entities.ProvinceResponse{}, actual)
	})
}

func TestFindAllCity(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		citySamples := _regionRepository.CityCollection
		regionRepositoryMock := _regionRepository.NewRegionRepositoryMock(&mock.Mock{})
		regionRepositoryMock.Mock.On("FindAllCity").Return(citySamples, nil)

		regionService := _regionService.NewRegionService(regionRepositoryMock)
		actual, err := regionService.FindAllCity(1, []map[string]interface{}{})

		expected := []entities.CityResponse{}
		copier.Copy(&expected, &citySamples)

		assert.Nil(t, err)
		assert.Equal(t, actual, expected)
	})
	t.Run("error", func(t *testing.T) {
		regionRepositoryMock := _regionRepository.NewRegionRepositoryMock(&mock.Mock{})
		regionRepositoryMock.Mock.On("FindAllCity").Return([]entities.City{}, web.WebError{})

		regionService := _regionService.NewRegionService(regionRepositoryMock)
		actual, err := regionService.FindAllCity(1, []map[string]interface{}{})

		assert.Error(t, err)
		assert.Equal(t, []entities.CityResponse{}, actual)
	})
}

func TestFindAllDistrict(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		districtSamples := _regionRepository.DistrictCollection
		regionRepositoryMock := _regionRepository.NewRegionRepositoryMock(&mock.Mock{})
		regionRepositoryMock.Mock.On("FindAllDistrict").Return(districtSamples, nil)
		
		citySample := _regionRepository.CityCollection[0]
		regionRepositoryMock.Mock.On("FindCity").Return(citySample, nil)

		regionService := _regionService.NewRegionService(regionRepositoryMock)
		actual, err := regionService.FindAllDistrict(int(citySample.CityID), int(citySample.ProvID), []map[string]interface{}{})

		expected := []entities.DistrictResponse{}
		copier.Copy(&expected, &districtSamples)

		assert.Nil(t, err)
		assert.Equal(t, actual, expected)
	})
	t.Run("invalid-city", func(t *testing.T) {
		districtSamples := _regionRepository.DistrictCollection
		regionRepositoryMock := _regionRepository.NewRegionRepositoryMock(&mock.Mock{})
		regionRepositoryMock.Mock.On("FindAllDistrict").Return(districtSamples, nil)
		
		citySample := _regionRepository.CityCollection[0]
		regionRepositoryMock.Mock.On("FindCity").Return(entities.City{}, web.WebError{})

		regionService := _regionService.NewRegionService(regionRepositoryMock)
		actual, err := regionService.FindAllDistrict(int(citySample.CityID), int(citySample.ProvID), []map[string]interface{}{})

		assert.Error(t, err)
		assert.Equal(t, []entities.DistrictResponse{}, actual)
	})
	t.Run("invalid-province", func(t *testing.T) {
		districtSamples := _regionRepository.DistrictCollection
		regionRepositoryMock := _regionRepository.NewRegionRepositoryMock(&mock.Mock{})
		regionRepositoryMock.Mock.On("FindAllDistrict").Return(districtSamples, nil)
		
		citySample := _regionRepository.CityCollection[0]
		regionRepositoryMock.Mock.On("FindCity").Return(citySample, nil)

		regionService := _regionService.NewRegionService(regionRepositoryMock)
		actual, err := regionService.FindAllDistrict(int(citySample.CityID), 999, []map[string]interface{}{})

		expected := []entities.DistrictResponse{}
		copier.Copy(&expected, &districtSamples)

		assert.Error(t, err)
		assert.Equal(t, []entities.DistrictResponse{}, actual)
	})
	t.Run("invalid-province", func(t *testing.T) {
		districtSamples := _regionRepository.DistrictCollection
		regionRepositoryMock := _regionRepository.NewRegionRepositoryMock(&mock.Mock{})
		regionRepositoryMock.Mock.On("FindAllDistrict").Return([]entities.District{}, web.WebError{})
		
		citySample := _regionRepository.CityCollection[0]
		regionRepositoryMock.Mock.On("FindCity").Return(citySample, nil)

		regionService := _regionService.NewRegionService(regionRepositoryMock)
		actual, err := regionService.FindAllDistrict(int(citySample.CityID), int(citySample.ProvID), []map[string]interface{}{})

		expected := []entities.DistrictResponse{}
		copier.Copy(&expected, &districtSamples)

		assert.Error(t, err)
		assert.Equal(t, []entities.DistrictResponse{}, actual)
	})
}