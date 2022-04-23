package handlers

import (
	"bringeee-capstone/configs"
	"bringeee-capstone/deliveries/helpers"
	middleware "bringeee-capstone/deliveries/middlewares"
	"bringeee-capstone/entities"
	"bringeee-capstone/entities/web"
	_authService "bringeee-capstone/services/auth"
	"net/http"
	"reflect"

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

		// return error response khusus jika err termasuk webError
		if reflect.TypeOf(err).String() == "web.WebError" {
			webErr := err.(web.WebError)
			return c.JSON(webErr.Code, web.ErrorResponse{
				Status: "ERROR",
				Code:   webErr.Code,
				Error:  webErr.Error(),
				Links:  links,
			})
		}

		// return error 500 jika bukan webError
		return c.JSON(http.StatusInternalServerError, web.ErrorResponse{
			Status: "ERROR",
			Code:   http.StatusInternalServerError,
			Error:  err.Error(),
			Links:  links,
		})
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

func (handler AuthHandler) CustomerMe(c echo.Context) error {

	// Token and Read Token
	token := c.Get("user")
	Id, _, err := middleware.ReadToken(token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helpers.ResponseFailed("Bad Request"))
	}

	// Define link
	links := map[string]string{"self": configs.Get().App.BaseURL + "/api/auth/me"}

	// Memanggil service auth me
	authRes, err := handler.authService.Me(Id, token)
	if err != nil {

		// return error response khusus jika err termasuk webError
		if reflect.TypeOf(err).String() == "web.WebError" {
			webErr := err.(web.WebError)
			return c.JSON(webErr.Code, web.ErrorResponse{
				Status: "ERROR",
				Code:   webErr.Code,
				Error:  webErr.Error(),
				Links:  links,
			})
		}

		// return error 500 jika bukan webError
		return c.JSON(http.StatusInternalServerError, web.ErrorResponse{
			Status: "ERROR",
			Code:   http.StatusInternalServerError,
			Error:  err.Error(),
			Links:  links,
		})
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
