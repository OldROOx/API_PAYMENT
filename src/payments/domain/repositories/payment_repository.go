package repositories

import (
	"app.payment/src/payments/domain/entities"
)

type PaymentRepository interface {
	Save(payment *entities.Payment) error
	FindByID(id string) (*entities.Payment, error)
	FindByOrderID(orderID uint) (*entities.Payment, error)
	UpdateStatus(id string, status entities.PaymentStatus, transactionID string) error
	UpdatePaymentFailed(id string, errorMessage string) error
}
