package application

import (
	"app.payment/src/payments/domain/entities"
)

type GatewayResponse struct {
	Success       bool   `json:"success"`
	TransactionID string `json:"transaction_id,omitempty"`
	ErrorMessage  string `json:"error_message,omitempty"`
}

type PaymentGateway interface {
	ProcessPayment(payment *entities.Payment) (GatewayResponse, error)
}
