package auth_test

import (
	"bringeee-capstone/entities"
	"bringeee-capstone/entities/web"
	_userRepository "bringeee-capstone/repositories/user"
	_authService "bringeee-capstone/services/auth"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)


func TestLogin(t *testing.T) {
	t.Run("driver-success", func(t *testing.T) {
		userSample := _userRepository.UserCollection[0]
		userSample.Role = "driver"
		userSample.Password = "$2a$12$iu2L7bKpW4Rpe5yPGt3KPOm5N229fSuMlkHYu5l25dIwgvvW6oQYO" // pass: password
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindByCustomer").Return(userSample, nil)

		driverSample := _userRepository.DriverCollection[0]
		driverSample.UserID = userSample.ID
		driverSample.AccountStatus = "VERIFIED"
		userRepositoryMock.Mock.On("FindByDriver").Return(driverSample, nil)
		
		authService := _authService.NewAuthService(userRepositoryMock)
		actual, err := authService.Login(entities.AuthRequest{
			Email: userSample.Email,
			Password: "password",
		})
		
		assert.Nil(t, err)
		assert.NotEqual(t, "", actual.(entities.DriverAuthResponse).Token)
	})
	t.Run("invalid-email", func(t *testing.T) {
		userSample := _userRepository.UserCollection[0]
		userSample.Role = "driver"
		userSample.Password = "$2a$12$iu2L7bKpW4Rpe5yPGt3KPOm5N229fSuMlkHYu5l25dIwgvvW6oQYO" // pass: password
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindByCustomer").Return(entities.User{}, web.WebError{})

		driverSample := _userRepository.DriverCollection[0]
		driverSample.UserID = userSample.ID
		driverSample.AccountStatus = "VERIFIED"
		userRepositoryMock.Mock.On("FindByDriver").Return(driverSample, nil)
		
		authService := _authService.NewAuthService(userRepositoryMock)
		_, err := authService.Login(entities.AuthRequest{
			Email: userSample.Email + "wrongwrongwrong",
			Password: "password",
		})

		assert.Error(t, err)
	})
	t.Run("invalid-password", func(t *testing.T) {
		userSample := _userRepository.UserCollection[0]
		userSample.Role = "driver"
		userSample.Password = "$2a$12$iu2L7bKpW4Rpe5yPGt3KPOm5N229fSuMlkHYu5l25dIwgvvW6oQYO" // pass: password
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindByCustomer").Return(userSample, nil)

		driverSample := _userRepository.DriverCollection[0]
		driverSample.UserID = userSample.ID
		driverSample.AccountStatus = "VERIFIED"
		userRepositoryMock.Mock.On("FindByDriver").Return(driverSample, nil)
		
		authService := _authService.NewAuthService(userRepositoryMock)
		_, err := authService.Login(entities.AuthRequest{
			Email: userSample.Email,
			Password: "invalidpasswordhere",
		})
		
		assert.Error(t, err)
	})
	t.Run("driver-need-confirmation", func(t *testing.T) {
		userSample := _userRepository.UserCollection[0]
		userSample.Role = "driver"
		userSample.Password = "$2a$12$iu2L7bKpW4Rpe5yPGt3KPOm5N229fSuMlkHYu5l25dIwgvvW6oQYO" // pass: password
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindByCustomer").Return(userSample, nil)

		driverSample := _userRepository.DriverCollection[0]
		driverSample.UserID = userSample.ID
		driverSample.AccountStatus = "PENDING"
		userRepositoryMock.Mock.On("FindByDriver").Return(driverSample, nil)
		
		authService := _authService.NewAuthService(userRepositoryMock)
		actual, err := authService.Login(entities.AuthRequest{
			Email: userSample.Email,
			Password: "password",
		})
		
		assert.Error(t, err)
		assert.NotEqual(t, entities.DriverAuthResponse{}, actual.(entities.DriverAuthResponse).Token)
	})
	t.Run("admin-success", func(t *testing.T) {
		userSample := _userRepository.UserCollection[0]
		userSample.Role = "admin"
		userSample.Password = "$2a$12$iu2L7bKpW4Rpe5yPGt3KPOm5N229fSuMlkHYu5l25dIwgvvW6oQYO" // pass: password
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindByCustomer").Return(userSample, nil)
		
		authService := _authService.NewAuthService(userRepositoryMock)
		actual, err := authService.Login(entities.AuthRequest{
			Email: userSample.Email,
			Password: "password",
		})
		
		assert.Nil(t, err)
		assert.NotEqual(t, "", actual.(entities.AdminAuthResponse).Token)
	})
	t.Run("customer-success", func(t *testing.T) {
		userSample := _userRepository.UserCollection[0]
		userSample.Role = "customer"
		userSample.Password = "$2a$12$iu2L7bKpW4Rpe5yPGt3KPOm5N229fSuMlkHYu5l25dIwgvvW6oQYO" // pass: password
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindByCustomer").Return(userSample, nil)
		
		authService := _authService.NewAuthService(userRepositoryMock)
		actual, err := authService.Login(entities.AuthRequest{
			Email: userSample.Email,
			Password: "password",
		})
		
		assert.Nil(t, err)
		assert.NotEqual(t, "", actual.(entities.CustomerAuthResponse).Token)
	})
}

func TestMe(t *testing.T) {
	t.Run("customer-success", func(t *testing.T) {
		userSample := _userRepository.UserCollection[0]
		userSample.Role = "customer"
		userSample.Password = "$2a$12$iu2L7bKpW4Rpe5yPGt3KPOm5N229fSuMlkHYu5l25dIwgvvW6oQYO" // pass: password
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindCustomer").Return(userSample, nil)

		jwt := jwt.Token{
			Raw: "",
			Method: jwt.SigningMethodHS256,
			Claims: jwt.MapClaims{},
		}
		authService := _authService.NewAuthService(userRepositoryMock)
		_, err := authService.Me(int(userSample.ID), &jwt)

		assert.Nil(t, err)
	})
	t.Run("admin-success", func(t *testing.T) {
		userSample := _userRepository.UserCollection[0]
		userSample.Role = "admin"
		userSample.Password = "$2a$12$iu2L7bKpW4Rpe5yPGt3KPOm5N229fSuMlkHYu5l25dIwgvvW6oQYO" // pass: password
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindCustomer").Return(userSample, nil)

		jwt := jwt.Token{
			Raw: "",
			Method: jwt.SigningMethodHS256,
			Claims: jwt.MapClaims{},
		}
		authService := _authService.NewAuthService(userRepositoryMock)
		_, err := authService.Me(int(userSample.ID), &jwt)

		assert.Nil(t, err)
	})
	t.Run("driver-success", func(t *testing.T) {
		userSample := _userRepository.UserCollection[0]
		userSample.Role = "driver"
		userSample.Password = "$2a$12$iu2L7bKpW4Rpe5yPGt3KPOm5N229fSuMlkHYu5l25dIwgvvW6oQYO" // pass: password
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindCustomer").Return(userSample, nil)

		driverSample := _userRepository.DriverCollection[0]
		userRepositoryMock.Mock.On("FindByDriver").Return(driverSample, nil)

		jwt := jwt.Token{
			Raw: "",
			Method: jwt.SigningMethodHS256,
			Claims: jwt.MapClaims{},
		}
		authService := _authService.NewAuthService(userRepositoryMock)
		_, err := authService.Me(int(userSample.ID), &jwt)

		assert.Nil(t, err)
	})
}