package region

import (
	"bringeee-capstone/entities"
	regionRepository "bringeee-capstone/repositories/region"
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
func (service RegionService) FindAllProvince(sort []map[string]interface{}) ([]entities.Province, error) {
	provinces, err := service.regionRepository.FindAllProvince(sort)
	if err != nil {
		return []entities.Province{}, err
	} 
	return provinces, nil
}
/*
 * Find All City
 * -------------------------------
 * Mengambil semua data kota berdasarkan provinsi
 *
 * @var sort		sort data, { field, sort[bool] }
 * @return City		list kota dalam entity domain
 */
func (service RegionService) FindAllCity(provinceID int, sort []map[string]interface{}) ([]entities.City, error) {
	cities, err := service.regionRepository.FindAllCity(provinceID, sort)
	if err != nil {
		return []entities.City{}, err
	} 
	return cities, nil
}
/*
 * Find All District
 * -------------------------------
 * Mengambil semua data kecamatan berdasarkan kota
 *
 * @var sort		sort data, { field, sort[bool] }
 * @return District	list kecamatan dalam entity domain 
 */
func (service RegionService) FindAllDistrict(cityID int, sort []map[string]interface{}) ([]entities.District, error) {
	districts, err := service.regionRepository.FindAllDistrict(cityID, sort)
	if err != nil {
		return []entities.District{}, err
	} 
	return districts, nil
}