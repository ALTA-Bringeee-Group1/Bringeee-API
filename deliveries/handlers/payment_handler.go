package handlers

import (
	"bringeee-capstone/configs"
	"bringeee-capstone/deliveries/helpers"
	"bringeee-capstone/entities/web"
	orderService "bringeee-capstone/services/order"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type PaymentHandler struct {
	orderService orderService.OrderServiceInterface
}

func NewPaymentHandler(orderService orderService.OrderServiceInterface) *PaymentHandler {
	return &PaymentHandler{
		orderService: orderService,
	}
}

type MidtransHookRequest struct {
	TransactionTime string `form:"transaction_time" json:"transaction_time"`
	TransactionStatus string `form:"transaction_status" json:"transaction_status"`
	OrderID string `form:"order_id" json:"order_id"`
	MerchantID string `form:"merchant_id" json:"merchant_id"`
	GrossAmount string `form:"gross_amount" json:"gross_amount"`
	FraudStatus string `form:"fraud_status" json:"fraud_status"`
	Currency string `form:"currency" json:"currency"`
} 


/*
 * Midtrans Payment notification Webhook
 * -------------------------------
 * Payment Webhook notification, dikirimkan oleh layanan pihak ketiga
 * referensi: https://docs.midtrans.com/en/after-payment/http-notification
 * endpoint: POST /api/payments/midtrans_webhook
 */
func (handler PaymentHandler) MidtransWebhook(c echo.Context) error {
	links := map[string]string{}	
	links["self"] = fmt.Sprintf("%s/api/payments/midtrans_webhook", configs.Get().App.BaseURL)

	midtransRequest := MidtransHookRequest{}
	c.Bind(&midtransRequest)
	if midtransRequest.TransactionStatus == "" {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Status: "ERROR",
			Code: http.StatusBadRequest,
			Error: "transaction status must be provided",
			Links: links,
		})
	}
	orderID, err := strconv.Atoi(midtransRequest.OrderID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Status: "ERROR",
			Code: http.StatusBadRequest,
			Error: "order id request body is invalid",
			Links: links,
		})
	}

	// service call
	err = handler.orderService.PaymentWebhook(orderID, midtransRequest.TransactionStatus)
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	}

	return c.JSON(http.StatusOK, web.SuccessResponse{
		Status: "OK",
		Code: http.StatusOK,
		Error: nil,
		Links: links,
		Data: map[string]interface{}{
			"id": orderID,
		},
	})
}