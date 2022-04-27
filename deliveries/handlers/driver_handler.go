package handlers

import (
	"bringeee-capstone/configs"
	middleware "bringeee-capstone/deliveries/middlewares"
	"bringeee-capstone/entities"
	"bringeee-capstone/entities/web"
	userService "bringeee-capstone/services/user"
	"mime/multipart"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
)

type DriverHandler struct {
	userService *userService.UserService
}

func NewDriverHandler(service *userService.UserService) *DriverHandler {
	return &DriverHandler{
		userService: service,
	}
}

/*
 * User Handler - Create
 * -------------------------------
 * Registrasi User kedalam sistem dan
 * mengembalikan token
 */
func (handler DriverHandler) CreateDriver(c echo.Context) error {

	// Bind request ke user request
	driverReq := entities.CreateDriverRequest{}
	c.Bind(&driverReq)

	// Define links (hateoas)
	links := map[string]string{"self": configs.Get().App.BaseURL + "/api/drivers"}

	// Read files
	files := map[string]*multipart.FileHeader{}
	avatar, _ := c.FormFile("avatar")
	stnk_file, _ := c.FormFile("stnk_file")
	ktp_file, _ := c.FormFile("ktp_file")
	driver_license_file, _ := c.FormFile("driver_license_file")
	vehicle_picture, _ := c.FormFile("vehicle_picture")
	files["avatar"] = avatar
	files["stnk_file"] = stnk_file
	files["ktp_file"] = ktp_file
	files["driver_license_file"] = driver_license_file
	files["vehicle_picture"] = vehicle_picture

	// registrasi user via call user service
	userRes, err := handler.userService.CreateDriver(driverReq, files)
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

func (handler DriverHandler) UpdateDriver(c echo.Context) error {

	// Bind request to user request
	userReq := entities.UpdateDriverRequest{}
	c.Bind(&userReq)

	// Get token
	token := c.Get("user")
	tokenId, role, err := middleware.ReadToken(token)
	links := map[string]string{"self": configs.Get().App.BaseURL + "/api/customers"}
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Code:   http.StatusBadRequest,
			Status: "ERROR",
			Error:  "bad request",
			Links:  links,
		})
	}
	if role != "driver" {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Code:   http.StatusBadRequest,
			Status: "ERROR",
			Error:  "bad request",
			Links:  links,
		})
	}

	// avatar
	files := map[string]*multipart.FileHeader{}
	avatar, _ := c.FormFile("avatar")
	if avatar != nil {
		files["avatar"] = avatar
	}

	// Update via user service call
	userRes, err := handler.userService.UpdateDriver(userReq, tokenId, files)
	if err != nil {
		if reflect.TypeOf(err).String() == "web.WebError" {
			webErr := err.(web.WebError)
			return c.JSON(webErr.Code, web.ErrorResponse{
				Code:   webErr.Code,
				Status: "ERROR",
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

	}

	// response
	return c.JSON(200, web.SuccessResponse{
		Status: "OK",
		Code:   200,
		Error:  nil,
		Links:  links,
		Data:   userRes,
	})
}
