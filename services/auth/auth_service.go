package auth

import (
	middleware "bringeee-capstone/deliveries/middlewares"
	userRepository "bringeee-capstone/repositories/user"

	"bringeee-capstone/entities"
	web "bringeee-capstone/entities/web"

	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo userRepository.UserRepositoryInterface
}

func NewAuthService(userRepo userRepository.UserRepositoryInterface) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

/*
 * Auth Service - Login
 * -------------------------------
 * Mencari user berdasarkan ID
 */
func (service AuthService) LoginCustomer(authReq entities.AuthRequest) (entities.CustomerAuthResponse, error) {

	// Get user by username via repository
	user, err := service.userRepo.FindBy("email", authReq.Email)
	if err != nil {
		return entities.CustomerAuthResponse{}, web.WebError{Code: 401, Message: "Invalid credential"}
	}

	// Verify password
	match := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authReq.Password))
	if match != nil {
		return entities.CustomerAuthResponse{}, web.WebError{Code: 401, Message: "Invalid password"}
	}

	// Konversi menjadi user response
	userRes := entities.CustomerResponse{}
	copier.Copy(&userRes, &user)

	// Create token
	token, err := middleware.CreateToken(int(userRes.ID), userRes.Name, userRes.Role)
	if err != nil {
		return entities.CustomerAuthResponse{}, web.WebError{Code: 500, Message: "Error create token"}
	}

	return entities.CustomerAuthResponse{
		Token: token,
		User:  userRes,
	}, nil
}

/*
 * Auth Service - Me
 * -------------------------------
 * Mendapatkan userdata yang sedang login
 */
func (service AuthService) CustomerMe(Id int) (entities.CustomerAuthResponse, error) {

	// Get userdata via repository
	user, err := service.userRepo.Find(Id)

	// Konversi user ke user response
	userRes := entities.UserResponse{}
	copier.Copy(&userRes, &user)

	// Bentuk auth response
	authRes := entities.CustomerAuthResponse{
		Token: userJWT.Raw,
		User:  userRes,
	}

	return authRes, err
}
