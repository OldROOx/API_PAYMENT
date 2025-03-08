package gateways

import (
	"app.payment/src/payments/application"
	"app.payment/src/payments/domain/entities"
	"errors"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

// MockPaymentGateway simula un procesador de pagos para pruebas
type MockPaymentGateway struct {
	failureRate float64
}

func NewMockPaymentGateway(failureRate float64) *MockPaymentGateway {
	return &MockPaymentGateway{
		failureRate: failureRate,
	}
}

func (g *MockPaymentGateway) ProcessPayment(payment *entities.Payment) (application.GatewayResponse, error) {
	// Simular un pequeño retraso
	time.Sleep(time.Millisecond * 500)

	// Inicializar generador de números aleatorios
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Simular falla según la tasa configurada
	if r.Float64() < g.failureRate {
		return application.GatewayResponse{
			Success:      false,
			ErrorMessage: "payment processor error: card declined",
		}, errors.New("payment processor error: card declined")
	}

	// Simular pago exitoso
	return application.GatewayResponse{
		Success:       true,
		TransactionID: uuid.New().String(),
	}, nil
}
