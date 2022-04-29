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

type AdminHandler struct {
	userService *userService.UserService
	orderService orderService.OrderServiceInterface
}

func NewAdminHandler(service *userService.UserService, orderService orderService.OrderServiceInterface) *AdminHandler {
	return &AdminHandler{
		userService: service,
		orderService: orderService,
	}
}

func (handler AdminHandler) DeleteDriver(c echo.Context) error {

	id, tx := strconv.Atoi(c.Param("id"))
	token := c.Get("user")
	_, role, err := middleware.ReadToken(token)
	links := map[string]string{"self": configs.Get().App.BaseURL + "/api/drivers/" + c.Param("id")}
	if tx != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Code:   http.StatusBadRequest,
			Status: "ERROR",
			Error:  "invalid parameter",
			Links:  links,
		})
	}
	if role != "admin" {
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
	err = handler.userService.DeleteDriver(id)
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

func (handler AdminHandler) UpdateDriverByAdmin(c echo.Context) error {

	// Bind request to user request
	userReq := entities.UpdateDriverByAdminRequest{}
	c.Bind(&userReq)
	id, tx := strconv.Atoi(c.Param("id"))

	// Get token
	token := c.Get("user")
	_, role, err := middleware.ReadToken(token)
	links := map[string]string{"self": configs.Get().App.BaseURL + "/api/drivers/" + c.Param("id")}
	if tx != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Code:   http.StatusBadRequest,
			Status: "ERROR",
			Error:  "invalid parameter",
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
	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Code:   http.StatusUnauthorized,
			Status: "ERROR",
			Error:  "unauthorized",
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

func (handler AdminHandler) GetAllDriver(c echo.Context) error {

	// Get token
	token := c.Get("user")
	_, role, err := middleware.ReadToken(token)

	// Translate query param to map of filters
	filters := []map[string]string{}
	name := c.QueryParam("name")
	if name != "" {
		filters = append(filters, map[string]string{
			"field":    "name",
			"operator": "=",
			"value":    name,
		})
	}
	gender := c.QueryParam("gender")
	if gender != "" {
		filters = append(filters, map[string]string{
			"field":    "gender",
			"operator": "=",
			"value":    gender,
		})
	}
	status := c.QueryParam("status")
	if status != "" {
		filters = append(filters, map[string]string{
			"field":    "status",
			"operator": "=",
			"value":    status,
		})
	}
	account_status := c.QueryParam("account_status")
	if account_status != "" {
		filters = append(filters, map[string]string{
			"field":    "account_status",
			"operator": "=",
			"value":    account_status,
		})
	}
	truck_type := c.QueryParam("truck_type")
	if truck_type != "" {
		filters = append(filters, map[string]string{
			"field":    "truck_type_id",
			"operator": "=",
			"value":    truck_type,
		})
	}

	// Sort parameter
	sorts := []map[string]interface{}{}
	sortName := c.QueryParam("sortName")
	sorts = append(sorts, map[string]interface{}{
		"field": "name",
		"desc":  map[string]bool{"1": true, "0": false}[sortName],
	})

	sortAge := c.QueryParam("sortAge")
	if sortAge != "" {
		switch sortAge {
		case "1":
			sorts = append(sorts, map[string]interface{}{
				"field": "age",
				"desc":  true,
			})
		case "0":
			sorts = append(sorts, map[string]interface{}{
				"field": "age",
				"desc":  false,
			})
		}
	}
	links := map[string]string{"self": configs.Get().App.BaseURL + "/api/drivers?limit=" + c.QueryParam("limit") + "&page=" + c.QueryParam("page")}
	if err != nil {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Code:   http.StatusUnauthorized,
			Status: "ERROR",
			Error:  "unauthorized",
			Links:  links,
		})
	}
	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Code:   http.StatusUnauthorized,
			Status: "ERROR",
			Error:  "unauthorized",
			Links:  links,
		})
	}
	// pagination param
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return c.JSON(400, helpers.MakeErrorResponse("ERROR", 400, "Limit Parameter format is invalid", links))
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		links := map[string]string{"self": configs.Get().App.BaseURL}
		return c.JSON(400, helpers.MakeErrorResponse("ERROR", 400, "page Parameter format is invalid", links))
	}
	links["self"] = configs.Get().App.BaseURL + "/api/drivers?limit=" + c.QueryParam("limit") + "&page=" + c.QueryParam("page")

	// Get all drivers
	driversRes, err := handler.userService.FindAllDriver(limit, page, filters, sorts)
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	}

	// Get pagination data
	pagination, err := handler.userService.GetPaginationDriver(limit, page, filters)
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	}

	links["first"] = configs.Get().App.BaseURL + "/api/drivers?limit=" + c.QueryParam("limit") + "&page=1"
	links["last"] = configs.Get().App.BaseURL + "/api/drivers?limit=" + c.QueryParam("limit") + "&page=" + strconv.Itoa(pagination.TotalPages)
	if pagination.Page > 1 {
		links["prev"] = configs.Get().App.BaseURL + "/api/drivers?limit=" + c.QueryParam("limit") + "&page=" + strconv.Itoa(pagination.Page-1)
	}
	if pagination.Page < pagination.TotalPages {
		links["next"] = configs.Get().App.BaseURL + "/api/drivers?limit=" + c.QueryParam("limit") + "&page=" + strconv.Itoa(pagination.Page+1)
	}

	// success response
	return c.JSON(200, web.SuccessListResponse{
		Status:     "OK",
		Code:       200,
		Error:      nil,
		Links:      links,
		Data:       driversRes,
		Pagination: pagination,
	})
}

func (handler AdminHandler) GetSingleDriver(c echo.Context) error {
	// Get param and token
	id, tx := strconv.Atoi(c.Param("id"))
	links := map[string]string{"self": configs.Get().App.BaseURL + "/api/drivers/" + c.Param("id")}
	if tx != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Code:   http.StatusBadRequest,
			Status: "ERROR",
			Error:  "invalid parameter",
			Links:  links,
		})
	}
	token := c.Get("user")
	_, role, err := middleware.ReadToken(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Code:   http.StatusUnauthorized,
			Status: "ERROR",
			Error:  "unauthorized",
			Links:  links,
		})
	}
	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Code:   http.StatusUnauthorized,
			Status: "ERROR",
			Error:  "unauthorized",
			Links:  links,
		})
	}

	// Get eventdata
	driver, err := handler.userService.FindDriver(id)
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	}
	// response
	return c.JSON(200, web.SuccessResponse{
		Status: "OK",
		Code:   200,
		Error:  nil,
		Links:  links,
		Data:   driver,
	})
}

func (handler AdminHandler) GetSingleCustomer(c echo.Context) error {
	// Get param and token
	id, tx := strconv.Atoi(c.Param("id"))
	links := map[string]string{"self": configs.Get().App.BaseURL + "/api/customers/" + c.Param("id")}
	if tx != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Code:   http.StatusBadRequest,
			Status: "ERROR",
			Error:  "invalid parameter",
			Links:  links,
		})
	}
	token := c.Get("user")
	_, role, err := middleware.ReadToken(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Code:   http.StatusUnauthorized,
			Status: "ERROR",
			Error:  "unauthorized",
			Links:  links,
		})
	}
	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Code:   http.StatusUnauthorized,
			Status: "ERROR",
			Error:  "unauthorized",
			Links:  links,
		})
	}

	// Get eventdata
	user, err := handler.userService.FindCustomer(id)
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	} else if user.ID == 0 {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Code:   http.StatusBadRequest,
			Status: "ERROR",
			Error:  "bad request",
			Links:  links,
		})
	}
	// response
	return c.JSON(200, web.SuccessResponse{
		Status: "OK",
		Code:   200,
		Error:  nil,
		Links:  links,
		Data:   user,
	})
}


/* 
 * Admin - Detail Order - Get Histories
 * ---------------------------------
 * List history tracking dari satu detail order tunggal
 * GET /api/orders/{orderID}/histories
 */
func (handler AdminHandler) DetailOrderHistory(c echo.Context) error {
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

	_, err = handler.orderService.Find(orderID)
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
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

/*
 * List Order (Admin)
 * ------------------------------------
 * Mendapatkan list order berdasarkan query param tersedia
 * GET /api/orders
 */
func (handler AdminHandler) ListOrders(c echo.Context) error {
	
	links := map[string]string{}
	_, role, _ := middleware.ReadToken(c.Get("user"))

	// pagination param
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 50
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
	}
	links["self"] = configs.Get().App.BaseURL + "/api/orders?page=" + strconv.Itoa(page)

	// Reject if not admin
	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Status: "ERROR",
			Code: http.StatusUnauthorized,
			Error: "Unauthorized user",
			Links: links,
		})
	}

	// filters & sorts
	filters := []map[string]interface{} {}
	sorts := []map[string]interface{}{
		{ "field": "total_volume", 	"desc": map[string]bool{"1": true, "0": false}[c.QueryParam("sortVolume")]},
		{ "field": "total_weight", 	"desc": map[string]bool{"1": true, "0": false}[c.QueryParam("sortWeight")]},
		{ "field": "distance", 		"desc": map[string]bool{"1": true, "0": false}[c.QueryParam("sortDistance")]},
	}

	// Multi status filters
	statusQuery := c.QueryParam("status")
	if statusQuery != "" {
		statusArr := strings.Split(statusQuery, ",")
		statusFilters := []map[string]string{}
		for _, status := range statusArr {
			statusFilters = append(statusFilters, map[string]string{
				"field": "status", 
				"operator": "=", 
				"value": status,
			})
		}
		filters = append(filters, map[string]interface{}{
			"or": statusFilters,
		})
	}
	truckTypeID := c.QueryParam("truck_type")
	if truckTypeID != "" {
		filters = append(filters, map[string]interface{}{ "field": "truck_type_id", "operator": "=", "value": truckTypeID })
	}

	// call order service
	ordersRes, err := handler.orderService.FindAll(limit, page, filters, sorts)
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	}

	// make pagination data & formatting pagination links
	paginationRes, err := handler.orderService.GetPagination(page, limit, filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, web.ErrorResponse{
			Status: "ERROR",
			Code: http.StatusInternalServerError,
			Error: err.Error(),
			Links: links,
		})
	}
	pageUrl := fmt.Sprintf("%s/api/orders?page=", configs.Get().App.BaseURL)
	links["first"] = pageUrl + "1"
	links["last"] = pageUrl + strconv.Itoa(paginationRes.TotalPages)
	if paginationRes.Page > 1 {
		links["previous"] = pageUrl + strconv.Itoa(page - 1)
	}
	if paginationRes.Page < paginationRes.TotalPages {
		links["previous"] = pageUrl + strconv.Itoa(page + 1)
	}

	// Success list response
	return c.JSON(http.StatusOK, web.SuccessListResponse{
		Status: "OK",
		Error: nil,
		Code: http.StatusOK,
		Links: links,
		Data: ordersRes,
		Pagination: paginationRes,
	})
}

/*
 * Admin - Detail order
 * ------------------------------------
 * Mendapatkan detail order
 * GET /api/order/{orderID}
 */
func (handler AdminHandler) DetailOrder(c echo.Context) error {
	links := map[string]string{}
	_, role, err  := middleware.ReadToken(c.Get("user"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Status: "OK",
			Code: http.StatusBadRequest,
			Error: "Order ID parameter is invalid",
			Links: links,
		})
	}

	// orderID param
	orderID, err := strconv.Atoi(c.Param("orderID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Status: "OK",
			Code: http.StatusBadRequest,
			Error: "Order ID parameter is invalid",
			Links: links,
		})
	}

	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Status: "ERROR",
			Code: http.StatusUnauthorized,
			Error: "Unauthorized user",
			Links: links,
		})
	}

	// set self links and filters
	links["self"] = fmt.Sprintf("%s/api/orders/%s", configs.Get().App.BaseURL,strconv.Itoa(orderID))
	
	// call service order
	orderRes, err := handler.orderService.Find(orderID)
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	}

	return c.JSON(http.StatusOK, web.SuccessResponse{
		Status: "OK",
		Code: http.StatusOK,
		Error: nil,
		Links: links,
		Data: orderRes,
	})
}

/*
 * Admin - Set fixed price
 * ------------------------------------
 * Mendapatkan detail order
 * PATCH /api/order/{orderID}
 */
func (handler AdminHandler) SetFixedPrice(c echo.Context) error {
	links := map[string]string{}
	orderID, err := strconv.Atoi(c.Param("orderID"))
	links["self"] = fmt.Sprintf("%s/api/orders/%s", configs.Get().App.BaseURL, c.Param("orderID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Status: "ERROR",
			Code: http.StatusBadRequest,
			Error: "Invalid order id parameter",
			Links: links,
		})
	}
	_, role, err := middleware.ReadToken(c.Get("user"))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Status: "ERROR",
			Code: http.StatusUnauthorized,
			Error: "Unauthorized user",
			Links: links,
		})
	}

	// reject if not admin
	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Status: "ERROR",
			Code: http.StatusUnauthorized,
			Error: "Unauthorized user",
			Links: links,
		})
	}

	// service action
	setPriceReq := entities.AdminSetPriceOrderRequest{}
	c.Bind(&setPriceReq)
	err = handler.orderService.SetFixOrder(orderID, setPriceReq)
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	}

	return c.JSON(http.StatusOK, web.SuccessResponse{
		Status: "OK",
		Code: http.StatusOK,
		Error: nil,
		Links: links,
		Data: map[string]interface{} {
			"id": orderID,
		},
	})
}