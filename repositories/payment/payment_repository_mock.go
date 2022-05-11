package payment

import (
	"bringeee-capstone/entities"
	"time"

	"github.com/stretchr/testify/mock"
)

type PaymentRepositoryMock struct {
	Mock *mock.Mock
}

func NewPaymentRepositoryMock(mock *mock.Mock) *PaymentRepositoryMock {
	return &PaymentRepositoryMock {
		Mock: mock,
	}
}

var PaymentResponseCollection = []entities.PaymentResponse {
	{
		OrderID: "1",
		TransactionID: "xya-7721d-ma",
		PaymentMethod: "BANK_TRANSFER_BNI",
		BillNumber: "5571230851238",
		Bank: "bni",
		GrossAmount: 880000,
		TransactionTime: time.Now(),
		TransactionExpire: time.Now().Add(time.Hour * 24),
	},
	{
		OrderID: "2",
		TransactionID: "xzc-7722d-ca",
		PaymentMethod: "BANK_TRANSFER_BCA",
		BillNumber: "5571230851238",
		Bank: "bca",
		GrossAmount: 880000,
		TransactionTime: time.Now(),
		TransactionExpire: time.Now().Add(time.Hour * 24),
	},
	{
		OrderID: "3",
		TransactionID: "xyd-7723d-ff",
		PaymentMethod: "BANK_TRANSFER_BRI",
		BillNumber: "5571230851238",
		Bank: "bri",
		GrossAmount: 880000,
		TransactionTime: time.Now(),
		TransactionExpire: time.Now().Add(time.Hour * 24),
	},
	{
		OrderID: "4",
		TransactionID: "xze-7724d-xu",
		PaymentMethod: "BANK_TRANSFER_PERMATA",
		BillNumber: "5571230851238",
		Bank: "permata",
		GrossAmount: 880000,
		TransactionTime: time.Now(),
		TransactionExpire: time.Now().Add(time.Hour * 24),
	},
	{
		OrderID: "5",
		TransactionID: "xyl-77215-de",
		PaymentMethod: "BANK_TRANSFER_MANDIRI",
		BillNumber: "5571230851238",
		Bank: "mandiri",
		GrossAmount: 880000,
		TransactionTime: time.Now(),
		TransactionExpire: time.Now().Add(time.Hour * 24),
	},
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
func (repo PaymentRepositoryMock) CreateBankTransferBCA(order entities.Order) (entities.PaymentResponse, error) {
	param := repo.Mock.Called()
	return param.Get(0).(entities.PaymentResponse), param.Error(1)
}
/*
 * Create Bank Transfer BNI
 * -------------------------------
 * Buat pembayaran untuk order tertentu menggunakan BNI
 *
 * @var order	Entity domain order yang dibuatkan pembayaran
 * @return any	Response dari layanan pihak ketiga
 */
func (repo PaymentRepositoryMock) CreateBankTransferBNI(order entities.Order) (entities.PaymentResponse, error) {
	param := repo.Mock.Called()
	return param.Get(0).(entities.PaymentResponse), param.Error(1)
}
/*
 * Create Bank Transfer BRI
 * -------------------------------
 * Buat pembayaran untuk order tertentu menggunakan BNI
 *
 * @var order	Entity domain order yang dibuatkan pembayaran
 * @return any	Response dari layanan pihak ketiga
 */
func (repo PaymentRepositoryMock) CreateBankTransferBRI(order entities.Order) (entities.PaymentResponse, error) {
	param := repo.Mock.Called()
	return param.Get(0).(entities.PaymentResponse), param.Error(1)
}
/*
 * Create Bank Transfer BRI
 * -------------------------------
 * Buat pembayaran untuk order tertentu menggunakan Mandiri
 *
 * @var order	Entity domain order yang dibuatkan pembayaran
 * @return any	Response dari layanan pihak ketiga
 */
func (repo PaymentRepositoryMock) CreateBankTransferMandiri(order entities.Order) (entities.PaymentResponse, error) {
	param := repo.Mock.Called()
	return param.Get(0).(entities.PaymentResponse), param.Error(1)
}
/*
 * Create Bank Transfer BRI
 * -------------------------------
 * Buat pembayaran untuk order tertentu menggunakan Permata
 *
 * @var order	Entity domain order yang dibuatkan pembayaran
 * @return any	Response dari layanan pihak ketiga
 */
func (repo PaymentRepositoryMock) CreateBankTransferPermata(order entities.Order) (entities.PaymentResponse, error) {
	param := repo.Mock.Called()
	return param.Get(0).(entities.PaymentResponse), param.Error(1)
}
/*
* Get Payment detail
* -------------------------------
* Mengambil data transaksi berdasarkan `transaction_id`
*
* @var transaction_id		Transaction ID
* @return PaymentResponse	Response
*/
func (repo PaymentRepositoryMock) GetPaymentStatus(transactionID string, paymentMethod string) (entities.PaymentResponse, error) {
	param := repo.Mock.Called()
	return param.Get(0).(entities.PaymentResponse), param.Error(1)
}
/*
* Cancel payment
* -------------------------------
* Membatalkan data transaksi berdasarkan `transaction_id`
*
* @var transaction_id		Transaction ID
* @return PaymentResponse	Response
*/
func (repo PaymentRepositoryMock) CancelPayment(transactionID string, paymentMethod string) error {
	param := repo.Mock.Called()
	return param.Error(0)
}

