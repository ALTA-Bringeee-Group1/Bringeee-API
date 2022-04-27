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
	"strconv"

	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	userService *userService.UserService
}

func NewAdminHandler(service *userService.UserService) *AdminHandler {
	return &AdminHandler{
		userService: service,
	}
}

func (handler AdminHandler) DeleteDriver(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))
	token := c.Get("user")
	_, role, err := middleware.ReadToken(token)
	links := map[string]string{"self": configs.Get().App.BaseURL + "/api/drivers/" + c.Param("id")}
	if role != "admin" {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Code:   http.StatusBadRequest,
			Status: "ERROR",
			Error:  "bad request",
			Links:  links,
		})
	}
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Code:   http.StatusBadRequest,
			Status: "ERROR",
			Error:  "bad request",
			Links:  links,
		})
	}

	// call delete service
	err = handler.userService.DeleteDriver(id)
	if err != nil {
		if reflect.TypeOf(err).String() == "web.WebError" {
			webErr := err.(web.WebError)
			return c.JSON(webErr.Code, web.ErrorResponse{
				Code:   webErr.Code,
				Status: "ERROR",
				Error:  webErr.Error(),
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
		Data: map[string]interface{}{
			"id": id,
		},
	})
}

func (handler AdminHandler) UpdateDriverByAdmin(c echo.Context) error {

	// Bind request to user request
	userReq := entities.UpdateDriverByAdminRequest{}
	c.Bind(&userReq)
	id, _ := strconv.Atoi(c.Param("id"))

	// Get token
	token := c.Get("user")
	_, role, err := middleware.ReadToken(token)
	links := map[string]string{"self": configs.Get().App.BaseURL + "/api/drivers/" + c.Param("id")}
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Code:   http.StatusBadRequest,
			Status: "ERROR",
			Error:  "bad request",
			Links:  links,
		})
	}
	if role != "admin" {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Code:   http.StatusBadRequest,
			Status: "ERROR",
			Error:  "bad request",
			Links:  links,
		})
	}

	// avatar
	files := map[string]*multipart.FileHeader{}
	stnk_file, _ := c.FormFile("stnk_file")
	if stnk_file != nil {
		files["stnk_file"] = stnk_file
	}

	// Update via user service call
	userRes, err := handler.userService.UpdateDriverByAdmin(userReq, id, files)
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
