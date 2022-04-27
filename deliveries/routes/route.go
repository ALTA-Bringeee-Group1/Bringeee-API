package routes

import (
	"bringeee-capstone/deliveries/handlers"
	middleware "bringeee-capstone/deliveries/middlewares"

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
func RegisterAuthRoute(e *echo.Echo, authHandler *handlers.AuthHandler) {
	e.POST("/api/auth", authHandler.Login)
	e.GET("/api/auth/me", authHandler.Me, middleware.JWTMiddleware())
}

func RegisterRegionHandler(e *echo.Echo, regionHandler *handlers.RegionHandler) {
	e.GET("/api/provinces", regionHandler.IndexProvince)
	e.GET("/api/provinces/:provinceID/cities", regionHandler.IndexCity)
	e.GET("/api/provinces/:provinceID/cities/:cityID/districts", regionHandler.IndexDistrict)
}