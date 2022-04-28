package handlers

import (
	"bringeee-capstone/configs"
	"bringeee-capstone/deliveries/helpers"
	middleware "bringeee-capstone/deliveries/middlewares"
	"bringeee-capstone/entities"
	"bringeee-capstone/entities/web"
	orderService "bringeee-capstone/services/order"
	userService "bringeee-capstone/services/user"
	"fmt"
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService *userService.UserService
	orderService orderService.OrderServiceInterface
}

func NewUserHandler(service *userService.UserService,orderService orderService.OrderServiceInterface) *UserHandler {
	return &UserHandler{
		userService: service,
		orderService: orderService,
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

	// Read files
	files := map[string]*multipart.FileHeader{}
	avatar, _ := c.FormFile("avatar")
	if avatar != nil {
		files["avatar"] = avatar
	}

	// registrasi user via call user service
	userRes, err := handler.userService.CreateCustomer(userReq, files)
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

func (handler UserHandler) UpdateCustomer(c echo.Context) error {

	// Bind request to user request
	userReq := entities.UpdateCustomerRequest{}
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
	if role == "driver" {
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
	userRes, err := handler.userService.UpdateCustomer(userReq, tokenId, files)
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

func (handler UserHandler) DeleteCustomer(c echo.Context) error {

	token := c.Get("user")
	tokenId, role, err := middleware.ReadToken(token)
	links := map[string]string{"self": configs.Get().App.BaseURL + "/api/customers"}
	if role == "driver" {
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
	err = handler.userService.DeleteCustomer(tokenId)
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
			"id": tokenId,
		},
	})
}


/* 
 * Customer - Detail Order - Get Histories
 * ---------------------------------
 * List history tracking dari satu detail order tunggal
 * GET /api/customers/orders/{orderID}/histories
 */
func (handler UserHandler) DetailOrderHistory(c echo.Context) error {
	userID, _, _ := middleware.ReadToken(c.Get("user"))
	links := map[string]string{}
	orderID, err := strconv.Atoi(c.Param("orderID"))
	links["self"] = fmt.Sprintf("%s/api/customer/orders/%s/histories", configs.Get().App.BaseURL, c.Param("orderID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Status: "ERROR",
			Code: 400,
			Error: "Order ID parameter is invalid",
			Links: links,
		})
	}

	// get single order
	order, err := handler.orderService.Find(orderID)
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	}

	// Reject if order doesn't belong to currently authenticated user
	if order.CustomerID != uint(userID) {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Status: "ERROR",
			Code: http.StatusUnauthorized,
			Error: "order doesn't belong to currently authenticated customer",
		})
	}

	// Get order tracking histories
	histories, err := handler.orderService.FindAllHistory(orderID, []map[string]interface{}{
		{ "field": "created_at", "desc": true },
	})
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	}

	return c.JSON(http.StatusOK, web.SuccessResponse{
		Status: "OK",
		Code: http.StatusOK,
		Error: nil,
		Links: links,
		Data: histories,
	})
}