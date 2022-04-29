package order

import (
	"bringeee-capstone/deliveries/helpers"
	"bringeee-capstone/deliveries/validations"
	"bringeee-capstone/entities"
	"bringeee-capstone/entities/web"
	orderRepository "bringeee-capstone/repositories/order"
	orderHistoryRepository "bringeee-capstone/repositories/order_history"
	userRepository "bringeee-capstone/repositories/user"
	"mime/multipart"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type OrderService struct {
	orderRepository        orderRepository.OrderRepositoryInterface
	orderHistoryRepository orderHistoryRepository.OrderHistoryRepositoryInterface
	userRepository         userRepository.UserRepositoryInterface
	validate               *validator.Validate
}

func NewOrderService(repository orderRepository.OrderRepositoryInterface, orderHistoryRepository orderHistoryRepository.OrderHistoryRepositoryInterface, userRepository userRepository.UserRepositoryInterface) *OrderService {
	return &OrderService{
		orderRepository:        repository,
		orderHistoryRepository: orderHistoryRepository,
		userRepository:         userRepository,
		validate:               validator.New(),
	}
}

/*
 * Find All
 * -------------------------------
 * Mengambil data order berdasarkan filters dan sorts
 *
 * @var limit 	batas limit hasil query
 * @var offset 	offset hasil query
 * @var filters	query untuk penyaringan data, { field, operator, value }
 * @var sorts	pengurutan data, { field, value[bool] }
 * @return order	list order dalam bentuk entity domain
 * @return error	error
 */
func (service OrderService) FindAll(limit int, page int, filters []map[string]interface{}, sorts []map[string]interface{}) ([]entities.OrderResponse, error) {

	offset := (page - 1) * limit

	// Repository action find all order
	orders, err := service.orderRepository.FindAll(limit, offset, filters, sorts)
	if err != nil {
		return []entities.OrderResponse{}, err
	}

	// Konversi ke order response
	ordersRes := []entities.OrderResponse{}
	copier.Copy(&ordersRes, &orders)
	for i, order := range orders {
		copier.Copy(&ordersRes[i], &order.Destination)
		copier.Copy(&ordersRes[i].Driver, &order.Driver.User)
		ordersRes[i].ID = order.ID // fix: overlap destinationID vs orderID
	}
	return ordersRes, nil
}

/*
 * Get Pagination
 * -------------------------------
 * Mengambil data pagination berdasarkan filters dan sorts
 *
 * @var limit 	batas limit hasil query
 * @var page 	halaman sekarang diakses
 * @var filters	query untuk penyaringan data, { field, operator, value }
 * @return order	response pagination
 * @return error	error
 */
func (service OrderService) GetPagination(limit int, page int, filters []map[string]interface{}) (web.Pagination, error) {
	totalRows, err := service.orderRepository.CountAll(filters)
	if err != nil {
		return web.Pagination{}, err
	}
	var totalPages int64 = 1
	if limit > 0 {
		totalPages = totalRows / int64(limit)
	}
	if totalPages <= 0 {
		totalPages = 1
	}
	return web.Pagination{
		Page:         page,
		Limit:        limit,
		TotalPages:   int(totalPages),
		TotalRecords: int(totalRows),
	}, nil
}

/*
 * Find
 * -------------------------------
 * Mengambil data order tunggal berdasarkan ID
 *
 * @var id 		id order
 * @return order	order tunggal dalam bentuk response
 * @return error	error
 */
func (service OrderService) Find(id int) (entities.OrderResponse, error) {
	order, err := service.orderRepository.Find(id)
	if err != nil {
		return entities.OrderResponse{}, err
	}

	// convert to response
	orderRes := entities.OrderResponse{}
	copier.Copy(&orderRes, &order)
	copier.Copy(&orderRes, &order.Destination)
	copier.Copy(&orderRes.Driver, &order.Driver.User)
	orderRes.ID = order.ID // fix: overlap destinationID vs orderID

	return orderRes, nil
}

/*
 * Find First
 * -------------------------------
 * Mengambil order pertama berdasarkan filter yang telah di tentukan pada parameter
 * dan mengambil data pertama sebagai data order tunggal
 * @var filter
 * @return OrderResponse	order response dalam bentuk tunggal
 * @return error			error
 */
func (service OrderService) FindFirst(filters []map[string]interface{}) (entities.OrderResponse, error) {
	// Repository call
	order, err := service.orderRepository.FindFirst(filters)
	if err != nil {
		return entities.OrderResponse{}, err
	}

	// Convert to response
	orderRes := entities.OrderResponse{}
	copier.Copy(&orderRes, &order)
	copier.Copy(&orderRes, &order.Destination)
	copier.Copy(&orderRes.Driver, &order.Driver.User)
	orderRes.ID = order.ID // fix: overlap destinationID vs orderID

	return orderRes, nil
}

/*
 * Customer Create order
 * -------------------------------
 * Membuat order baru berdasarkan user yang sedang login
 * @var orderRequest		request create order oleh customer
 * @var files				list file request untuk diteruskan ke validation dan upload
 * @return OrderResponse	order response
 */
func (service OrderService) Create(orderRequest entities.CustomerCreateOrderRequest, files map[string]*multipart.FileHeader, userID int) (entities.OrderResponse, error) {
	// validation
	err := validations.ValidateCustomerCreateOrderRequest(service.validate, orderRequest, files)
	if err != nil {
		return entities.OrderResponse{}, err
	}

	// convert request to domain
	order := entities.Order{}
	destination := entities.Destination{}
	copier.Copy(&order, &orderRequest)
	copier.Copy(&destination, &orderRequest)
	order.CustomerID = uint(userID)
	order.Status = "REQUESTED"

	// Upload file to cloud storage
	for field, file := range files {
		switch field {
		case "order_picture":
			fileFile, err := file.Open()
			if err != nil {
				return entities.OrderResponse{}, web.WebError{Code: 500, Message: "Cannot process the requested file"}
			}
			defer fileFile.Close()

			fileName := uuid.New().String() + file.Filename
			fileUrl, err := helpers.UploadFileToS3("orders/order_picture/"+fileName, fileFile)
			if err != nil {
				return entities.OrderResponse{}, web.WebError{Code: 500, ProductionMessage: "Cannot upload requested file", DevelopmentMessage: err.Error()}
			}
			order.OrderPicture = fileUrl
		}
	}

	// repository call
	order, err = service.orderRepository.Store(order, destination)
	if err != nil {
		return entities.OrderResponse{}, err
	}

	// Log
	service.orderHistoryRepository.Create(int(order.ID), "Order dibuat dan diajukan oleh customer", "customer")

	// get newly order data
	orderRes, err := service.Find(int(order.ID))
	if err != nil {
		return entities.OrderResponse{}, web.WebError{
			Code:               500,
			DevelopmentMessage: "Cannot get newly inserted data: " + err.Error(),
			ProductionMessage:  "Cannot get newly data",
		}
	}
	return orderRes, nil
}

/*
 * Admin Set fixed price order
 * -------------------------------
 * Set fixed price order oleh admin untuk diteruskan kembali ke user agar di konfirmasi/cancel
 * @var orderID				orderID
 * @var orderRequest		request create order oleh customer
 * @return OrderResponse	order response
 */
func (service OrderService) SetFixOrder(orderID int, setPriceRequest entities.AdminSetPriceOrderRequest) error  {
	err := validations.ValidateAdminSetPriceOrderRequest(service.validate, setPriceRequest)
	if err != nil {
		return err
	}

	// get order
	order, err := service.orderRepository.Find(orderID)
	if err != nil {
		return err
	}

	// reject if status other than requested
	if order.Status != "REQUESTED" {
		return web.WebError{
			Code: 400,
			Message: "Status order was already confirmed or cancelled",
		}
	}
	
	// Update via repository
	order.FixPrice = setPriceRequest.FixedPrice
	_, err = service.orderRepository.Update(order, orderID)
	if err != nil {
		return err
	}
	service.orderHistoryRepository.Create(orderID, "Perubahan detail order oleh admin", "admin")
	return nil
}

/*
 * Confirm Order
 * -------------------------------
 * Confirm order yang sudah dibuat
 * @var orderID				ID Order yang akan di cancel
 * @return OrderResponse	order response
 */
func (service OrderService) ConfirmOrder(orderID int, userID int, isAdmin bool) error  {
	// get order
	order, err := service.orderRepository.Find(orderID)
	if err != nil {
		return err
	}

	// set fix price = estimated price if fix price empty
	if order.FixPrice <= 0 {
		order.FixPrice = order.EstimatedPrice
	}

	// reject if status other than requested
	if order.Status != "REQUESTED" {
		return web.WebError{
			Code: 400,
			Message: "Status order was already confirmed or cancelled",
		}
	}

	// reject if order doesn't belong to customer (except admin)
	if !isAdmin {
		if order.CustomerID != uint(userID) {
			return web.WebError{
				Code: 400,
				Message: "Order doesn't belong to currently authenticated user",
			}
		}
	}

	// Update via repository
	order.Status = "CONFIRMED"
	_, err = service.orderRepository.Update(order, orderID)
	if err != nil {
		return err
	}
	service.orderHistoryRepository.Create(orderID, "Order dikonfirmasi", map[bool]string{true: "admin", false: "customer"}[isAdmin])
	return nil
}

/*
 * Cancel Order
 * -------------------------------
 * Cancel order yang sudah dibuat
 * @var orderID				ID Order yang akan di cancel
 * @return OrderResponse	order response
 */
func (service OrderService) CancelOrder(orderID int, userID int, isAdmin bool) error  {
	// get order
	order, err := service.orderRepository.Find(orderID)
	if err != nil {
		return err
	}

	// reject if status other than requested
	if order.Status != "REQUESTED" && order.Status != "CONFIRMED" {
		return web.WebError{
			Code: 400,
			Message: "Cannot cancel order: status order was already " + order.Status,
		}
	}

	// reject if order doesn't belong to customer (except admin)
	if !isAdmin {
		if order.CustomerID != uint(userID) {
			return web.WebError{
				Code: 400,
				Message: "Order doesn't belong to currently authenticated user",
			}
		}
	}

	// Update via repository
	order.Status = "CANCELLED"
	_, err = service.orderRepository.Update(order, orderID)
	if err != nil {
		return err
	}
	service.orderHistoryRepository.Create(orderID, "Order dibatalkan", map[bool]string{true: "admin", false: "customer"}[isAdmin])
	return nil
}

/*
 * Create Payment
 * -------------------------------
 * Buat pembayaran untuk order tertentu ke layanan pihak ketiga
 * @var orderID					ID Order yang akan di cancel
 * @var createPaymentRequest	request payment
 * @return PaymentResponse		response payment
 */
func (service OrderService) CreatePayment(orderID int, createPaymentRequest entities.CreatePaymentRequest) (entities.PaymentResponse, error) {
	panic("implement me")
}

/*
 * Get Payment
 * -------------------------------
 * Mengambil data pembayaran yang sudah ada berdasarkan transaction_id yang sudah di set pada order
 * @var orderID					ID Order yang akan di cancel
 * @var createPaymentRequest	request payment
 * @return PaymentResponse		response payment
 */
func (service OrderService) GetPayment(orderID int, createPaymentRequest entities.CreatePaymentRequest) (entities.PaymentResponse, error) {
	panic("implement me")
}

/*
 * Find All History
 * -------------------------------
 * Mengambil data order berdasarkan filters dan sorts
 *
 * @var limit 	batas limit hasil query
 * @var offset 	offset hasil query
 * @var filters	query untuk penyaringan data, { field, operator, value }
 * @var sorts	pengurutan data, { field, value[bool] }
 * @return order	list order dalam bentuk entity domain
 * @return error	error
 */
func (service OrderService) FindAllHistory(orderID int, sorts []map[string]interface{}) ([]entities.OrderHistoryResponse, error) {
	histories, err := service.orderHistoryRepository.FindAll(orderID, sorts)
	if err != nil {
		return []entities.OrderHistoryResponse{}, err
	}
	historiesRes := []entities.OrderHistoryResponse{}
	copier.Copy(&historiesRes, &histories)
	return historiesRes, nil
}

/*
 * Webhook
 * -------------------------------
 * Payment Webhook notification, dikirimkan oleh layanan pihak ketiga
 * referensi: https://docs.midtrans.com/en/after-payment/http-notification
 *
 * @var limit 	batas limit hasil query
 * @var offset 	offset hasil query
 * @var filters	query untuk penyaringan data, { field, operator, value }
 * @var sorts	pengurutan data, { field, value[bool] }
 * @return order	list order dalam bentuk entity domain
 * @return error	error
 */
func (service OrderService) PaymentWebhook(orderID int) error {

	// if status settlement, set order to MANIFESTED
	// if status is deny, cancel, expire, set to CANCELLED

	panic("implement me")
}

/*
 * Take order for shipping
 * -------------------------------
 * Pengambilan order oleh driver untuk di set statusnya menjadi ON_PROCESS
 * @var orderID 	order id terkait
 * @var userID		authenticated user (role: driver)
 */
func (service OrderService) TakeOrder(orderID int, userID int) error {
	order, err := service.orderRepository.Find(orderID)
	if err != nil {
		return web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "The requested ID doesn't match with any record"}
	}
	driver, err := service.userRepository.FindByDriver("user_id", strconv.Itoa(userID))
	if err != nil {
		return web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "The requested Driver ID doesn't match with any record"}
	}
	if order.DriverID != 0 {
		return web.WebError{Code: 400, ProductionMessage: "bad request", DevelopmentMessage: "This order already taken by someone else"}
	}
	if order.TruckTypeID != driver.TruckTypeID {
		return web.WebError{Code: 400, ProductionMessage: "bad request", DevelopmentMessage: "Cannot take order that doesn't match with your truck type"}
	}
	if order.Status != "CONFIRMED" {
		return web.WebError{Code: 400, ProductionMessage: "bad request", DevelopmentMessage: "This order is not confirmed by admin"}
	}
	order.DriverID = driver.ID
	order.Status = "MANIFESTED"
	order, err = service.orderRepository.Update(order, userID)
	if err != nil {
		return web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "Cannot update current order"}
	}
	// Log
	service.orderHistoryRepository.Create(int(order.ID), "Order diambil oleh "+driver.User.Name, "driver")
	return err
}

/*
 * Finish Order
 * -------------------------------
 * penyelesaian order oleh driver untuk di set statusnya menjadi DELIVERED
 * @var orderID 	order id terkait
 * @var userID		authenticated user (role: driver)
 */
func (service OrderService) FinishOrder(orderID int, userID int, files map[string]*multipart.FileHeader) error {
	// validation
	err := validations.ValidateUpdateOrderRequest(files)
	if err != nil {
		return web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: err.Error()}
	}
	order, err := service.orderRepository.Find(orderID)
	if err != nil {
		return web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "The requested ID doesn't match with any record"}
	}
	driver, err := service.userRepository.FindByDriver("user_id", strconv.Itoa(userID))
	if err != nil {
		return web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "The requested Driver ID doesn't match with any record"}
	}
	if order.Status == "DELIVERED" {
		return web.WebError{Code: 400, ProductionMessage: "bad request", DevelopmentMessage: "Already finish!"}
	}
	if order.DriverID == 0 {
		return web.WebError{Code: 400, ProductionMessage: "bad request", DevelopmentMessage: "The current order has not belong to somenone else, take the ordrer first"}
	}
	if order.DriverID != driver.ID {
		return web.WebError{Code: 400, ProductionMessage: "bad request", DevelopmentMessage: "Cannot finish order that belong to someone else"}
	}
	if len(files) == 0 {
		return web.WebError{Code: 400, ProductionMessage: "bad request", DevelopmentMessage: "Arrived picture must be filled"}
	}
	// Upload file to cloud storage
	for field, file := range files {
		switch field {
		case "arrived_picture":
			fileFile, err := file.Open()
			if err != nil {
				return web.WebError{Code: 500, Message: "Cannot process the requested file"}
			}
			defer fileFile.Close()

			fileName := uuid.New().String() + file.Filename
			fileUrl, err := helpers.UploadFileToS3("orders/arrived_picture/"+fileName, fileFile)
			if err != nil {
				return web.WebError{Code: 500, ProductionMessage: "Cannot upload requested file", DevelopmentMessage: err.Error()}
			}
			order.OrderPicture = fileUrl
		}
	}
	order.Status = "DELIVERED"
	order, err = service.orderRepository.Update(order, userID)
	if err != nil {
		return web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "Cannot update current order"}
	}
	// Log
	service.orderHistoryRepository.Create(int(order.ID), "Order yang diantar "+driver.User.Name+" telah sampai di tujuan", "driver")

	return err
}
