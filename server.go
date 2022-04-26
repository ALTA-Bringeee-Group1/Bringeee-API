package main

import (
	"bringeee-capstone/configs"
	"bringeee-capstone/deliveries/handlers"
	"bringeee-capstone/deliveries/routes"
	truckRepository "bringeee-capstone/repositories/truck_type"
	truckTypeRepository "bringeee-capstone/repositories/truck_type"
	userRepository "bringeee-capstone/repositories/user"
	authService "bringeee-capstone/services/auth"
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
	truckRepository := truckRepository.NewTruckTypeRepository(db)
	userService := userService.NewUserService(userRepository, truckRepository)
	userHandler := handlers.NewUserHandler(userService)
	routes.RegisterDriverRoute(e, userHandler)
	routes.RegisterCustomerRoute(e, userHandler)

	truckTypeRepository := truckTypeRepository.NewTruckTypeRepository(db)
	truckTypeService := truckTypeService.NewTruckTypeService(*truckTypeRepository)
	truckTypeHandler := handlers.NewTruckTypeHandler(*truckTypeService)
	routes.RegisterTruckTypeRoute(e, truckTypeHandler)
	authService := authService.NewAuthService(userRepository)
	authHandler := handlers.NewAuthHandler(authService)
	routes.RegisterAuthRoute(e, authHandler)

	e.Logger.Fatal(e.Start(":" + config.App.Port))
}
