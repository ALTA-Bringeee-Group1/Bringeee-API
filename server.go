package main

import (
	"bringeee-capstone/configs"
	"bringeee-capstone/deliveries/handlers"
	"bringeee-capstone/deliveries/routes"
	orderRepository "bringeee-capstone/repositories/order"
	orderHistoryRepository "bringeee-capstone/repositories/order_history"
	regionRepository "bringeee-capstone/repositories/region"
	truckRepository "bringeee-capstone/repositories/truck_type"
	userRepository "bringeee-capstone/repositories/user"
	authService "bringeee-capstone/services/auth"
	orderService "bringeee-capstone/services/order"
	regionService "bringeee-capstone/services/region"
	truckTypeService "bringeee-capstone/services/truck_type"
	userService "bringeee-capstone/services/user"
	"bringeee-capstone/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config := configs.Get()
	db := utils.NewMysqlGorm(config)
	utils.Migrate(db)

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS},
	}))

	userRepository := userRepository.NewUserRepository(db)
	truckTypeRepository := truckRepository.NewTruckTypeRepository(db)
	authService := authService.NewAuthService(userRepository)
	regionRepository := regionRepository.NewRegionRepository(db)
	orderRepository := orderRepository.NewOrderRepository(db)
	orderHistoryRepository := orderHistoryRepository.NewOrderHistoryRepository(db)

	orderService := orderService.NewOrderService(orderRepository, orderHistoryRepository)
	userService := userService.NewUserService(userRepository, truckTypeRepository)
	regionService := regionService.NewRegionService(regionRepository)
	truckTypeService := truckTypeService.NewTruckTypeService(*truckTypeRepository)
	
	userHandler := handlers.NewUserHandler(userService)
	orderHandler := handlers.NewOrderHandler(orderService, userService)
	regionHandler := handlers.NewRegionHandler(regionService)
	authHandler := handlers.NewAuthHandler(authService)
	truckTypeHandler := handlers.NewTruckTypeHandler(*truckTypeService)
	adminHandler := handlers.NewAdminHandler(userService, orderService)
	driverHandler := handlers.NewDriverHandler(userService)

	routes.RegisterTruckTypeRoute(e, truckTypeHandler)
	routes.RegisterAuthRoute(e, authHandler)
	routes.RegisterDriverRoute(e, driverHandler)
	routes.RegisterAdminRoute(e, adminHandler)
	routes.RegisterRegionHandler(e, regionHandler)
	routes.RegisterCustomerRoute(e, userHandler, orderHandler)
	
	e.Logger.Fatal(e.Start(":" + config.App.Port))
}