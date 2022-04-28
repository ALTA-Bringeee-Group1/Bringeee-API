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
	"net/url"
	"strconv"
	"strings"
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
		return entities.CustomerAuthResponse{}, web.WebError{Code: 400, ProductionMessage: "bad request", DevelopmentMessage: "date of birth format is invalid"}
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
				return entities.CustomerAuthResponse{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "Cannot process file image"}
			}
			defer fileFile.Close()

			// Upload file to S3
			filename := uuid.New().String() + file.Filename
			fileURL, err := helpers.UploadFileToS3("users/avatar/"+filename, fileFile)
			if err != nil {
				return entities.CustomerAuthResponse{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "Cannot upload file image"}
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
		return entities.DriverAuthResponse{}, web.WebError{Code: 400, ProductionMessage: "server error", DevelopmentMessage: "date of birth format is invalid"}
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
				return entities.DriverAuthResponse{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "Cannot process file image"}
			}
			defer fileFile.Close()

			// Upload file to S3
			filename := uuid.New().String() + file.Filename
			fileURL, err := helpers.UploadFileToS3("users/avatar/"+filename, fileFile)
			if err != nil {
				return entities.DriverAuthResponse{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: err.Error()}
			}
			user.Avatar = fileURL

		case "ktp_file":
			fileFile, err := file.Open()
			if err != nil {
				return entities.DriverAuthResponse{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "Cannot process file image"}
			}
			defer fileFile.Close()

			// Upload file to S3
			filename := uuid.New().String() + file.Filename
			fileURL, err := helpers.UploadFileToS3("drivers/ktp/"+filename, fileFile)
			if err != nil {
				return entities.DriverAuthResponse{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: err.Error()}
			}
			driver.KtpFile = fileURL

		case "stnk_file":
			fileFile, err := file.Open()
			if err != nil {
				return entities.DriverAuthResponse{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "Cannot process file image"}
			}
			defer fileFile.Close()

			// Upload file to S3
			filename := uuid.New().String() + file.Filename
			fileURL, err := helpers.UploadFileToS3("drivers/stnk/"+filename, fileFile)
			if err != nil {
				return entities.DriverAuthResponse{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: err.Error()}
			}
			driver.StnkFile = fileURL

		case "driver_license_file":
			fileFile, err := file.Open()
			if err != nil {
				return entities.DriverAuthResponse{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "Cannot process file image"}
			}
			defer fileFile.Close()

			// Upload file to S3
			filename := uuid.New().String() + file.Filename
			fileURL, err := helpers.UploadFileToS3("drivers/driver_license/"+filename, fileFile)
			if err != nil {
				return entities.DriverAuthResponse{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: err.Error()}
			}
			driver.DriverLicenseFile = fileURL

		case "vehicle_picture":
			fileFile, err := file.Open()
			if err != nil {
				return entities.DriverAuthResponse{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "Cannot process file image"}
			}
			defer fileFile.Close()

			// Upload file to S3
			filename := uuid.New().String() + file.Filename
			fileURL, err := helpers.UploadFileToS3("drivers/vehicle_picture/"+filename, fileFile)
			if err != nil {
				return entities.DriverAuthResponse{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: err.Error()}
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
	copier.Copy(&userRes, &user)
	copier.Copy(&userRes, &driverRes)
	// generate token

	// Buat auth response untuk dimasukkan token dan driver
	authRes := entities.DriverAuthResponse{
		Token: "",
		User:  userRes,
	}
	return authRes, nil
}

func (service UserService) UpdateCustomer(customerRequest entities.UpdateCustomerRequest, id int, files map[string]*multipart.FileHeader) (entities.CustomerResponse, error) {

	// validation
	err := validations.ValidateUpdateCustomerRequest(files)
	if err != nil {
		return entities.CustomerResponse{}, err
	}

	// Get user by ID via repository
	user, err := service.userRepo.FindCustomer(id)
	if err != nil {
		return entities.CustomerResponse{}, web.WebError{Code: 400, ProductionMessage: "server error", DevelopmentMessage: "The requested ID doesn't match with any record"}
	}
	// Avatar
	for field, file := range files {
		switch field {
		case "avatar":
			// Delete avatar lama jika ada yang baru
			if user.Avatar != "" {
				u, _ := url.Parse(user.Avatar)
				objectPathS3 := strings.TrimPrefix(u.Path, "/")
				helpers.DeleteFromS3(objectPathS3)
			}
			fileFile, err := file.Open()
			if err != nil {
				return entities.CustomerResponse{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "Cannot process file image"}
			}
			defer fileFile.Close()

			// Upload file to S3
			filename := uuid.New().String() + file.Filename
			fileURL, err := helpers.UploadFileToS3("users/avatar/"+filename, fileFile)
			if err != nil {
				return entities.CustomerResponse{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: err.Error()}
			}
			user.Avatar = fileURL
		}
	}

	// Konversi dari request ke domain entities user - mengabaikan nilai kosong pada request
	copier.CopyWithOption(&user, &customerRequest, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	// Hanya hash password jika password juga diganti (tidak kosong)
	if customerRequest.Password != "" {
		hashedPassword, _ := helpers.HashPassword(user.Password)
		user.Password = hashedPassword
	}

	// Update via repository
	user, err = service.userRepo.UpdateCustomer(user, id)
	// Konversi user domain menjadi user response
	userRes := entities.CustomerResponse{}
	copier.Copy(&userRes, &user)

	return userRes, err
}

func (service UserService) UpdateDriver(driverRequest entities.UpdateDriverRequest, id int, files map[string]*multipart.FileHeader) (entities.DriverResponse, error) {

	// validation
	err := validations.ValidateUpdateDriverRequest(files)
	if err != nil {
		return entities.DriverResponse{}, err
	}

	// Get user by ID via repository
	user, err := service.userRepo.FindCustomer(id)
	if err != nil {
		return entities.DriverResponse{}, web.WebError{Code: 400, ProductionMessage: "server error", DevelopmentMessage: "The requested ID doesn't match with any record"}
	}
	// Avatar
	for field, file := range files {
		switch field {
		case "avatar":
			// Delete avatar lama jika ada yang baru
			if user.Avatar != "" {
				u, _ := url.Parse(user.Avatar)
				objectPathS3 := strings.TrimPrefix(u.Path, "/")
				helpers.DeleteFromS3(objectPathS3)
			}
			fileFile, err := file.Open()
			if err != nil {
				return entities.DriverResponse{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "Cannot process file image"}
			}
			defer fileFile.Close()

			// Upload file to S3
			filename := uuid.New().String() + file.Filename
			fileURL, err := helpers.UploadFileToS3("users/avatar/"+filename, fileFile)
			if err != nil {
				return entities.DriverResponse{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: err.Error()}
			}
			user.Avatar = fileURL
		}
	}

	// Konversi dari request ke domain entities user - mengabaikan nilai kosong pada request
	copier.CopyWithOption(&user, &driverRequest, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	// Hanya hash password jika password juga diganti (tidak kosong)
	if driverRequest.Password != "" {
		hashedPassword, _ := helpers.HashPassword(user.Password)
		user.Password = hashedPassword
	}

	// Update via repository
	user, err = service.userRepo.UpdateCustomer(user, id)
	if err != nil {
		return entities.DriverResponse{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: err.Error()}
	}
	// find driver
	driver, err := service.userRepo.FindByDriver("user_id", strconv.Itoa(int(user.ID)))
	if err != nil {
		return entities.DriverResponse{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: err.Error()}
	}
	copier.CopyWithOption(&driver, &driverRequest, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	driver, err = service.userRepo.UpdateDriver(driver, int(driver.ID))
	// Konversi user domain menjadi user response
	userRes := entities.DriverResponse{}
	copier.Copy(&userRes, &driver.User)
	copier.Copy(&userRes, &driver)
	copier.Copy(&userRes.TruckType, &driver.TruckType)

	return userRes, err
}
func (service UserService) UpdateDriverByAdmin(driverRequest entities.UpdateDriverByAdminRequest, id int, files map[string]*multipart.FileHeader) (entities.DriverResponse, error) {

	// validation
	err := validations.ValidateUpdateDriverRequest(files)
	if err != nil {
		return entities.DriverResponse{}, err
	}
	// find driver
	driver, err := service.userRepo.FindDriver(id)
	if err != nil {
		return entities.DriverResponse{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: err.Error()}
	}
	// Avatar
	for field, file := range files {
		switch field {
		case "ktp_file":
			// Delete avatar lama jika ada yang baru
			if driver.KtpFile != "" {
				u, _ := url.Parse(driver.KtpFile)
				objectPathS3 := strings.TrimPrefix(u.Path, "/")
				helpers.DeleteFromS3(objectPathS3)
			}
			fileFile, err := file.Open()
			if err != nil {
				return entities.DriverResponse{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "Cannot process file image"}
			}
			defer fileFile.Close()

			// Upload file to S3
			filename := uuid.New().String() + file.Filename
			fileURL, err := helpers.UploadFileToS3("drivers/ktp/"+filename, fileFile)
			if err != nil {
				return entities.DriverResponse{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: err.Error()}
			}
			driver.KtpFile = fileURL

		case "stnk_file":
			// Delete avatar lama jika ada yang baru
			if driver.StnkFile != "" {
				u, _ := url.Parse(driver.StnkFile)
				objectPathS3 := strings.TrimPrefix(u.Path, "/")
				helpers.DeleteFromS3(objectPathS3)
			}
			fileFile, err := file.Open()
			if err != nil {
				return entities.DriverResponse{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "Cannot process file image"}
			}
			defer fileFile.Close()

			// Upload file to S3
			filename := uuid.New().String() + file.Filename
			fileURL, err := helpers.UploadFileToS3("drivers/stnk/"+filename, fileFile)
			if err != nil {
				return entities.DriverResponse{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: err.Error()}
			}
			driver.StnkFile = fileURL

		case "driver_license_file":
			// Delete avatar lama jika ada yang baru
			if driver.DriverLicenseFile != "" {
				u, _ := url.Parse(driver.DriverLicenseFile)
				objectPathS3 := strings.TrimPrefix(u.Path, "/")
				helpers.DeleteFromS3(objectPathS3)
			}
			fileFile, err := file.Open()
			if err != nil {
				return entities.DriverResponse{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "Cannot process file image"}
			}
			defer fileFile.Close()

			// Upload file to S3
			filename := uuid.New().String() + file.Filename
			fileURL, err := helpers.UploadFileToS3("drivers/driver_license/"+filename, fileFile)
			if err != nil {
				return entities.DriverResponse{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: err.Error()}
			}
			driver.DriverLicenseFile = fileURL

		case "vehicle_picture":
			// Delete avatar lama jika ada yang baru
			if driver.VehiclePicture != "" {
				u, _ := url.Parse(driver.VehiclePicture)
				objectPathS3 := strings.TrimPrefix(u.Path, "/")
				helpers.DeleteFromS3(objectPathS3)
			}
			fileFile, err := file.Open()
			if err != nil {
				return entities.DriverResponse{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "Cannot process file image"}
			}
			defer fileFile.Close()

			// Upload file to S3
			filename := uuid.New().String() + file.Filename
			fileURL, err := helpers.UploadFileToS3("drivers/vehicle_picture/"+filename, fileFile)
			if err != nil {
				return entities.DriverResponse{}, web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: err.Error()}
			}
			driver.VehiclePicture = fileURL
		}
	}

	copier.CopyWithOption(&driver, &driverRequest, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	driver, err = service.userRepo.UpdateDriver(driver, int(driver.ID))
	// Konversi user domain menjadi user response
	userRes := entities.DriverResponse{}
	copier.Copy(&userRes, &driver.User)
	copier.Copy(&userRes, &driver)
	copier.Copy(&userRes.TruckType, &driver.TruckType)

	return userRes, err
}

func (service UserService) DeleteCustomer(id int) error {

	// Cari user berdasarkan ID via repo
	user, err := service.userRepo.FindCustomer(id)
	if err != nil {
		return web.WebError{Code: 400, ProductionMessage: "server error", DevelopmentMessage: "The requested ID doesn't match with any record"}
	}

	// Delete avatar lama jika ada yang baru
	if user.Avatar != "" {
		u, _ := url.Parse(user.Avatar)
		objectPathS3 := strings.TrimPrefix(u.Path, "/")
		helpers.DeleteFromS3(objectPathS3)
	}

	// Delete via repository
	err = service.userRepo.DeleteCustomer(id)
	return err
}

func (service UserService) FindDriver(id int) (entities.DriverResponse, error) {

	driver, err := service.userRepo.FindDriver(id)
	if err != nil {
		return entities.DriverResponse{}, err
	}
	user, err := service.userRepo.FindByCustomer("id", strconv.Itoa(int(driver.UserID)))
	driverRes := entities.DriverResponse{}
	copier.Copy(&driverRes, &user)
	copier.Copy(&driverRes, &driver)

	return driverRes, err
}

func (service UserService) FindCustomer(id int) (entities.CustomerResponse, error) {

	user, err := service.userRepo.FindCustomer(id)
	if err != nil {
		return entities.CustomerResponse{}, err
	} else if user.Role == "driver" {
		return entities.CustomerResponse{}, err
	}
	userRes := entities.CustomerResponse{}
	copier.Copy(&userRes, &user)

	return userRes, err
}

func (service UserService) GetPaginationCustomer(limit, page int, filters []map[string]string) (web.Pagination, error) {
	totalRows, err := service.userRepo.CountAllCustomer(filters)
	if err != nil {
		return web.Pagination{}, err
	}
	if limit <= 0 {
		limit = 1
	}
	totalPages := totalRows / int64(limit)
	if totalRows%int64(limit) > 0 {
		totalPages++
	}

	return web.Pagination{
		Page:       page,
		Limit:      limit,
		TotalPages: int(totalPages),
	}, nil
}

func (service UserService) GetPaginationDriver(limit, page int, filters []map[string]string) (web.Pagination, error) {
	totalRows, err := service.userRepo.CountAllDriver(filters)
	if err != nil {
		return web.Pagination{}, err
	}
	if limit <= 0 {
		limit = 1
	}
	totalPages := totalRows / int64(limit)
	if totalRows%int64(limit) > 0 {
		totalPages++
	}

	return web.Pagination{
		Page:       page,
		Limit:      limit,
		TotalPages: int(totalPages),
	}, nil
}

func (service UserService) FindAllCustomer(limit, page int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.CustomerResponse, error) {

	offset := (page - 1) * limit

	usersRes := []entities.CustomerResponse{}
	users, err := service.userRepo.FindAllDriver(limit, offset, filters, sorts)

	copier.Copy(&usersRes, &users)

	return usersRes, err
}

func (service UserService) FindAllDriver(limit, page int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.DriverResponse, error) {

	offset := (page - 1) * limit

	driversRes := []entities.DriverResponse{}
	drivers, err := service.userRepo.FindAllDriver(limit, offset, filters, sorts)
	copier.Copy(&driversRes, &drivers)
	for key, value := range drivers {
		user, _ := service.userRepo.FindByCustomer("id", strconv.Itoa(int(value.UserID)))
		driversRes[key].Name = user.Name
		driversRes[key].Email = user.Email
		driversRes[key].DOB = user.DOB
		driversRes[key].Gender = user.Gender
		driversRes[key].Address = user.Address
		driversRes[key].PhoneNumber = user.PhoneNumber
		driversRes[key].Avatar = user.Avatar
		driversRes[key].Role = user.Role
	}

	return driversRes, err
}

func (service UserService) DeleteDriver(id int) error {

	// Cari user berdasarkan ID via repo
	driver, err := service.userRepo.FindDriver(id)
	if err != nil {
		return web.WebError{Code: 400, ProductionMessage: "bad request", DevelopmentMessage: "The requested ID doesn't match with any record"}
	}

	user, err := service.userRepo.FindByCustomer("id", strconv.Itoa(int(driver.UserID)))
	if err != nil {
		return web.WebError{Code: 400, ProductionMessage: "bad request", DevelopmentMessage: "The requested ID user doesn't match with any record"}
	}
	// Delete file di s3
	if user.Avatar != "" {
		u, _ := url.Parse(user.Avatar)
		objectPathS3 := strings.TrimPrefix(u.Path, "/")
		helpers.DeleteFromS3(objectPathS3)
	}
	if driver.KtpFile != "" {
		u, _ := url.Parse(driver.KtpFile)
		objectPathS3 := strings.TrimPrefix(u.Path, "/")
		helpers.DeleteFromS3(objectPathS3)
	}
	if driver.StnkFile != "" {
		u, _ := url.Parse(driver.StnkFile)
		objectPathS3 := strings.TrimPrefix(u.Path, "/")
		helpers.DeleteFromS3(objectPathS3)
	}
	if driver.DriverLicenseFile != "" {
		u, _ := url.Parse(driver.DriverLicenseFile)
		objectPathS3 := strings.TrimPrefix(u.Path, "/")
		helpers.DeleteFromS3(objectPathS3)
	}
	if driver.VehiclePicture != "" {
		u, _ := url.Parse(driver.VehiclePicture)
		objectPathS3 := strings.TrimPrefix(u.Path, "/")
		helpers.DeleteFromS3(objectPathS3)
	}

	// Delete via repository
	service.userRepo.DeleteDriver(int(driver.ID))
	err = service.userRepo.DeleteCustomer(int(driver.UserID))
	return err
}

func (service UserService) FindByDriver(field string, value string) (entities.DriverResponse, error) {

	driver, err := service.userRepo.FindByDriver(field, value)
	if err != nil {
		return entities.DriverResponse{}, err
	}
	user, err := service.FindByCustomer("id", strconv.Itoa(int(driver.UserID)))
	if err != nil {
		return entities.DriverResponse{}, err
	}
	driverRes := entities.DriverResponse{}
	copier.Copy(&driverRes, &user)
	copier.Copy(&driverRes, &driver)

	return driverRes, err
}

func (service UserService) FindByCustomer(field string, value string) (entities.CustomerResponse, error) {

	user, err := service.userRepo.FindByCustomer(field, value)
	if err != nil {
		return entities.CustomerResponse{}, err
	}
	userRes := entities.CustomerResponse{}
	copier.Copy(&userRes, &user)

	return userRes, err
}
