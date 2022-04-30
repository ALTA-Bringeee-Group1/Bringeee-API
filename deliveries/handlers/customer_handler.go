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
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type CustomerHandler struct {
	userService  *userService.UserService
	orderService orderService.OrderServiceInterface
}

func NewCustomerHandler(service *userService.UserService, orderService orderService.OrderServiceInterface) *CustomerHandler {
	return &CustomerHandler{
		userService:  service,
		orderService: orderService,
	}
}

/*
 * User Handler - Create
 * -------------------------------
 * Registrasi User kedalam sistem dan
 * mengembalikan token
 */
func (handler CustomerHandler) CreateCustomer(c echo.Context) error {

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
		return helpers.WebErrorResponse(c, err, links)
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

func (handler CustomerHandler) UpdateCustomer(c echo.Context) error {

	// Bind request to user request
	userReq := entities.UpdateCustomerRequest{}
	c.Bind(&userReq)

	// Get token
	token := c.Get("user")
	tokenId, role, err := middleware.ReadToken(token)
	links := map[string]string{"self": configs.Get().App.BaseURL + "/api/customers"}
	if err != nil {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Code:   http.StatusUnauthorized,
			Status: "ERROR",
			Error:  "unauthorized",
			Links:  links,
		})
	}
	if role == "driver" {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Code:   http.StatusUnauthorized,
			Status: "ERROR",
			Error:  "unauthorized",
			Links:  links,
		})
	}

	// avatar
	files := map[string]*multipart.FileHeader{}
	avatar, _ := c.FormFile("avatar")
	if avatar != nil {
		files["avatar"] = avatar
	}

	if len(files) == 0 && userReq.Name == "" &&
		userReq.Address == "" && userReq.DOB == "" &&
		userReq.Email == "" && userReq.Gender == "" &&
		userReq.Password == "" && userReq.PhoneNumber == "" {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Code:   http.StatusBadRequest,
			Status: "ERROR",
			Error:  "no such data filled",
			Links:  links,
		})
	}
	// Update via user service call
	userRes, err := handler.userService.UpdateCustomer(userReq, tokenId, files)
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
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

func (handler CustomerHandler) DeleteCustomer(c echo.Context) error {

	token := c.Get("user")
	tokenId, role, err := middleware.ReadToken(token)
	links := map[string]string{"self": configs.Get().App.BaseURL + "/api/customers"}
	if role == "driver" {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Code:   http.StatusUnauthorized,
			Status: "ERROR",
			Error:  "unauthorized",
			Links:  links,
		})
	}
	if err != nil {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Code:   http.StatusUnauthorized,
			Status: "ERROR",
			Error:  "unauthorized",
			Links:  links,
		})
	}

	// call delete service
	err = handler.userService.DeleteCustomer(tokenId)
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
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
 * Customer - List Order
 * ------------------------------------
 * Mendapatkan list order berdasarkan
 * query param yang telah ditentukan
 */
func (handler CustomerHandler) ListOrders(c echo.Context) error {

	userID, role, _ := middleware.ReadToken(c.Get("user"))
	links := map[string]string{}
	filters := []map[string]interface{}{}

	// pagination param
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 50
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
	}

	// reject if not customer
	if role != "customer" {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Status: "ERROR",
			Code:   http.StatusUnauthorized,
			Error:  "Unauthorized user",
			Links:  links,
		})
	}

	// Multi status filters
	statusQuery := c.QueryParam("status")
	if statusQuery != "" {
		statusArr := strings.Split(statusQuery, ",")
		statusFilters := []map[string]string{}
		for _, status := range statusArr {
			statusFilters = append(statusFilters, map[string]string{
				"field":    "status",
				"operator": "=",
				"value":    status,
			})
		}
		filters = append(filters, map[string]interface{}{
			"or": statusFilters,
		})
	}

	// set self links and filters
	links["self"] = configs.Get().App.BaseURL + "/api/customers/orders?page=" + strconv.Itoa(page)
	filters = append(filters, map[string]interface{}{
		"field":    "customer_id",
		"operator": "=",
		"value":    strconv.Itoa(userID),
	})

	// get authenticated userdata
	_, err = handler.userService.FindCustomer(userID)
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	}

	// call order service
	ordersRes, err := handler.orderService.FindAll(0, 0, filters, []map[string]interface{}{
		{"field": "updated_at", "desc": true},
	})
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	}

	// make pagination data & formatting pagination links
	paginationRes, err := handler.orderService.GetPagination(page, limit, filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, web.ErrorResponse{
			Status: "ERROR",
			Code:   http.StatusInternalServerError,
			Error:  err.Error(),
			Links:  links,
		})
	}
	pageUrl := fmt.Sprintf("%s/api/customers/orders?page=", configs.Get().App.BaseURL)
	links["first"] = pageUrl + "1"
	links["last"] = pageUrl + strconv.Itoa(paginationRes.TotalPages)
	if paginationRes.Page > 1 {
		links["previous"] = pageUrl + strconv.Itoa(page-1)
	}
	if paginationRes.Page < paginationRes.TotalPages {
		links["previous"] = pageUrl + strconv.Itoa(page+1)
	}

	// Success list response
	return c.JSON(http.StatusOK, web.SuccessListResponse{
		Status:     "OK",
		Error:      nil,
		Code:       http.StatusOK,
		Links:      links,
		Data:       ordersRes,
		Pagination: paginationRes,
	})
}

/*
 * Customer - Detail Order
 * ------------------------------------
 * Mendapatkan detail order customer
 * yang hanya dimiliki customer
 */
func (handler CustomerHandler) DetailOrder(c echo.Context) error {
	userID, role, _ := middleware.ReadToken(c.Get("user"))
	links := map[string]string{}
	orderRes := entities.OrderResponse{}

	// orderID param
	orderID, err := strconv.Atoi(c.Param("orderID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Status: "OK",
			Code:   http.StatusBadRequest,
			Error:  "Order ID parameter is invalid",
			Links:  links,
		})
	}

	// reject if not customer
	if role != "customer" {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Status: "ERROR",
			Code:   http.StatusUnauthorized,
			Error:  "Unauthorized user",
			Links:  links,
		})
	}

	// set self links and filters
	links["self"] = fmt.Sprintf("%s/api/customers/orders/%s", configs.Get().App.BaseURL, strconv.Itoa(orderID))

	// get authenticated userdata
	customer, err := handler.userService.FindCustomer(userID)
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	}

	// call service order
	orderRes, err = handler.orderService.Find(orderID)
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	}

	// reject if customer_id doesn't match
	if orderRes.CustomerID != customer.ID {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Status: "ERROR",
			Code:   http.StatusBadRequest,
			Error:  "order doesn't belong to currently authenticated customer",
			Links:  links,
		})
	}

	return c.JSON(http.StatusOK, web.SuccessResponse{
		Status: "OK",
		Code:   http.StatusOK,
		Error:  nil,
		Links:  links,
		Data:   orderRes,
	})
}

/*
 * Customer - Detail Order - Get Histories
 * ---------------------------------
 * List history tracking dari satu detail order tunggal
 * GET /api/customers/orders/{orderID}/histories
 */
func (handler CustomerHandler) DetailOrderHistory(c echo.Context) error {
	userID, _, _ := middleware.ReadToken(c.Get("user"))
	links := map[string]string{}
	orderID, err := strconv.Atoi(c.Param("orderID"))
	links["self"] = fmt.Sprintf("%s/api/customer/orders/%s/histories", configs.Get().App.BaseURL, c.Param("orderID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Status: "ERROR",
			Code:   400,
			Error:  "Order ID parameter is invalid",
			Links:  links,
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
			Code:   http.StatusUnauthorized,
			Error:  "order doesn't belong to currently authenticated customer",
		})
	}

	// Get order tracking histories
	histories, err := handler.orderService.FindAllHistory(orderID, []map[string]interface{}{
		{"field": "created_at", "desc": true},
	})
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	}

	return c.JSON(http.StatusOK, web.SuccessResponse{
		Status: "OK",
		Code:   http.StatusOK,
		Error:  nil,
		Links:  links,
		Data:   histories,
	})
}

/*
 * Customer - Create Order
 * ---------------------------------
 * Create order
 * GET /api/customers/orders
 */
func (handler CustomerHandler) CreateOrder(c echo.Context) error {
	links := map[string]string{}
	links["self"] = fmt.Sprintf("%s/api/customers/orders", configs.Get().App.BaseURL)

	filesReq := map[string]*multipart.FileHeader{}
	userID, role, err := middleware.ReadToken(c.Get("user"))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Status: "ERROR",
			Code:   http.StatusUnauthorized,
			Error:  "Unauthorized user",
			Links:  links,
		})
	}

	// check authenticated user
	if role != "customer" {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Status: "ERROR",
			Code:   http.StatusUnauthorized,
			Error:  "Unauthorized user",
			Links:  links,
		})
	}

	// populate request
	orderReq := entities.CustomerCreateOrderRequest{}
	c.Bind(&orderReq)
	filesReq["order_picture"], _ = c.FormFile("order_picture")

	orderRes, err := handler.orderService.Create(orderReq, filesReq, userID)
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	}

	return c.JSON(http.StatusOK, web.SuccessResponse{
		Status: "OK",
		Error:  nil,
		Code:   http.StatusOK,
		Links:  links,
		Data: map[string]interface{}{
			"id": orderRes.ID,
		},
	})
}

func (handler CustomerHandler) ConfirmOrder(c echo.Context) error {
	links := map[string]string{}
	orderID, err := strconv.Atoi(c.Param("orderID"))
	links["self"] = fmt.Sprintf("%s/api/customers/orders/%s", configs.Get().App.BaseURL, c.Param("orderID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Status: "ERROR",
			Code:   http.StatusBadRequest,
			Error:  "invalid parameter",
			Links:  links,
		})
	}
	tokenID, role, err := middleware.ReadToken(c.Get("user"))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Status: "ERROR",
			Code:   http.StatusUnauthorized,
			Error:  "unauthorized",
			Links:  links,
		})
	}

	if role != "customer" {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Status: "ERROR",
			Code:   http.StatusUnauthorized,
			Error:  "unauthorized",
			Links:  links,
		})
	}
	err = handler.orderService.ConfirmOrder(orderID, tokenID, map[string]bool{"admin": true, "": false}[role])
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	}

	return c.JSON(http.StatusOK, web.SuccessResponse{
		Status: "OK",
		Code:   http.StatusOK,
		Error:  nil,
		Links:  links,
		Data: map[string]interface{}{
			"id": orderID,
		},
	})
}

func (handler CustomerHandler) CancelOrder(c echo.Context) error {
	links := map[string]string{}
	orderID, err := strconv.Atoi(c.Param("orderID"))
	links["self"] = fmt.Sprintf("%s/api/customers/orders/%s/cancel", configs.Get().App.BaseURL, c.Param("orderID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Status: "ERROR",
			Code:   http.StatusBadRequest,
			Error:  "invalid parameter",
			Links:  links,
		})
	}
	tokenID, role, err := middleware.ReadToken(c.Get("user"))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Status: "ERROR",
			Code:   http.StatusUnauthorized,
			Error:  "unauthorized",
			Links:  links,
		})
	}

	if role != "customer" {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Status: "ERROR",
			Code:   http.StatusUnauthorized,
			Error:  "unauthorized",
			Links:  links,
		})
	}

	err = handler.orderService.CancelOrder(orderID, tokenID, map[string]bool{"admin": true, "": false}[role])
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	}

	return c.JSON(http.StatusOK, web.SuccessResponse{
		Status: "OK",
		Code:   http.StatusOK,
		Error:  nil,
		Links:  links,
		Data: map[string]interface{}{
			"id": orderID,
		},
	})
}
/* 
 * Customer - Create payment
 * ---------------------------------
 * Membuat detail pembayaran agar order untuk dibayarkan
 * POST /api/customers/orders/{orderID}/payment
 */
func (handler CustomerHandler) CreatePayment(c echo.Context) error {
	orderID, err := strconv.Atoi(c.Param("orderID"))
	links := map[string]string{}
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Status: "ERROR",
			Code: http.StatusBadRequest,
			Error: "Invalid order id parameter",
			Links: links,
		})
	}
	links["self"] = fmt.Sprintf("%s/api/customers/orders/%s/payment", configs.Get().App.BaseURL, c.Param("orderID"))

	// check if authenticated user is a customer
	userID, role, err := middleware.ReadToken(c.Get("user"))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Status: "ERROR",
			Code: http.StatusUnauthorized,
			Error: "Unauthorized user",
			Links: links,
		})
	}
	if role != "customer" {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Status: "ERROR",
			Code: http.StatusUnauthorized,
			Error: "Unauthorized user",
			Links: links,
		})
	}

	// reject if order belong to someone else
	order, err := handler.orderService.Find(orderID)
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	}
	if int(order.CustomerID) != userID {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Status: "ERROR",
			Code: http.StatusUnauthorized,
			Error: "This order doesn't belong to currently authenticated user",
			Links: links,
		})
	}

	// service payment action
	paymentReq := entities.CreatePaymentRequest{}
	c.Bind(&paymentReq)
	paymentRes, err := handler.orderService.CreatePayment(orderID, paymentReq)
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	}
	return c.JSON(http.StatusOK, web.SuccessResponse{
		Status: "OK",
		Code: http.StatusOK,
		Error: nil,
		Links: links,
		Data: paymentRes,
	})
}

/* 
 * Customer - Get payment detail
 * ---------------------------------
 * Mendapatkan detail pembayaran untuk order yang sudah dibuatkan payment
 * GET /api/customers/orders/{orderID}/payment
 */
func (handler CustomerHandler) GetPayment(c echo.Context) error {
	orderID, err := strconv.Atoi(c.Param("orderID"))
	links := map[string]string{}
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Status: "ERROR",
			Code: http.StatusBadRequest,
			Error: "Invalid order id parameter",
			Links: links,
		})
	}
	links["self"] = fmt.Sprintf("%s/api/customers/orders/%s/payment", configs.Get().App.BaseURL, c.Param("orderID"))

	// check if authenticated user is a customer
	userID, role, err := middleware.ReadToken(c.Get("user"))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Status: "ERROR",
			Code: http.StatusUnauthorized,
			Error: "Unauthorized user",
			Links: links,
		})
	}
	if role != "customer" {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Status: "ERROR",
			Code: http.StatusUnauthorized,
			Error: "Unauthorized user",
			Links: links,
		})
	}
	// reject if order belong to someone else
	order, err := handler.orderService.Find(orderID)
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	}
	if int(order.CustomerID) != userID {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Status: "ERROR",
			Code: http.StatusUnauthorized,
			Error: "This order doesn't belong to currently authenticated user",
			Links: links,
		})
	}
	// service payment action
	paymentRes, err := handler.orderService.GetPayment(orderID)
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	}
	return c.JSON(http.StatusOK, web.SuccessResponse{
		Status: "OK",
		Code: http.StatusOK,
		Error: nil,
		Links: links,
		Data: paymentRes,
	})
}


/* 
 * Customer - Cancel payment order
 * ---------------------------------
 * Mendapatkan detail pembayaran untuk order yang sudah dibuatkan payment
 * GET /api/customers/orders/{orderID}/payment/cancel
 */
func (handler CustomerHandler) CancelPayment(c echo.Context) error {
	orderID, err := strconv.Atoi(c.Param("orderID"))
	links := map[string]string{}
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Status: "ERROR",
			Code: http.StatusBadRequest,
			Error: "Invalid order id parameter",
			Links: links,
		})
	}
	links["self"] = fmt.Sprintf("%s/api/customers/orders/%s/payment/cancel", configs.Get().App.BaseURL, c.Param("orderID"))

	// check if authenticated user is a customer
	userID, role, err := middleware.ReadToken(c.Get("user"))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Status: "ERROR",
			Code: http.StatusUnauthorized,
			Error: "Unauthorized user",
			Links: links,
		})
	}
	if role != "customer" {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Status: "ERROR",
			Code: http.StatusUnauthorized,
			Error: "Unauthorized user",
			Links: links,
		})
	}
	// reject if order belong to someone else
	order, err := handler.orderService.Find(orderID)
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	}
	if int(order.CustomerID) != userID {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Status: "ERROR",
			Code: http.StatusUnauthorized,
			Error: "This order doesn't belong to currently authenticated user",
			Links: links,
		})
	}
	// service payment action
	err = handler.orderService.CancelPayment(orderID)
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	}
	return c.JSON(http.StatusOK, web.SuccessResponse{
		Status: "OK",
		Code: http.StatusOK,
		Error: nil,
		Links: links,
		Data: map[string]interface{}{
			"id": orderID,
		},
	})
}
