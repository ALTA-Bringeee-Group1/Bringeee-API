package region

import (
	"bringeee-capstone/entities"

	"github.com/stretchr/testify/mock"
)

type RegionRepositoryMock struct {
	Mock *mock.Mock
}

func NewRegionRepositoryMock(mock *mock.Mock) *RegionRepositoryMock {
	return &RegionRepositoryMock{
		Mock: mock,
	}
}

var ProvinceCollection = []entities.Province{
	{
		ProvID: 1,
		ProvName: "Banten",
		LocationID: "",
		Status: 1,
	},
	{
		ProvID: 2,
		ProvName: "Jawa Barat",
		LocationID: "",
		Status: 1,
	},
	{
		ProvID: 3,
		ProvName: "Jawa Tengah",
		LocationID: "",
		Status: 1,
	},
	{
		ProvID: 4,
		ProvName: "Jawa Timur",
		LocationID: "",
		Status: 1,
	},
	{
		ProvID: 5,
		ProvName: "Sulawesi Tengah",
		LocationID: "",
		Status: 1,
	},
	{
		ProvID: 6,
		ProvName: "Sulawesi Tenggara",
		LocationID: "",
		Status: 1,
	},
	{
		ProvID: 7,
		ProvName: "Kalimantan Timur",
		LocationID: "",
		Status: 1,
	},
}
var CityCollection = []entities.City{
	{
		CityID: 1,
		CityName: "Malang",
		ProvID: 4,
	},
	{
		CityID: 2,
		CityName: "Semarang",
		ProvID: 3,
	},
}
var DistrictCollection = []entities.District{
	{
		DisID: 1,
		DisName: "Ampelgading",
		CityID: 1,
	},
	{
		DisID: 1,
		DisName: "Blimbing",
		CityID: 2,
	},
	{
		DisID: 1,
		DisName: "Klojen",
		CityID: 3,
	},
	{
		DisID: 5,
		DisName: "Lowokwaru",
		CityID: 4,
	},
	{
		DisID: 6,
		DisName: "Sukun",
		CityID: 1,
	},
	
}

/*
 * Find All Province
 * -------------------------------
 * Mengambil semua data provinsi
 *
 * @var sort		sort data, { field, sort[bool] }
 * @return Province	list provinsi dalam entity response
 */
func (repo RegionRepositoryMock) FindAllProvince(sort []map[string]interface{}) ([]entities.Province, error) {
	args := repo.Mock.Called()
	return args.Get(0).([]entities.Province), args.Error(1)
}
/*
 * Find All City
 * -------------------------------
 * Mengambil semua data kota berdasarkan provinsi
 *
 * @var sort		sort data, { field, sort[bool] }
 * @return City		list kota dalam entity response
 */
func (repo RegionRepositoryMock) FindAllCity(provinceID int, sort []map[string]interface{}) ([]entities.City, error) {
	args := repo.Mock.Called()
	return args.Get(0).([]entities.City), args.Error(1)
}

/*
 * Find City
 * -------------------------------
 * Mencari data kota tunggal berdasarkan ID
 *
 * @var id 		data id
 */
func (repo RegionRepositoryMock) FindCity(id int) (entities.City, error) {
	args := repo.Mock.Called()
	return args.Get(0).(entities.City), args.Error(1)
}

/*
 * Find All District
 * -------------------------------
 * Mengambil semua data kecamatan berdasarkan kota
 *
 * @var sort		sort data, { field, sort[bool] }
 * @return District	list kecamatan dalam entity response 
 */
func (repo RegionRepositoryMock) FindAllDistrict(cityID int, sort []map[string]interface{}) ([]entities.District, error) {
	args := repo.Mock.Called()
	return args.Get(0).([]entities.District), args.Error(1)
}