package user

import (
	"bringeee-capstone/deliveries/helpers"
	_middleware "bringeee-capstone/deliveries/middlewares"
	"bringeee-capstone/deliveries/validations"
	"bringeee-capstone/entities"
	web "bringeee-capstone/entities/web"
	userRepository "bringeee-capstone/repositories/user"
	"mime/multipart"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type UserService struct {
	userRepo userRepository.UserRepositoryInterface
	validate *validator.Validate
}

func NewUserService(repository userRepository.UserRepositoryInterface) *UserService {
	return &UserService{
		userRepo: repository,
		validate: validator.New(),
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
			fileURL, err := helpers.UploadFileToS3("users/file/"+filename, fileFile)
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
