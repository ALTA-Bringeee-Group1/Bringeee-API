package routes

import (
	"bringeee-capstone/deliveries/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterCustomerRoute(e *echo.Echo, userHandler *handlers.UserHandler) {
	group := e.Group("/api/customers")
	group.POST("", userHandler.CreateCustomer) // Registration customer
}

func RegisterDriverRoute(e *echo.Echo, driverHandler *handlers.UserHandler) {
	group := e.Group("/api/drivers")
	group.POST("", driverHandler.CreateDriver) // Registration driver
}

func RegisterTruckTypeRoute(e *echo.Echo, truckTypeHandler *handlers.TruckTypeHandler) {
	e.GET("/api/truck_types", truckTypeHandler.Index)
}