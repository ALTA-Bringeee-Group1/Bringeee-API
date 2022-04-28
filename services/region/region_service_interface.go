package region

import "bringeee-capstone/entities"

type RegionServiceInterface interface {
	/*
	 * Find All Province
	 * -------------------------------
	 * Mengambil semua data provinsi
	 *
	 * @var sort		sort data, { field, sort[bool] }
	 * @return Province	list provinsi dalam entity response
	 */
	FindAllProvince(sort []map[string]interface{}) ([]entities.ProvinceResponse, error)

	/*
	 * Find All City
	 * -------------------------------
	 * Mengambil semua data kota berdasarkan provinsi
	 *
	 * @var sort		sort data, { field, sort[bool] }
	 * @return City		list kota dalam entity response
	 */
	FindAllCity(provinceID int, sort []map[string]interface{}) ([]entities.CityResponse, error)

	/*
	 * Find All District
	 * -------------------------------
	 * Mengambil semua data kecamatan berdasarkan kota
	 *
	 * @var sort		sort data, { field, sort[bool] }
	 * @return District	list kecamatan dalam entity response 
	 */
	FindAllDistrict(cityID int, provinceID int, sort []map[string]interface{}) ([]entities.DistrictResponse, error)
}