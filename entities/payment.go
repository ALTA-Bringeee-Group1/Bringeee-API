package entities


type MidtransBankTransferBCAResponse struct {
	Currency string `json:"currency"`
	FraudStatus string `json:"fraud_status"`
	GrossAmount string `json:"gross_amount"`
	MerchantID string `json:"merchant_id"`
	OrderID string `json:"order_id"`
	PaymentType string `json:"payment_type"`
	StatusCode string `json:"status_code"`
	StatusMessage string `json:"status_message"`
	TransactionID string `json:"transaction_id"`
	TransactionStatus string `json:"transaction_status"`
	TransactionTime string `json:"transaction_time"`
	VaNumbers []struct {
		Bank string `json:"bank"`
		VaNumber string `json:"va_number"`
	} `json:"va_numbers"`
}

type MidtransBankTransferBNIResponse struct {
	Currency string `json:"currency"`
	FraudStatus string `json:"fraud_status"`
	GrossAmount string `json:"gross_amount"`
	MerchantID string `json:"merchant_id"`
	OrderID string `json:"order_id"`
	PaymentType string `json:"payment_type"`
	StatusCode string `json:"status_code"`
	StatusMessage string `json:"status_message"`
	TransactionID string `json:"transaction_id"`
	TransactionStatus string `json:"transaction_status"`
	TransactionTime string `json:"transaction_time"`
	VaNumbers []struct {
		Bank string `json:"bank"`
		VaNumber string `json:"va_number"`
	} `json:"va_numbers"`
}

type MidtransBankTransferBRIResponse struct {
	Currency string `json:"currency"`
	FraudStatus string `json:"fraud_status"`
	GrossAmount string `json:"gross_amount"`
	MerchantID string `json:"merchant_id"`
	OrderID string `json:"order_id"`
	PaymentType string `json:"payment_type"`
	StatusCode string `json:"status_code"`
	StatusMessage string `json:"status_message"`
	TransactionID string `json:"transaction_id"`
	TransactionStatus string `json:"transaction_status"`
	TransactionTime string `json:"transaction_time"`
	VaNumbers []struct {
		Bank string `json:"bank"`
		VaNumber string `json:"va_number"`
	} `json:"va_numbers"`
}

type MidtransBankTransferMandiriResponse struct {
	BillKey string `json:"bill_key"`
	BillerCode string `json:"biller_code"`
	Currency string `json:"currency"`
	FraudStatus string `json:"fraud_status"`
	GrossAmount string `json:"gross_amount"`
	MerchantID string `json:"merchant_id"`
	OrderID string `json:"order_id"`
	PaymentType string `json:"payment_type"`
	StatusCode string `json:"status_code"`
	StatusMessage string `json:"status_message"`
	TransactionID string `json:"transaction_id"`
	TransactionStatus string `json:"transaction_status"`
	TransactionTime string `json:"transaction_time"`
}

type MidtransBankTransferPermataResponse struct {
	Currency string `json:"currency"`
	FraudStatus string `json:"fraud_status"`
	GrossAmount string `json:"gross_amount"`
	MerchantID string `json:"merchant_id"`
	OrderID string `json:"order_id"`
	PaymentType string `json:"payment_type"`
	StatusCode string `json:"status_code"`
	StatusMessage string `json:"status_message"`
	TransactionID string `json:"transaction_id"`
	TransactionStatus string `json:"transaction_status"`
	TransactionTime string `json:"transaction_time"`
	PermataVaNumber string `json:"permata_va_number"`
}