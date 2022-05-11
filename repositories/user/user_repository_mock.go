package user

import (
	"bringeee-capstone/entities"
	truckTypeRepository "bringeee-capstone/repositories/truck_type"
	"time"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type UserRepositoryMock struct {
	Mock *mock.Mock
}

func NewUserRepositoryMock(mock *mock.Mock) *UserRepositoryMock {
	return &UserRepositoryMock{
		Mock: mock,
	}
}

var UserCollection = []entities.User{
    {
        Model:       gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
        Name:        "test1",
        Email:       "test1@mail.com",
        Password:    "test",
        PhoneNumber: "081222333444",
        Gender:      "male",
        DOB:         time.Now(),
        Address:     "jl. reformasi",
        Avatar:      "some avatar",
        Role:        "customer",
    },
    {
        Model:       gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
        Name:        "test2",
        Email:       "test2@mail.com",
        Password:    "test",
        PhoneNumber: "082333444555",
        Gender:      "male",
        DOB:         time.Now(),
        Address:     "jl. reformasi",
        Avatar:      "some avatar",
        Role:        "admin",
    },
    {
        Model:       gorm.Model{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()},
        Name:        "test3",
        Email:       "test3@mail.com",
        Password:    "test",
        PhoneNumber: "083444555666",
        Gender:      "male",
        DOB:         time.Now(),
        Address:     "jl. reformasi",
        Avatar:      "some avatar",
        Role:        "driver",
    },
    {
        Model:       gorm.Model{ID: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()},
        Name:        "test4",
        Email:       "test4@mail.com",
        Password:    "test",
        PhoneNumber: "084555666777",
        Gender:      "male",
        DOB:         time.Now(),
        Address:     "jl. reformasi",
        Avatar:      "some avatar",
        Role:        "driver",
    },
}
var DriverCollection = []entities.Driver{
    {Model: gorm.Model{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()},
        UserID:            3,
        TruckTypeID:       1,
        KtpFile:           "some ktp file",
        StnkFile:          "some stnk file",
        DriverLicenseFile: "some SIM file",
        Age:               29,
        VehicleIdentifier: "DP 1111 CC",
        NIK:               "123456789",
        AccountStatus:     "PENDING",
        Status:            "IDLE",
        VehiclePicture:    "some vehicle picture",
        User:              UserCollection[2],
        TruckType:         truckTypeRepository.TruckTypeCollection[0],
    },
    {Model: gorm.Model{ID: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()},
        UserID:            4,
        TruckTypeID:       2,
        KtpFile:           "some ktp file",
        StnkFile:          "some stnk file",
        DriverLicenseFile: "some SIM file",
        Age:               39,
        VehicleIdentifier: "DP 2222 CC",
        NIK:               "987654321",
        AccountStatus:     "VERIFIED",
        Status:            "BUSY",
        VehiclePicture:    "some vehicle picture",
        User:              UserCollection[3],
        TruckType:         truckTypeRepository.TruckTypeCollection[1],
    },
}

func (repo UserRepositoryMock) FindAllCustomer(limit int, offset int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.User, error) {
	args := repo.Mock.Called()
	return args.Get(0).([]entities.User), args.Error(1)
}
func (repo UserRepositoryMock) FindCustomer(id int) (entities.User, error) {
	args := repo.Mock.Called()
	return args.Get(0).(entities.User), args.Error(1)
}
func (repo UserRepositoryMock) FindByCustomer(field string, value string) (entities.User, error) {
	args := repo.Mock.Called()
	return args.Get(0).(entities.User), args.Error(1)
}
func (repo UserRepositoryMock) CountAllCustomer(filters []map[string]string) (int64, error) {
	args := repo.Mock.Called()
	return int64(args.Int(0)), args.Error(1)
}
func (repo UserRepositoryMock) StoreCustomer(user entities.User) (entities.User, error) {
	args := repo.Mock.Called()
	return args.Get(0).(entities.User), args.Error(1)
}
func (repo UserRepositoryMock) UpdateCustomer(user entities.User, id int) (entities.User, error) {
	args := repo.Mock.Called()
	return args.Get(0).(entities.User), args.Error(1)
}
func (repo UserRepositoryMock) DeleteCustomer(id int) error {
	args := repo.Mock.Called()
	return args.Error(1)
}
func (repo UserRepositoryMock) FindAllDriver(limit int, offset int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.Driver, error) {
	args := repo.Mock.Called()
	return args.Get(0).([]entities.Driver), args.Error(1)
}
func (repo UserRepositoryMock) FindDriver(id int) (entities.Driver, error) {
	args := repo.Mock.Called()
	return args.Get(0).(entities.Driver), args.Error(1)
}
func (repo UserRepositoryMock) FindByDriver(field string, value string) (entities.Driver, error) {
	args := repo.Mock.Called()
	return args.Get(0).(entities.Driver), args.Error(1)
}
func (repo UserRepositoryMock) CountAllDriver(filters []map[string]string) (int64, error) {
	args := repo.Mock.Called()
	return int64(args.Int(0)), args.Error(1)
}
func (repo UserRepositoryMock) StoreDriver(driver entities.Driver) (entities.Driver, error) {
	args := repo.Mock.Called()
	return args.Get(0).(entities.Driver), args.Error(1)
}
func (repo UserRepositoryMock) UpdateDriver(driver entities.Driver, id int) (entities.Driver, error) {
	args := repo.Mock.Called()
	return args.Get(0).(entities.Driver), args.Error(1)
}
func (repo UserRepositoryMock) DeleteDriver(id int) error {
	args := repo.Mock.Called()
	return args.Error(1)
}