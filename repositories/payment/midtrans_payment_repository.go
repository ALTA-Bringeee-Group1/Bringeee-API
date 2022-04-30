package payment

import (
	"bringeee-capstone/configs"
	"bringeee-capstone/entities"
	"bringeee-capstone/entities/web"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type MidtransPaymentRepository struct {
	baseURL string
	client *http.Client
}

func NewMidtransPaymentRepository() *MidtransPaymentRepository {
	baseUrl := ""
	if configs.Get().App.ENV == "production" {
		baseUrl = configs.Get().Payment.MidtransBaseURLProduction
	} else {
		baseUrl = configs.Get().Payment.MidtransBaseURLDevelopment
	}
	return &MidtransPaymentRepository{
		baseURL: baseUrl,
		client: &http.Client{},
	}
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
func (repository MidtransPaymentRepository) CreateBankTransferBCA(order entities.Order) (entities.PaymentResponse, error) {
	// prepare request
	requestBody, _ := json.Marshal(map[string]interface{}{
		"payment_type": "bank_transfer",
		"transaction_details": map[string]interface{} {
			"order_id": order.ID,
			"gross_amount": order.FixPrice,
		},
		"bank_transfer": map[string]interface{} {
			"bank": "bca",
		},
	})
	request, err := http.NewRequest(http.MethodPost, repository.baseURL + "/charge", bytes.NewBuffer(requestBody))
	if err != nil {
		return entities.PaymentResponse{}, web.WebError{ Code: 500, ProductionMessage: "Payment server error", DevelopmentMessage: err.Error() }
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Basic " + base64.StdEncoding.EncodeToString([]byte(configs.Get().Payment.MidtransServerKey + ":")))
	request.Header.Set("Content-Type", "application/json")

	// do request
	response, err := repository.client.Do(request)
	if err != nil {
		return entities.PaymentResponse{}, web.WebError{ Code: 500, ProductionMessage: "Payment server error", DevelopmentMessage: "HTTP Response failed: " + err.Error() }
	}
	defer response.Body.Close()

	// parse response
	var data entities.MidtransBankTransferBCAResponse
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return entities.PaymentResponse{}, web.WebError{ Code: 500, ProductionMessage: "Payment server error", DevelopmentMessage: "Parsing response failed: " + err.Error() }
	}
	if data.StatusCode != "201" {
		return entities.PaymentResponse{}, web.WebError{ Code: 500, Message: "Error creating payment transaction" }
	}

	// translate response
	grossAmount, _ := strconv.ParseFloat(data.GrossAmount, 64)
	trTime, _ := time.Parse("2006-01-02 15:04:05", data.TransactionTime)
	paymentRes := entities.PaymentResponse {
		OrderID: strconv.Itoa(int(order.ID)),
		TransactionID: data.TransactionID,
		PaymentMethod: "BANK_TRANSFER_BCA",
		BillNumber: data.VaNumbers[0].VaNumber,
		Bank: data.VaNumbers[0].Bank,
		GrossAmount: int64(grossAmount),
		TransactionTime: trTime,
		TransactionExpire: trTime.Add(time.Hour * 24),
	}
	return paymentRes, nil
}

/*
 * Create Bank Transfer BNI
 * -------------------------------
 * Buat pembayaran untuk order tertentu menggunakan BNI
 *
 * @var order	Entity domain order yang dibuatkan pembayaran
 * @return any	Response dari layanan pihak ketiga
 */
func (repository MidtransPaymentRepository) CreateBankTransferBNI(order entities.Order) (entities.PaymentResponse, error) {
	// prepare request
	requestBody, _ := json.Marshal(map[string]interface{}{
		"payment_type": "bank_transfer",
		"transaction_details": map[string]interface{} {
			"order_id": order.ID,
			"gross_amount": order.FixPrice,
		},
		"bank_transfer": map[string]interface{} {
			"bank": "bni",
		},
	})
	request, err := http.NewRequest(http.MethodPost, repository.baseURL + "/charge", bytes.NewBuffer(requestBody))
	if err != nil {
		return entities.PaymentResponse{}, web.WebError{ Code: 500, ProductionMessage: "Payment server error", DevelopmentMessage: err.Error() }
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Basic " + base64.StdEncoding.EncodeToString([]byte(configs.Get().Payment.MidtransServerKey + ":")))
	request.Header.Set("Content-Type", "application/json")

	// do request
	response, err := repository.client.Do(request)
	if err != nil {
		return entities.PaymentResponse{}, web.WebError{ Code: 500, ProductionMessage: "Payment server error", DevelopmentMessage: "HTTP Response failed: " + err.Error() }
	}
	defer response.Body.Close()

	// parse response
	var data entities.MidtransBankTransferBNIResponse
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return entities.PaymentResponse{}, web.WebError{ Code: 500, ProductionMessage: "Payment server error", DevelopmentMessage: "Parsing response failed: " + err.Error() }
	}
	if data.StatusCode != "201" {
		return entities.PaymentResponse{}, web.WebError{ Code: 500, Message: "Error creating payment transaction" }
	}

	// translate response
	grossAmount, _ := strconv.ParseFloat(data.GrossAmount, 64)
	trTime, _ := time.Parse("2006-01-02 15:04:05", data.TransactionTime)
	paymentRes := entities.PaymentResponse {
		OrderID: strconv.Itoa(int(order.ID)),
		TransactionID: data.TransactionID,
		PaymentMethod: "BANK_TRANSFER_BNI",
		BillNumber: data.VaNumbers[0].VaNumber,
		Bank: data.VaNumbers[0].Bank,
		GrossAmount: int64(grossAmount),
		TransactionTime: trTime,
		TransactionExpire: trTime.Add(time.Hour * 24),
	}
	return paymentRes, nil
}

/*
 * Create Bank Transfer BRI
 * -------------------------------
 * Buat pembayaran untuk order tertentu menggunakan BNI
 *
 * @var order	Entity domain order yang dibuatkan pembayaran
 * @return any	Response dari layanan pihak ketiga
 */
func (repository MidtransPaymentRepository) CreateBankTransferBRI(order entities.Order) (entities.PaymentResponse, error) {
	// prepare request
	requestBody, _ := json.Marshal(map[string]interface{}{
		"payment_type": "bank_transfer",
		"transaction_details": map[string]interface{} {
			"order_id": order.ID,
			"gross_amount": order.FixPrice,
		},
		"bank_transfer": map[string]interface{} {
			"bank": "bri",
		},
	})
	request, err := http.NewRequest(http.MethodPost, repository.baseURL + "/charge", bytes.NewBuffer(requestBody))
	if err != nil {
		return entities.PaymentResponse{}, web.WebError{ Code: 500, ProductionMessage: "Payment server error", DevelopmentMessage: err.Error() }
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Basic " + base64.StdEncoding.EncodeToString([]byte(configs.Get().Payment.MidtransServerKey + ":")))
	request.Header.Set("Content-Type", "application/json")

	// do request
	response, err := repository.client.Do(request)
	if err != nil {
		return entities.PaymentResponse{}, web.WebError{ Code: 500, ProductionMessage: "Payment server error", DevelopmentMessage: "HTTP Response failed: " + err.Error() }
	}
	defer response.Body.Close()

	// parse response
	var data entities.MidtransBankTransferBRIResponse
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return entities.PaymentResponse{}, web.WebError{ Code: 500, ProductionMessage: "Payment server error", DevelopmentMessage: "Parsing response failed: " + err.Error() }
	}
	if data.StatusCode != "201" {
		return entities.PaymentResponse{}, web.WebError{ Code: 500, Message: "Error creating payment transaction" }
	}

	// translate response
	grossAmount, _ := strconv.ParseFloat(data.GrossAmount, 64)
	trTime, _ := time.Parse("2006-01-02 15:04:05", data.TransactionTime)
	paymentRes := entities.PaymentResponse {
		OrderID: strconv.Itoa(int(order.ID)),
		TransactionID: data.TransactionID,
		PaymentMethod: "BANK_TRANSFER_BRI",
		BillNumber: data.VaNumbers[0].VaNumber,
		Bank: data.VaNumbers[0].Bank,
		GrossAmount: int64(grossAmount),
		TransactionTime: trTime,
		TransactionExpire: trTime.Add(time.Hour * 24),
	}
	return paymentRes, nil
}

/*
 * Create Bank Transfer Mandiri
 * -------------------------------
 * Buat pembayaran untuk order tertentu menggunakan Mandiri
 *
 * @var order	Entity domain order yang dibuatkan pembayaran
 * @return any	Response dari layanan pihak ketiga
 */
func (repository MidtransPaymentRepository) CreateBankTransferMandiri(order entities.Order) (entities.PaymentResponse, error) {
	// prepare request
	requestBody, _ := json.Marshal(map[string]interface{}{
		"payment_type": "echannel",
		"transaction_details": map[string]interface{} {
			"order_id": order.ID,
			"gross_amount": order.FixPrice,
		},
		"echannel": map[string]interface{} {
			"bill_info1": "Payment:",
			"bill_info2": "Online purchase",
		},
	})
	request, err := http.NewRequest(http.MethodPost, repository.baseURL + "/charge", bytes.NewBuffer(requestBody))
	if err != nil {
		return entities.PaymentResponse{}, web.WebError{ Code: 500, ProductionMessage: "Payment server error", DevelopmentMessage: err.Error() }
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Basic " + base64.StdEncoding.EncodeToString([]byte(configs.Get().Payment.MidtransServerKey + ":")))
	request.Header.Set("Content-Type", "application/json")

	// do request
	response, err := repository.client.Do(request)
	if err != nil {
		return entities.PaymentResponse{}, web.WebError{ Code: 500, ProductionMessage: "Payment server error", DevelopmentMessage: "HTTP Response failed: " + err.Error() }
	}
	defer response.Body.Close()

	// parse response
	var data entities.MidtransBankTransferMandiriResponse
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return entities.PaymentResponse{}, web.WebError{ Code: 500, ProductionMessage: "Payment server error", DevelopmentMessage: "Parsing response failed: " + err.Error() }
	}
	if data.StatusCode != "201" {
		return entities.PaymentResponse{}, web.WebError{ Code: 500, Message: "Error creating payment transaction" }
	}

	// translate response
	grossAmount, _ := strconv.ParseFloat(data.GrossAmount, 64)
	trTime, _ := time.Parse("2006-01-02 15:04:05", data.TransactionTime)
	paymentRes := entities.PaymentResponse {
		OrderID: strconv.Itoa(int(order.ID)),
		TransactionID: data.TransactionID,
		PaymentMethod: "BANK_TRANSFER_MANDIRI",
		BillNumber: data.BillKey,
		Bank: "mandiri",
		GrossAmount: int64(grossAmount),
		TransactionTime: trTime,
		TransactionExpire: trTime.Add(time.Hour * 24),
	}
	return paymentRes, nil
}

/*
 * Create Bank Transfer Permata
 * -------------------------------
 * Buat pembayaran untuk order tertentu menggunakan Permata
 *
 * @var order	Entity domain order yang dibuatkan pembayaran
 * @return any	Response dari layanan pihak ketiga
 */
func (repository MidtransPaymentRepository) CreateBankTransferPermata(order entities.Order) (entities.PaymentResponse, error) {
	// prepare request
	requestBody, _ := json.Marshal(map[string]interface{}{
		"payment_type": "permata",
		"transaction_details": map[string]interface{} {
			"order_id": order.ID,
			"gross_amount": order.FixPrice,
		},
	})
	request, err := http.NewRequest(http.MethodPost, repository.baseURL + "/charge", bytes.NewBuffer(requestBody))
	if err != nil {
		return entities.PaymentResponse{}, web.WebError{ Code: 500, ProductionMessage: "Payment server error", DevelopmentMessage: err.Error() }
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Basic " + base64.StdEncoding.EncodeToString([]byte(configs.Get().Payment.MidtransServerKey + ":")))
	request.Header.Set("Content-Type", "application/json")

	// do request
	response, err := repository.client.Do(request)
	if err != nil {
		return entities.PaymentResponse{}, web.WebError{ Code: 500, ProductionMessage: "Payment server error", DevelopmentMessage: "HTTP Response failed: " + err.Error() }
	}
	defer response.Body.Close()

	// parse response
	var data entities.MidtransBankTransferPermataResponse
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return entities.PaymentResponse{}, web.WebError{ Code: 500, ProductionMessage: "Payment server error", DevelopmentMessage: "Parsing response failed: " + err.Error() }
	}
	if data.StatusCode != "201" {
		return entities.PaymentResponse{}, web.WebError{ Code: 500, Message: "Error creating payment transaction" }
	}

	// translate response
	grossAmount, _ := strconv.ParseFloat(data.GrossAmount, 64)
	trTime, _ := time.Parse("2006-01-02 15:04:05", data.TransactionTime)
	paymentRes := entities.PaymentResponse {
		OrderID: strconv.Itoa(int(order.ID)),
		TransactionID: data.TransactionID,
		PaymentMethod: "BANK_TRANSFER_PERMATA",
		BillNumber: data.PermataVaNumber,
		Bank: "permata",
		GrossAmount: int64(grossAmount),
		TransactionTime: trTime,
		TransactionExpire: trTime.Add(time.Hour * 24),
	}
	return paymentRes, nil
}

/*
 * Get Payment detail
 * -------------------------------
 * Mengambil data transaksi berdasarkan `transaction_id`
 *
 * @var transaction_id		Transaction ID
 * @return PaymentResponse	Response
 */
func (repository MidtransPaymentRepository) GetPaymentStatus(transactionID string) (entities.PaymentResponse, error) {
	// req, err := http.NewRequest(http.MethodGet, repository.baseURL + "/" + transactionID + "/status", nil)
	panic("as")
}