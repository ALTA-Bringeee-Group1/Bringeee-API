package region

import (
	"bringeee-capstone/entities"
	"bringeee-capstone/entities/web"
	regionRepository "bringeee-capstone/repositories/region"

	"github.com/jinzhu/copier"
)

type RegionService struct {
	regionRepository regionRepository.RegionRepositoryInterface
}

func NewRegionService(repository regionRepository.RegionRepositoryInterface) *RegionService {
	return &RegionService{
		regionRepository: repository,
	}
}

/*
 * Find All Province
 * -------------------------------
 * Mengambil semua data provinsi
 *
 * @var sort		sort data, { field, sort[bool] }
 * @return Province	list provinsi dalam entity domain
 */
func (service RegionService) FindAllProvince(sort []map[string]interface{}) ([]entities.ProvinceResponse, error) {
	provinces, err := service.regionRepository.FindAllProvince(sort)
	if err != nil {
		return []entities.ProvinceResponse{}, err
	} 
	provRes := []entities.ProvinceResponse{}
	copier.Copy(&provRes, &provinces)
	return provRes, nil
}
/*
 * Find All City
 * -------------------------------
 * Mengambil semua data kota berdasarkan provinsi
 *
 * @var sort		sort data, { field, sort[bool] }
 * @return City		list kota dalam entity domain
 */
func (service RegionService) FindAllCity(provinceID int, sort []map[string]interface{}) ([]entities.CityResponse, error) {
	cities, err := service.regionRepository.FindAllCity(provinceID, sort)
	if err != nil {
		return []entities.CityResponse{}, err
	}
	citiesRes := []entities.CityResponse{}
	copier.Copy(&citiesRes, &cities)
	return citiesRes, nil
}
/*
 * Find All District
 * -------------------------------
 * Mengambil semua data kecamatan berdasarkan kota
 *
 * @var sort		sort data, { field, sort[bool] }
 * @return District	list kecamatan dalam entity domain 
 */
func (service RegionService) FindAllDistrict(cityID int, provinceID int, sort []map[string]interface{}) ([]entities.DistrictResponse, error) {
	city, err := service.regionRepository.FindCity(cityID)
	if err != nil {
		return []entities.DistrictResponse{}, err
	}
	if city.ProvID != uint(provinceID) {
		return []entities.DistrictResponse{}, web.WebError{Code: 400, Message: "invalid city & province param combination"}
	}
	districts, err := service.regionRepository.FindAllDistrict(cityID, sort)
	if err != nil {
		return []entities.DistrictResponse{}, err
	} 
	districtsRes := []entities.DistrictResponse{}
	copier.Copy(&districtsRes, &districts)
	return districtsRes, nil
}