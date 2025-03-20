package application

import (
	"app.payment/src/payments/domain/entities"
)

// GatewayResponse representa la respuesta de un procesador de pagos
type GatewayResponse struct {
	Success       bool   `json:"success"`
	TransactionID string `json:"transaction_id,omitempty"`
	ErrorMessage  string `json:"error_message,omitempty"`
}

// PaymentGateway define las operaciones para procesar pagos
type PaymentGateway interface {
	ProcessPayment(payment *entities.Payment) (GatewayResponse, error)
}
