package routes

import (
	"bringeee-capstone/deliveries/handlers"
	middleware "bringeee-capstone/deliveries/middlewares"

	"github.com/labstack/echo/v4"
)

func RegisterCustomerRoute(e *echo.Echo, customerHandler *handlers.CustomerHandler, orderHandler *handlers.OrderHandler) {
	group := e.Group("/api/customers")
	group.POST("", customerHandler.CreateCustomer)                               // Registration customer
	group.PUT("", customerHandler.UpdateCustomer, middleware.JWTMiddleware())    // Edit customer profile
	group.DELETE("", customerHandler.DeleteCustomer, middleware.JWTMiddleware()) // delete customer

	order := e.Group("/api/customers/orders", middleware.JWTMiddleware())
	order.GET("", customerHandler.ListOrders)
	order.POST("", customerHandler.CreateOrder)
	order.GET("/:orderID", customerHandler.DetailOrder)
	order.GET("/:orderID/histories", customerHandler.DetailOrderHistory)
	order.POST("/:orderID/confirm", customerHandler.ConfirmOrder)
	order.POST("/:orderID/cancel", customerHandler.CancelOrder)
	order.POST("/:orderID/payment", customerHandler.CreatePayment)
	order.GET("/:orderID/payment", customerHandler.GetPayment)
	order.POST("/:orderID/payment/cancel", customerHandler.CancelPayment)
}

func RegisterDriverRoute(e *echo.Echo, driverHandler *handlers.DriverHandler) {
	group := e.Group("/api/drivers")
	group.POST("", driverHandler.CreateDriver)                            // Registration driver
	group.PUT("", driverHandler.UpdateDriver, middleware.JWTMiddleware()) // Edit driver profile

	order := e.Group("/api/drivers/orders", middleware.JWTMiddleware())
	order.GET("", driverHandler.ListOrders)
	order.POST("/:id/take_order", driverHandler.TakeOrder, middleware.JWTMiddleware())
	order.POST("/:id/finish_order", driverHandler.FinishOrder, middleware.JWTMiddleware())
	e.GET("/api/drivers/current_order", driverHandler.CurrentOrder, middleware.JWTMiddleware())
}

func RegisterAdminRoute(e *echo.Echo, AdminHandler *handlers.AdminHandler) {
	e.POST("api/drivers/:id/confirm", AdminHandler.VerifiedDriverAccount, middleware.JWTMiddleware())
	e.GET("api/drivers", AdminHandler.GetAllDriver, middleware.JWTMiddleware())
	e.GET("api/drivers/:id", AdminHandler.GetSingleDriver, middleware.JWTMiddleware())
	e.GET("api/customers/:id", AdminHandler.GetSingleCustomer, middleware.JWTMiddleware())
	e.GET("api/customers", AdminHandler.GetAllCustomer, middleware.JWTMiddleware())
	e.PUT("/api/drivers/:id", AdminHandler.UpdateDriverByAdmin, middleware.JWTMiddleware()) // edit driver credential
	e.DELETE("api/drivers/:id", AdminHandler.DeleteDriver, middleware.JWTMiddleware())      // delete driver

	order := e.Group("/api/orders", middleware.JWTMiddleware())
	order.GET("", AdminHandler.ListOrders)
	order.GET("/:orderID", AdminHandler.DetailOrder)
	order.PATCH("/:orderID", AdminHandler.SetFixedPrice)
	order.POST("/:orderID/confirm", AdminHandler.ConfirmOrder)
	order.POST("/:orderID/cancel", AdminHandler.CancelOrder)
	order.GET("/:orderID/histories", AdminHandler.DetailOrderHistory)
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
