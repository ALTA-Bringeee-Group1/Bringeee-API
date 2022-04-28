package handlers

import (
	"bringeee-capstone/configs"
	"bringeee-capstone/deliveries/helpers"
	middleware "bringeee-capstone/deliveries/middlewares"
	"bringeee-capstone/entities"
	"bringeee-capstone/entities/web"
	userRepository "bringeee-capstone/repositories/user"
	orderService "bringeee-capstone/services/order"
	userService "bringeee-capstone/services/user"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	orderService orderService.OrderServiceInterface
	userService userService.UserServiceInterface
	userRepository userRepository.UserRepositoryInterface
}

func NewOrderHandler(
	service orderService.OrderServiceInterface,
	userService userService.UserServiceInterface,
	userRepository userRepository.UserRepositoryInterface,
) *OrderHandler {
	return &OrderHandler{
		orderService: service,
		userService: userService,
		userRepository: userRepository,
	}
}

/*
 * List Order (Customer, Admin, Driver)
 * ------------------------------------
 * Mendapatkan list order untuk
 * role spesific tertentu (customer, admin & driver)
 */
func (handler OrderHandler) Index(c echo.Context) error {
	
	userID, role, _ := middleware.ReadToken(c.Get("user"))
	links := map[string]string{}
	filters := []map[string]interface{} {}
	ordersRes := []entities.OrderResponse{}

	// pagination param
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 50
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
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

	switch role {
	case "customer":
		
		// set self links and filters
		links["self"] = configs.Get().App.BaseURL + "/api/customers/orders?page=" + strconv.Itoa(page)
		filters = append(filters, map[string]interface{}{
			"field": "customer_id", 
			"operator": "=", 
			"value": strconv.Itoa(userID),
		})

		// get authenticated userdata
		_, err := handler.userService.FindCustomer(userID)
		if err != nil {
			return helpers.WebErrorResponse(c, err, links)
		}
		
		// call order service
		ordersRes, err = handler.orderService.FindAll(0, 0, filters, []map[string]interface{}{
			{ "field": "updated_at", "desc": true },
		})
		if err != nil {
			return helpers.WebErrorResponse(c, err, links)
		}
	case "driver":
		// find userdata driver
		driver, err := handler.userRepository.FindByDriver("user_id", strconv.Itoa(userID))
		if err != nil {
			return c.JSON(http.StatusUnauthorized, web.ErrorResponse{ 
				Status: "ERROR", 
				Code: http.StatusUnauthorized,
				Error: "Unauthorized user",  
				Links: links,
			})
		}

		// set self links and filters
		links["self"] = configs.Get().App.BaseURL + "/api/drivers/orders?page=" + strconv.Itoa(page)
		filters = append(filters, map[string]interface{}{
			"field": "truck_type_id", 
			"operator": "=", 
			"value": driver.ID,
		})
		filters = append(filters, map[string]interface{}{
			"field": "status", 
			"operator": "=", 
			"value": "MANIFESTED",
		})

		// sorts
		sorts := []map[string]interface{}{
			{ "field": "updated_at", "desc": true },
		}
		sorts = append(sorts, map[string]interface{}{"field": "total_volume", 	"desc": map[string]bool{"1": true, "0": false}[c.QueryParam("sortVolume")]})
		sorts = append(sorts, map[string]interface{}{"field": "total_weight", 	"desc": map[string]bool{"1": true, "0": false}[c.QueryParam("sortWeight")]})
		sorts = append(sorts, map[string]interface{}{"field": "total_distance", "desc": map[string]bool{"1": true, "0": false}[c.QueryParam("sortDistance")]})

		// call order service
		ordersRes, err = handler.orderService.FindAll(0, 0, filters, sorts)
		if err != nil {
			return helpers.WebErrorResponse(c, err, links)
		}

	case "admin":
		// Admin order list
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
	pageUrl := fmt.Sprintf("%s/api/customers/orders?page=", configs.Get().App.BaseURL)
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
 * Detail Order (Customer, Admin, Driver)
 * ------------------------------------
 * Mendapatkan detail order untuk
 * role spesific tertentu (customer, admin & driver)
 */
func (handler OrderHandler) Show(c echo.Context) error {
	userID, role, _ := middleware.ReadToken(c.Get("user"))
	links := map[string]string{}
	orderRes := entities.OrderResponse{}

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

	switch role {
	case "customer":
		// set self links and filters
		links["self"] = fmt.Sprintf("%s/api/customers/orders/%s", configs.Get().App.BaseURL,strconv.Itoa(orderID))

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
				Code: http.StatusBadRequest,
				Error: "order doesn't belong to currently authenticated customer",
				Links: links,
			})
		}
	}

	return c.JSON(http.StatusOK, web.SuccessResponse{
		Status: "OK",
		Code: http.StatusOK,
		Error: nil,
		Links: links,
		Data: orderRes,
	})
}