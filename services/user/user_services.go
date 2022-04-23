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

func (service UserService) CreateCustomer(userRequest entities.CreateCustomerRequest, avatar *multipart.FileHeader) (entities.CustomerAuthResponse, error) {

	// Validation
	userFiles := []*multipart.FileHeader{}
	if avatar != nil {
		userFiles = append(userFiles, avatar)
	}
	err := validations.ValidateCreateCustomerRequest(service.validate, userRequest, userFiles)
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
	if avatar != nil {
		avatarFile, err := avatar.Open()
		if err != nil {
			return entities.CustomerAuthResponse{}, web.WebError{Code: 500, Message: "Cannot process avatar image"}
		}
		defer avatarFile.Close()

		// Upload avatar to S3
		filename := uuid.New().String() + avatar.Filename
		avatarURL, err := helpers.UploadFileToS3("avatar/"+filename, avatarFile)
		if err != nil {
			return entities.CustomerAuthResponse{}, web.WebError{Code: 500, Message: err.Error()}
		}
		user.Avatar = avatarURL
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
