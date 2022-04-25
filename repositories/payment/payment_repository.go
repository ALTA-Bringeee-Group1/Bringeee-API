package payment

import (
	"bringeee-capstone/entities"
)

type PaymentRepository struct {}

func NewPaymentRepository() *PaymentRepository {
	return &PaymentRepository{}
}

/*
 * Create Bank Transfer BCA
 * -------------------------------
 * Buat pembayaran untuk order tertentu menggunakan BCA
 * Contoh referensi: https://docs.midtrans.com/en/core-api/bank-transfer
 *
 * @var order	Entity domain order yang dibuatkan pembayaran
 * @return any	Response dari layanan pihak ketiga
 */
func (repository PaymentRepository) CreateBankTransferBCA(order entities.Order) (interface{}, error) {
	panic("implement me")
}

/*
 * Create Bank Transfer BNI
 * -------------------------------
 * Buat pembayaran untuk order tertentu menggunakan BNI
 *
 * @var order	Entity domain order yang dibuatkan pembayaran
 * @return any	Response dari layanan pihak ketiga
 */
func (repository PaymentRepository) CreateBankTransferBNI(order entities.Order) (interface{}, error) {
	panic("implement me")
}

/*
 * Create Bank Transfer BRI
 * -------------------------------
 * Buat pembayaran untuk order tertentu menggunakan BNI
 *
 * @var order	Entity domain order yang dibuatkan pembayaran
 * @return any	Response dari layanan pihak ketiga
 */
func (repository PaymentRepository) CreateBankTransferBRI(order entities.Order) (interface{}, error) {
	panic("implement me")
}

/*
 * Create Bank Transfer BRI
 * -------------------------------
 * Buat pembayaran untuk order tertentu menggunakan Mandiri
 *
 * @var order	Entity domain order yang dibuatkan pembayaran
 * @return any	Response dari layanan pihak ketiga
 func (repository PaymentRepository) */
func (repository PaymentRepository) CreateBankTransferMandiri(order entities.Order) (interface{}, error) {
	panic("implement me")
}

/*
 * Create Bank Transfer BRI
 * -------------------------------
 * Buat pembayaran untuk order tertentu menggunakan Permata
 *
 * @var order	Entity domain order yang dibuatkan pembayaran
 * @return any	Response dari layanan pihak ketiga
 func (repository PaymentRepository) */
func (repository PaymentRepository) CreateBankTransferPermata(order entities.Order) (interface{}, error) {
	panic("implement me")
}