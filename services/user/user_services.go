package user

import (
	"bringeee-capstone/deliveries/helpers"
	_middleware "bringeee-capstone/deliveries/middlewares"
	"bringeee-capstone/deliveries/validations"
	"bringeee-capstone/entities"
	web "bringeee-capstone/entities/web"
	truckRepository "bringeee-capstone/repositories/truck_type"
	userRepository "bringeee-capstone/repositories/user"
	"mime/multipart"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type UserService struct {
	userRepo  userRepository.UserRepositoryInterface
	truckRepo truckRepository.TruckTypeRepositoryInterface
	validate  *validator.Validate
}

func NewUserService(repository userRepository.UserRepositoryInterface, truckRepo truckRepository.TruckTypeRepositoryInterface) *UserService {
	return &UserService{
		userRepo:  repository,
		truckRepo: truckRepo,
		validate:  validator.New(),
	}
}

func (service UserService) CreateCustomer(userRequest entities.CreateCustomerRequest, files map[string]*multipart.FileHeader) (entities.CustomerAuthResponse, error) {

	// Validation
	err := validations.ValidateCreateCustomerRequest(service.validate, userRequest, files)
	if err != nil {
		return entities.CustomerAuthResponse{}, err
	}

	// Konversi user request menjadi domain untuk diteruskan ke repository
	user := entities.User{}
	copier.Copy(&user, &userRequest)

	// Konversi datetime untuk field datetime (dob)
	dob, err := time.Parse("2006-01-02", userRequest.DOB)
	if err != nil {
		return entities.CustomerAuthResponse{}, web.WebError{Code: 400, Message: "date of birth format is invalid"}
	}
	user.DOB = dob

	// Password hashing menggunakan bcrypt
	hashedPassword, _ := helpers.HashPassword(user.Password)
	user.Password = hashedPassword

	// Upload avatar if exists
	for field, file := range files {
		switch field {
		case "avatar":
			fileFile, err := file.Open()
			if err != nil {
				return entities.CustomerAuthResponse{}, web.WebError{Code: 500, Message: "Cannot process file image"}
			}
			defer fileFile.Close()

			// Upload file to S3
			filename := uuid.New().String() + file.Filename
			fileURL, err := helpers.UploadFileToS3("users/avatar/"+filename, fileFile)
			if err != nil {
				return entities.CustomerAuthResponse{}, web.WebError{Code: 500, Message: err.Error()}
			}
			user.Avatar = fileURL
		}
	}
	user.Role = "customer"

	// Insert ke sistem melewati repository
	user, err = service.userRepo.StoreCustomer(user)
	if err != nil {
		return entities.CustomerAuthResponse{}, err
	}

	// Konversi hasil repository menjadi user response
	userRes := entities.CustomerResponse{}
	copier.Copy(&userRes, &user)

	// generate token
	token, err := _middleware.CreateToken(int(user.ID), user.Name, user.Role)
	if err != nil {
		return entities.CustomerAuthResponse{}, err
	}

	// Buat auth response untuk dimasukkan token dan user
	authRes := entities.CustomerAuthResponse{
		Token: token,
		User:  userRes,
	}
	return authRes, nil
}

func (service UserService) CreateDriver(driverRequest entities.CreateDriverRequest, files map[string]*multipart.FileHeader) (entities.DriverAuthResponse, error) {

	// Validation
	err := validations.ValidateCreateDriverRequest(service.validate, driverRequest, files)
	if err != nil {
		return entities.DriverAuthResponse{}, err
	}

	// Konversi user request menjadi domain untuk diteruskan ke repository
	user := entities.User{}
	driver := entities.Driver{}
	copier.Copy(&user, &driverRequest)
	copier.Copy(&driver, &driverRequest)

	// Konversi datetime untuk field datetime (dob)
	dob, err := time.Parse("2006-01-02", driverRequest.DOB)
	if err != nil {
		return entities.DriverAuthResponse{}, web.WebError{Code: 400, Message: "date of birth format is invalid"}
	}
	user.DOB = dob

	// Password hashing menggunakan bcrypt
	hashedPassword, _ := helpers.HashPassword(user.Password)
	user.Password = hashedPassword

	// Upload avatar if exists
	for field, file := range files {
		switch field {
		case "avatar":
			fileFile, err := file.Open()
			if err != nil {
				return entities.DriverAuthResponse{}, web.WebError{Code: 500, Message: "Cannot process file image"}
			}
			defer fileFile.Close()

			// Upload file to S3
			filename := uuid.New().String() + file.Filename
			fileURL, err := helpers.UploadFileToS3("users/avatar/"+filename, fileFile)
			if err != nil {
				return entities.DriverAuthResponse{}, web.WebError{Code: 500, Message: err.Error()}
			}
			user.Avatar = fileURL

		case "ktp_file":
			fileFile, err := file.Open()
			if err != nil {
				return entities.DriverAuthResponse{}, web.WebError{Code: 500, Message: "Cannot process file image"}
			}
			defer fileFile.Close()

			// Upload file to S3
			filename := uuid.New().String() + file.Filename
			fileURL, err := helpers.UploadFileToS3("drivers/ktp/"+filename, fileFile)
			if err != nil {
				return entities.DriverAuthResponse{}, web.WebError{Code: 500, Message: err.Error()}
			}
			driver.KtpFile = fileURL

		case "stnk_file":
			fileFile, err := file.Open()
			if err != nil {
				return entities.DriverAuthResponse{}, web.WebError{Code: 500, Message: "Cannot process file image"}
			}
			defer fileFile.Close()

			// Upload file to S3
			filename := uuid.New().String() + file.Filename
			fileURL, err := helpers.UploadFileToS3("drivers/stnk/"+filename, fileFile)
			if err != nil {
				return entities.DriverAuthResponse{}, web.WebError{Code: 500, Message: err.Error()}
			}
			driver.StnkFile = fileURL

		case "driver_license_file":
			fileFile, err := file.Open()
			if err != nil {
				return entities.DriverAuthResponse{}, web.WebError{Code: 500, Message: "Cannot process file image"}
			}
			defer fileFile.Close()

			// Upload file to S3
			filename := uuid.New().String() + file.Filename
			fileURL, err := helpers.UploadFileToS3("drivers/driver_license/"+filename, fileFile)
			if err != nil {
				return entities.DriverAuthResponse{}, web.WebError{Code: 500, Message: err.Error()}
			}
			driver.DriverLicenseFile = fileURL

		case "vehicle_picture":
			fileFile, err := file.Open()
			if err != nil {
				return entities.DriverAuthResponse{}, web.WebError{Code: 500, Message: "Cannot process file image"}
			}
			defer fileFile.Close()

			// Upload file to S3
			filename := uuid.New().String() + file.Filename
			fileURL, err := helpers.UploadFileToS3("drivers/vehicle_picture/"+filename, fileFile)
			if err != nil {
				return entities.DriverAuthResponse{}, web.WebError{Code: 500, Message: err.Error()}
			}
			driver.VehiclePicture = fileURL
		}
	}
	user.Role = "driver"

	// Insert ke sistem melewati repository
	user, err = service.userRepo.StoreCustomer(user)
	if err != nil {
		return entities.DriverAuthResponse{}, err
	}

	// Insert ke database driver repo
	driver.AccountStatus = "PENDING"
	driver.Status = "IDLE"
	driver.UserID = user.ID
	driver, err = service.userRepo.StoreDriver(driver)
	if err != nil {
		return entities.DriverAuthResponse{}, err
	}
	driverRes, _ := service.userRepo.FindDriver(int(driver.ID))
	// Konversi hasil repository menjadi driver response
	userRes := entities.DriverResponse{}
	copier.Copy(&userRes, &driverRes)
	// generate token
	token, err := _middleware.CreateToken(int(user.ID), user.Name, user.Role)
	if err != nil {
		return entities.DriverAuthResponse{}, err
	}

	// Buat auth response untuk dimasukkan token dan driver
	authRes := entities.DriverAuthResponse{
		Token: token,
		User:  userRes,
	}
	return authRes, nil
}
