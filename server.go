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

	regionRepository := regionRepository.NewRegionRepository(db)
	userRepository := userRepository.NewUserRepository(db)
	truckTypeRepository := truckRepository.NewTruckTypeRepository(db)
	orderRepository := orderRepository.NewOrderRepository(db)
	orderHistoryRepository := orderHistoryRepository.NewOrderHistoryRepository(db)

	authService := authService.NewAuthService(userRepository)
	userService := userService.NewUserService(userRepository, truckTypeRepository)
	orderService := orderService.NewOrderService(orderRepository, orderHistoryRepository)
	regionService := regionService.NewRegionService(regionRepository)
	truckTypeService := truckTypeService.NewTruckTypeService(*truckTypeRepository)

	authHandler := handlers.NewAuthHandler(authService)
<<<<<<< HEAD
	customerHandler := handlers.NewCustomerHandler(userService, orderService)
	driverHandler := handlers.NewDriverHandler(userService)
=======
	userHandler := handlers.NewUserHandler(userService, orderService)
	driverHandler := handlers.NewDriverHandler(userService, orderService)
>>>>>>> refactor: driver list order handler
	adminHandler := handlers.NewAdminHandler(userService, orderService)
	truckTypeHandler := handlers.NewTruckTypeHandler(*truckTypeService)
	regionHandler := handlers.NewRegionHandler(regionService)
	_ = handlers.NewOrderHandler(orderService, userService)

	routes.RegisterDriverRoute(e, driverHandler)
	routes.RegisterAdminRoute(e, adminHandler)
	routes.RegisterTruckTypeRoute(e, truckTypeHandler)
	routes.RegisterAuthRoute(e, authHandler)
	routes.RegisterRegionHandler(e, regionHandler)
	routes.RegisterCustomerRoute(e, customerHandler, orderHandler)

<<<<<<< HEAD
=======
	routes.RegisterCustomerRoute(e, userHandler)
	
>>>>>>> refactor: customer detail order handler
	e.Logger.Fatal(e.Start(":" + config.App.Port))
}
