package region

import "bringeee-capstone/entities"

type RegionRepositoryInterface interface {
	
	/*
	 * Find All Province
	 * -------------------------------
	 * Mengambil semua data provinsi
	 *
	 * @var sort		sort data, { field, sort[bool] }
	 * @return Province	list provinsi dalam entity domain
	 */
	FindAllProvince(sort []map[string]interface{}) ([]entities.Province, error)

	/*
	 * Find All City
	 * -------------------------------
	 * Mengambil semua data kota berdasarkan provinsi
	 *
	 * @var sort		sort data, { field, sort[bool] }
	 * @return City		list kota dalam entity domain
	 */
	FindAllCity(provinceID int, sort []map[string]interface{}) ([]entities.City, error)

	/*
	 * Find All District
	 * -------------------------------
	 * Mengambil semua data kecamatan berdasarkan kota
	 *
	 * @var sort		sort data, { field, sort[bool] }
	 * @return District	list kecamatan dalam entity domain 
	 */
	FindAllDistrict(cityID int, sort []map[string]interface{}) ([]entities.District, error)
}