package payment

import "bringeee-capstone/entities"

type PaymentRepositoryInterface interface {


	/*
	 * Create Bank Transfer BCA
	 * -------------------------------
	 * Buat pembayaran untuk order tertentu menggunakan BCA
	 * Contoh referensi: https://docs.midtrans.com/en/core-api/bank-transfer
	 *
	 * @var order	Entity domain order yang dibuatkan pembayaran
	 * @return any	Response dari layanan pihak ketiga
	 */
	CreateBankTransferBCA(order entities.Order) (entities.PaymentResponse, error)

	/*
	 * Create Bank Transfer BNI
	 * -------------------------------
	 * Buat pembayaran untuk order tertentu menggunakan BNI
	 *
	 * @var order	Entity domain order yang dibuatkan pembayaran
	 * @return any	Response dari layanan pihak ketiga
	 */
	CreateBankTransferBNI(order entities.Order) (entities.PaymentResponse, error)

	/*
	 * Create Bank Transfer BRI
	 * -------------------------------
	 * Buat pembayaran untuk order tertentu menggunakan BNI
	 *
	 * @var order	Entity domain order yang dibuatkan pembayaran
	 * @return any	Response dari layanan pihak ketiga
	 */
	CreateBankTransferBRI(order entities.Order) (entities.PaymentResponse, error)

	/*
	 * Create Bank Transfer BRI
	 * -------------------------------
	 * Buat pembayaran untuk order tertentu menggunakan Mandiri
	 *
	 * @var order	Entity domain order yang dibuatkan pembayaran
	 * @return any	Response dari layanan pihak ketiga
	 */
	CreateBankTransferMandiri(order entities.Order) (entities.PaymentResponse, error)

	/*
	 * Create Bank Transfer BRI
	 * -------------------------------
	 * Buat pembayaran untuk order tertentu menggunakan Permata
	 *
	 * @var order	Entity domain order yang dibuatkan pembayaran
	 * @return any	Response dari layanan pihak ketiga
	 */
	CreateBankTransferPermata(order entities.Order) (entities.PaymentResponse, error)

	/*
	* Get Payment detail
	* -------------------------------
	* Mengambil data transaksi berdasarkan `transaction_id`
	*
	* @var transaction_id		Transaction ID
	* @return PaymentResponse	Response
	*/
	GetPaymentStatus(transactionID string, paymentMethod string) (entities.PaymentResponse, error)

	/*
	* Cancel payment
	* -------------------------------
	* Membatalkan data transaksi berdasarkan `transaction_id`
	*
	* @var transaction_id		Transaction ID
	* @return PaymentResponse	Response
	*/
	CancelPayment(transactionID string, paymentMethod string) error
}