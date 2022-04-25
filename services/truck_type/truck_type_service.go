package truck_type

import (
	"bringeee-capstone/entities"
	truckTypeRepository "bringeee-capstone/repositories/truck_type"

	"github.com/jinzhu/copier"
)

type TruckTypeService struct {
	truckTypeRepo truckTypeRepository.TruckTypeRepositoryInterface
}

func NewTruckTypeService(repo truckTypeRepository.TruckTypeRepositoryInterface) *TruckTypeService {
	return &TruckTypeService{
		truckTypeRepo: repo,
	}
}


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
func (service TruckTypeService) FindAll(limit int, page int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.TruckTypeResponse, error) {
	offset := (page - 1) * limit

	// Repository action find all truckType
	truckTypes, err := service.truckTypeRepo.FindAll(limit, offset, filters, sorts)
	if err != nil {
		return []entities.TruckTypeResponse{}, err
	}

	// Konversi ke truckType response
	truckTypesRes := []entities.TruckTypeResponse{}
	copier.Copy(&truckTypesRes, &truckTypes)
	
	return truckTypesRes, nil
}