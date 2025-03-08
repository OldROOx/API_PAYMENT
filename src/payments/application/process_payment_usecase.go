package application

import (
	"app.payment/src/payments/domain/entities"
	"app.payment/src/payments/domain/repositories"
	"errors"
	"github.com/google/uuid"
	"time"
)

type ProcessPaymentUseCase struct {
	paymentRepo    repositories.PaymentRepository
	eventPublisher repositories.EventPublisher
	paymentGateway PaymentGateway
}

func NewProcessPaymentUseCase(
	paymentRepo repositories.PaymentRepository,
	eventPublisher repositories.EventPublisher,
	paymentGateway PaymentGateway,
) *ProcessPaymentUseCase {
	return &ProcessPaymentUseCase{
		paymentRepo:    paymentRepo,
		eventPublisher: eventPublisher,
		paymentGateway: paymentGateway,
	}
}

func (uc *ProcessPaymentUseCase) Execute(orderID uint, amount float64, method string) (*entities.Payment, error) {
	// Check if payment already exists for this order
	existingPayment, err := uc.paymentRepo.FindByOrderID(orderID)
	if err == nil && existingPayment != nil {
		// Payment already exists
		return existingPayment, errors.New("payment already exists for this order")
	}

	// Create new payment
	paymentID := uuid.New().String()
	payment := &entities.Payment{
		ID:        paymentID,
		OrderID:   orderID,
		Amount:    amount,
		Method:    method,
		Status:    entities.PaymentStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save to repository
	if err := uc.paymentRepo.Save(payment); err != nil {
		return nil, err
	}

	// Process payment with gateway
	gatewayResponse, err := uc.paymentGateway.ProcessPayment(payment)
	if err != nil {
		// Update payment status to failed
		uc.paymentRepo.UpdatePaymentFailed(paymentID, err.Error())

		// Publish payment.failed event
		failedEvent := entities.Event{
			ID:        uuid.New().String(),
			Type:      "payment.failed",
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"payment_id": payment.ID,
				"order_id":   payment.OrderID,
				"reason":     err.Error(),
			},
		}
		uc.eventPublisher.PublishEvent(failedEvent)

		return payment, err
	}

	// Update payment with transaction ID
	err = uc.paymentRepo.UpdateStatus(paymentID, entities.PaymentStatusCompleted, gatewayResponse.TransactionID)
	if err != nil {
		return nil, err
	}

	// Update payment object
	payment.Status = entities.PaymentStatusCompleted
	payment.TransactionID = gatewayResponse.TransactionID

	// Publish payment.completed event
	completedEvent := entities.Event{
		ID:        uuid.New().String(),
		Type:      "payment.completed",
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"payment_id":     payment.ID,
			"order_id":       payment.OrderID,
			"transaction_id": payment.TransactionID,
		},
	}

	uc.eventPublisher.PublishEvent(completedEvent)

	return payment, nil
}
