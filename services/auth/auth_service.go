package auth

import (
	"bringeee-capstone/deliveries/helpers"
	middleware "bringeee-capstone/deliveries/middlewares"
	userRepository "bringeee-capstone/repositories/user"
	"bringeee-capstone/utils"
	"fmt"

	"bringeee-capstone/entities"
	web "bringeee-capstone/entities/web"

	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/copier"
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
func (service AuthService) Login(authReq entities.AuthRequest) (interface{}, error) {

	// Get user by username via repository
	user, err := service.userRepo.FindByCustomer("email", authReq.Email)
	if err != nil {
		return entities.CustomerAuthResponse{}, web.WebError{Code: 401, Message: "Invalid credential"}
	}

	// Verify password
	if helpers.CheckPasswordHash(authReq.Password, user.Password) {
		return entities.CustomerAuthResponse{}, web.WebError{Code: 401, Message: "Invalid password"}
	}

	// if role == driver
	if user.Role == "driver" {
		// Konversi menjadi driver response
		driver, _ := service.userRepo.FindByDriver("user_id", int(user.ID))
		if driver.AccountStatus != "VERIFIED" {
			return entities.DriverAuthResponse{}, web.WebError{Code: 403, Message: "Waiting for admin confirmation"}
		}
		userRes := entities.DriverResponse{}
		copier.Copy(&userRes, &driver.User)
		copier.Copy(&userRes, &driver)
		copier.Copy(&userRes.TruckType, &driver.TruckType)

		fmt.Println(utils.JsonEncode(driver))

		// Create token
		token, err := middleware.CreateToken(int(user.ID), userRes.Name, userRes.Role)
		if err != nil {
			return entities.DriverAuthResponse{}, web.WebError{Code: 500, Message: "Error create token"}
		}

		return entities.DriverAuthResponse{
			Token: token,
			User:  userRes,
		}, nil
	}

	// if role == admin
	if user.Role == "admin" {
		// Konversi menjadi admin response
		userRes := entities.AdminResponse{}
		copier.Copy(&userRes, &user)

		// Create token
		token, err := middleware.CreateToken(int(userRes.ID), userRes.Name, userRes.Role)
		if err != nil {
			return entities.AdminAuthResponse{}, web.WebError{Code: 500, Message: "Error create token"}
		}

		return entities.AdminAuthResponse{
			Token: token,
			User:  userRes,
		}, nil
	}

	// Konversi menjadi customer response
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
func (service AuthService) Me(Id int, token interface{}) (interface{}, error) {

	userJWT := token.(*jwt.Token)
	// Get userdata via repository
	user, err := service.userRepo.FindCustomer(Id)

	// Konversi user ke user response
	if user.Role == "driver" {
		driver, _ := service.userRepo.FindByDriver("user_id", Id)
		userRes := entities.DriverResponse{}
		copier.Copy(&userRes, &user)
		copier.Copy(&userRes, &driver)
		authRes := entities.DriverAuthResponse{
			Token: userJWT.Raw,
			User:  userRes,
		}
		return authRes, err
	}
	if user.Role == "admin" {
		userRes := entities.AdminResponse{}
		copier.Copy(&userRes, &user)
		authRes := entities.AdminAuthResponse{
			Token: userJWT.Raw,
			User:  userRes,
		}
		return authRes, err
	}
	userRes := entities.CustomerResponse{}
	copier.Copy(&userRes, &user)

	// Bentuk auth response
	authRes := entities.CustomerAuthResponse{
		Token: userJWT.Raw,
		User:  userRes,
	}

	return authRes, err
}
