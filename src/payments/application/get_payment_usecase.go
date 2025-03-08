package application

import (
	"app.payment/src/payments/domain/entities"
	"app.payment/src/payments/domain/repositories"
)

type GetPaymentUseCase struct {
	paymentRepo repositories.PaymentRepository
}

func NewGetPaymentUseCase(paymentRepo repositories.PaymentRepository) *GetPaymentUseCase {
	return &GetPaymentUseCase{paymentRepo: paymentRepo}
}

func (uc *GetPaymentUseCase) ExecuteByID(id string) (*entities.Payment, error) {
	return uc.paymentRepo.FindByID(id)
}

func (uc *GetPaymentUseCase) ExecuteByOrderID(orderID uint) (*entities.Payment, error) {
	return uc.paymentRepo.FindByOrderID(orderID)
}
