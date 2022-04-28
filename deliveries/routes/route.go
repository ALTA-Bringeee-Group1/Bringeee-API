package routes

import (
	"bringeee-capstone/deliveries/handlers"
	middleware "bringeee-capstone/deliveries/middlewares"

	"github.com/labstack/echo/v4"
)

func RegisterCustomerRoute(e *echo.Echo, userHandler *handlers.UserHandler, orderHandler *handlers.OrderHandler) {
	group := e.Group("/api/customers")
	group.POST("", userHandler.CreateCustomer)                               // Registration customer
	group.PUT("", userHandler.UpdateCustomer, middleware.JWTMiddleware())    // Edit customer profile
	group.DELETE("", userHandler.DeleteCustomer, middleware.JWTMiddleware()) // delete customer

	order := e.Group("/api/customers/orders", middleware.JWTMiddleware())
	order.GET("", orderHandler.Index)
	order.GET("/:orderID", orderHandler.Show)
	order.GET("/:orderID/histories", userHandler.DetailOrderHistory)
}

func RegisterDriverRoute(e *echo.Echo, driverHandler *handlers.DriverHandler) {
	group := e.Group("/api/drivers")
	group.POST("", driverHandler.CreateDriver)                            // Registration driver
	group.PUT("", driverHandler.UpdateDriver, middleware.JWTMiddleware()) // Edit driver profile
}

func RegisterAdminRoute(e *echo.Echo, AdminHandler *handlers.AdminHandler) {
	e.PUT("/api/drivers/:id", AdminHandler.UpdateDriverByAdmin, middleware.JWTMiddleware()) // edit driver credential
	e.DELETE("api/drivers/:id", AdminHandler.DeleteDriver, middleware.JWTMiddleware())      // delete driver
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
