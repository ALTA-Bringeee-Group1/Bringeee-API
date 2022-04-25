package region

import (
	"bringeee-capstone/entities"

	"gorm.io/gorm"
)

type RegionRepository struct {
	db *gorm.DB
}

func NewRegionRepository(db *gorm.DB) *RegionRepository {
	return &RegionRepository{
		db: db,
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
func (repository RegionRepository) findAllProvince(sort []map[string]interface{}) ([]entities.Province, error) {
	panic("implement me")
}
/*
 * Find All City
 * -------------------------------
 * Mengambil semua data kota berdasarkan provinsi
 *
 * @var sort		sort data, { field, sort[bool] }
 * @return City		list kota dalam entity domain
 */
func (repository RegionRepository) findAllCity(provinceID int, sort []map[string]interface{}) ([]entities.City, error) {
	panic("implement me")
}
/*
 * Find All District
 * -------------------------------
 * Mengambil semua data kecamatan berdasarkan kota
 *
 * @var sort		sort data, { field, sort[bool] }
 * @return District	list kecamatan dalam entity domain 
 */
func (repository RegionRepository) findAllDistrict(cityID int, sort []map[string]interface{}) ([]entities.District, error) {
	panic("implement me")
}