package handlers

import (
	"bringeee-capstone/configs"
	middleware "bringeee-capstone/deliveries/middlewares"
	"bringeee-capstone/entities/web"
	orderService "bringeee-capstone/services/order"
	userService "bringeee-capstone/services/user"
	"bringeee-capstone/utils"
	"fmt"
	"net/http"
	"reflect"
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
 * List Order (Customer, Admin, Driver)
 * ------------------------------------
 * Mendapatkan list order untuk
 * role spesific tertentu (customer, admin & driver)
 */
func (handler OrderHandler) Index(c echo.Context) error {
	
	// hateoas links
	links := map[string]string{}
	
	// get token
	userID, role, _ := middleware.ReadToken(c.Get("user"))
	
	switch role {
	case "customer":
		
		// set links
		links["self"] = fmt.Sprintf("%s/api/customers/orders", configs.Get().App.BaseURL)

		// get authenticated userdata
		_, err := handler.userService.FindCustomer(userID)
		if err != nil {
			if reflect.TypeOf(err).String() == "web.WebError" {
				webErr := err.(web.WebError)
				return c.JSON(webErr.Code, web.ErrorResponse{
					Status: "ERROR",
					Code: webErr.Code,
					Error: webErr.Error(),
					Links: links,
				})
			}
			return c.JSON(http.StatusInternalServerError, web.ErrorResponse{
				Status: "ERROR",
				Code: http.StatusInternalServerError,
				Error: "Server Error",
				Links: links,
			})
		}
		
		// call order service
		filters := []map[string]interface{} {
			{ "field": "customer_id", "operator": "=", "value": strconv.Itoa(userID) },
			{
				"or": []map[string]string {
					{ "field": "status", "operator": "=", "value": "CONFIRMED" },
					{ "field": "status", "operator": "=", "value": "PENDING", "logic": "OR" },
					{ "field": "status", "operator": "=", "value": "MANIFESTED", "logic": "OR" },
				},
			},
		}
		orders, err := handler.orderService.FindAll(0, 0, filters, []map[string]interface{}{})
		fmt.Println(utils.JsonEncode(orders))
		return c.String(200, "hooh")
	}
	return c.String(200, "hooh")
}