package handlers

import (
	"bringeee-capstone/configs"
	"bringeee-capstone/deliveries/helpers"
	middleware "bringeee-capstone/deliveries/middlewares"
	"bringeee-capstone/entities"
	"bringeee-capstone/entities/web"
	authService "bringeee-capstone/services/auth"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *authService.AuthService
}

func NewAuthHandler(service *authService.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: service,
	}
}

func (handler AuthHandler) Login(c echo.Context) error {
	// Populate request input
	authReq := entities.AuthRequest{
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
	}

	// define link hateoas
	links := map[string]string{"self": configs.Get().App.BaseURL + "/api/auth"}

	// call auth service login
	authRes, err := handler.authService.Login(authReq)
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	}

	// send response
	return c.JSON(200, web.SuccessResponse{
		Status: "OK",
		Code:   200,
		Error:  nil,
		Links:  links,
		Data:   authRes,
	})
}

func (handler AuthHandler) Me(c echo.Context) error {

	// Token and Read Token
	token := c.Get("user")
	Id, _, err := middleware.ReadToken(token)

	// Define link
	links := map[string]string{"self": configs.Get().App.BaseURL + "/api/auth/me"}
	if err != nil {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Code:   http.StatusUnauthorized,
			Status: "ERROR",
			Error:  "unauthorized",
			Links:  links,
		})
	}

	// Memanggil service auth me
	authRes, err := handler.authService.Me(Id, token)
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	}

	// Response
	return c.JSON(200, web.SuccessResponse{
		Status: "OK",
		Code:   200,
		Error:  nil,
		Links:  links,
		Data:   authRes,
	})
}
