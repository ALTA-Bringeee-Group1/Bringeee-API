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

	"github.com/labstack/echo/v4"
)

type DriverHandler struct {
	userService  *userService.UserService
	orderService orderService.OrderServiceInterface
}

func NewDriverHandler(service *userService.UserService, orderService orderService.OrderServiceInterface) *DriverHandler {
	return &DriverHandler{
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

func (handler DriverHandler) UpdateDriver(c echo.Context) error {

	// Bind request to user request
	userReq := entities.UpdateDriverRequest{}
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
	if role != "driver" {
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
		userReq.Address == "" && userReq.Age == 0 && userReq.DOB == "" &&
		userReq.Email == "" && userReq.Gender == "" && userReq.NIK == "" &&
		userReq.Password == "" && userReq.PhoneNumber == "" && userReq.VehicleIdentifier == "" {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Code:   http.StatusBadRequest,
			Status: "ERROR",
			Error:  "no such data filled",
			Links:  links,
		})
	}

	// Update via user service call
	userRes, err := handler.userService.UpdateDriver(userReq, tokenId, files)
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

/*
 * List Order (Driver)
 * ------------------------------------
 * Mendapatkan list order untuk
 * role spesific tertentu (customer, admin & driver)
 */
func (handler DriverHandler) ListOrders(c echo.Context) error {

	links := map[string]string{}
	userID, role, _ := middleware.ReadToken(c.Get("user"))

	// pagination param
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 50
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
	}
	links["self"] = configs.Get().App.BaseURL + "/api/drivers/orders?page=" + strconv.Itoa(page)

	// Reject if not driver
	if role != "driver" {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Status: "ERROR",
			Code:   http.StatusUnauthorized,
			Error:  "Unauthorized user",
			Links:  links,
		})
	}

	// find userdata driver
	driver, err := handler.userService.FindByDriver("user_id", strconv.Itoa(userID))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Status: "ERROR",
			Code:   http.StatusUnauthorized,
			Error:  "Unauthorized user",
			Links:  links,
		})
	}

	// filters & sorts
	filters := []map[string]interface{}{
		{
			"field":    "truck_type_id",
			"operator": "=",
			"value":    strconv.Itoa(int(driver.TruckTypeID)),
		},
		{
			"field":    "status",
			"operator": "=",
			"value":    "MANIFESTED",
		},
	}
	sorts := []map[string]interface{}{
		{"field": "updated_at", "desc": true},
		{"field": "total_volume", "desc": map[string]bool{"1": true, "0": false}[c.QueryParam("sortVolume")]},
		{"field": "total_weight", "desc": map[string]bool{"1": true, "0": false}[c.QueryParam("sortWeight")]},
		{"field": "distance", "desc": map[string]bool{"1": true, "0": false}[c.QueryParam("sortDistance")]},
	}

	// call order service
	ordersRes, err := handler.orderService.FindAll(limit, 0, filters, sorts)
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
 * Driver - Current order
 * ------------------------------------
 * Mendapatkan detail order customer
 * yang hanya dimiliki customer
 * GET /api/drivers/current_order
 */
func (handler DriverHandler) CurrentOrder(c echo.Context) error {
	userID, role, _ := middleware.ReadToken(c.Get("user"))
	links := map[string]string{}

	// reject if not customer
	if role != "driver" {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Status: "ERROR",
			Code:   http.StatusUnauthorized,
			Error:  "Unauthorized user",
			Links:  links,
		})
	}

	// set self links and filters
	links["self"] = fmt.Sprintf("%s/api/drivers/current_order", configs.Get().App.BaseURL)

	// find userdata driver
	driver, err := handler.userService.FindByDriver("user_id", strconv.Itoa(userID))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Status: "ERROR",
			Code:   http.StatusUnauthorized,
			Error:  "Unauthorized user",
			Links:  links,
		})
	}

	// service call
	orderRes, err := handler.orderService.FindFirst([]map[string]interface{}{
		{"field": "driver_id", "operator": "=", "value": strconv.Itoa(int(driver.ID))},
		{"field": "status", "operator": "=", "value": "ON_PROCESS"},
	})
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	}

	return c.JSON(http.StatusOK, web.SuccessResponse{
		Status: "OK",
		Code:   http.StatusOK,
		Error:  nil,
		Links:  links,
		Data:   orderRes,
	})
}

func (handler DriverHandler) TakeOrder(c echo.Context) error {
	// Get param and token
	id, tx := strconv.Atoi(c.Param("id"))
	links := map[string]string{"self": configs.Get().App.BaseURL + "/api/drivers/orders/" + c.Param("id") + "/take_order"}
	if tx != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Code:   http.StatusBadRequest,
			Status: "ERROR",
			Error:  "invalid parameter",
			Links:  links,
		})
	}
	token := c.Get("user")
	tokenID, role, err := middleware.ReadToken(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Code:   http.StatusUnauthorized,
			Status: "ERROR",
			Error:  "unauthorized",
			Links:  links,
		})
	}
	if role != "driver" {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Code:   http.StatusUnauthorized,
			Status: "ERROR",
			Error:  "unauthorized",
			Links:  links,
		})
	}

	err = handler.orderService.TakeOrder(id, tokenID)
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
			"id": id,
		},
	})
}

func (handler DriverHandler) FinishOrder(c echo.Context) error {
	// Get param and token
	id, tx := strconv.Atoi(c.Param("id"))
	links := map[string]string{"self": configs.Get().App.BaseURL + "/api/drivers/orders/" + c.Param("id") + "/take_order"}
	if tx != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Code:   http.StatusBadRequest,
			Status: "ERROR",
			Error:  "invalid parameter",
			Links:  links,
		})
	}
	token := c.Get("user")
	tokenID, role, err := middleware.ReadToken(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Code:   http.StatusUnauthorized,
			Status: "ERROR",
			Error:  "unauthorized",
			Links:  links,
		})
	}
	if role != "driver" {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Code:   http.StatusUnauthorized,
			Status: "ERROR",
			Error:  "unauthorized",
			Links:  links,
		})
	}
	// Read files
	files := map[string]*multipart.FileHeader{}
	arrived_picture, _ := c.FormFile("arrived_picture")
	files["arrived_picture"] = arrived_picture

	err = handler.orderService.FinishOrder(id, tokenID, files)
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
			"id": id,
		},
	})
}
