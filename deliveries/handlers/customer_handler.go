package handlers

import (
	"bringeee-capstone/configs"
	"bringeee-capstone/entities"
	"bringeee-capstone/entities/web"
	userService "bringeee-capstone/services/user"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService *userService.UserService
}

func NewUserHandler(service *userService.UserService) *UserHandler {
	return &UserHandler{
		userService: service,
	}
}

/*
 * User Handler - Create
 * -------------------------------
 * Registrasi User kedalam sistem dan
 * mengembalikan token
 */
func (handler UserHandler) CreateCustomer(c echo.Context) error {

	// Bind request ke user request
	userReq := entities.CreateCustomerRequest{}
	c.Bind(&userReq)

	// Define links (hateoas)
	links := map[string]string{"self": configs.Get().App.BaseURL + "/api/customers"}

	// Read file avatar
	avatar, _ := c.FormFile("avatar")

	// registrasi user via call user service
	userRes, err := handler.userService.CreateCustomer(userReq, avatar)
	if err != nil {

		// return error response khusus jika err termasuk webError / ValidationError
		if reflect.TypeOf(err).String() == "web.WebError" {
			webErr := err.(web.WebError)
			return c.JSON(webErr.Code, web.ErrorResponse{
				Status: "ERROR",
				Code:   webErr.Code,
				Error:  webErr.Error(),
				Links:  links,
			})
		} else if reflect.TypeOf(err).String() == "web.ValidationError" {
			valErr := err.(web.ValidationError)
			return c.JSON(valErr.Code, web.ValidationErrorResponse{
				Status: "ERROR",
				Code:   valErr.Code,
				Error:  valErr.Error(),
				Errors: valErr.Errors,
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

	// response
	return c.JSON(http.StatusCreated, web.SuccessResponse{
		Status: "OK",
		Code:   http.StatusCreated,
		Error:  nil,
		Links:  links,
		Data:   userRes,
	})
}
