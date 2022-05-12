package main

import (
	"bringeee-capstone/configs"
	"bringeee-capstone/deliveries/handlers"
	"bringeee-capstone/deliveries/routes"
	distanceMatrixRepository "bringeee-capstone/repositories/distance_matrix"
	orderRepository "bringeee-capstone/repositories/order"
	orderHistoryRepository "bringeee-capstone/repositories/order_history"
	paymentRepository "bringeee-capstone/repositories/payment"
	regionRepository "bringeee-capstone/repositories/region"
	truckRepository "bringeee-capstone/repositories/truck_type"
	userRepository "bringeee-capstone/repositories/user"
	authService "bringeee-capstone/services/auth"
	orderService "bringeee-capstone/services/order"
	regionService "bringeee-capstone/services/region"
	storageProvider "bringeee-capstone/services/storage"
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
	midtransPaymentRepository := paymentRepository.NewMidtransPaymentRepository()
	distanceMatrixRepository := distanceMatrixRepository.NewDistanceMatrixRepository()

	authService := authService.NewAuthService(userRepository)
	userService := userService.NewUserService(userRepository, truckTypeRepository, orderRepository)
	orderService := orderService.NewOrderService(orderRepository, orderHistoryRepository, userRepository, midtransPaymentRepository, distanceMatrixRepository, truckTypeRepository)
	regionService := regionService.NewRegionService(regionRepository)
	truckTypeService := truckTypeService.NewTruckTypeService(*truckTypeRepository)

	s3 := storageProvider.NewS3()

	authHandler := handlers.NewAuthHandler(authService)
	customerHandler := handlers.NewCustomerHandler(userService, orderService, s3)
	driverHandler := handlers.NewDriverHandler(userService, orderService, s3)
	adminHandler := handlers.NewAdminHandler(userService, orderService, truckTypeService, s3)
	truckTypeHandler := handlers.NewTruckTypeHandler(*truckTypeService)
	regionHandler := handlers.NewRegionHandler(regionService)
	orderHandler := handlers.NewOrderHandler(orderService, userService)
	paymentHandler := handlers.NewPaymentHandler(orderService)

	routes.RegisterDriverRoute(e, driverHandler)
	routes.RegisterAdminRoute(e, adminHandler)
	routes.RegisterTruckTypeRoute(e, truckTypeHandler)
	routes.RegisterAuthRoute(e, authHandler)
	routes.RegisterRegionHandler(e, regionHandler)
	routes.RegisterCustomerRoute(e, customerHandler, orderHandler)
	routes.RegisterCustomerRoute(e, customerHandler, orderHandler)
	routes.RegisterPaymentRoute(e, paymentHandler)

	e.Logger.Fatal(e.Start(":" + config.App.Port))
}
