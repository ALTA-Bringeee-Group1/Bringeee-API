package order

import (
	"bringeee-capstone/deliveries/helpers"
	"bringeee-capstone/deliveries/validations"
	"bringeee-capstone/entities"
	"bringeee-capstone/entities/web"
	orderRepository "bringeee-capstone/repositories/order"
	orderHistoryRepository "bringeee-capstone/repositories/order_history"
	paymentRepository "bringeee-capstone/repositories/payment"
	userRepository "bringeee-capstone/repositories/user"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gopkg.in/guregu/null.v4"
)

type OrderService struct {
	orderRepository        orderRepository.OrderRepositoryInterface
	orderHistoryRepository orderHistoryRepository.OrderHistoryRepositoryInterface
	userRepository         userRepository.UserRepositoryInterface
	paymentRepository      paymentRepository.PaymentRepositoryInterface
	validate               *validator.Validate
}

func NewOrderService(
	repository orderRepository.OrderRepositoryInterface, 
	orderHistoryRepository orderHistoryRepository.OrderHistoryRepositoryInterface,
	userRepository         userRepository.UserRepositoryInterface,
	paymentRepository paymentRepository.PaymentRepositoryInterface,
) *OrderService {
	return &OrderService{
		orderRepository:        repository,
		orderHistoryRepository: orderHistoryRepository,
		userRepository:         userRepository,
		paymentRepository: 		paymentRepository,
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
func (service OrderService) SetFixOrder(orderID int, setPriceRequest entities.AdminSetPriceOrderRequest) error {
	err := validations.ValidateAdminSetPriceOrderRequest(service.validate, setPriceRequest)
	if err != nil {
		return err
	}

	// get order
	order, err := service.orderRepository.Find(orderID)
	if err != nil {
		return err
	}

	// reject if status other than requested and already fixed
	if order.Status != "REQUESTED" && order.Status != "NEED_CUSTOMER_CONFIRM" {
		return web.WebError{
			Code:    400,
			Message: "Status order was already confirmed or cancelled",
		}
	}

	// Update via repository
	order.FixPrice = setPriceRequest.FixedPrice
	order.Status = "NEED_CUSTOMER_CONFIRM"
	order.TransactionID = ""
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
func (service OrderService) ConfirmOrder(orderID int, userID int, isAdmin bool) error {
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
	if order.Status != "REQUESTED" && order.Status != "NEED_CUSTOMER_CONFIRM" {
		return web.WebError{
			Code:    400,
			Message: "Status order was already confirmed or cancelled",
		}
	}

	// reject if order doesn't belong to customer (except admin)
	if !isAdmin {
		if order.CustomerID != uint(userID) {
			return web.WebError{
				Code:    400,
				Message: "Order doesn't belong to currently authenticated user",
			}
		} else if order.Status != "NEED_CUSTOMER_CONFIRM" {
			return web.WebError{
				Code:    400,
				Message: "Waiting for admin response",
			}
		}
	}
	if isAdmin {
		if order.Status == "NEED_CUSTOMER_CONFIRM" {
			return web.WebError{
				Code:    400,
				Message: "Waiting for customer response",
			}
		}
	}

	// Update via repository
	order.Status = "CONFIRMED"
	order.DriverID = null.IntFromPtr(nil)
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
func (service OrderService) CancelOrder(orderID int, userID int, isAdmin bool) error {
	// get order
	order, err := service.orderRepository.Find(orderID)
	if err != nil {
		return err
	}

	// reject if status other than requested
	if order.Status != "REQUESTED" && order.Status != "CONFIRMED" && order.Status != "NEED_CUSTOMER_CONFIRM" {
		return web.WebError{
			Code:    400,
			Message: "Cannot cancel order: status order was already " + order.Status,
		}
	}

	// reject if order doesn't belong to customer (except admin)
	if !isAdmin {
		if order.CustomerID != uint(userID) {
			return web.WebError{
				Code:    400,
				Message: "Order doesn't belong to currently authenticated user",
			}
		}
	}

	// Update via repository
	order.Status = "CANCELLED"
	order.DriverID = null.IntFromPtr(nil)
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
	// validation
	err := validations.ValidateSimpleStruct(service.validate, createPaymentRequest)
	if err != nil {
		return entities.PaymentResponse{}, err
	}

	paymentRes := entities.PaymentResponse{}

	// order detail
	order, err := service.orderRepository.Find(orderID)
	if err != nil {
		return entities.PaymentResponse{}, web.WebError{ Code: 400, ProductionMessage: "Cannot get order details", DevelopmentMessage: "Error order detail: " + err.Error() }
	}

	// reject if status is other than confirmed
	if order.Status != "CONFIRMED" {
		return entities.PaymentResponse{}, web.WebError{ Code: 400, Message: "Transaction hasn't been confirmed or already been paid"}
	}

	// reject if order has transaction
	if order.TransactionID != "" {
		return entities.PaymentResponse{}, web.WebError{ Code: 400, Message: "Transaction has been set for this order" }
	}

	// Process payment by method
	switch strings.ToUpper(createPaymentRequest.PaymentMethod) {
	case "BANK_TRANSFER_BCA":
		paymentRes, err = service.paymentRepository.CreateBankTransferBCA(order)
	case "BANK_TRANSFER_BNI":
		paymentRes, err = service.paymentRepository.CreateBankTransferBNI(order)
	case "BANK_TRANSFER_BRI":
		paymentRes, err = service.paymentRepository.CreateBankTransferBRI(order)
	case "BANK_TRANSFER_MANDIRI":
		paymentRes, err = service.paymentRepository.CreateBankTransferMandiri(order)
	case "BANK_TRANSFER_PERMATA":
		paymentRes, err = service.paymentRepository.CreateBankTransferPermata(order)
	default:
		return entities.PaymentResponse{}, web.WebError{ Code: 400, Message: "Invalid payment method" }
	}
	if err != nil {
		return entities.PaymentResponse{}, err
	}

	// set order transaction
	order.DriverID = null.IntFromPtr(nil)
	order.TransactionID = paymentRes.TransactionID
	order.PaymentMethod = createPaymentRequest.PaymentMethod
	_, err = service.orderRepository.Update(order, int(order.ID))
	if err != nil {
		return entities.PaymentResponse{}, err
	}
	service.orderHistoryRepository.Create(orderID, "Pesanan sudah di bayarkan", "customer")

	return paymentRes, nil
}

/*
 * Get Payment
 * -------------------------------
 * Mengambil data pembayaran yang sudah ada berdasarkan transaction_id yang sudah di set pada order
 * @var orderID					ID Order yang akan di cancel
 * @var createPaymentRequest	request payment
 * @return PaymentResponse		response payment
 */
func (service OrderService) GetPayment(orderID int) (entities.PaymentResponse, error) {
	// get order
	order, err := service.orderRepository.Find(orderID)
	if err != nil {
		return entities.PaymentResponse{}, nil
	}
	
	// get payment
	paymentRes, err := service.paymentRepository.GetPaymentStatus(order.TransactionID, order.PaymentMethod)
	if err != nil {
		return entities.PaymentResponse{}, err
	}
	return paymentRes, nil
}

/*
 * Cancel Order payment
 * -------------------------------
 * Cancel payment order yang sudah dibuat
 * @var orderID 	order id terkait
 * @return 			error
 */
func (service OrderService) CancelPayment(orderID int) error {
	// get order
	order, err := service.orderRepository.Find(orderID)
	if err != nil {
		return nil
	}
	// reject if status is other than confirmed
	if order.Status != "CONFIRMED" {
		return web.WebError{ Code: 400, Message: "Order hasn't been confirmed or already been paid"}
	}
	// reject if transaction_id is empty
	if order.TransactionID == "" {
		return web.WebError{
			Code:    400,
			Message: "Payment transaction hasn't been created or already been cancelled",
		}
	}

	// repository action
	err = service.paymentRepository.CancelPayment(order.TransactionID, order.PaymentMethod)
	if err != nil {
		return err
	}
	// update order
	order.DriverID = null.IntFromPtr(nil)
	order.PaymentMethod = ""
	order.TransactionID = ""
	_, err = service.orderRepository.Update(order, orderID)
	if err != nil {
		return web.WebError{ Code: 500, ProductionMessage: "Failed to update order data", DevelopmentMessage: "update order fail:" + err.Error() }
	}
	return nil
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
 * @var order 		order id
 * @var status 		payment status
 * @return error	error
 */
func (service OrderService) PaymentWebhook(orderID int, status string) error {
	order, err := service.orderRepository.Find(orderID)
	if err != nil {
		return err
	}
	if status == "settlement" {
		// if status settlement, set order to MANIFESTED
	 	order.Status = "MANIFESTED"
	} 
	order.DriverID = null.IntFromPtr(nil)
	_, err = service.orderRepository.Update(order, orderID)
	if err != nil {
		return err
	}
	return nil
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
		return web.WebError{Code: 500, ProductionMessage: "The requested Driver ID doesn't match with any record", DevelopmentMessage: err.Error()}
	}
	if driver.Status == "BUSY" {
		return web.WebError{Code: 400, Message: "Finish your current order first"}
	}
	if !order.DriverID.Valid {
		return web.WebError{Code: 400, Message: "This order already taken by someone else"}
	}
	if order.TruckTypeID != driver.TruckTypeID {
		return web.WebError{Code: 400, Message: "Cannot take order that doesn't match with your truck type"}
	}
	if order.Status != "MANIFESTED" {
		return web.WebError{Code: 400, ProductionMessage: "This order hasn't been paid for by the customer"}
	}
	driverID := int64(driver.ID)
	order.DriverID = null.IntFromPtr(&driverID)
	order.Status = "ON_PROCESS"
	order, err = service.orderRepository.Update(order, userID)
	if err != nil {
		return web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "Cannot update order: " + err.Error()}
	}
	driver.Status = "BUSY"
	driver, err = service.userRepository.UpdateDriver(driver, userID)
	if err != nil {
		return web.WebError{Code: 500, ProductionMessage: "server error", DevelopmentMessage: "Cannot update driver:" + err.Error()}
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
		return err
	}
	order, err := service.orderRepository.Find(orderID)
	if err != nil {
		return err
	}
	driver, err := service.userRepository.FindByDriver("user_id", strconv.Itoa(userID))
	if err != nil {
		return err
	}
	if order.Status != "ON_PROCESS" {
		return web.WebError{Code: 400, Message: "Order wasn't on process"}
	}
	if !order.DriverID.Valid {
		return web.WebError{Code: 400, Message: "The current order has not belong to any driver, take the order first"}
	}
	if uint(order.DriverID.Int64) != driver.ID {
		return web.WebError{Code: 400, Message: "Cannot finish order that belong to someone else"}
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
			order.ArrivedPicture = fileUrl
		}
	}

	order.Status = "DELIVERED"
	order, err = service.orderRepository.Update(order, userID)
	if err != nil {
		return web.WebError{Code: 500, ProductionMessage: "Cannot update current order", DevelopmentMessage: "Update order error: " + err.Error()}
	}
	driver.Status = "IDLE"
	driver, err = service.userRepository.UpdateDriver(driver, userID)
	if err != nil {
		return web.WebError{Code: 500, ProductionMessage: "Cannot update current driver", DevelopmentMessage: "Update driver error: " + err.Error()}
	}
	// Log
	service.orderHistoryRepository.Create(int(order.ID), "Order yang diantar "+driver.User.Name+" telah sampai di tujuan", "driver")

	return err
}
