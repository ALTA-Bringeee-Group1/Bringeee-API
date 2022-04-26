package handlers

import (
	"bringeee-capstone/configs"
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

func NewDriverHandler(service *userService.UserService) *UserHandler {
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
func (handler UserHandler) CreateDriver(c echo.Context) error {

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
	driver_lisence_file, _ := c.FormFile("driver_lisence_file")
	vehicle_picture, _ := c.FormFile("vehicle_picture")
	files["avatar"] = avatar
	files["stnk_file"] = stnk_file
	files["ktp_file"] = ktp_file
	files["driver_lisence_file"] = driver_lisence_file
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
