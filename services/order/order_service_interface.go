package order

import (
	"bringeee-capstone/entities"
	"bringeee-capstone/entities/web"
	"mime/multipart"
)

type OrderServiceInterface interface {

	/*
	 * Find All
	 * -------------------------------
	 * Mengambil data order berdasarkan filters dan sorts
	 *
	 * @var limit 	batas limit hasil query
	 * @var offset 	offset hasil query
	 * @var filters	query untuk penyaringan data, { field, operator, value }
	 * @var sorts	pengurutan data, { field, value[bool] }
	 * @return order			list order dalam bentuk entity response
	 * @return error			error
	 */
	FindAll(limit int, offset int, filters []map[string]interface{}, sorts []map[string]interface{}) ([]entities.OrderResponse, error)

	/*
	 * Find First
	 * -------------------------------
	 * Mengambil order pertama berdasarkan filter yang telah di tentukan pada parameter
	 * dan mengambil data pertama sebagai data order tunggal
	 * @var filter 
	 * @return OrderResponse	order response dalam bentuk tunggal
	 * @return error			error
	 */
	FindFirst(filters []map[string]interface{}) (entities.OrderResponse, error)

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
	GetPagination(limit int, page int, filters []map[string]interface{}) (web.Pagination, error)

	/*
	 * Find
	 * -------------------------------
	 * Mengambil data order tunggal berdasarkan ID
	 *
	 * @var id 			id order
	 * @return order	order tunggal dalam bentuk response
	 * @return error	error
	 */
	Find(id int) (entities.OrderResponse, error)

	/*
	 * Customer Create order
	 * -------------------------------
	 * Membuat order baru berdasarkan user yang sedang login
	 * @var orderRequest		request create order oleh customer
	 * @var files				list file request untuk diteruskan ke validation dan upload
	 * @return OrderResponse	order response 
	 */
	Create(orderRequest entities.CustomerCreateOrderRequest, files map[string]*multipart.FileHeader, userID int) (entities.OrderResponse, error)

	/*
	 * Admin Set fixed price order
	 * -------------------------------
	 * Set fixed price order oleh admin untuk diteruskan kembali ke user agar di konfirmasi/cancel
	 * @var orderRequest		request create order oleh customer
	 * @var orderID				orderID
	 * @return OrderResponse	order response 
	 */
	SetFixOrder(orderID int, setPriceRequest entities.AdminSetPriceOrderRequest) error 

	/*
	 * Confirm Order
	 * -------------------------------
	 * Confirm order yang sudah dibuat
	 * @var orderID				ID Order yang akan di cancel
	 * @var userID 				authenticated user id (role: customer, admin)
	 * @return OrderResponse	order response 
	 */
	ConfirmOrder(orderID int, userID int, isAdmin bool) error 

	/*
	 * Cancel Order
	 * -------------------------------
	 * Cancel order yang sudah dibuat
	 * @var orderID				ID Order yang akan di cancel
	 * @var userID				Authenticated user id (role: customer, admin)
	 * @return OrderResponse	order response 
	 */
	CancelOrder(orderID int, userID int, isAdmin bool) error 

	/*
	 * Create Payment
	 * -------------------------------
	 * Buat pembayaran untuk order tertentu ke layanan pihak ketiga
	 * @var orderID					ID Order yang akan di cancel
	 * @var createPaymentRequest	request payment
	 * @return PaymentResponse		response payment 
	 */
	CreatePayment(orderID int, createPaymentRequest entities.CreatePaymentRequest) (entities.PaymentResponse, error) 

	/*
	 * Get Payment
	 * -------------------------------
	 * Mengambil data pembayaran yang sudah ada berdasarkan transaction_id yang sudah di set pada order
	 * @var orderID					ID Order yang akan di cancel
	 * @return PaymentResponse		response payment 
	 */
	GetPayment(orderID int) (entities.PaymentResponse, error)
	
	/*
	 * Find All History
	 * -------------------------------
	 * Mengambil data order berdasarkan filters dan sorts
	 *
	 * @var limit 	batas limit hasil query
	 * @var offset 	offset hasil query
	 * @var sorts	pengurutan data, { field, value[bool] }
	 * @return order	list order dalam bentuk entity response
	 * @return error	error
	 */
	FindAllHistory(orderID int, sorts []map[string]interface{}) ([]entities.OrderHistoryResponse, error)

	/*
	 * Webhook
	 * -------------------------------
	 * Payment Webhook notification, dikirimkan oleh layanan pihak ketiga
	 * referensi: https://docs.midtrans.com/en/after-payment/http-notification
	 *
	 * @var orer 		order id
	 * @return error	error
	 */
	PaymentWebhook(orderID int) error

	/*
	 * Take order for shipping
	 * -------------------------------
	 * Pengambilan order oleh driver untuk di set statusnya menjadi ON_PROCESS
	 * @var orderID 	order id terkait
	 * @var userID		authenticated user (role: driver)
	 */
	TakeOrder(orderID int, userID int) error

	/*
	 * Finish Order
	 * -------------------------------
	 * penyelesaian order oleh driver untuk di set statusnya menjadi DELIVERED
	 * @var orderID 	order id terkait
	 * @var userID		authenticated user (role: driver)
	 * @var files		list file request untuk diteruskan ke validation dan upload
	 */
	FinishOrder(orderID int, userID int, files map[string]*multipart.FileHeader) error
}