package region

import (
	"bringeee-capstone/entities"
	"bringeee-capstone/entities/web"
	"bringeee-capstone/utils"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
func (repository RegionRepository) FindAllProvince(sorts []map[string]interface{}) ([]entities.Province, error) {
	provinces := []entities.Province{}
	builder := repository.db

	// OrderBy Filters
	for _, sort := range sorts {
		fmt.Println(utils.JsonEncode(sort))
		builder.Order(clause.OrderByColumn{Column: clause.Column{Name: sort["field"].(string)}, Desc: sort["desc"].(bool)})
	}

	// Operation
	tx := builder.Find(&provinces)
	if tx.Error != nil {
		return []entities.Province{}, web.WebError{Code: 500, DevelopmentMessage: "Server data error", ProductionMessage: tx.Error.Error()}
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
func (repository RegionRepository) FindAllCity(provinceID int, sorts []map[string]interface{}) ([]entities.City, error) {
	city := []entities.City{}
	builder := repository.db.Model(&entities.City{})

	builder.Where("prov_id = ?", provinceID)

	// OrderBy Filters
	for _, sort := range sorts {
		builder.Order(clause.OrderByColumn{Column: clause.Column{Name: sort["field"].(string)}, Desc: sort["desc"].(bool)})
	}

	// Operation
	tx := builder.Find(&city)
	if tx.Error != nil {
		return []entities.City{}, web.WebError{Code: 500, DevelopmentMessage: "Server data error", ProductionMessage: tx.Error.Error()}
	}
	return city, nil
}

/*
 * Find All District
 * -------------------------------
 * Mengambil semua data kecamatan berdasarkan kota
 *
 * @var sort		sort data, { field, sort[bool] }
 * @return District	list kecamatan dalam entity domain
 */
func (repository RegionRepository) FindAllDistrict(cityID int, sorts []map[string]interface{}) ([]entities.District, error) {
	districts := []entities.District{}
	builder := repository.db.Model(&entities.District{})

	builder.Where("city_id = ?", cityID)

	// OrderBy Filters
	for _, sort := range sorts {
		builder.Order(clause.OrderByColumn{Column: clause.Column{Name: sort["field"].(string)}, Desc: sort["desc"].(bool)})
	}

	// Operation
	tx := builder.Find(&districts)
	if tx.Error != nil {
		return []entities.District{}, web.WebError{Code: 500, DevelopmentMessage: "Server data error", ProductionMessage: tx.Error.Error()}
	}
	return districts, nil
}
