package user_test

import (
	"bringeee-capstone/entities"
	"bringeee-capstone/entities/web"
	_orderRepository "bringeee-capstone/repositories/order"
	_truckTypeRepository "bringeee-capstone/repositories/truck_type"
	_userRepository "bringeee-capstone/repositories/user"
	_storageProvider "bringeee-capstone/services/storage"
	_userService "bringeee-capstone/services/user"
	"mime/multipart"
	"testing"

	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCustomer(t *testing.T) {
	sampleCustomerCentral := _userRepository.UserCollection[0]
	sampleRequestCentral := entities.CreateCustomerRequest{}
	copier.Copy(&sampleRequestCentral, &sampleCustomerCentral)
	sampleRequestCentral.DOB = "1999-12-12"
	sampleFileRequestCentral := map[string]*multipart.FileHeader {
		"avatar": {
			Filename: "avatar.jpg",
			Size: 155 * 1024,
		},
	}
	t.Run("success", func(t *testing.T) {
		sampleCustomer := sampleCustomerCentral
		sampleRequest := sampleRequestCentral
		sampleFileRequest := sampleFileRequestCentral

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("UploadFromRequest").Return("example.com/images.png", nil)
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("StoreCustomer").Return(sampleCustomer, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.CreateCustomer(sampleRequest, sampleFileRequest, storageProvider)

		expected := entities.CustomerResponse{}
		copier.Copy(&expected, &sampleCustomer)

		assert.Nil(t, err)
		assert.NotEqual(t, "", actual.Token)
		assert.Equal(t, expected, actual.User)
	})
	t.Run("validation-fail", func(t *testing.T) {
		sampleCustomer := sampleCustomerCentral
		sampleRequest := sampleRequestCentral
		sampleFileRequest := sampleFileRequestCentral

		sampleRequest.Name = ""
		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("UploadFromRequest").Return("example.com/images.png", nil)
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("StoreCustomer").Return(sampleCustomer, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.CreateCustomer(sampleRequest, sampleFileRequest, storageProvider)

		expected := entities.CustomerResponse{}
		copier.Copy(&expected, &sampleCustomer)

		assert.Error(t, err)
		assert.Equal(t, entities.CustomerAuthResponse{}, actual)
	})
	t.Run("invalid-dob", func(t *testing.T) {
		sampleCustomer := sampleCustomerCentral
		sampleRequest := sampleRequestCentral
		sampleFileRequest := sampleFileRequestCentral

		sampleRequest.DOB = "2022222222"
		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("UploadFromRequest").Return("example.com/images.png", nil)
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("StoreCustomer").Return(sampleCustomer, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.CreateCustomer(sampleRequest, sampleFileRequest, storageProvider)

		expected := entities.CustomerResponse{}
		copier.Copy(&expected, &sampleCustomer)

		assert.Error(t, err)
		assert.Equal(t, entities.CustomerAuthResponse{}, actual)
	})
	t.Run("upload-fail", func(t *testing.T) {
		sampleCustomer := sampleCustomerCentral
		sampleRequest := sampleRequestCentral
		sampleFileRequest := sampleFileRequestCentral

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("UploadFromRequest").Return("", web.WebError{})
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("StoreCustomer").Return(sampleCustomer, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.CreateCustomer(sampleRequest, sampleFileRequest, storageProvider)

		expected := entities.CustomerAuthResponse{}
		assert.Error(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("store-fail", func(t *testing.T) {
		sampleRequest := sampleRequestCentral
		sampleFileRequest := sampleFileRequestCentral

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("UploadFromRequest").Return("domain.com/images.jpg", nil)
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("StoreCustomer").Return(entities.User{}, web.WebError{})

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.CreateCustomer(sampleRequest, sampleFileRequest, storageProvider)

		expected := entities.CustomerAuthResponse{}
		assert.Error(t, err)
		assert.Equal(t, expected, actual)
	})
}


func TestCreateDriver(t *testing.T) {
	sampleDriverCentral := _userRepository.DriverCollection[0]
	sampleUserCentral := sampleDriverCentral.User
	sampleRequestCentral := entities.CreateDriverRequest{}
	copier.Copy(&sampleRequestCentral, &sampleDriverCentral)
	copier.Copy(&sampleRequestCentral, &sampleDriverCentral.User)
	sampleRequestCentral.DOB = "1999-12-12"
	sampleFileRequestCentral := map[string]*multipart.FileHeader{
		"avatar": { Filename: "avatar.jpg", Size: 800 * 1024 },
		"ktp_file": { Filename: "ktp_file.jpg", Size: 800 * 1024 },
		"stnk_file": { Filename: "stnk_file.jpg", Size: 800 * 1024 },
		"driver_license_file": { Filename: "driver_license_file.jpg", Size: 800 * 1024 },
		"vehicle_picture": { Filename: "vehicle_picture.jpg", Size: 800 * 1024 },
	}

	t.Run("success", func(t *testing.T) {
		sampleDriver := sampleDriverCentral
		sampleUser := sampleUserCentral
		sampleRequest := sampleRequestCentral
		sampleFileRequest := sampleFileRequestCentral
		
		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("UploadFromRequest").Return("domain.com/images.jpg", nil)
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("StoreCustomer").Return(sampleUser, nil)
		userRepositoryMock.Mock.On("StoreDriver").Return(sampleDriver, nil)
		userRepositoryMock.Mock.On("FindDriver").Return(sampleDriver, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.CreateDriver(sampleRequest, sampleFileRequest, storageProvider)

		expected := entities.DriverResponse{}
		copier.Copy(&expected, &sampleUser) // code order matters
		copier.Copy(&expected, &sampleDriver)

		assert.Nil(t, err)
		assert.Equal(t, expected, actual.User)
	})
	t.Run("validation-fail", func(t *testing.T) {
		sampleDriver := sampleDriverCentral
		sampleUser := sampleUserCentral
		sampleRequest := sampleRequestCentral
		sampleFileRequest := sampleFileRequestCentral
		sampleRequest.Name = ""
		
		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("UploadFromRequest").Return("domain.com/images.jpg", nil)
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("StoreCustomer").Return(sampleUser, nil)
		userRepositoryMock.Mock.On("StoreDriver").Return(sampleDriver, nil)
		userRepositoryMock.Mock.On("FindDriver").Return(sampleDriver, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.CreateDriver(sampleRequest, sampleFileRequest, storageProvider)

		assert.Error(t, err)
		assert.Equal(t, entities.DriverAuthResponse{}, actual)
	})
	t.Run("dob-fail", func(t *testing.T) {
		sampleDriver := sampleDriverCentral
		sampleUser := sampleUserCentral
		sampleRequest := sampleRequestCentral
		sampleFileRequest := sampleFileRequestCentral
		sampleRequest.DOB = "99999"
		
		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("UploadFromRequest").Return("domain.com/images.jpg", nil)
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("StoreCustomer").Return(sampleUser, nil)
		userRepositoryMock.Mock.On("StoreDriver").Return(sampleDriver, nil)
		userRepositoryMock.Mock.On("FindDriver").Return(sampleDriver, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.CreateDriver(sampleRequest, sampleFileRequest, storageProvider)

		assert.Error(t, err)
		assert.Equal(t, entities.DriverAuthResponse{}, actual)
	})
	t.Run("upload-fail", func(t *testing.T) {
		sampleDriver := sampleDriverCentral
		sampleUser := sampleUserCentral
		sampleRequest := sampleRequestCentral
		sampleFileRequest := sampleFileRequestCentral
		
		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("UploadFromRequest").Return("", web.WebError{})
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("StoreCustomer").Return(sampleUser, nil)
		userRepositoryMock.Mock.On("StoreDriver").Return(sampleDriver, nil)
		userRepositoryMock.Mock.On("FindDriver").Return(sampleDriver, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.CreateDriver(sampleRequest, sampleFileRequest, storageProvider)

		assert.Error(t, err)
		assert.Equal(t, entities.DriverAuthResponse{}, actual)
	})
	t.Run("user-repo-fail", func(t *testing.T) {
		sampleDriver := sampleDriverCentral
		sampleRequest := sampleRequestCentral
		sampleFileRequest := sampleFileRequestCentral
		
		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("UploadFromRequest").Return("domain.com/images.jpg", nil)
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("StoreCustomer").Return(entities.User{}, web.WebError{})
		userRepositoryMock.Mock.On("StoreDriver").Return(sampleDriver, nil)
		userRepositoryMock.Mock.On("FindDriver").Return(sampleDriver, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.CreateDriver(sampleRequest, sampleFileRequest, storageProvider)

		assert.Error(t, err)
		assert.Equal(t, entities.DriverAuthResponse{}, actual)
	})
	t.Run("driver-repo-fail", func(t *testing.T) {
		sampleDriver := sampleDriverCentral
		sampleUser := sampleUserCentral
		sampleRequest := sampleRequestCentral
		sampleFileRequest := sampleFileRequestCentral
		
		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("UploadFromRequest").Return("domain.com/images.jpg", nil)
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("StoreCustomer").Return(sampleUser, nil)
		userRepositoryMock.Mock.On("StoreDriver").Return(entities.Driver{}, web.WebError{})
		userRepositoryMock.Mock.On("FindDriver").Return(sampleDriver, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.CreateDriver(sampleRequest, sampleFileRequest, storageProvider)

		assert.Error(t, err)
		assert.Equal(t, entities.DriverAuthResponse{}, actual)
	})
}


func TestUpdateCustomer(t *testing.T) {
	sampleCustomerCentral := _userRepository.UserCollection[0]
	sampleRequestCentral := entities.UpdateCustomerRequest{}
	copier.Copy(&sampleRequestCentral, &sampleCustomerCentral)
	sampleRequestCentral.DOB = "1999-12-12"
	sampleFileRequestCentral := map[string]*multipart.FileHeader {
		"avatar": {
			Filename: "avatar.jpg",
			Size: 155 * 1024,
		},
	}
	t.Run("success", func(t *testing.T) {
		sampleRequest := sampleRequestCentral
		sampleFileRequest := sampleFileRequestCentral
		sampleCustomer := sampleCustomerCentral

		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindCustomer").Return(sampleCustomer, nil)

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)
		storageProvider.Mock.On("UploadFromRequest").Return("domain.com/images.jpg", nil)

		customerOutput := sampleCustomer
		copier.CopyWithOption(&customerOutput, &sampleRequest, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		userRepositoryMock.Mock.On("UpdateCustomer").Return(customerOutput, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.UpdateCustomer(sampleRequest, int(sampleCustomer.ID), sampleFileRequest, storageProvider)
		expected := entities.CustomerResponse{}
		copier.Copy(&expected, &customerOutput)

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("validation-fail", func(t *testing.T) {
		sampleRequest := sampleRequestCentral
		sampleFileRequest := sampleFileRequestCentral
		sampleCustomer := sampleCustomerCentral
		sampleFileRequest["avatar"].Size = 1024 * 2048

		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindCustomer").Return(sampleCustomer, nil)

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)
		storageProvider.Mock.On("UploadFromRequest").Return("domain.com/images.jpg", nil)

		customerOutput := sampleCustomer
		copier.CopyWithOption(&customerOutput, &sampleRequest, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		userRepositoryMock.Mock.On("UpdateCustomer").Return(customerOutput, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.UpdateCustomer(sampleRequest, int(sampleCustomer.ID), sampleFileRequest, storageProvider)
		assert.Error(t, err)
		assert.Equal(t, entities.CustomerResponse{}, actual)
	})
	t.Run("find-customer-fail", func(t *testing.T) {
		sampleRequest := sampleRequestCentral
		sampleFileRequest := map[string]*multipart.FileHeader{}
		sampleCustomer := sampleCustomerCentral

		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindCustomer").Return(entities.User{}, web.WebError{})

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)
		storageProvider.Mock.On("UploadFromRequest").Return("domain.com/images.jpg", nil)

		customerOutput := sampleCustomer
		copier.CopyWithOption(&customerOutput, &sampleRequest, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		userRepositoryMock.Mock.On("UpdateCustomer").Return(customerOutput, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.UpdateCustomer(sampleRequest, int(sampleCustomer.ID), sampleFileRequest, storageProvider)
		assert.Error(t, err)
		assert.Equal(t, entities.CustomerResponse{}, actual)
	})
	t.Run("upload-fail", func(t *testing.T) {
		sampleRequest := sampleRequestCentral
		sampleFileRequest := map[string]*multipart.FileHeader{
			"avatar": { Filename: "avatar.jpg", Size: 512 * 1024 },
		}
		sampleCustomer := sampleCustomerCentral

		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindCustomer").Return(sampleCustomer, nil)

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)
		storageProvider.Mock.On("UploadFromRequest").Return("", web.WebError{})

		customerOutput := sampleCustomer
		copier.CopyWithOption(&customerOutput, &sampleRequest, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		userRepositoryMock.Mock.On("UpdateCustomer").Return(customerOutput, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.UpdateCustomer(sampleRequest, int(sampleCustomer.ID), sampleFileRequest, storageProvider)
		assert.Error(t, err)
		assert.Equal(t, entities.CustomerResponse{}, actual)
	})
}

func TestUpdateDriver(t *testing.T) {
	sampleDriverCentral := _userRepository.DriverCollection[0]
	sampleUserCentral := _userRepository.DriverCollection[0].User
	sampleRequestCentral := entities.UpdateDriverRequest{}
	copier.Copy(&sampleRequestCentral, &sampleDriverCentral)
	copier.Copy(&sampleRequestCentral, &sampleUserCentral)

	t.Run("success", func(t *testing.T) {
		sampleUser := sampleUserCentral
		sampleDriver := sampleDriverCentral
		sampleRequest := sampleRequestCentral
		sampleFile := map[string]*multipart.FileHeader{
			"avatar": { Filename: "avatar.jpg", Size: 1024 * 88 },
		}
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindCustomer").Return(sampleUser, nil)

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)
		storageProvider.Mock.On("UploadFromRequest").Return("domain.com/avatar.jpg", nil)

		outputCustomer := sampleUser
		outputDriver := sampleDriver
		copier.CopyWithOption(&outputCustomer, &sampleRequest, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		copier.CopyWithOption(&outputDriver, &sampleRequest, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		userRepositoryMock.Mock.On("UpdateCustomer").Return(outputCustomer, nil)
		userRepositoryMock.Mock.On("FindByDriver").Return(sampleDriver, nil)
		userRepositoryMock.Mock.On("UpdateDriver").Return(outputDriver, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.UpdateDriver(sampleRequest, int(sampleUser.ID), sampleFile, storageProvider)
		expected := entities.DriverResponse{}
		copier.Copy(&expected, &outputDriver.User)
		copier.Copy(&expected, &outputDriver)
		copier.Copy(&expected.TruckType, &outputDriver.TruckType)

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("validation-fail", func(t *testing.T) {
		sampleUser := sampleUserCentral
		sampleDriver := sampleDriverCentral
		sampleRequest := sampleRequestCentral
		sampleFile := map[string]*multipart.FileHeader{
			"avatar": { Filename: "avatar.jpg", Size: 1024 * 9999 },
		}
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindCustomer").Return(sampleUser, nil)

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)
		storageProvider.Mock.On("UploadFromRequest").Return("domain.com/avatar.jpg", nil)

		outputCustomer := sampleUser
		outputDriver := sampleDriver
		copier.CopyWithOption(&outputCustomer, &sampleRequest, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		copier.CopyWithOption(&outputDriver, &sampleRequest, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		userRepositoryMock.Mock.On("UpdateCustomer").Return(outputCustomer, nil)
		userRepositoryMock.Mock.On("FindByDriver").Return(sampleDriver, nil)
		userRepositoryMock.Mock.On("UpdateDriver").Return(outputDriver, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.UpdateDriver(sampleRequest, int(sampleUser.ID), sampleFile, storageProvider)

		assert.Error(t, err)
		assert.Equal(t, entities.DriverResponse{}, actual)
	})
	t.Run("find-customer-fail", func(t *testing.T) {
		sampleUser := sampleUserCentral
		sampleDriver := sampleDriverCentral
		sampleRequest := sampleRequestCentral
		sampleFile := map[string]*multipart.FileHeader{
			"avatar": { Filename: "avatar.jpg", Size: 1024 * 55 },
		}
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindCustomer").Return(entities.User{}, web.WebError{})

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)
		storageProvider.Mock.On("UploadFromRequest").Return("domain.com/avatar.jpg", nil)

		outputCustomer := sampleUser
		outputDriver := sampleDriver
		copier.CopyWithOption(&outputCustomer, &sampleRequest, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		copier.CopyWithOption(&outputDriver, &sampleRequest, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		userRepositoryMock.Mock.On("UpdateCustomer").Return(outputCustomer, nil)
		userRepositoryMock.Mock.On("FindByDriver").Return(sampleDriver, nil)
		userRepositoryMock.Mock.On("UpdateDriver").Return(outputDriver, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.UpdateDriver(sampleRequest, int(sampleUser.ID), sampleFile, storageProvider)

		assert.Error(t, err)
		assert.Equal(t, entities.DriverResponse{}, actual)
	})
	t.Run("upload-fail", func(t *testing.T) {
		sampleUser := sampleUserCentral
		sampleDriver := sampleDriverCentral
		sampleRequest := sampleRequestCentral
		sampleFile := map[string]*multipart.FileHeader{
			"avatar": { Filename: "avatar.jpg", Size: 1024 * 55 },
		}
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindCustomer").Return(sampleUser, nil)

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)
		storageProvider.Mock.On("UploadFromRequest").Return("", web.WebError{})

		outputCustomer := sampleUser
		outputDriver := sampleDriver
		copier.CopyWithOption(&outputCustomer, &sampleRequest, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		copier.CopyWithOption(&outputDriver, &sampleRequest, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		userRepositoryMock.Mock.On("UpdateCustomer").Return(outputCustomer, nil)
		userRepositoryMock.Mock.On("FindByDriver").Return(sampleDriver, nil)
		userRepositoryMock.Mock.On("UpdateDriver").Return(outputDriver, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.UpdateDriver(sampleRequest, int(sampleUser.ID), sampleFile, storageProvider)

		assert.Error(t, err)
		assert.Equal(t, entities.DriverResponse{}, actual)
	})
	t.Run("update-user-repo-fail", func(t *testing.T) {
		sampleUser := sampleUserCentral
		sampleDriver := sampleDriverCentral
		sampleRequest := sampleRequestCentral
		sampleFile := map[string]*multipart.FileHeader{
			"avatar": { Filename: "avatar.jpg", Size: 1024 * 55 },
		}
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindCustomer").Return(sampleUser, nil)

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)
		storageProvider.Mock.On("UploadFromRequest").Return("domain.com/avatar.jpg", nil)

		outputCustomer := sampleUser
		outputDriver := sampleDriver
		copier.CopyWithOption(&outputCustomer, &sampleRequest, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		copier.CopyWithOption(&outputDriver, &sampleRequest, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		userRepositoryMock.Mock.On("UpdateCustomer").Return(entities.User{}, web.WebError{})
		userRepositoryMock.Mock.On("FindByDriver").Return(sampleDriver, nil)
		userRepositoryMock.Mock.On("UpdateDriver").Return(outputDriver, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.UpdateDriver(sampleRequest, int(sampleUser.ID), sampleFile, storageProvider)

		assert.Error(t, err)
		assert.Equal(t, entities.DriverResponse{}, actual)
	})
	t.Run("find-driver-repo-fail", func(t *testing.T) {
		sampleUser := sampleUserCentral
		sampleDriver := sampleDriverCentral
		sampleRequest := sampleRequestCentral
		sampleFile := map[string]*multipart.FileHeader{
			"avatar": { Filename: "avatar.jpg", Size: 1024 * 55 },
		}
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindCustomer").Return(sampleUser, nil)

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)
		storageProvider.Mock.On("UploadFromRequest").Return("domain.com/avatar.jpg", nil)

		outputCustomer := sampleUser
		outputDriver := sampleDriver
		copier.CopyWithOption(&outputCustomer, &sampleRequest, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		copier.CopyWithOption(&outputDriver, &sampleRequest, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		userRepositoryMock.Mock.On("UpdateCustomer").Return(outputCustomer, nil)
		userRepositoryMock.Mock.On("FindByDriver").Return(entities.Driver{}, web.WebError{})
		userRepositoryMock.Mock.On("UpdateDriver").Return(outputDriver, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.UpdateDriver(sampleRequest, int(sampleUser.ID), sampleFile, storageProvider)

		assert.Error(t, err)
		assert.Equal(t, entities.DriverResponse{}, actual)
	})
}

func TestUpdateDriverByAdmin(t *testing.T) {
	sampleDriverCentral := _userRepository.DriverCollection[0]
	sampleUserCentral := _userRepository.DriverCollection[0].User
	sampleRequestCentral := entities.UpdateDriverByAdminRequest{}
	copier.Copy(&sampleRequestCentral, &sampleDriverCentral)
	copier.Copy(&sampleRequestCentral, &sampleUserCentral)

	t.Run("success", func(t *testing.T) {
		sampleUser := sampleUserCentral
		sampleDriver := sampleDriverCentral
		sampleRequest := sampleRequestCentral
		sampleFile := map[string]*multipart.FileHeader{
			"avatar": { Filename: "avatar.jpg", Size: 800 * 1024 },
			"ktp_file": { Filename: "ktp_file.jpg", Size: 800 * 1024 },
			"stnk_file": { Filename: "stnk_file.jpg", Size: 800 * 1024 },
			"driver_license_file": { Filename: "driver_license_file.jpg", Size: 800 * 1024 },
			"vehicle_picture": { Filename: "vehicle_picture.jpg", Size: 800 * 1024 },
		}
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindDriver").Return(sampleDriver, nil)

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)
		storageProvider.Mock.On("UploadFromRequest").Return("domain.com/avatar.jpg", nil)

		outputDriver := sampleDriver
		copier.CopyWithOption(&outputDriver, &sampleRequest, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		userRepositoryMock.Mock.On("UpdateDriver").Return(outputDriver, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.UpdateDriverByAdmin(sampleRequest, int(sampleUser.ID), sampleFile, storageProvider)
		expected := entities.DriverResponse{}
		copier.Copy(&expected, &outputDriver.User)
		copier.Copy(&expected, &outputDriver)
		copier.Copy(&expected.TruckType, &outputDriver.TruckType)

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
}


func TestDeleteCustomer(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		sampleCustomer := _userRepository.UserCollection[0]
		sampleCustomer.Role = "customer"
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindCustomer").Return(sampleCustomer, nil)

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)

		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("DeleteBatch").Return(nil)

		userRepositoryMock.Mock.On("DeleteCustomer").Return(nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			orderRepositoryMock,
		)
		err := userService.DeleteCustomer(int(sampleCustomer.ID), storageProvider)
		assert.Nil(t, err)
	})
	t.Run("repo-fail", func(t *testing.T) {
		sampleCustomer := _userRepository.UserCollection[0]
		sampleCustomer.Role = "customer"
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindCustomer").Return(entities.User{}, web.WebError{})

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)

		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("DeleteBatch").Return(nil)

		userRepositoryMock.Mock.On("DeleteCustomer").Return(nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			orderRepositoryMock,
		)
		err := userService.DeleteCustomer(int(sampleCustomer.ID), storageProvider)
		assert.Error(t, err)
	})
	t.Run("success", func(t *testing.T) {
		sampleCustomer := _userRepository.UserCollection[0]
		sampleCustomer.Role = "admin"
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindCustomer").Return(sampleCustomer, nil)

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)

		orderRepositoryMock := _orderRepository.NewOrderRepositoryMock(&mock.Mock{})
		orderRepositoryMock.Mock.On("DeleteBatch").Return(nil)

		userRepositoryMock.Mock.On("DeleteCustomer").Return(nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			orderRepositoryMock,
		)
		err := userService.DeleteCustomer(int(sampleCustomer.ID), storageProvider)
		assert.Error(t, err)
	})
}

func TestFindDriver(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		sampleDriver := _userRepository.DriverCollection[0]
		sampleUser := _userRepository.DriverCollection[0].User
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindDriver").Return(sampleDriver, nil)
		userRepositoryMock.Mock.On("FindByCustomer").Return(sampleUser, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)

		actual, err := userService.FindDriver(int(sampleDriver.ID))
		expected := entities.DriverResponse{}
		copier.Copy(&expected, &sampleUser)
		copier.Copy(&expected, &sampleDriver)

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("repo-fail", func(t *testing.T) {
		sampleDriver := _userRepository.DriverCollection[0]
		sampleUser := _userRepository.DriverCollection[0].User
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindDriver").Return(entities.Driver{}, web.WebError{})
		userRepositoryMock.Mock.On("FindByCustomer").Return(sampleUser, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)

		actual, err := userService.FindDriver(int(sampleDriver.ID))

		assert.Error(t, err)
		assert.Equal(t, entities.DriverResponse{}, actual)
	})
}

func TestFindCustomer(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		sampleCustomer := _userRepository.UserCollection[0]
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindCustomer").Return(sampleCustomer, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)

		actual, err := userService.FindCustomer(int(sampleCustomer.ID))
		expected := entities.CustomerResponse{}
		copier.Copy(&expected, &sampleCustomer)

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("repo-fail", func(t *testing.T) {
		sampleCustomer := _userRepository.UserCollection[0]
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindCustomer").Return(entities.User{}, web.WebError{})

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)

		actual, err := userService.FindCustomer(int(sampleCustomer.ID))
		assert.Error(t, err)
		assert.Equal(t, entities.CustomerResponse{}, actual)
	})
	t.Run("user-driver", func(t *testing.T) {
		sampleCustomer := _userRepository.UserCollection[0]
		sampleCustomer.Role = "driver"
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindCustomer").Return(sampleCustomer, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)

		actual, err := userService.FindCustomer(int(sampleCustomer.ID))

		assert.Nil(t, err)
		assert.Equal(t, entities.CustomerResponse{}, actual)
	})
}


func TestGetPaginationCustomer(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("CountAllCustomer").Return(20, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.GetPaginationCustomer(5, 1, []map[string]string{})

		expected := web.Pagination {
			Page:       1,
			Limit:      5,
			TotalPages: int(4),
		}
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("repo-fail", func(t *testing.T) {
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("CountAllCustomer").Return(0, web.WebError{})

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.GetPaginationCustomer(5, 1, []map[string]string{})
		assert.Error(t, err)
		assert.Equal(t, web.Pagination{}, actual)
	})
	t.Run("limit-zero", func(t *testing.T) {
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("CountAllCustomer").Return(20, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.GetPaginationCustomer(0, 1, []map[string]string{})

		expected := web.Pagination {
			Page:       1,
			Limit:      1,
			TotalPages: int(20),
		}
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("added-page-on-active-module", func(t *testing.T) {
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("CountAllCustomer").Return(22, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.GetPaginationCustomer(5, 1, []map[string]string{})

		expected := web.Pagination {
			Page:       1,
			Limit:      5,
			TotalPages: int(5),
		}
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestGetPaginationDriver(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("CountAllDriver").Return(20, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.GetPaginationDriver(5, 1, []map[string]string{})

		expected := web.Pagination {
			Page:       1,
			Limit:      5,
			TotalPages: int(4),
		}
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("repo-fail", func(t *testing.T) {
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("CountAllDriver").Return(0, web.WebError{})

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.GetPaginationDriver(5, 1, []map[string]string{})
		assert.Error(t, err)
		assert.Equal(t, web.Pagination{}, actual)
	})
	t.Run("limit-zero", func(t *testing.T) {
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("CountAllDriver").Return(20, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.GetPaginationDriver(0, 1, []map[string]string{})

		expected := web.Pagination {
			Page:       1,
			Limit:      1,
			TotalPages: int(20),
		}
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("added-page-on-active-module", func(t *testing.T) {
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("CountAllDriver").Return(22, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.GetPaginationDriver(5, 1, []map[string]string{})

		expected := web.Pagination {
			Page:       1,
			Limit:      5,
			TotalPages: int(5),
		}
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestFindAllCustomer(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userSamples := _userRepository.UserCollection
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindAllCustomer").Return(userSamples, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)

		actual, err := userService.FindAllCustomer(0, 0, []map[string]string{}, []map[string]interface{}{})
		expected := []entities.CustomerResponse{}
		copier.Copy(&expected, &actual)
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("fail", func(t *testing.T) {
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindAllCustomer").Return([]entities.User{}, web.WebError{})

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)

		actual, err := userService.FindAllCustomer(0, 0, []map[string]string{}, []map[string]interface{}{})
		assert.Error(t, err)
		assert.Equal(t, []entities.CustomerResponse{}, actual)
	})
}

func TestFindAllDriver(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		driverSamples := _userRepository.DriverCollection
		userSample := _userRepository.UserCollection[0]
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindAllDriver").Return(driverSamples, nil)
		userRepositoryMock.Mock.On("FindByCustomer").Return(userSample, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		_, err := userService.FindAllDriver(0, 0, []map[string]string{}, []map[string]interface{}{})
		assert.Nil(t, err)
	})
	t.Run("fail", func(t *testing.T) {
		userSample := _userRepository.UserCollection[0]
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindAllDriver").Return([]entities.Driver{}, web.WebError{})
		userRepositoryMock.Mock.On("FindByCustomer").Return(userSample, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		_, err := userService.FindAllDriver(0, 0, []map[string]string{}, []map[string]interface{}{})
		assert.Error(t, err)
	})
}

func TestDeleteDriver(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		driverSample := _userRepository.DriverCollection[0]
		userSample := driverSample.User
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindDriver").Return(driverSample, nil)
		userRepositoryMock.Mock.On("FindByCustomer").Return(userSample, nil)

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)

		userRepositoryMock.Mock.On("DeleteDriver").Return(nil)
		userRepositoryMock.Mock.On("DeleteCustomer").Return(nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)

		err := userService.DeleteDriver(int(driverSample.ID), storageProvider)
		assert.Nil(t, err)
	})
	t.Run("find-driver-fail", func(t *testing.T) {
		driverSample := _userRepository.DriverCollection[0]
		userSample := driverSample.User
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindDriver").Return(entities.Driver{}, web.WebError{})
		userRepositoryMock.Mock.On("FindByCustomer").Return(userSample, nil)

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)

		userRepositoryMock.Mock.On("DeleteDriver").Return(nil)
		userRepositoryMock.Mock.On("DeleteCustomer").Return(nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)

		err := userService.DeleteDriver(int(driverSample.ID), storageProvider)
		assert.Error(t, err)
	})
	t.Run("find-user-fail", func(t *testing.T) {
		driverSample := _userRepository.DriverCollection[0]
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindDriver").Return(driverSample, nil)
		userRepositoryMock.Mock.On("FindByCustomer").Return(entities.User{}, web.WebError{})

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)

		userRepositoryMock.Mock.On("DeleteDriver").Return(nil)
		userRepositoryMock.Mock.On("DeleteCustomer").Return(nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)

		err := userService.DeleteDriver(int(driverSample.ID), storageProvider)
		assert.Error(t, err)
	})
}


func TestFindByDriver(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		driverSample := _userRepository.DriverCollection[0]
		userSample := _userRepository.DriverCollection[0].User
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindByDriver").Return(driverSample, nil)
		userRepositoryMock.Mock.On("FindByCustomer").Return(userSample, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.FindByDriver("id", "1")
		expected := entities.DriverResponse{}
		copier.Copy(&expected, &actual)
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("driver-repo-fail", func(t *testing.T) {
		userSample := _userRepository.DriverCollection[0].User
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindByDriver").Return(entities.Driver{}, web.WebError{})
		userRepositoryMock.Mock.On("FindByCustomer").Return(userSample, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.FindByDriver("id", "1")
		expected := entities.DriverResponse{}
		copier.Copy(&expected, &actual)
		assert.Error(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("user-repo-fail", func(t *testing.T) {
		driverSample := _userRepository.DriverCollection[0]
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindByDriver").Return(driverSample, nil)
		userRepositoryMock.Mock.On("FindByCustomer").Return(entities.User{}, web.WebError{})

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.FindByDriver("id", "1")
		expected := entities.DriverResponse{}
		copier.Copy(&expected, &actual)
		assert.Error(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestVerifiedDriverAccount(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		driverSample := _userRepository.DriverCollection[0]
		driverSample.AccountStatus = "PENDING"
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindDriver").Return(driverSample, nil)

		driverOutput := driverSample
		driverOutput.AccountStatus = "VERIFIED"
		userRepositoryMock.Mock.On("UpdateDriver").Return(driverOutput, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		err := userService.VerifiedDriverAccount(int(driverSample.ID))

		assert.Nil(t, err)
	})
	t.Run("find-driver-fail", func(t *testing.T) {
		driverSample := _userRepository.DriverCollection[0]
		driverSample.AccountStatus = "PENDING"
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindDriver").Return(entities.Driver{}, web.WebError{})

		driverOutput := driverSample
		driverOutput.AccountStatus = "VERIFIED"
		userRepositoryMock.Mock.On("UpdateDriver").Return(driverOutput, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		err := userService.VerifiedDriverAccount(int(driverSample.ID))

		assert.Error(t, err)
	})
	t.Run("verified-driver", func(t *testing.T) {
		driverSample := _userRepository.DriverCollection[0]
		driverSample.AccountStatus = "VERIFIED"
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindDriver").Return(driverSample, nil)

		driverOutput := driverSample
		userRepositoryMock.Mock.On("UpdateDriver").Return(driverOutput, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		err := userService.VerifiedDriverAccount(int(driverSample.ID))

		assert.Error(t, err)
	})
}

func TestCountCustomer(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("CountAllCustomer").Return(20, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.CountCustomer([]map[string]string{})
		assert.Nil(t, err)
		assert.Equal(t, 20, actual)
	})
	t.Run("repo-error", func(t *testing.T) {
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("CountAllCustomer").Return(0, web.WebError{})

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.CountCustomer([]map[string]string{})
		assert.Error(t, err)
		assert.Equal(t, 0, actual)
	})
}

func TestCountDriver(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("CountAllDriver").Return(20, nil)

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.CountDriver([]map[string]string{})
		assert.Nil(t, err)
		assert.Equal(t, 20, actual)
	})
	t.Run("repo-error", func(t *testing.T) {
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("CountAllDriver").Return(0, web.WebError{})

		userService := _userService.NewUserService(
			userRepositoryMock,
			_truckTypeRepository.NewTruckTypeRepositoryMock(&mock.Mock{}),
			_orderRepository.NewOrderRepositoryMock(&mock.Mock{}),
		)
		actual, err := userService.CountDriver([]map[string]string{})
		assert.Error(t, err)
		assert.Equal(t, 0, actual)
	})
}