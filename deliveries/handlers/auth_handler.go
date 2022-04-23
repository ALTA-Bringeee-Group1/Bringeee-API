package handlers

import (
	"bringeee-capstone/deliveries/helpers"
	_authService "bringeee-capstone/services/auth"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

type AuthHandler struct {
	authService _authService.AuthServiceInterface
}

func NewAuthHandler(auth _authService.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		authService: auth,
	}
}

func (ah *AuthHandler) LoginHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		type loginData struct {
			Identifier string `json:"identifier"`
			Password   string `json:"password"`
		}
		var login loginData
		err := c.Bind(&login)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helpers.ResponseFailed("error bind data"))
		}
		token, errorLogin := ah.authService.Login(login.Identifier, login.Password)
		if errorLogin != nil {
			return c.JSON(http.StatusBadRequest, helpers.ResponseFailed(fmt.Sprintf("%v", errorLogin)))
		}
		responseToken := map[string]interface{}{
			"token": token,
		}
		return c.JSON(200, web.SuccessResponse{
			Status: "OK",
			Code:   200,
			Error:  nil,
			Links:  links,
			Data:   authRes,
		})
	}
}
