package payment

import (
	"bringeee-capstone/configs"
	"bringeee-capstone/entities"
	"bringeee-capstone/entities/web"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
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
func (repository MidtransPaymentRepository) CreateBankTransferBCA(order entities.Order) (interface{}, error) {
	// prepare request
	requestBody, _ := json.Marshal(map[string]interface{}{
		"payment_type": "bank_transfer",
		"transaction_details": map[string]interface{} {
			"order_id": order.ID,
			"gross_amount": 90000,
		},
		"bank_transfer": map[string]interface{} {
			"bank": "bca",
		},
	})
	request, err := http.NewRequest(http.MethodPost, repository.baseURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, web.WebError{ Code: 500, ProductionMessage: "Payment server error", DevelopmentMessage: err.Error() }
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Basic " + base64.StdEncoding.EncodeToString([]byte(configs.Get().Payment.MidtransServerKey + ":")))
	request.Header.Set("Content-Type", "application/json")

	// do request
	response, err := repository.client.Do(request)
	if err != nil {
		return nil, web.WebError{ Code: 500, ProductionMessage: "Payment server error", DevelopmentMessage: "HTTP Response failed: " + err.Error() }
	}
	defer response.Body.Close()

	// parse response
	var data interface{}
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return nil, web.WebError{ Code: 500, ProductionMessage: "Payment server error", DevelopmentMessage: "Parsing response failed: " + err.Error() }
	}

	return data, nil
}

/*
 * Create Bank Transfer BNI
 * -------------------------------
 * Buat pembayaran untuk order tertentu menggunakan BNI
 *
 * @var order	Entity domain order yang dibuatkan pembayaran
 * @return any	Response dari layanan pihak ketiga
 */
func (repository MidtransPaymentRepository) CreateBankTransferBNI(order entities.Order) (interface{}, error) {
	// prepare request
	requestBody, _ := json.Marshal(map[string]interface{}{
		"payment_type": "bank_transfer",
		"transaction_details": map[string]interface{} {
			"order_id": order.ID,
			"gross_amount": 90000,
		},
		"bank_transfer": map[string]interface{} {
			"bank": "bni",
		},
	})
	request, err := http.NewRequest(http.MethodPost, repository.baseURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, web.WebError{ Code: 500, ProductionMessage: "Payment server error", DevelopmentMessage: err.Error() }
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Basic " + base64.StdEncoding.EncodeToString([]byte(configs.Get().Payment.MidtransServerKey + ":")))
	request.Header.Set("Content-Type", "application/json")

	// do request
	response, err := repository.client.Do(request)
	if err != nil {
		return nil, web.WebError{ Code: 500, ProductionMessage: "Payment server error", DevelopmentMessage: "HTTP Response failed: " + err.Error() }
	}
	defer response.Body.Close()

	// parse response
	var data interface{}
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return nil, web.WebError{ Code: 500, ProductionMessage: "Payment server error", DevelopmentMessage: "Parsing response failed: " + err.Error() }
	}

	return data, nil
}

/*
 * Create Bank Transfer BRI
 * -------------------------------
 * Buat pembayaran untuk order tertentu menggunakan BNI
 *
 * @var order	Entity domain order yang dibuatkan pembayaran
 * @return any	Response dari layanan pihak ketiga
 */
func (repository MidtransPaymentRepository) CreateBankTransferBRI(order entities.Order) (interface{}, error) {
	// prepare request
	requestBody, _ := json.Marshal(map[string]interface{}{
		"payment_type": "bank_transfer",
		"transaction_details": map[string]interface{} {
			"order_id": order.ID,
			"gross_amount": 90000,
		},
		"bank_transfer": map[string]interface{} {
			"bank": "bri",
		},
	})
	request, err := http.NewRequest(http.MethodPost, repository.baseURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, web.WebError{ Code: 500, ProductionMessage: "Payment server error", DevelopmentMessage: err.Error() }
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Basic " + base64.StdEncoding.EncodeToString([]byte(configs.Get().Payment.MidtransServerKey + ":")))
	request.Header.Set("Content-Type", "application/json")

	// do request
	response, err := repository.client.Do(request)
	if err != nil {
		return nil, web.WebError{ Code: 500, ProductionMessage: "Payment server error", DevelopmentMessage: "HTTP Response failed: " + err.Error() }
	}
	defer response.Body.Close()

	// parse response
	var data interface{}
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return nil, web.WebError{ Code: 500, ProductionMessage: "Payment server error", DevelopmentMessage: "Parsing response failed: " + err.Error() }
	}

	return data, nil
}

/*
 * Create Bank Transfer BRI
 * -------------------------------
 * Buat pembayaran untuk order tertentu menggunakan Mandiri
 *
 * @var order	Entity domain order yang dibuatkan pembayaran
 * @return any	Response dari layanan pihak ketiga
 */
func (repository MidtransPaymentRepository) CreateBankTransferMandiri(order entities.Order) (interface{}, error) {
	// prepare request
	requestBody, _ := json.Marshal(map[string]interface{}{
		"payment_type": "echannel",
		"transaction_details": map[string]interface{} {
			"order_id": order.ID,
			"gross_amount": 90000,
		},
		"echannel": map[string]interface{} {
			"bill_info1": "Payment:",
			"bill_info2": "Online purchase",
		},
	})
	request, err := http.NewRequest(http.MethodPost, repository.baseURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, web.WebError{ Code: 500, ProductionMessage: "Payment server error", DevelopmentMessage: err.Error() }
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Basic " + base64.StdEncoding.EncodeToString([]byte(configs.Get().Payment.MidtransServerKey + ":")))
	request.Header.Set("Content-Type", "application/json")

	// do request
	response, err := repository.client.Do(request)
	if err != nil {
		return nil, web.WebError{ Code: 500, ProductionMessage: "Payment server error", DevelopmentMessage: "HTTP Response failed: " + err.Error() }
	}
	defer response.Body.Close()

	// parse response
	var data interface{}
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return nil, web.WebError{ Code: 500, ProductionMessage: "Payment server error", DevelopmentMessage: "Parsing response failed: " + err.Error() }
	}

	return data, nil
}

/*
 * Create Bank Transfer BRI
 * -------------------------------
 * Buat pembayaran untuk order tertentu menggunakan Permata
 *
 * @var order	Entity domain order yang dibuatkan pembayaran
 * @return any	Response dari layanan pihak ketiga
 */
func (repository MidtransPaymentRepository) CreateBankTransferPermata(order entities.Order) (interface{}, error) {
	// prepare request
	requestBody, _ := json.Marshal(map[string]interface{}{
		"payment_type": "permata",
		"transaction_details": map[string]interface{} {
			"order_id": order.ID,
			"gross_amount": 90000,
		},
	})
	request, err := http.NewRequest(http.MethodPost, repository.baseURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, web.WebError{ Code: 500, ProductionMessage: "Payment server error", DevelopmentMessage: err.Error() }
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Basic " + base64.StdEncoding.EncodeToString([]byte(configs.Get().Payment.MidtransServerKey + ":")))
	request.Header.Set("Content-Type", "application/json")

	// do request
	response, err := repository.client.Do(request)
	if err != nil {
		return nil, web.WebError{ Code: 500, ProductionMessage: "Payment server error", DevelopmentMessage: "HTTP Response failed: " + err.Error() }
	}
	defer response.Body.Close()

	// parse response
	var data interface{}
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return nil, web.WebError{ Code: 500, ProductionMessage: "Payment server error", DevelopmentMessage: "Parsing response failed: " + err.Error() }
	}

	return data, nil
}