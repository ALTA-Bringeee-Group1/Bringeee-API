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
	truckTypeService "bringeee-capstone/services/truck_type"
	userService "bringeee-capstone/services/user"
	"bringeee-capstone/utils"
	"fmt"

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
	orderService := orderService.NewOrderService(orderRepository, orderHistoryRepository, userRepository, midtransPaymentRepository)
	regionService := regionService.NewRegionService(regionRepository)
	truckTypeService := truckTypeService.NewTruckTypeService(*truckTypeRepository)

	authHandler := handlers.NewAuthHandler(authService)
	customerHandler := handlers.NewCustomerHandler(userService, orderService)
	driverHandler := handlers.NewDriverHandler(userService, orderService)
	adminHandler := handlers.NewAdminHandler(userService, orderService, truckTypeService)
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

	data, _ := distanceMatrixRepository.EstimateShortest(
		"-6.129634", 
		"106.827312", 
		"-7.795688637022531", 
		"110.3653103342137",
	)
	fmt.Println(utils.JsonEncode(data))

	fmt.Println("Awok?")

	// e.Logger.Fatal(e.Start(":" + config.App.Port))
}
