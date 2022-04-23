package routes

import (
	"bringeee-capstone/deliveries/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterCustomerRoute(e *echo.Echo, userHandler *handlers.UserHandler) {
	group := e.Group("/api/customers")
	group.POST("", userHandler.CreateCustomer) // Registration
}
