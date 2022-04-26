package truck_type

import "bringeee-capstone/entities"


type TruckTypeServiceInterface interface {

	/*
	 * Find All
	 * -------------------------------
	 * Mengambil data truckType berdasarkan filters dan sorts
	 *
	 * @var limit 			batas limit hasil query
	 * @var offset 			offset hasil query
	 * @var filters			query untuk penyaringan data, { field, operator, value }
	 * @var sorts			pengurutan data, { field, value[bool] }
	 * @return truckType	list truckType dalam bentuk entity domain
	 * @return error		error
	 */
	FindAll(limit int, offset int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.TruckTypeResponse, error)
}