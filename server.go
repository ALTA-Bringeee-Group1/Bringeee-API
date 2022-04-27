package main

import (
	"bringeee-capstone/configs"
	"bringeee-capstone/deliveries/handlers"
	"bringeee-capstone/deliveries/routes"
	orderRepository "bringeee-capstone/repositories/order"
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
	userService := userService.NewUserService(userRepository, truckTypeRepository)
	userHandler := handlers.NewUserHandler(userService)
	routes.RegisterDriverRoute(e, userHandler)
	routes.RegisterCustomerRoute(e, userHandler)

	truckTypeService := truckTypeService.NewTruckTypeService(*truckTypeRepository)
	truckTypeHandler := handlers.NewTruckTypeHandler(*truckTypeService)
	routes.RegisterTruckTypeRoute(e, truckTypeHandler)
	authService := authService.NewAuthService(userRepository)
	authHandler := handlers.NewAuthHandler(authService)
	routes.RegisterAuthRoute(e, authHandler)

	regionRepository := regionRepository.NewRegionRepository(db)
	regionService := regionService.NewRegionService(*regionRepository)
	regionHandler := handlers.NewRegionHandler(*regionService)
	routes.RegisterRegionHandler(e, regionHandler)

	orderRepository := orderRepository.NewOrderRepository(db)
	_ = orderService.NewOrderService(orderRepository)

	e.Logger.Fatal(e.Start(":" + config.App.Port))
}
