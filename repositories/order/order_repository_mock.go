package order

import (
	"bringeee-capstone/entities"
	truckTypeRepository "bringeee-capstone/repositories/truck_type"
	userRepository "bringeee-capstone/repositories/user"
	"time"

	"github.com/stretchr/testify/mock"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type OrderRepositoryMock struct {
	Mock *mock.Mock
}

var OrderCollection = []entities.Order{
	{
		Model: gorm.Model{ ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now() },
		DriverID: null.IntFrom(3),
		Driver: userRepository.DriverCollection[0],
		CustomerID: userRepository.UserCollection[0].ID,
		Customer: userRepository.UserCollection[0],
		DestinationID: 1,
		Destination: entities.Destination{
			Model: gorm.Model{ ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now() },
			DestinationStartProvince: "JAWA TIMUR",
			DestinationStartCity: "MALANG",
			DestinationStartDistrict: "GONDANGLEGI",
			DestinationStartAddress: "Jl. Melati 03/01, Putat Lor",
			DestinationStartPostal: "65174",
			DestinationStartLat: "-8.16281415182338",
			DestinationStartLong: "112.65624115364912",
			DestinationEndProvince: "JAWA TENGAH",
			DestinationEndCity: "SEMARANG",
			DestinationEndDistrict: "GUBENG",
			DestinationEndAddress: "Jl. Mawar, 03/03",
			DestinationEndPostal: "67712",
			DestinationEndLat: "-7.006076535110995",
			DestinationEndLong: "110.43177925041718",
		},
		TruckTypeID: int(truckTypeRepository.TruckTypeCollection[0].ID),
		TruckType: truckTypeRepository.TruckTypeCollection[0],
		OrderPicture: "https://source.unsplash.com/600x600/?cargo",
		TotalVolume: 16000330,
		TotalWeight: 1000,
		Distance: 140,
		EstimatedPrice: 3000000,
		FixPrice: 3250000,
		TransactionID: "xya-7721d-ma",
		PaymentMethod: "BANK_TRANSFER_BNI",
		Status: "DELIVERED",
		ArrivedPicture: "",
	},
	{
		Model: gorm.Model{ ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now() },
		DriverID: null.IntFromPtr(nil),
		CustomerID: userRepository.UserCollection[0].ID,
		Customer: userRepository.UserCollection[0],
		DestinationID: 3,
		Destination: entities.Destination{
			Model: gorm.Model{ ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now() },
			DestinationStartProvince: "DI YOGYAKARTA",
			DestinationStartCity: "YOGYAKARTA",
			DestinationStartDistrict: "DANUREJAN",
			DestinationStartAddress: "Suryatmajan, Kec. Danurejan, Kota Yogyakarta, Daerah Istimewa Yogyakarta",
			DestinationStartPostal: "55213",
			DestinationStartLat: "-7.793050394271023",
			DestinationStartLong: "110.36756312713727",
			DestinationEndProvince: "JAWA TENGAH",
			DestinationEndCity: "SURAKARTA (SOLO)",
			DestinationEndDistrict: "JEBRES",
			DestinationEndAddress: "Tegalharjo, Kec. Jebres, Kota Surakarta, Jawa Tengah",
			DestinationEndPostal: "57129",
			DestinationEndLat: "-7.561160260537138",
			DestinationEndLong: "110.83655443176414",
		},
		TruckTypeID: int(truckTypeRepository.TruckTypeCollection[0].ID),
		TruckType: truckTypeRepository.TruckTypeCollection[0],
		OrderPicture: "https://source.unsplash.com/600x600/?cargo",
		TotalVolume: 16000330,
		TotalWeight: 1000,
		Distance: 140,
		EstimatedPrice: 3000000,
		FixPrice: 0,
		TransactionID: "xzc-7722d-ca",
		PaymentMethod: "BANK_TRANSFER_BCA",
		Status: "NEED_CUSTOMER_CONFIRM",
		ArrivedPicture: "",
	},
	{
		Model: gorm.Model{ ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now() },
		DriverID: null.IntFromPtr(nil),
		CustomerID: userRepository.UserCollection[0].ID,
		Customer: userRepository.UserCollection[0],
		DestinationID: 3,
		Destination: entities.Destination{
			Model: gorm.Model{ ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now() },
			DestinationStartProvince: "DI YOGYAKARTA",
			DestinationStartCity: "YOGYAKARTA",
			DestinationStartDistrict: "DANUREJAN",
			DestinationStartAddress: "Suryatmajan, Kec. Danurejan, Kota Yogyakarta, Daerah Istimewa Yogyakarta",
			DestinationStartPostal: "55213",
			DestinationStartLat: "-7.793050394271023",
			DestinationStartLong: "110.36756312713727",
			DestinationEndProvince: "JAWA TENGAH",
			DestinationEndCity: "SURAKARTA (SOLO)",
			DestinationEndDistrict: "JEBRES",
			DestinationEndAddress: "Tegalharjo, Kec. Jebres, Kota Surakarta, Jawa Tengah",
			DestinationEndPostal: "57129",
			DestinationEndLat: "-7.561160260537138",
			DestinationEndLong: "110.83655443176414",
		},
		TruckTypeID: int(truckTypeRepository.TruckTypeCollection[0].ID),
		TruckType: truckTypeRepository.TruckTypeCollection[0],
		OrderPicture: "https://source.unsplash.com/600x600/?cargo",
		TotalVolume: 16000330,
		TotalWeight: 1000,
		Distance: 140,
		EstimatedPrice: 3000000,
		FixPrice: 0,
		Status: "REQUESTED",
		ArrivedPicture: "",
	},
}

func NewOrderRepositoryMock(mock *mock.Mock) *OrderRepositoryMock {
	return &OrderRepositoryMock{
		Mock: mock,
	}
}

/*
 * Find All
 * -------------------------------
 * Mengambil data order berdasarkan filters dan sorts
 *
 * @var limit 	batas limit hasil query
 * @var offset 	offset hasil query
 * @var filters	query untuk penyaringan data, { field, operator, value }
 * @var sorts	pengurutan data, { field, value[bool] }
 * @return order	list order dalam bentuk entity domain
 * @return error	error
 */
func (repo OrderRepositoryMock) FindAll(limit int, offset int, filters []map[string]interface{}, sorts []map[string]interface{}) ([]entities.Order, error) {
	args := repo.Mock.Called(limit, offset, filters, sorts)
	return args.Get(0).([]entities.Order), args.Error(1)
}

/*
 * Find
 * -------------------------------
 * Mencari order tunggal berdasarkan ID
 *
 * @var id 		data id
 * @return order	single order dalam bentuk entity domain
 * @return error	error
 */
func (repository OrderRepositoryMock) Find(id int) (entities.Order, error) {
	args := repository.Mock.Called(id)
	return args.Get(0).(entities.Order), args.Error(1)
}

/*
 * Find User
 * -------------------------------
 * Mencari order berdasarkan field tertentu
 *
 * @var field 	kolom data
 * @var value	nilai data
 * @return order	single order dalam bentuk entity domain
 * @return error	error
 */
func (repository OrderRepositoryMock) FindBy(field string, value string) (entities.Order, error) {
	args := repository.Mock.Called(field, value)
	return args.Get(0).(entities.Order), args.Error(1)
}

/*
 * Find First
 * -------------------------------
 * Mengambil data order tunggal berdasarkan filter
 *
 * @var filters		query untuk penyaringan data, { field, operator, value }
 * @return order	order dalam bentuk entity domain
 * @return error	error
 */
func (repository OrderRepositoryMock) FindFirst(filters []map[string]interface{}) (entities.Order, error) {
	args := repository.Mock.Called(filters)
	return args.Get(0).(entities.Order), args.Error(1)
}

/*
 * CountAll
 * -------------------------------
 * Menghitung semua orders (ini digunakan untuk pagination di service)
 *
 * @return order	single order dalam bentuk entity domain
 * @return error	error
 */
func (repository OrderRepositoryMock) CountAll(filters []map[string]interface{}) (int64, error) {
	args := repository.Mock.Called(filters)
	return int64(args.Int(0)), args.Error(1)
}

/*
 * Store
 * -------------------------------
 * Menambahkan data order kedalam database, beserta data destination
 *
 * @var order		single order entity
 * @var destination	single destination entity
 * @return order	single order dalam bentuk entity domain
 */
func (repository OrderRepositoryMock) Store(order entities.Order, destination entities.Destination) (entities.Order, error) {
	args := repository.Mock.Called()
	return args.Get(0).(entities.Order), args.Error(1)
}

/*
 * Update
 * -------------------------------
 * Mengupdate data order berdasarkan ID
 *
 * @var order		single order entity
 * @return order	single order dalam bentuk entity domain
 * @return error	error
 */
func (repository OrderRepositoryMock) Update(order entities.Order, id int) (entities.Order, error) {
	args := repository.Mock.Called(order)
	return args.Get(0).(entities.Order), args.Error(1)
}

/*
 * Delete
 * -------------------------------
 * Delete order berdasarkan ID
 *
 * @return error	error
 */
func (repository OrderRepositoryMock) Delete(id int, destinationID int) error {
	args := repository.Mock.Called(id)
	return args.Error(1)
}

/*
 * Delete Batch
 * -------------------------------
 * Delete multiple order berdasarkan filter tertentu
 *
 * @var filters	query untuk penyaringan data, { field, operator, value }
 * @return error	error
 */
func (repository OrderRepositoryMock) DeleteBatch(filters []map[string]interface{}) error {
	args := repository.Mock.Called(filters)
	return args.Error(1)
}

func (repository OrderRepositoryMock) FindByDate(day int) ([]map[string]interface{}, error) {
	args := repository.Mock.Called(day)
	return args.Get(0).([]map[string]interface{}), nil
}

func (repository OrderRepositoryMock) FindByMonth(month int, year int) ([]entities.Order, error) {
	args := repository.Mock.Called(month, year)
	return args.Get(0).([]entities.Order), nil
}