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
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	orderService orderService.OrderServiceInterface
	userService userService.UserServiceInterface
}

func NewOrderHandler(
	service orderService.OrderServiceInterface,
	userService userService.UserServiceInterface,
) *OrderHandler {
	return &OrderHandler{
		orderService: service,
		userService: userService,
	}
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