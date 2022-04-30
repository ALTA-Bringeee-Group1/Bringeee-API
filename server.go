package main

import (
	"bringeee-capstone/configs"
	"bringeee-capstone/deliveries/handlers"
	"bringeee-capstone/deliveries/routes"
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

	authService := authService.NewAuthService(userRepository)
	userService := userService.NewUserService(userRepository, truckTypeRepository)
	orderService := orderService.NewOrderService(orderRepository, orderHistoryRepository, userRepository, midtransPaymentRepository)
	regionService := regionService.NewRegionService(regionRepository)
	truckTypeService := truckTypeService.NewTruckTypeService(*truckTypeRepository)

	authHandler := handlers.NewAuthHandler(authService)
	customerHandler := handlers.NewCustomerHandler(userService, orderService)
	driverHandler := handlers.NewDriverHandler(userService, orderService)
	adminHandler := handlers.NewAdminHandler(userService, orderService)
	truckTypeHandler := handlers.NewTruckTypeHandler(*truckTypeService)
	regionHandler := handlers.NewRegionHandler(regionService)
	orderHandler := handlers.NewOrderHandler(orderService, userService)

	routes.RegisterDriverRoute(e, driverHandler)
	routes.RegisterAdminRoute(e, adminHandler)
	routes.RegisterTruckTypeRoute(e, truckTypeHandler)
	routes.RegisterAuthRoute(e, authHandler)
	routes.RegisterRegionHandler(e, regionHandler)
	routes.RegisterCustomerRoute(e, customerHandler, orderHandler)
	routes.RegisterCustomerRoute(e, customerHandler, orderHandler)

	// err := orderService.ConfirmOrder(1, 2, false)
	// fmt.Println(err)

	// data, err := orderService.CreatePayment(5, entities.CreatePaymentRequest{ PaymentMethod: "BANK_TRANSFER_MANDIRI" })
	// if err != nil {
	// 	fmt.Println(utils.JsonEncode(err))
	// }
	// fmt.Println(utils.JsonEncode(data))
	// fmt.Println(utils.JsonEncode(configs.Get().Payment.MidtransServerKey))	
	
	// paymentRes, _ := midtransPaymentRepository.GetPaymentStatus("4da46a8a-b133-4ec2-b23c-76f795245018", "BANK_TRANSFER_MANDIRI")
	// fmt.Println(utils.JsonEncode(paymentRes))
	// paymentRes, _ = midtransPaymentRepository.GetPaymentStatus("80da96ea-ce20-45eb-a89f-320759d41272", "BANK_TRANSFER_BRI")
	// fmt.Println(utils.JsonEncode(paymentRes))
	// paymentRes, _ = midtransPaymentRepository.GetPaymentStatus("83d81c00-8eb5-4de1-8450-7d57bfb52ecc", "BANK_TRANSFER_BNI")
	// fmt.Println(utils.JsonEncode(paymentRes))
	// paymentRes, _ = midtransPaymentRepository.GetPaymentStatus("01072be0-cd5a-4dfc-a63a-d9c199f0cb94", "BANK_TRANSFER_PERMATA")
	// fmt.Println(utils.JsonEncode(paymentRes))
	// paymentRes, _ = midtransPaymentRepository.GetPaymentStatus("51e3a2a7-fe69-44be-875e-0d75df142299", "BANK_TRANSFER_BCA")
	// fmt.Println(utils.JsonEncode(paymentRes))
	
	e.Logger.Fatal(e.Start(":" + config.App.Port))
}
